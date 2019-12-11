package mbslave

import (
	"bytes"
	"github.com/goburrow/serial"
	"github.com/schnack/gotest"
	"sync"
	"testing"
)

func TestRtuFrameDelay(t *testing.T) {
	if err := gotest.Expect(RtuFrameDelay(9600).String()).Eq("3.645ms"); err != nil {
		t.Error(err)
	}
}

func TestRtuTransport_Listen(t *testing.T) {
	port := NewFixturePort([]byte{0x01, 0x05, 0x00, 0x01, 0xff, 0x00, 0xdd, 0xfa}, false, nil)
	rt := &RtuTransport{
		Config:     &serial.Config{},
		SlaveId:    0x01,
		FrameDelay: RtuFrameDelay(9600),
		Port:       port,
		handler: func(request Request) Response {
			return &RtuResponse{
				SlaveId:  request.GetSlaveId(),
				Function: request.GetFunction(),
				Address:  request.GetAddress(),
				Data:     request.GetData(),
			}
		},
	}

	if err := gotest.Expect(rt.Listen()).Error("EOF"); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(port.WriteBuff.Bytes()).Eq([]byte{0x01, 0x05, 0x00, 0x01, 0xff, 0x00, 0xdd, 0xfa}); err != nil {
		t.Error(err)
	}

}

func TestRtuTransport_getFrame(t *testing.T) {
	var mu sync.Mutex
	buff := bytes.NewBuffer([]byte{0x01, 0x02})

	if err := gotest.Expect((&RtuTransport{}).getFrame(buff, mu)).Eq([]byte{1, 2}); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(buff.Bytes()).Eq([]byte{}); err != nil {
		t.Error(err)
	}
}
