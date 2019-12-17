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
	"fmt"
	"math"
	"sync"
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
	SetDiscreteInputs(address uint16, value bool) error
	SetCoils(address uint16, value bool) error
	SetHoldingRegisters(address uint16, value uint16) error
	SetInputRegisters(address uint16, value uint16) error
	GetDiscreteInputs(address uint16) bool
	GetCoils(address uint16) bool
	GetHoldingRegisters(address uint16) uint16
	GetInputRegisters(address uint16) uint16
	LengthDiscreteInputs() int
	LengthCoils() int
	LengthInputRegisters() int
	LengthHoldingRegisters() int
}

type DefaultDataModel struct {
	SlaveId            uint8
	discreteInputs     []bool
	coils              []bool
	inputRegisters     []uint16
	holdingRegisters   []uint16
	function           [256]func(Request, Response)
	muDiscreteInputs   sync.RWMutex
	muCoils            sync.RWMutex
	muInputRegisters   sync.RWMutex
	muHoldingRegisters sync.RWMutex
}

//TODO tests

func (dm *DefaultDataModel) LengthDiscreteInputs() int {
	return len(dm.discreteInputs)
}

func (dm *DefaultDataModel) LengthCoils() int {
	return len(dm.coils)
}

func (dm *DefaultDataModel) LengthInputRegisters() int {
	return len(dm.inputRegisters)
}

func (dm *DefaultDataModel) LengthHoldingRegisters() int {
	return len(dm.holdingRegisters)
}

func (dm *DefaultDataModel) SetDiscreteInputs(address uint16, value bool) error {
	dm.muDiscreteInputs.Lock()
	defer dm.muDiscreteInputs.Unlock()
	if len(dm.discreteInputs) > int(address) {
		return fmt.Errorf("there is no register at this address")
	}
	dm.discreteInputs[int(address)] = value
	return nil
}

func (dm *DefaultDataModel) SetCoils(address uint16, value bool) error {
	dm.muCoils.Lock()
	defer dm.muCoils.Unlock()
	if len(dm.coils) > int(address) {
		return fmt.Errorf("there is no register at this address")
	}
	dm.coils[int(address)] = value
	return nil
}

func (dm *DefaultDataModel) SetHoldingRegisters(address uint16, value uint16) error {
	dm.muHoldingRegisters.Lock()
	defer dm.muHoldingRegisters.Unlock()
	if len(dm.holdingRegisters) > int(address) {
		return fmt.Errorf("there is no register at this address")
	}
	dm.holdingRegisters[int(address)] = value
	return nil
}

func (dm *DefaultDataModel) SetInputRegisters(address uint16, value uint16) error {
	dm.muInputRegisters.Lock()
	defer dm.muInputRegisters.Unlock()
	if len(dm.inputRegisters) > int(address) {
		return fmt.Errorf("there is no register at this address")
	}
	dm.inputRegisters[int(address)] = value
	return nil
}

func (dm *DefaultDataModel) GetDiscreteInputs(address uint16) bool {
	dm.muDiscreteInputs.RLock()
	defer dm.muDiscreteInputs.RUnlock()
	if len(dm.discreteInputs) > int(address) {
		return false
	}
	return dm.discreteInputs[int(address)]
}

func (dm *DefaultDataModel) GetCoils(address uint16) bool {
	dm.muCoils.RLock()
	defer dm.muCoils.RUnlock()
	if len(dm.coils) > int(address) {
		return false
	}
	return dm.coils[int(address)]
}

func (dm *DefaultDataModel) GetHoldingRegisters(address uint16) uint16 {
	dm.muHoldingRegisters.RLock()
	defer dm.muHoldingRegisters.RUnlock()
	if len(dm.holdingRegisters) > int(address) {
		return 0
	}
	return dm.holdingRegisters[int(address)]
}

func (dm *DefaultDataModel) GetInputRegisters(address uint16) uint16 {
	dm.muInputRegisters.Lock()
	defer dm.muInputRegisters.Unlock()
	if len(dm.inputRegisters) > int(address) {
		return 0
	}
	return dm.inputRegisters[int(address)]
}

func NewDefaultDataModel(slaveId uint8) *DefaultDataModel {
	return &DefaultDataModel{
		SlaveId: slaveId,
	}
}

func (dm *DefaultDataModel) Init() {
	dm.discreteInputs = make([]bool, math.MaxUint16)
	dm.coils = make([]bool, math.MaxUint16)
	dm.inputRegisters = make([]uint16, math.MaxUint16)
	dm.holdingRegisters = make([]uint16, math.MaxUint16)
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
	dm.muCoils.RLock()
	defer dm.muCoils.RUnlock()
	dm.read1bit(dm.coils, request, resp)
}

func (dm *DefaultDataModel) ReadDiscreteInputs(request Request, resp Response) {
	dm.muDiscreteInputs.RLock()
	defer dm.muDiscreteInputs.RUnlock()
	dm.read1bit(dm.discreteInputs, request, resp)
}

func (dm *DefaultDataModel) ReadHoldingRegisters(request Request, resp Response) {
	dm.muHoldingRegisters.RLock()
	defer dm.muHoldingRegisters.RUnlock()
	dm.read16bit(dm.holdingRegisters, request, resp)
}

func (dm *DefaultDataModel) ReadInputRegisters(request Request, resp Response) {
	dm.muInputRegisters.RLock()
	defer dm.muInputRegisters.RUnlock()
	dm.read16bit(dm.inputRegisters, request, resp)
}

func (dm *DefaultDataModel) WriteSingleCoil(request Request, resp Response) {
	dm.muCoils.Lock()
	defer dm.muCoils.Unlock()
	if binary.BigEndian.Uint16(request.GetData()) != 0 {
		dm.coils[int(request.GetAddress())] = true
	} else {
		dm.coils[int(request.GetAddress())] = false
	}
	resp.SetSingleWrite(request.GetAddress(), request.GetData())
}

func (dm *DefaultDataModel) WriteSingleRegister(request Request, resp Response) {
	dm.muHoldingRegisters.Lock()
	defer dm.muHoldingRegisters.Unlock()
	dm.holdingRegisters[int(request.GetAddress())] = binary.BigEndian.Uint16(request.GetData())
	resp.SetSingleWrite(request.GetAddress(), request.GetData())
}

func (dm *DefaultDataModel) WriteMultipleCoils(request Request, resp Response) {
	dm.muCoils.Lock()
	defer dm.muCoils.Unlock()
	endAddress := uint32(request.GetAddress()) + uint32(request.GetQuantity())
	if endAddress >= uint32(len(dm.coils)) {
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
				dm.coils[targetAddress] = true
			} else {
				dm.coils[targetAddress] = false
			}
		}
	}
	resp.SetMultiWrite(request.GetAddress(), request.GetQuantity())
}

func (dm *DefaultDataModel) WriteMultipleRegisters(request Request, resp Response) {
	dm.muHoldingRegisters.Lock()
	defer dm.muHoldingRegisters.Unlock()
	endAddress := uint32(request.GetAddress()) + uint32(request.GetQuantity())
	if endAddress >= uint32(len(dm.coils)) {
		resp.SetError(ErrorAddress)
		return
	}
	if len(request.GetData())%2 != 0 {
		resp.SetError(ErrorData)
	}

	for i := 0; i <= int(request.GetQuantity()); i++ {
		dm.holdingRegisters[int(request.GetAddress())+i] = binary.BigEndian.Uint16(request.GetData()[i*2 : (i+1)*2])
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
