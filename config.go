package mbslave

import (
	"go.bug.st/serial"
	"time"
)

const (
	// NoParity disable parity control (default)
	NoParity = serial.NoParity
	// OddParity enable odd-parity check
	OddParity = serial.OddParity
	// EvenParity enable even-parity check
	EvenParity = serial.EvenParity
	// MarkParity enable mark-parity (always 1) check
	MarkParity = serial.MarkParity
	// SpaceParity enable space-parity (always 0) check
	SpaceParity = serial.SpaceParity
)

const (
	// OneStopBit sets 1 stop bit (default)
	OneStopBit = serial.OneStopBit
	// OnePointFiveStopBits sets 1.5 stop bits
	OnePointFiveStopBits = serial.OnePointFiveStopBits
	// TwoStopBits sets 2 stop bits
	TwoStopBits = serial.TwoStopBits
)

type Config struct {
	Port     string
	BaudRate int
	DataBits int
	Parity   serial.Parity
	StopBits serial.StopBits
	// Интервал между adu
	SilentInterval time.Duration

	SlaveId              uint8
	SizeDiscreteInputs   uint16
	SizeCoils            uint16
	SizeInputRegisters   uint16
	SizeHoldingRegisters uint16
}
