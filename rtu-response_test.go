package mbslave

import (
	"github.com/schnack/gotest"
	"testing"
)

func TestRtuResponse_GetRaw(t *testing.T) {
	response := RtuResponse{
		SlaveId:  0x0A,
		Function: FuncReadCoils,
		Address:  0,
		Err:      0,
		Data:     []byte{0x01},
	}

	data, err := response.GetRaw()
	if err := gotest.Expect(err).Nil(); err != nil {
		t.Fatal(err)
	}

	if err := gotest.Expect(data).Eq([]byte{0x0A, 0x01, 0x01, 0x01, 0x92, 0x6C}); err != nil {
		t.Error(err)
	}

	response.Function = FuncWriteSingleCoil
	response.Address = 0x0001
	response.Data = []byte{0xff, 0x00}

	data, err = response.GetRaw()
	if err := gotest.Expect(err).Nil(); err != nil {
		t.Fatal(err)
	}

	if err := gotest.Expect(data).Eq([]byte{0x0A, 0x05, 0x00, 0x01, 0xff, 0x00, 0xDC, 0x81}); err != nil {
		t.Error(err)
	}

	response.Function = ExceptionReadCoils
	response.Err = 0x02

	data, err = response.GetRaw()
	if err := gotest.Expect(err).Nil(); err != nil {
		t.Fatal(err)
	}

	if err := gotest.Expect(data).Eq([]byte{0x0A, 0x81, 0x02, 0xb0, 0x53}); err != nil {
		t.Error(err)
	}
}
