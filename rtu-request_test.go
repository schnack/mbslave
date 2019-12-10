package mbslave

import (
	"github.com/schnack/gotest"
	"testing"
)

func TestRtuRequest_GetSlaveId(t *testing.T) {
	rtu := RtuRequest{}
	slaveId, err := rtu.GetSlaveId()

	if err := gotest.Expect(err).Error("slaveId not found"); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(slaveId).Eq(uint8(0)); err != nil {
		t.Error(err)
	}

	rtu.raw = []byte{0x01, 0x02, 0x03}
	slaveId, err = rtu.GetSlaveId()
	if err := gotest.Expect(err).Nil(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(slaveId).Eq(uint8(0x01)); err != nil {
		t.Error(err)
	}
}

func TestRtuRequest_GetFunction(t *testing.T) {
	rtu := RtuRequest{}
	function, err := rtu.GetFunction()

	if err := gotest.Expect(err).Error("function not found"); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(function).Eq(uint8(0)); err != nil {
		t.Error(err)
	}

	rtu.raw = []byte{0x01, 0x02, 0x03}
	function, err = rtu.GetFunction()
	if err := gotest.Expect(err).Nil(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(function).Eq(uint8(0x02)); err != nil {
		t.Error(err)
	}

	rtu.raw = []byte{0x01, 0x00, 0x03}
	function, err = rtu.GetFunction()
	if err := gotest.Expect(err).Error("function not found"); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(function).Eq(uint8(0x00)); err != nil {
		t.Error(err)
	}
}

func TestRtuRequest_GetAddress(t *testing.T) {
	rtu := RtuRequest{}
	address, err := rtu.GetAddress()

	if err := gotest.Expect(err).Error("address not found"); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(address).Eq(uint16(0)); err != nil {
		t.Error(err)
	}

	rtu.raw = []byte{0x01, 0x02, 0x00, 0x02}
	address, err = rtu.GetAddress()
	if err := gotest.Expect(err).Nil(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(address).Eq(uint16(2)); err != nil {
		t.Error(err)
	}
}

func TestRtuRequest_GetQuantity(t *testing.T) {
	rtu := RtuRequest{}
	quantity, err := rtu.GetQuantity()

	if err := gotest.Expect(err).Error("function not found"); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(quantity).Eq(uint16(0)); err != nil {
		t.Error(err)
	}

	rtu.raw = []byte{0x01, 0x02, 0x00, 0x02, 0x00, 0x03}
	quantity, err = rtu.GetQuantity()
	if err := gotest.Expect(err).Nil(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(quantity).Eq(uint16(3)); err != nil {
		t.Error(err)
	}

	rtu.raw = []byte{0x01, 0x05, 0x00, 0x02, 0x00, 0x03}
	quantity, err = rtu.GetQuantity()
	if err := gotest.Expect(err).Nil(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(quantity).Eq(uint16(1)); err != nil {
		t.Error(err)
	}
}

func TestRtuRequest_GetCountByte(t *testing.T) {
	rtu := RtuRequest{}
	quantity, err := rtu.GetCountByte()

	if err := gotest.Expect(err).Error("function not found"); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(quantity).Eq(uint8(0)); err != nil {
		t.Error(err)
	}

	rtu.raw = []byte{0x01, 0x01, 0x00, 0x02, 0x00, 0x03, 0x01}
	quantity, err = rtu.GetCountByte()
	if err := gotest.Expect(err).Error("function not support count byte"); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(quantity).Eq(uint8(0)); err != nil {
		t.Error(err)
	}

	rtu.raw = []byte{0x01, 0x0f, 0x00, 0x00, 0x00, 0x01, 0x01}
	quantity, err = rtu.GetCountByte()
	if err := gotest.Expect(err).Nil(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(quantity).Eq(uint8(1)); err != nil {
		t.Error(err)
	}

}

func TestRtuRequest_GetData(t *testing.T) {
	rtu := RtuRequest{}
	data, err := rtu.GetData()

	if err := gotest.Expect(err).Error("function not found"); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(data).Nil(); err != nil {
		t.Error(err)
	}

	rtu.raw = []byte{0x01, 0x01, 0x00, 0x01}
	data, err = rtu.GetData()
	if err := gotest.Expect(err).Nil(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(data).Eq([]byte{}); err != nil {
		t.Error(err)
	}

	rtu.raw = []byte{0x01, 0x05, 0x00, 0x00, 0xff, 0x00}
	data, err = rtu.GetData()
	if err := gotest.Expect(err).Nil(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(data).Eq([]byte{0xff, 0x00}); err != nil {
		t.Error(err)
	}

	rtu.raw = []byte{0x01, 0x10, 0x00, 0x00, 0x00, 0x02, 0x04, 0x01, 0x02, 0x03, 0x04}
	data, err = rtu.GetData()
	if err := gotest.Expect(err).Nil(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(data).Eq([]byte{0x01, 0x02, 0x03, 0x04}); err != nil {
		t.Error(err)
	}
}

func TestRtuRequest_GetCrc(t *testing.T) {
	rtu := RtuRequest{}
	crc, err := rtu.GetCrc()

	if err := gotest.Expect(err).Error("function not found"); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(crc).Eq(uint16(0)); err != nil {
		t.Error(err)
	}

	rtu.raw = []byte{0x01, 0x01, 0x00, 0x01, 0x01, 0x02, 0xf0, 0x0f}
	crc, err = rtu.GetCrc()
	if err := gotest.Expect(err).Nil(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(crc).Eq(uint16(4080)); err != nil {
		t.Error(err)
	}

	rtu.raw = []byte{0x01, 0x05, 0x00, 0x00, 0xff, 0x00, 0xf0, 0x0f}
	crc, err = rtu.GetCrc()
	if err := gotest.Expect(err).Nil(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(crc).Eq(uint16(4080)); err != nil {
		t.Error(err)
	}

	rtu.raw = []byte{0x01, 0x10, 0x00, 0x00, 0x00, 0x02, 0x04, 0x01, 0x02, 0x03, 0x04, 0xf0, 0x0f}
	crc, err = rtu.GetCrc()
	if err := gotest.Expect(err).Nil(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(crc).Eq(uint16(4080)); err != nil {
		t.Error(err)
	}
}

func TestRtuRequest_Validate(t *testing.T) {
	rtu := RtuRequest{}

	if err := gotest.Expect(rtu.Validate()).False(); err != nil {
		t.Error(err)
	}

	rtu.raw = []byte{0x11, 0x02, 0x00, 0xC4, 0x00, 0x16, 0xba, 0xa9}
	if err := gotest.Expect(rtu.Validate()).True(); err != nil {
		t.Error(err)
	}

}

func TestRtuRequest_GetRaw(t *testing.T) {
	if err := gotest.Expect((&RtuRequest{raw: []byte{0x11, 0x02, 0x00, 0xC4, 0x00, 0x16, 0xba, 0xa9}}).GetRaw()).Eq([]byte{0x11, 0x02, 0x00, 0xC4, 0x00, 0x16, 0xba, 0xa9}); err != nil {
		t.Error(err)
	}
}
