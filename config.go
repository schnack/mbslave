package mbslave

import "go.bug.st/serial"

type Config struct {
	Port string
	serial.Mode

	SlaveId              uint8
	SizeDiscreteInputs   uint16
	SizeCoils            uint16
	SizeInputRegisters   uint16
	SizeHoldingRegisters uint16
}
