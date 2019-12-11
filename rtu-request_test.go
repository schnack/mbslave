package mbslave

import (
	"github.com/schnack/gotest"
	"testing"
)

func TestNewRtuRequestWriteSingle(t *testing.T) {
	rtu, err := NewRtuRequest([]byte{})
	if err := gotest.Expect(err).Error("frame damaged"); err != nil {
		t.Error(err)
	}

	rtu, err = NewRtuRequest([]byte{0x01, 0x05, 0x00, 0x01, 0xff, 0x00, 0xdd, 0xfa})
	if err := gotest.Expect(err).NotError(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(rtu.GetSlaveId()).Eq(uint8(0x01)); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(rtu.GetFunction()).Eq(uint8(0x05)); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(rtu.GetAddress()).Eq(uint16(0x0001)); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(rtu.GetQuantity()).Eq(uint16(0x0001)); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(rtu.GetData()).Eq([]byte{0xff, 0x00}); err != nil {
		t.Error(err)
	}
}

func TestNewRtuRequestWriteMultiple(t *testing.T) {

	rtu, err := NewRtuRequest([]byte{0x01, 0x0f, 0x00, 0x01, 0x00, 0x07, 0x01, 0xff, 0xb3, 0x16})
	if err := gotest.Expect(err).Nil(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(rtu.GetSlaveId()).Eq(uint8(0x01)); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(rtu.GetFunction()).Eq(uint8(0x0f)); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(rtu.GetAddress()).Eq(uint16(0x0001)); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(rtu.GetQuantity()).Eq(uint16(0x0007)); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(rtu.GetData()).Eq([]byte{0xff}); err != nil {
		t.Error(err)
	}
}

func TestNewRtuRequestRead(t *testing.T) {

	rtu, err := NewRtuRequest([]byte{0x01, 0x02, 0x00, 0x01, 0x00, 0x07, 0x68, 0x08})
	if err := gotest.Expect(err).NotError(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(rtu.GetSlaveId()).Eq(uint8(0x01)); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(rtu.GetFunction()).Eq(uint8(0x02)); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(rtu.GetAddress()).Eq(uint16(0x0001)); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(rtu.GetQuantity()).Eq(uint16(0x0007)); err != nil {
		t.Error(err)
	}
}

func TestRtuRequest_GetSlaveId(t *testing.T) {
	if err := gotest.Expect((&RtuRequest{SlaveId: 0x01}).GetSlaveId()).Eq(uint8(0x01)); err != nil {
		t.Error(err)
	}
}

func TestRtuRequest_GetFunction(t *testing.T) {
	if err := gotest.Expect((&RtuRequest{Function: 0x01}).GetFunction()).Eq(uint8(0x01)); err != nil {
		t.Error(err)
	}
}

func TestRtuRequest_GetAddress(t *testing.T) {
	if err := gotest.Expect((&RtuRequest{Address: 0x0001}).GetAddress()).Eq(uint16(0x0001)); err != nil {
		t.Error(err)
	}
}

func TestRtuRequest_GetQuantity(t *testing.T) {
	if err := gotest.Expect((&RtuRequest{Quantity: 0x0001}).GetQuantity()).Eq(uint16(0x0001)); err != nil {
		t.Error(err)
	}
}

func TestRtuRequest_GetCountByte(t *testing.T) {
	if err := gotest.Expect((&RtuRequest{CountByte: 0x01}).GetCountByte()).Eq(uint8(0x01)); err != nil {
		t.Error(err)
	}
}

func TestRtuRequest_GetData(t *testing.T) {
	if err := gotest.Expect((&RtuRequest{Data: []byte{0x01, 0x02}}).GetData()).Eq([]byte{0x01, 0x02}); err != nil {
		t.Error(err)
	}
}

func TestRtuRequest_GetCrc(t *testing.T) {
	if err := gotest.Expect((&RtuRequest{CRC: 0x0001}).GetCrc()).Eq(uint16(0x0001)); err != nil {
		t.Error(err)
	}
}

func TestRtuRequest_Validate(t *testing.T) {
	rtu, err := NewRtuRequest([]byte{0x11, 0x02, 0x00, 0xC4, 0x00, 0x16, 0xba, 0xa9})

	if err := gotest.Expect(err).Nil(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(rtu.Validate()).Nil(); err != nil {
		t.Error(err)
	}
}

func TestRtuRequest_GetADU(t *testing.T) {
	if err := gotest.Expect((&RtuRequest{raw: []byte{0x11, 0x02, 0x00, 0xC4, 0x00, 0x16, 0xba, 0xa9}}).GetADU()).Eq([]byte{0x11, 0x02, 0x00, 0xC4, 0x00, 0x16, 0xba, 0xa9}); err != nil {
		t.Error(err)
	}
}
