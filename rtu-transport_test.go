package mbslave

import (
	"bytes"
	"github.com/schnack/gotest"
	"github.com/sirupsen/logrus"
	"sync"
	"testing"
	"time"
)

func setupRtuTransport() {
	InoutSerialPort.Load()
}

func teardownRtuTransport() {
	InoutSerialPort.Unload()
}

func TestRtuTransport_Listen(t *testing.T) {
	setupRtuTransport()
	defer teardownRtuTransport()
	config := &Config{
		Port:           "com",
		BaudRate:       9600,
		SilentInterval: 2 * time.Hour,
	}
	InoutSerialPort.GetOut(config.Port).Write([]byte{0x01, 0x05, 0x00, 0x01, 0xff, 0x00, 0xdd, 0xfa})

	port, _ := OpenSerialPort(config)
	rt := &RtuTransport{
		Config: config,
		Port:   port,
		handler: func(request Request, resp Response) {
			_ = request.Parse()
			resp.SetSingleWrite(request.GetAddress(), request.GetData())
		},
		Log: logrus.StandardLogger(),
	}

	if err := gotest.Expect(rt.Listen()).Error("EOF"); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(InoutSerialPort.GetIn(config.Port).Bytes()).Eq([]byte{0x01, 0x05, 0x00, 0x01, 0xff, 0x00, 0xdd, 0xfa}); err != nil {
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
		Config: &Config{
			BaudRate: 9600,
		},
	}
	if err := gotest.Expect(rt.SilentInterval().String()).Eq("3.645ms"); err != nil {
		t.Error(err)
	}
}
