//Пример наследования и переопределения функций
//
//		type UserDataModel struct {
//			DefaultDataModel
//		}
//
//		func (ud *UserDataModel)Init() error {
//			ud.DefaultDataModel.Init()
//			ud.SetFunction(0x01, ud.ReadCoils)
//			return nil
//		}
package mbslave

import (
	"math"
)

const (
	DiscretesInputAddress   = 0
	CoilsAddress            = DiscretesInputAddress + (math.MaxUint16+1)/8
	InputRegistersAddress   = CoilsAddress + (math.MaxUint16+1)/8
	HoldingRegistersAddress = InputRegistersAddress + (math.MaxUint16+1)*2
	ModBusDataEnd           = HoldingRegistersAddress + (math.MaxUint16+1)*2

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
	Init()
	Handler(req Request) Response
}

type DefaultDataModel struct {
	SlaveId  uint8
	raw      [ModBusDataEnd]byte
	function [256]func(Request) Response
}

func (dm *DefaultDataModel) Init() {
	dm.SetFunction(FuncReadCoils, dm.ReadCoils)
	dm.SetFunction(FuncReadDiscreteInputs, dm.ReadDiscreteInputs)
	dm.SetFunction(FuncReadHoldingRegisters, dm.ReadHoldingRegisters)
	dm.SetFunction(FuncReadInputRegisters, dm.ReadInputRegisters)
	dm.SetFunction(FuncWriteSingleCoil, dm.WriteSingleCoil)
	dm.SetFunction(FuncWriteSingleRegister, dm.WriteSingleRegister)
	dm.SetFunction(FuncWriteMultipleCoils, dm.WriteMultipleCoils)
	dm.SetFunction(FuncWriteMultipleRegisters, dm.WriteMultipleRegisters)
}

func (dm *DefaultDataModel) SetFunction(code uint8, f func(Request) Response) {
	dm.function[code] = f
}

func (dm *DefaultDataModel) Handler(req Request) (response Response) {

	if req.GetSlaveId() != dm.SlaveId && req.GetSlaveId() != 255 {
		return response
	}

	if err := req.Parse(); err != nil {
		return response
	}

	if dm.function[req.GetFunction()] != nil {
		response = dm.function[req.GetFunction()](req)
	} else {
		response = NewRtuResponse(req.GetSlaveId(), req.GetFunction(), 0, nil, ErrorFunction)
	}

	if req.GetSlaveId() == 255 {
		response = nil
	}

	return
}

func (dm *DefaultDataModel) ReadCoils(request Request) Response {
	return nil
}

func (dm *DefaultDataModel) ReadDiscreteInputs(request Request) Response {
	return nil
}

func (dm *DefaultDataModel) ReadHoldingRegisters(request Request) Response {
	return nil
}

func (dm *DefaultDataModel) ReadInputRegisters(request Request) Response {
	return nil
}

func (dm *DefaultDataModel) WriteSingleCoil(request Request) Response {
	return nil
}

func (dm *DefaultDataModel) WriteSingleRegister(request Request) Response {
	return nil
}

func (dm *DefaultDataModel) WriteMultipleCoils(request Request) Response {
	return nil
}

func (dm *DefaultDataModel) WriteMultipleRegisters(request Request) Response {
	return nil
}
