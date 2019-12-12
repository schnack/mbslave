package mbslave

import (
	"encoding/binary"
	"fmt"
)

type RtuResponse struct {
	SlaveId  uint8
	Function uint8
	Address  uint16
	Err      uint8
	Data     []byte
}

func NewRtuResponse(slaveId, function uint8, address uint16, data []byte, errCode uint8) Response {
	return &RtuResponse{
		SlaveId:  slaveId,
		Function: function,
		Address:  address,
		Err:      errCode,
		Data:     data,
	}
}

func (rr *RtuResponse) GetSlaveId() uint8 {
	return rr.SlaveId
}

func (rr *RtuResponse) GetFunction() uint8 {
	return rr.Function
}

func (rr *RtuResponse) GetAddress() uint16 {
	return rr.Address
}

func (rr *RtuResponse) GetError() uint8 {
	return rr.Err
}

func (rr *RtuResponse) GetData() []byte {
	return rr.Data
}

func (rr *RtuResponse) GetADU() (b []byte, err error) {
	address := make([]byte, 2)
	binary.BigEndian.PutUint16(address, rr.Address)

	b = append(b, rr.SlaveId)
	b = append(b, rr.Function)
	switch rr.Function {
	case FuncReadCoils, FuncReadDiscreteInputs, FuncReadInputRegisters, FuncReadHoldingRegisters:
		b = append(b, uint8(len(rr.Data)))
		if len(rr.Data) > 0 {
			b = append(b, rr.Data...)
		} else {
			return nil, fmt.Errorf("there is no data to answer")
		}
	case FuncWriteSingleCoil, FuncWriteSingleRegister, FuncWriteMultipleCoils, FuncWriteMultipleRegisters:
		b = append(b, address...)
		if len(rr.Data) > 1 {
			b = append(b, rr.Data[0:2]...)
		} else {
			return nil, fmt.Errorf("there is no data to answer")
		}
	case ExceptionFunction(FuncReadCoils),
		ExceptionFunction(FuncReadDiscreteInputs),
		ExceptionFunction(FuncReadHoldingRegisters),
		ExceptionFunction(FuncReadInputRegisters),
		ExceptionFunction(FuncWriteSingleCoil),
		ExceptionFunction(FuncWriteSingleRegister),
		ExceptionFunction(FuncWriteMultipleCoils),
		ExceptionFunction(FuncWriteMultipleRegisters):
		if rr.Err != 0 {
			b = append(b, rr.Err)
		} else {
			return nil, fmt.Errorf("the error cannot be 0")
		}
	default:
		b = append(b, rr.Data...)
	}

	crc := make([]byte, 2)
	binary.LittleEndian.PutUint16(crc, CalcCRC(b))
	b = append(b, crc...)
	return
}
