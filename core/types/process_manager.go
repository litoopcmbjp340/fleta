package types

// ProcessManager is a interface to give a chain data
type ProcessManager interface {
	Processes() []Process
	Process(id uint8) (Process, error)
	ProcessByName(name string) (Process, error)
}