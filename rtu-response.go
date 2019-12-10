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

func (rr *RtuResponse) GetRaw() (b []byte, err error) {
	address := make([]byte, 2)
	binary.BigEndian.PutUint16(address, rr.Address)

	b = append(b, rr.SlaveId)
	b = append(b, rr.Function)
	switch rr.Function {
	case FuncReadCoils, FuncReadDiscreteInputs, FuncReadInputRegisters, FuncReadHoldingRegisters:
		b = append(b, uint8(len(rr.Data)))
		b = append(b, rr.Data...)
	case FuncWriteSingleCoil, FuncWriteSingleRegister, FuncWriteMultipleCoils, FuncWriteMultipleRegisters:
		b = append(b, address...)
		b = append(b, rr.Data[0:2]...)
	case ExceptionReadCoils, ExceptionReadDiscreteInputs, ExceptionReadHoldingRegisters, ExceptionReadInputRegisters,
		ExceptionWriteSingleCoil, ExceptionWriteSingleRegister, ExceptionWriteMultipleCoils, ExceptionWriteMultipleRegisters:
		b = append(b, rr.Err)
	default:
		return nil, fmt.Errorf("sorry, this function is not supported")
	}

	crc := make([]byte, 2)
	binary.LittleEndian.PutUint16(crc, CalcCRC(b))
	b = append(b, crc...)
	return
}
