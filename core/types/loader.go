package types

import (
	"github.com/fletaio/fleta/common"
	"github.com/fletaio/fleta/common/hash"
)

// Loader is an interface to load state data from the target chain
type Loader interface {
	Name() string
	Version() uint16
	TargetHeight() uint32
	LastHash() hash.Hash256
	LastTimestamp() uint64
	Seq(addr common.Address) uint64
	Account(addr common.Address) (Account, error)
	AddressByName(Name string) (common.Address, error)
	IsExistAccount(addr common.Address) (bool, error)
	IsExistAccountName(Name string) (bool, error)
	AccountData(addr common.Address, name []byte) []byte
	AccountDataKeys(addr common.Address, Prefix []byte) ([][]byte, error)
	IsExistUTXO(id uint64) (bool, error)
	UTXO(id uint64) (*UTXO, error)
}

type internalLoader interface {
	Loader
	ProcessData(pid uint8, name []byte) []byte
	ProcessDataKeys(pid uint8, Prefix []byte) ([][]byte, error)
}

type emptyLoader struct {
	coord *common.Coordinate
}

// NewEmptyLoader is used for generating genesis state
func NewEmptyLoader() internalLoader {
	return &emptyLoader{}
}

// Name returns ""
func (st *emptyLoader) Name() string {
	return ""
}

// Version returns 0
func (st *emptyLoader) Version() uint16 {
	return 0
}

// TargetHeight returns 0
func (st *emptyLoader) TargetHeight() uint32 {
	return 0
}

// LastHash returns hash.Hash256{}
func (st *emptyLoader) LastHash() hash.Hash256 {
	return hash.Hash256{}
}

// LastTimestamp returns 0
func (st *emptyLoader) LastTimestamp() uint64 {
	return 0
}

// Seq returns 0
func (st *emptyLoader) Seq(addr common.Address) uint64 {
	return 0
}

// Account returns ErrNotExistAccount
func (st *emptyLoader) Account(addr common.Address) (Account, error) {
	return nil, ErrNotExistAccount
}

// AddressByName returns ErrNotExistAccount
func (st *emptyLoader) AddressByName(Name string) (common.Address, error) {
	return common.Address{}, ErrNotExistAccount
}

// IsExistAccount returns false
func (st *emptyLoader) IsExistAccount(addr common.Address) (bool, error) {
	return false, nil
}

// IsExistAccountName returns false
func (st *emptyLoader) IsExistAccountName(Name string) (bool, error) {
	return false, nil
}

// AccountDataKeys returns nil
func (st *emptyLoader) AccountDataKeys(addr common.Address, Prefix []byte) ([][]byte, error) {
	return nil, nil
}

// AccountData returns nil
func (st *emptyLoader) AccountData(addr common.Address, name []byte) []byte {
	return nil
}

// IsExistUTXO returns false
func (st *emptyLoader) IsExistUTXO(id uint64) (bool, error) {
	return false, nil
}

// UTXO returns ErrNotExistUTXO
func (st *emptyLoader) UTXO(id uint64) (*UTXO, error) {
	return nil, ErrNotExistUTXO
}

// ProcessDataKeys returns nil
func (st *emptyLoader) ProcessDataKeys(pid uint8, Prefix []byte) ([][]byte, error) {
	return nil, nil
}

// ProcessData returns nil
func (st *emptyLoader) ProcessData(pid uint8, name []byte) []byte {
	return nil
}