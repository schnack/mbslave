package mbslave

import (
	"encoding/binary"
	"fmt"
)

type RtuResponse struct {
	slaveId  uint8
	function uint8
	address  uint16
	err      uint8
	data     []byte
	crc      uint16

	unanswered bool
}

func NewRtuResponse(request Request) Response {
	return &RtuResponse{
		slaveId:  request.GetSlaveId(),
		function: request.GetFunction(),
	}
}

func (rr *RtuResponse) Unanswered(on bool) {
	rr.unanswered = on
}

func (rr *RtuResponse) SetError(err uint8) {
	rr.function = ExceptionFunction(rr.function)
	rr.err = err
}

func (rr *RtuResponse) SetRead(data []byte) {
	rr.data = data
}

func (rr *RtuResponse) SetSingleWrite(address uint16, data []byte) {
	rr.address = address
	rr.data = data
}

func (rr *RtuResponse) SetMultiWrite(address uint16, countReg uint16) {
	rr.address = address
	rr.data = make([]byte, 2)
	binary.BigEndian.PutUint16(rr.data, countReg)
}

func (rr *RtuResponse) GetSlaveId() uint8 {
	return rr.slaveId
}

func (rr *RtuResponse) GetFunction() uint8 {
	return rr.function
}

func (rr *RtuResponse) GetAddress() uint16 {
	return rr.address
}

func (rr *RtuResponse) GetError() uint8 {
	return rr.err
}

func (rr *RtuResponse) GetData() []byte {
	return rr.data
}

func (rr *RtuResponse) GetADU() (b []byte, err error) {
	if rr.unanswered {
		return nil, fmt.Errorf("not answer")
	}

	address := make([]byte, 2)
	binary.BigEndian.PutUint16(address, rr.address)

	b = append(b, rr.slaveId)
	b = append(b, rr.function)
	switch rr.function {
	case FuncReadCoils, FuncReadDiscreteInputs, FuncReadInputRegisters, FuncReadHoldingRegisters:
		b = append(b, uint8(len(rr.data)))
		if len(rr.data) > 0 {
			b = append(b, rr.data...)
		} else {
			return nil, fmt.Errorf("there is no data to answer")
		}
	case FuncWriteSingleCoil, FuncWriteSingleRegister, FuncWriteMultipleCoils, FuncWriteMultipleRegisters:
		b = append(b, address...)
		if len(rr.data) > 1 {
			b = append(b, rr.data[0:2]...)
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
		if rr.err != 0 {
			b = append(b, rr.err)
		} else {
			return nil, fmt.Errorf("the error cannot be 0")
		}
	default:
		b = append(b, rr.data...)
	}

	crc := make([]byte, 2)
	binary.LittleEndian.PutUint16(crc, CalcCRC(b))
	b = append(b, crc...)
	return
}
