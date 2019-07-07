package pof

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io/ioutil"
	"net"
	"sync"
	"time"

	"github.com/fletaio/fleta/common"
	"github.com/fletaio/fleta/common/queue"
	"github.com/fletaio/fleta/common/util"
	"github.com/fletaio/fleta/encoding"
	"github.com/fletaio/fleta/service/p2p"
)

// FormulatorPeer is a formulator peer
type FormulatorPeer struct {
	sync.Mutex
	id          string
	netAddr     string
	conn        net.Conn
	pubhash     common.PublicHash
	address     common.Address
	guessHeight uint32
	writeQueue  *queue.Queue
	isClose     bool
}

// NewFormulatorPeer returns a ormulatorPeer
func NewFormulatorPeer(conn net.Conn, pubhash common.PublicHash, address common.Address) *FormulatorPeer {
	p := &FormulatorPeer{
		id:         address.String(),
		netAddr:    conn.RemoteAddr().String(),
		conn:       conn,
		pubhash:    pubhash,
		address:    address,
		writeQueue: queue.NewQueue(),
	}
	go func() {
		defer p.conn.Close()

		for {
			if p.isClose {
				return
			}
			v := p.writeQueue.Pop()
			if v == nil {
				time.Sleep(50 * time.Millisecond)
				continue
			}
			bs := v.([]byte)
			var buffer bytes.Buffer
			buffer.Write(bs[:2])
			buffer.Write(make([]byte, 4))
			if len(bs) > 2 {
				zw := gzip.NewWriter(&buffer)
				zw.Write(bs[2:])
				zw.Flush()
				zw.Close()
			}
			wbs := buffer.Bytes()
			binary.LittleEndian.PutUint32(wbs[2:], uint32(len(wbs)-6))
			if err := p.conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
				return
			}
			_, err := p.conn.Write(wbs)
			if err != nil {
				return
			}
		}
	}()
	return p
}

// Close closes peer
func (p *FormulatorPeer) Close() {
	p.conn.Close()
	p.isClose = true
}

// ID returns the id of the peer
func (p *FormulatorPeer) ID() string {
	return p.id
}

// NetAddr returns the network address of the peer
func (p *FormulatorPeer) NetAddr() string {
	return p.netAddr
}

// Address returns the formulator address of the peer
func (p *FormulatorPeer) Address() common.Address {
	return p.address
}

// ReadMessageData returns a message data
func (p *FormulatorPeer) ReadMessageData() (interface{}, []byte, error) {
	var t uint16
	if v, _, err := p2p.ReadUint16(p.conn); err != nil {
		return nil, nil, err
	} else {
		t = v
	}

	if Len, _, err := p2p.ReadUint32(p.conn); err != nil {
		return nil, nil, err
	} else if Len == 0 {
		return nil, nil, nil
	} else {
		zbs := make([]byte, Len)
		if _, err := p2p.FillBytes(p.conn, zbs); err != nil {
			return nil, nil, err
		}
		zr, err := gzip.NewReader(bytes.NewReader(zbs))
		if err != nil {
			return nil, nil, err
		}
		defer zr.Close()

		fc := encoding.Factory("pof.message")
		m, err := fc.Create(t)
		if err != nil {
			return nil, nil, err
		}
		bs, err := ioutil.ReadAll(zr)
		if err != nil {
			return nil, nil, err
		}
		if err := encoding.Unmarshal(bs, &m); err != nil {
			return nil, nil, err
		}
		return m, bs, nil
	}
}

// Send sends a message to the peer
func (p *FormulatorPeer) Send(m interface{}) error {
	var buffer bytes.Buffer
	fc := encoding.Factory("pof.message")
	t, err := fc.TypeOf(m)
	if err != nil {
		return err
	}
	buffer.Write(util.Uint16ToBytes(t))
	enc := encoding.NewEncoder(&buffer)
	if err := enc.Encode(m); err != nil {
		return err
	}
	p.SendRaw(buffer.Bytes())
	return nil
}

// SendRaw sends bytes to the peer
func (p *FormulatorPeer) SendRaw(bs []byte) {
	p.writeQueue.Push(bs)
}

// UpdateGuessHeight updates the guess height of the peer
func (p *FormulatorPeer) UpdateGuessHeight(height uint32) {
	p.Lock()
	defer p.Unlock()

	p.guessHeight = height
}

// GuessHeight updates the guess height of the peer
func (p *FormulatorPeer) GuessHeight() uint32 {
	return p.guessHeight
}
