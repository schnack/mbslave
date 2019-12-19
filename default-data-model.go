package mbslave

import (
	"encoding/binary"
	"fmt"
	"math"
	"sync"
)

type Event int

const (
	EventRead = Event(iota)
	EventWrite
)

type DefaultDataModel struct {
	BaseDataModel
	discreteInputs   []bool
	coils            []bool
	inputRegisters   []uint16
	holdingRegisters []uint16

	callbackDiscreteInputs   []func(event Event, addr uint16, value bool)
	callbackCoils            []func(event Event, addr uint16, value bool)
	callbackInputRegisters   []func(event Event, addr uint16, value uint16)
	callbackHoldingRegisters []func(event Event, addr uint16, value uint16)

	muDiscreteInputs   sync.RWMutex
	muCoils            sync.RWMutex
	muInputRegisters   sync.RWMutex
	muHoldingRegisters sync.RWMutex
}

func NewDefaultDataModel(slaveId uint8) *DefaultDataModel {
	dm := &DefaultDataModel{
		discreteInputs:           make([]bool, math.MaxUint16),
		coils:                    make([]bool, math.MaxUint16),
		inputRegisters:           make([]uint16, math.MaxUint16),
		holdingRegisters:         make([]uint16, math.MaxUint16),
		callbackDiscreteInputs:   make([]func(event Event, addr uint16, value bool), math.MaxUint16),
		callbackCoils:            make([]func(event Event, addr uint16, value bool), math.MaxUint16),
		callbackInputRegisters:   make([]func(event Event, addr uint16, value uint16), math.MaxUint16),
		callbackHoldingRegisters: make([]func(event Event, addr uint16, value uint16), math.MaxUint16),
	}
	dm.SetSlaveId(slaveId)
	dm.SetFunction(FuncReadCoils, dm.ReadCoils)
	dm.SetFunction(FuncReadDiscreteInputs, dm.ReadDiscreteInputs)
	dm.SetFunction(FuncReadHoldingRegisters, dm.ReadHoldingRegisters)
	dm.SetFunction(FuncReadInputRegisters, dm.ReadInputRegisters)
	dm.SetFunction(FuncWriteSingleCoil, dm.WriteSingleCoil)
	dm.SetFunction(FuncWriteSingleRegister, dm.WriteSingleRegister)
	dm.SetFunction(FuncWriteMultipleCoils, dm.WriteMultipleCoils)
	dm.SetFunction(FuncWriteMultipleRegisters, dm.WriteMultipleRegisters)
	return dm
}

func (dm *DefaultDataModel) SetCallbackDiscreteInputs(addr uint16, f func(event Event, addr uint16, value bool)) {
	dm.callbackDiscreteInputs[addr] = f
}

func (dm *DefaultDataModel) SetCallbackCoils(addr uint16, f func(event Event, addr uint16, value bool)) {
	dm.callbackCoils[addr] = f
}

func (dm *DefaultDataModel) SetCallbackInputRegisters(addr uint16, f func(event Event, addr uint16, value uint16)) {
	dm.callbackInputRegisters[addr] = f
}

func (dm *DefaultDataModel) SetCallbackHoldingRegisters(addr uint16, f func(event Event, addr uint16, value uint16)) {
	dm.callbackHoldingRegisters[addr] = f
}

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
	if len(dm.discreteInputs) <= int(address) {
		return fmt.Errorf("there is no register at this address")
	}
	dm.discreteInputs[int(address)] = value
	if len(dm.callbackDiscreteInputs) > int(address) && dm.callbackDiscreteInputs[int(address)] != nil {
		go dm.callbackDiscreteInputs[int(address)](EventWrite, address, value)
	}
	return nil
}

func (dm *DefaultDataModel) SetCoils(address uint16, value bool) error {
	dm.muCoils.Lock()
	defer dm.muCoils.Unlock()
	if len(dm.coils) <= int(address) {
		return fmt.Errorf("there is no register at this address")
	}
	dm.coils[int(address)] = value
	if len(dm.callbackCoils) > int(address) && dm.callbackCoils[int(address)] != nil {
		go dm.callbackCoils[int(address)](EventWrite, address, value)
	}
	return nil
}

func (dm *DefaultDataModel) SetHoldingRegisters(address uint16, value uint16) error {
	dm.muHoldingRegisters.Lock()
	defer dm.muHoldingRegisters.Unlock()
	if len(dm.holdingRegisters) <= int(address) {
		return fmt.Errorf("there is no register at this address")
	}
	dm.holdingRegisters[int(address)] = value
	if len(dm.callbackHoldingRegisters) > int(address) && dm.callbackHoldingRegisters[int(address)] != nil {
		go dm.callbackHoldingRegisters[int(address)](EventWrite, address, value)
	}
	return nil
}

func (dm *DefaultDataModel) SetInputRegisters(address uint16, value uint16) error {
	dm.muInputRegisters.Lock()
	defer dm.muInputRegisters.Unlock()
	if len(dm.inputRegisters) <= int(address) {
		return fmt.Errorf("there is no register at this address")
	}
	dm.inputRegisters[int(address)] = value
	if len(dm.callbackInputRegisters) > int(address) && dm.callbackInputRegisters[int(address)] != nil {
		go dm.callbackInputRegisters[int(address)](EventWrite, address, value)
	}
	return nil
}

func (dm *DefaultDataModel) GetDiscreteInputs(address uint16) bool {
	dm.muDiscreteInputs.RLock()
	defer dm.muDiscreteInputs.RUnlock()
	if len(dm.discreteInputs) <= int(address) {
		return false
	}
	if len(dm.callbackDiscreteInputs) > int(address) && dm.callbackDiscreteInputs[int(address)] != nil {
		go dm.callbackDiscreteInputs[int(address)](EventRead, address, dm.discreteInputs[int(address)])
	}
	return dm.discreteInputs[int(address)]
}

func (dm *DefaultDataModel) GetCoils(address uint16) bool {
	dm.muCoils.RLock()
	defer dm.muCoils.RUnlock()
	if len(dm.coils) <= int(address) {
		return false
	}
	if len(dm.callbackCoils) > int(address) && dm.callbackCoils[int(address)] != nil {
		go dm.callbackCoils[int(address)](EventRead, address, dm.coils[int(address)])
	}
	return dm.coils[int(address)]
}

func (dm *DefaultDataModel) GetHoldingRegisters(address uint16) uint16 {
	dm.muHoldingRegisters.RLock()
	defer dm.muHoldingRegisters.RUnlock()
	if len(dm.holdingRegisters) <= int(address) {
		return 0
	}
	if len(dm.callbackHoldingRegisters) > int(address) && dm.callbackHoldingRegisters[int(address)] != nil {
		go dm.callbackHoldingRegisters[int(address)](EventRead, address, dm.holdingRegisters[int(address)])
	}
	return dm.holdingRegisters[int(address)]
}

func (dm *DefaultDataModel) GetInputRegisters(address uint16) uint16 {
	dm.muInputRegisters.RLock()
	defer dm.muInputRegisters.RUnlock()
	if len(dm.inputRegisters) <= int(address) {
		return 0
	}
	if len(dm.callbackInputRegisters) > int(address) && dm.callbackInputRegisters[int(address)] != nil {
		go dm.callbackInputRegisters[int(address)](EventRead, address, dm.inputRegisters[int(address)])
	}
	return dm.inputRegisters[int(address)]
}

func (dm *DefaultDataModel) ReadCoils(request Request, resp Response) {
	endAddress := uint32(request.GetAddress()) + uint32(request.GetQuantity())
	if endAddress >= uint32(dm.LengthCoils()) {
		resp.SetError(ErrorAddress)
		return
	}

	bufSize := request.GetQuantity() / 8
	if request.GetQuantity()%8 != 0 {
		bufSize++
	}
	buff := make([]byte, bufSize)

	for i := request.GetAddress(); i < uint16(endAddress); i++ {
		if dm.GetCoils(i) {
			index := i - request.GetAddress()
			buff[index/8] |= 1 << (index % 8)
		}
	}
	resp.SetRead(buff)
}

func (dm *DefaultDataModel) ReadDiscreteInputs(request Request, resp Response) {
	endAddress := uint32(request.GetAddress()) + uint32(request.GetQuantity())
	if endAddress >= uint32(dm.LengthCoils()) {
		resp.SetError(ErrorAddress)
		return
	}

	bufSize := request.GetQuantity() / 8
	if request.GetQuantity()%8 != 0 {
		bufSize++
	}
	buff := make([]byte, bufSize)

	for i := request.GetAddress(); i < uint16(endAddress); i++ {
		if dm.GetDiscreteInputs(i) {
			index := i - request.GetAddress()
			buff[index/8] |= 1 << (index % 8)
		}
	}
	resp.SetRead(buff)
}

func (dm *DefaultDataModel) ReadHoldingRegisters(request Request, resp Response) {
	endAddress := uint32(request.GetAddress()) + uint32(request.GetQuantity())
	if endAddress >= uint32(dm.LengthHoldingRegisters()) {
		resp.SetError(ErrorAddress)
		return
	}

	buff := make([]byte, request.GetQuantity()*2)
	for i := request.GetAddress(); i < uint16(endAddress); i++ {
		index := i - request.GetAddress()
		binary.BigEndian.PutUint16(buff[index*2:(index+1)*2], dm.GetHoldingRegisters(i))
	}

	resp.SetRead(buff)
}

func (dm *DefaultDataModel) ReadInputRegisters(request Request, resp Response) {

	endAddress := uint32(request.GetAddress()) + uint32(request.GetQuantity())
	if endAddress >= uint32(dm.LengthInputRegisters()) {
		resp.SetError(ErrorAddress)
		return
	}

	buff := make([]byte, request.GetQuantity()*2)
	for i := request.GetAddress(); i < uint16(endAddress); i++ {
		index := i - request.GetAddress()
		binary.BigEndian.PutUint16(buff[index*2:(index+1)*2], dm.GetInputRegisters(i))
	}

	resp.SetRead(buff)
}

func (dm *DefaultDataModel) WriteSingleCoil(request Request, resp Response) {
	if err := dm.SetCoils(request.GetAddress(), binary.BigEndian.Uint16(request.GetData()) != 0); err != nil {
		resp.SetError(ErrorAddress)
		return
	}
	resp.SetSingleWrite(request.GetAddress(), request.GetData())
}

func (dm *DefaultDataModel) WriteSingleRegister(request Request, resp Response) {
	if err := dm.SetHoldingRegisters(request.GetAddress(), binary.BigEndian.Uint16(request.GetData())); err != nil {
		resp.SetError(ErrorAddress)
		return
	}
	resp.SetSingleWrite(request.GetAddress(), request.GetData())
}

func (dm *DefaultDataModel) WriteMultipleCoils(request Request, resp Response) {
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
			if err := dm.SetCoils(uint16(targetAddress), value>>ii&0x01 == 1); err != nil {
				resp.SetError(ErrorAddress)
				return
			}
		}
	}
	resp.SetMultiWrite(request.GetAddress(), request.GetQuantity())
}

func (dm *DefaultDataModel) WriteMultipleRegisters(request Request, resp Response) {
	endAddress := uint32(request.GetAddress()) + uint32(request.GetQuantity())
	if endAddress >= uint32(len(dm.coils)) {
		resp.SetError(ErrorAddress)
		return
	}
	if len(request.GetData())%2 != 0 {
		resp.SetError(ErrorData)
		return
	}

	for i := 0; i <= int(request.GetQuantity()); i++ {
		if err := dm.SetHoldingRegisters(uint16(int(request.GetAddress())+i), binary.BigEndian.Uint16(request.GetData()[i*2:(i+1)*2])); err != nil {
			resp.SetError(ErrorAddress)
			return
		}
	}
	resp.SetMultiWrite(request.GetAddress(), request.GetQuantity())
}
