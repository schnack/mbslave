package mbslave

import (
	"encoding/binary"
	"fmt"
)

type RtuRequest struct {
	SlaveId   uint8
	Function  uint8
	Address   uint16
	Quantity  uint16
	CountByte uint8
	Data      []byte
	CRC       uint16
	raw       []byte
}

func NewRtuRequest(b []byte) Request {
	request := &RtuRequest{raw: b}
	return request
}

func (rr *RtuRequest) Parse() error {
	countFrame := len(rr.raw)
	if countFrame < 4 {
		return fmt.Errorf("frame damaged")
	}

	rr.SlaveId = rr.raw[0]
	rr.Function = rr.raw[1]

	switch rr.Function {
	case FuncWriteSingleCoil, FuncWriteSingleRegister:
		if countFrame != 8 {
			return fmt.Errorf("frame damaged")
		}
		rr.Address = binary.BigEndian.Uint16(rr.raw[2:4])
		rr.Quantity = 1
		rr.Data = rr.raw[4:6]
		rr.CRC = binary.LittleEndian.Uint16(rr.raw[6:8])

	case FuncWriteMultipleCoils, FuncWriteMultipleRegisters:
		if countFrame < 7 {
			return fmt.Errorf("frame damaged")
		}
		rr.Address = binary.BigEndian.Uint16(rr.raw[2:4])
		rr.Quantity = binary.BigEndian.Uint16(rr.raw[4:6])
		rr.CountByte = rr.raw[6]

		if countFrame != (9 + int(rr.CountByte)) {
			return fmt.Errorf("frame damaged")
		}
		rr.Data = rr.raw[7 : 7+int(rr.CountByte)]

		rr.CRC = binary.LittleEndian.Uint16(rr.raw[7+int(rr.CountByte) : 7+int(rr.CountByte)+2])

	case FuncReadDiscreteInputs, FuncReadCoils, FuncReadInputRegisters, FuncReadHoldingRegisters:
		if countFrame != 8 {
			return fmt.Errorf("frame damaged")
		}
		rr.Address = binary.BigEndian.Uint16(rr.raw[2:4])
		rr.Quantity = binary.BigEndian.Uint16(rr.raw[4:6])
		rr.CRC = binary.LittleEndian.Uint16(rr.raw[6:8])

	default:
		rr.Data = rr.raw[2 : len(rr.raw)-2]
		rr.CRC = binary.LittleEndian.Uint16(rr.raw[len(rr.raw)-2:])
	}

	if err := rr.Validate(); err != nil {
		return err
	}
	return nil
}

// GetSlaveId - returns the address of the device even if the ADU is not parsed
func (rr *RtuRequest) GetSlaveId() uint8 {
	if len(rr.raw) > 0 && rr.Function == 0 {
		return rr.raw[0]
	}
	return rr.SlaveId
}

// GetFunction - returns the function
func (rr *RtuRequest) GetFunction() uint8 {
	if len(rr.raw) > 1 && rr.Function == 0 {
		return rr.raw[1]
	}
	return rr.Function
}

func (rr *RtuRequest) GetAddress() uint16 {
	return rr.Address
}

func (rr *RtuRequest) GetQuantity() uint16 {
	return rr.Quantity
}

func (rr *RtuRequest) GetCountByte() uint8 {
	return rr.CountByte
}

func (rr *RtuRequest) GetData() []byte {
	return rr.Data
}

func (rr *RtuRequest) GetCrc() uint16 {
	return rr.CRC
}

func (rr *RtuRequest) Validate() error {
	calc := CalcCRC(rr.raw[:len(rr.raw)-2])
	if rr.GetCrc() != calc {
		return fmt.Errorf("crc: 0x%04x, calc: 0x%04x", rr.GetCrc(), calc)
	}
	return nil
}

func (rr *RtuRequest) GetADU() []byte {
	return rr.raw
}
