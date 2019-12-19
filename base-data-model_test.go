package mbslave

import (
	"github.com/schnack/gotest"
	"testing"
)

func TestBaseDataModel_SetSlaveId(t *testing.T) {
	bdm := BaseDataModel{}
	bdm.SetSlaveId(0x02)
	if err := gotest.Expect(bdm.SlaveId).Eq(uint8(0x02)); err != nil {
		t.Error(err)
	}
}

func TestBaseDataModel_SetFunction(t *testing.T) {
	bdm := BaseDataModel{}
	bdm.SetFunction(0x01, func(req Request, resp Response) {})
	if err := gotest.Expect(bdm.function[1]).NotNil(); err != nil {
		t.Error(err)
	}
}

func TestBaseDataModel_Handler(t *testing.T) {
	bdm := BaseDataModel{SlaveId: 0x01}
	bdm.SetFunction(0x01, func(req Request, resp Response) {
		resp.SetRead([]byte{0xff})
	})

	request := NewRtuRequest([]byte{0x01, 0x01, 0x00, 0x00, 0x00, 0x08, 0x3d, 0xcc})
	if err := gotest.Expect(request.Parse()).NotError(); err != nil {
		t.Error(err)
	}
	response := NewRtuResponse(request)

	bdm.Handler(request, response)

	if err := gotest.Expect(response.GetData()).Eq([]byte{0xff}); err != nil {
		t.Error(err)
	}
}
