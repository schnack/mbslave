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

func NewRtuRequest(b []byte) (Request, error) {
	request := &RtuRequest{raw: b}
	err := request.parse()
	return request, err
}

func (rr *RtuRequest) parse() error {
	countFrame := len(rr.raw)
	if countFrame < 4 {
		return fmt.Errorf("frame damaged")
	}

	rr.SlaveId = rr.raw[0]
	rr.Function = rr.raw[1]
	rr.Address = binary.BigEndian.Uint16(rr.raw[2:4])

	switch rr.Function {
	case FuncWriteSingleCoil, FuncWriteSingleRegister:
		if countFrame != 8 {
			return fmt.Errorf("frame damaged")
		}
		rr.Quantity = 1
		rr.Data = rr.raw[4:6]
		rr.CRC = binary.LittleEndian.Uint16(rr.raw[6:8])

	case FuncWriteMultipleCoils, FuncWriteMultipleRegisters:
		if countFrame < 7 {
			return fmt.Errorf("frame damaged")
		}
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
		rr.Quantity = binary.BigEndian.Uint16(rr.raw[4:6])
		rr.CRC = binary.LittleEndian.Uint16(rr.raw[6:8])

	default:
		return fmt.Errorf("function not found")
	}

	if err := rr.Validate(); err != nil {
		return err
	}
	return nil
}

func (rr *RtuRequest) GetSlaveId() uint8 {
	return rr.SlaveId
}

func (rr *RtuRequest) GetFunction() uint8 {
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
