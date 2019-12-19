package mbslave

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
