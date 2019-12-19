package mbslave

const (
	FuncReadCoils              = uint8(1)
	FuncReadDiscreteInputs     = uint8(2)
	FuncReadHoldingRegisters   = uint8(3)
	FuncReadInputRegisters     = uint8(4)
	FuncWriteSingleCoil        = uint8(5)
	FuncWriteSingleRegister    = uint8(6)
	FuncWriteMultipleCoils     = uint8(15)
	FuncWriteMultipleRegisters = uint8(16)

	ErrorFunction = uint8(1)
	ErrorAddress  = uint8(2)
	ErrorData     = uint8(3)
	ErrorFatal    = uint8(4)
	ErrorDelay    = uint8(5)
	ErrorWait     = uint8(6)
	ErrorFail     = uint8(7)
)

type DataModel interface {
	Handler(req Request, resp Response)
	SetFunction(code uint8, f func(Request, Response))
	SetSlaveId(uint8)
}
