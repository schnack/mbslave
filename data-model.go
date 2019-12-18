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

type BaseDataModel struct {
	SlaveId  uint8
	function [256]func(Request, Response)
}

func (bdm *BaseDataModel) SetSlaveId(id uint8) {
	bdm.SlaveId = id
}

func (bdm *BaseDataModel) SetFunction(code uint8, f func(Request, Response)) {
	bdm.function[code] = f
}

func (bdm *BaseDataModel) Handler(req Request, resp Response) {

	if req.GetSlaveId() != bdm.SlaveId && req.GetSlaveId() != 255 {
		resp.Unanswered(true)
		return
	}

	if err := req.Parse(); err != nil {
		resp.Unanswered(true)
		return
	}

	if bdm.function[req.GetFunction()] != nil {
		bdm.function[req.GetFunction()](req, resp)
	} else {
		resp.SetError(ErrorFunction)
	}

	if req.GetSlaveId() == 255 {
		resp.Unanswered(true)
	}
	return
}
