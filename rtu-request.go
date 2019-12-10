package mbslave

import (
	"encoding/binary"
	"fmt"
)

type RtuRequest struct {
	raw []byte
}

func (rr *RtuRequest) GetSlaveId() (uint8, error) {
	if len(rr.raw) > 0 {
		return rr.raw[0], nil
	}
	return 0, fmt.Errorf("slaveId not found")
}

func (rr *RtuRequest) GetFunction() (uint8, error) {
	if len(rr.raw) > 1 {
		switch rr.raw[1] {
		case FuncReadDiscreteInputs, FuncReadCoils, FuncWriteSingleCoil, FuncWriteMultipleCoils, FuncReadInputRegisters, FuncReadHoldingRegisters, FuncWriteSingleRegister, FuncWriteMultipleRegisters:
			return rr.raw[1], nil
		default:
			return 0, fmt.Errorf("function not found")
		}
	}
	return 0, fmt.Errorf("function not found")
}

func (rr *RtuRequest) GetAddress() (uint16, error) {
	if len(rr.raw) > 3 {
		return binary.BigEndian.Uint16(rr.raw[2:4]), nil
	}
	return 0, fmt.Errorf("address not found")
}

func (rr *RtuRequest) GetQuantity() (uint16, error) {
	f, err := rr.GetFunction()
	if err != nil {
		return 0, err
	}
	if len(rr.raw) > 5 {
		switch f {
		case FuncWriteSingleCoil, FuncWriteSingleRegister:
			return 1, nil
		case FuncReadDiscreteInputs, FuncReadCoils, FuncWriteMultipleCoils, FuncReadInputRegisters, FuncReadHoldingRegisters, FuncWriteMultipleRegisters:
			return binary.BigEndian.Uint16(rr.raw[4:6]), nil
		default:
			return 0, fmt.Errorf("function not support quantity")
		}

	}
	return 0, fmt.Errorf("quantity not found")
}

func (rr *RtuRequest) GetCountByte() (uint8, error) {
	f, err := rr.GetFunction()
	if err != nil {
		return 0, err
	}
	if len(rr.raw) > 6 {
		switch f {
		case FuncWriteMultipleCoils, FuncWriteMultipleRegisters:
			return rr.raw[6], nil
		default:
			return 0, fmt.Errorf("function not support count byte")
		}
	}
	return 0, fmt.Errorf("count bytes not found")
}

func (rr *RtuRequest) GetData() ([]byte, error) {
	f, err := rr.GetFunction()
	if err != nil {
		return nil, err
	}
	switch f {
	case FuncReadCoils, FuncReadDiscreteInputs, FuncReadInputRegisters, FuncReadHoldingRegisters:
		return []byte{}, nil
	case FuncWriteSingleCoil, FuncWriteSingleRegister:
		if len(rr.raw) > 5 {
			return rr.raw[4:6], nil
		} else {
			return nil, fmt.Errorf("data not found")
		}
	case FuncWriteMultipleCoils, FuncWriteMultipleRegisters:
		count, err := rr.GetCountByte()
		if err != nil {
			return nil, err
		}
		if len(rr.raw) > (6 + int(count)) {
			return rr.raw[7 : 7+int(count)], nil
		} else {
			return nil, fmt.Errorf("data not found")
		}
	default:
		return nil, fmt.Errorf("function not support quantity")
	}
}

func (rr *RtuRequest) GetCrc() (uint16, error) {
	f, err := rr.GetFunction()
	if err != nil {
		return 0, err
	}
	switch f {
	case FuncReadCoils, FuncReadDiscreteInputs, FuncReadInputRegisters, FuncReadHoldingRegisters:
		if len(rr.raw) > 7 {
			return binary.LittleEndian.Uint16(rr.raw[6:8]), nil
		} else {
			return 0, fmt.Errorf("crc not found")
		}

	case FuncWriteSingleCoil, FuncWriteSingleRegister:
		if len(rr.raw) > 7 {
			return binary.LittleEndian.Uint16(rr.raw[6:8]), nil
		} else {
			return 0, fmt.Errorf("crc not found")
		}
	case FuncWriteMultipleCoils, FuncWriteMultipleRegisters:
		count, err := rr.GetCountByte()
		if err != nil {
			return 0, err
		}
		if len(rr.raw) > (8 + int(count)) {
			return binary.LittleEndian.Uint16(rr.raw[7+int(count) : 7+int(count)+2]), nil
		} else {
			return 0, fmt.Errorf("crc not found")
		}
	default:
		return 0, fmt.Errorf("function not support quantity")
	}
}

func (rr *RtuRequest) Validate() bool {
	crc, err := rr.GetCrc()
	if err != nil {
		return false
	}
	return crc == CalcCRC(rr.raw[:len(rr.raw)-2])
}

func (rr *RtuRequest) GetRaw() []byte {
	return rr.raw
}
