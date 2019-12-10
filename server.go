package mbslave

import (
	"github.com/goburrow/serial"
	"io"
)

const (
	FuncReadCoils                   = uint8(1)
	FuncReadDiscreteInputs          = uint8(2)
	FuncReadHoldingRegisters        = uint8(3)
	FuncReadInputRegisters          = uint8(4)
	FuncWriteSingleCoil             = uint8(5)
	FuncWriteSingleRegister         = uint8(6)
	FuncWriteMultipleCoils          = uint8(15)
	FuncWriteMultipleRegisters      = uint8(16)
	ExceptionReadCoils              = FuncReadCoils | 1<<7
	ExceptionReadDiscreteInputs     = FuncReadDiscreteInputs | 1<<7
	ExceptionReadHoldingRegisters   = FuncReadHoldingRegisters | 1<<7
	ExceptionReadInputRegisters     = FuncReadInputRegisters | 1<<7
	ExceptionWriteSingleCoil        = FuncWriteSingleCoil | 1<<7
	ExceptionWriteSingleRegister    = FuncWriteSingleRegister | 1<<7
	ExceptionWriteMultipleCoils     = FuncWriteMultipleCoils | 1<<7
	ExceptionWriteMultipleRegisters = FuncWriteMultipleRegisters | 1<<7

	ErrorFunction = uint8(1)
	ErrorAddress  = uint8(1)
	ErrorData     = uint8(1)
	ErrorFatal    = uint8(4)
	ErrorDelay    = uint8(5)
	ErrorWait     = uint8(6)
	ErrorFail     = uint8(7)
)

type Server struct {
	config *serial.Config
	port   io.ReadWriteCloser

	DiscreteInputs   []byte
	Coils            []byte
	HoldingRegisters []uint16
	InputRegisters   []uint16
}
