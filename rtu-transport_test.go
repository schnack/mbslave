package mbslave

import (
	"bytes"
	"github.com/goburrow/serial"
	"github.com/schnack/gotest"
	"github.com/sirupsen/logrus"
	"sync"
	"testing"
)

func TestRtuTransport_Listen(t *testing.T) {
	port := NewFixturePort([]byte{0x01, 0x05, 0x00, 0x01, 0xff, 0x00, 0xdd, 0xfa}, false, nil)
	rt := &RtuTransport{
		Config: serial.Config{},
		Port:   port,
		handler: func(request Request, resp Response) {
			_ = request.Parse()
			resp.SetSingleWrite(request.GetAddress(), request.GetData())
		},
		Log: logrus.StandardLogger(),
	}

	if err := gotest.Expect(rt.Listen()).Error("unable to read data from serial port"); err != nil {
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

func TestRtuTransport_rtuFrameDelay(t *testing.T) {
	rt := &RtuTransport{
		Config: serial.Config{BaudRate: 9600},
	}
	if err := gotest.Expect(rt.rtuFrameDelay().String()).Eq("3.645ms"); err != nil {
		t.Error(err)
	}
}
