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
	"encoding/binary"
	"math"
)

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
	Init()
	Handler(req Request, resp Response)
}

type DefaultDataModel struct {
	SlaveId          uint8
	DiscreteInputs   []bool
	Coils            []bool
	InputRegisters   []uint16
	HoldingRegisters []uint16
	function         [256]func(Request, Response)
}

func NewDefaultDataModel(slaveId uint8) *DefaultDataModel {
	return &DefaultDataModel{
		SlaveId: slaveId,
	}
}

func (dm *DefaultDataModel) Init() {
	dm.DiscreteInputs = make([]bool, math.MaxUint16)
	dm.Coils = make([]bool, math.MaxUint16)
	dm.InputRegisters = make([]uint16, math.MaxUint16)
	dm.HoldingRegisters = make([]uint16, math.MaxUint16)
	dm.SetFunction(FuncReadCoils, dm.ReadCoils)
	dm.SetFunction(FuncReadDiscreteInputs, dm.ReadDiscreteInputs)
	dm.SetFunction(FuncReadHoldingRegisters, dm.ReadHoldingRegisters)
	dm.SetFunction(FuncReadInputRegisters, dm.ReadInputRegisters)
	dm.SetFunction(FuncWriteSingleCoil, dm.WriteSingleCoil)
	dm.SetFunction(FuncWriteSingleRegister, dm.WriteSingleRegister)
	dm.SetFunction(FuncWriteMultipleCoils, dm.WriteMultipleCoils)
	dm.SetFunction(FuncWriteMultipleRegisters, dm.WriteMultipleRegisters)
}

func (dm *DefaultDataModel) SetFunction(code uint8, f func(Request, Response)) {
	dm.function[code] = f
}

func (dm *DefaultDataModel) Handler(req Request, resp Response) {

	if req.GetSlaveId() != dm.SlaveId && req.GetSlaveId() != 255 {
		resp.Unanswered(true)
		return
	}

	if err := req.Parse(); err != nil {
		resp.Unanswered(true)
		return
	}

	if dm.function[req.GetFunction()] != nil {
		dm.function[req.GetFunction()](req, resp)
	} else {
		resp.SetError(ErrorFunction)
	}

	if req.GetSlaveId() == 255 {
		resp.Unanswered(true)
	}
	return
}

func (dm *DefaultDataModel) ReadCoils(request Request, resp Response) {
	dm.read1bit(dm.Coils, request, resp)
}

func (dm *DefaultDataModel) ReadDiscreteInputs(request Request, resp Response) {
	dm.read1bit(dm.DiscreteInputs, request, resp)
}

func (dm *DefaultDataModel) ReadHoldingRegisters(request Request, resp Response) {
	dm.read16bit(dm.HoldingRegisters, request, resp)
}

func (dm *DefaultDataModel) ReadInputRegisters(request Request, resp Response) {
	dm.read16bit(dm.InputRegisters, request, resp)
}

func (dm *DefaultDataModel) WriteSingleCoil(request Request, resp Response) {
	if binary.BigEndian.Uint16(request.GetData()) != 0 {
		dm.Coils[int(request.GetAddress())] = true
	} else {
		dm.Coils[int(request.GetAddress())] = false
	}
	resp.SetSingleWrite(request.GetAddress(), request.GetData())
}

func (dm *DefaultDataModel) WriteSingleRegister(request Request, resp Response) {
	dm.HoldingRegisters[int(request.GetAddress())] = binary.BigEndian.Uint16(request.GetData())
	resp.SetSingleWrite(request.GetAddress(), request.GetData())
}

func (dm *DefaultDataModel) WriteMultipleCoils(request Request, resp Response) {
	endAddress := uint32(request.GetAddress()) + uint32(request.GetQuantity())
	if endAddress >= uint32(len(dm.Coils)) {
		resp.SetError(ErrorAddress)
		return
	}

	for i, value := range request.GetData() {
		for ii := 0; ii < 8; ii++ {
			targetAddress := int(request.GetAddress()) + i*8 + ii
			if targetAddress > int(endAddress) {
				break
			}
			if value>>ii&0x01 == 1 {
				dm.Coils[targetAddress] = true
			} else {
				dm.Coils[targetAddress] = false
			}
		}
	}
	resp.SetMultiWrite(request.GetAddress(), request.GetQuantity())
}

func (dm *DefaultDataModel) WriteMultipleRegisters(request Request, resp Response) {
	endAddress := uint32(request.GetAddress()) + uint32(request.GetQuantity())
	if endAddress >= uint32(len(dm.Coils)) {
		resp.SetError(ErrorAddress)
		return
	}
	if len(request.GetData())%2 != 0 {
		resp.SetError(ErrorData)
	}

	for i := 0; i <= int(request.GetQuantity()); i++ {
		dm.HoldingRegisters[int(request.GetAddress())+i] = binary.BigEndian.Uint16(request.GetData()[i*2 : (i+1)*2])
	}
	resp.SetMultiWrite(request.GetAddress(), request.GetQuantity())
}

func (dm *DefaultDataModel) read1bit(data []bool, request Request, resp Response) {
	endAddress := uint32(request.GetAddress()) + uint32(request.GetQuantity())
	if endAddress >= uint32(len(data)) {
		resp.SetError(ErrorAddress)
		return
	}

	bufSize := request.GetQuantity() / 8
	if request.GetQuantity()%8 != 0 {
		bufSize++
	}
	buff := make([]byte, bufSize)

	for i, value := range data[request.GetAddress():endAddress] {
		if value {
			buff[i/8] |= 1 << (i % 8)
		}
	}
	resp.SetRead(buff)
}

func (dm *DefaultDataModel) read16bit(data []uint16, request Request, resp Response) {
	endAddress := uint32(request.GetAddress()) + uint32(request.GetQuantity())
	if endAddress >= uint32(len(data)) {
		resp.SetError(ErrorAddress)
		return
	}

	buff := make([]byte, request.GetQuantity()*2)
	for i, value := range data[request.GetAddress():endAddress] {
		binary.BigEndian.PutUint16(buff[i*2:(i+1)*2], value)
	}
	resp.SetRead(buff)
}
