package mbslave

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"sync"
	"time"
)

type RtuTransport struct {
	*Config
	handler        func(request Request, response Response)
	Port           serial.Port
	Log            logrus.FieldLogger
	silentInterval time.Duration
}

func NewRtuTransport(config *Config) *RtuTransport {
	return &RtuTransport{
		Config: config,
		Log:    logrus.StandardLogger(),
	}
}

func (rt *RtuTransport) SetHandler(f func(request Request, response Response)) {
	rt.handler = f
}

func (rt *RtuTransport) Listen() (exitError error) {
	rt.silentInterval = rt.SilentInterval()
	var err error
	if rt.Port, err = OpenSerialPort(rt.Config); err != nil {
		return err
	}
	defer rt.Port.Close()
	rt.Log.Debugf("start listing %s %d %d", rt.Config.Port, rt.BaudRate, rt.DataBits)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		// Буфер фрейма
		buff := new(bytes.Buffer)
		var muBuff sync.Mutex

		cb, ce := rt.readChan(rt.Port)

		for {
			select {
			case data := <-cb:
				muBuff.Lock()
				buff.WriteByte(data)
				muBuff.Unlock()
			case exitError = <-ce:
				// Обрабатываем финальный пакет удобно для тестов
				for data := range cb {
					muBuff.Lock()
					buff.WriteByte(data)
					muBuff.Unlock()
				}
				_ = rt.newFrame(buff, muBuff)
				return
			case <-time.After(rt.silentInterval):
				if err := rt.newFrame(buff, muBuff); err != nil {
					exitError = err
					return
				}
			}
		}
	}()
	wg.Wait()
	return
}

func (*RtuTransport) readChan(port serial.Port) (<-chan byte, <-chan error) {
	cb := make(chan byte, 256)
	ce := make(chan error)

	go func() {
		b := make([]byte, 1)
		defer close(cb)
		defer close(ce)
		for {
			n, err := port.Read(b)
			if err != nil {
				ce <- err
				return
			}
			if n != 0 {
				cb <- b[0]
			} else {
				ce <- fmt.Errorf("unable to read data from serial port")
				return
			}
		}
	}()
	return cb, ce
}

// getFrame - синхронизирует буфер
func (*RtuTransport) getFrame(buff *bytes.Buffer, mu sync.Mutex) []byte {
	mu.Lock()
	defer mu.Unlock()
	if buff.Len() == 0 {
		return nil
	}
	defer buff.Reset()
	return buff.Bytes()
}

func (rt *RtuTransport) newFrame(buff *bytes.Buffer, muBuff sync.Mutex) error {
	adu := rt.getFrame(buff, muBuff)
	if len(adu) == 0 {
		return nil
	}

	request := NewRtuRequest(adu)
	rt.Log.Debugf("<- in  raw(%03d): [% x]", len(adu), adu)

	response := NewRtuResponse(request)

	if rt.handler != nil {
		rt.handler(request, response)
	}
	rt.Log.Debugf("request   id: %02x func: %02x addr: %04x quat: %04x size: %02x data: [% x] crc: %04x",
		request.GetSlaveId(),
		request.GetFunction(),
		request.GetAddress(),
		request.GetQuantity(),
		request.GetCountByte(),
		request.GetData(),
		request.GetCrc(),
	)

	if adu, err := response.GetADU(); err == nil {
		rt.Log.Debugf("response  id: %02x func: %02x addr: %04x data: [% x] err: %02x",
			response.GetSlaveId(),
			response.GetFunction(),
			response.GetAddress(),
			response.GetData(),
			response.GetError(),
		)
		n, err := rt.Port.Write(adu)
		if err != nil {
			return err
		}
		rt.Log.Debugf("-> out raw(%03d): [% x]", n, adu)
	}
	return nil
}

func (rt *RtuTransport) SilentInterval() (frameDelay time.Duration) {
	if rt.Config.SilentInterval.Nanoseconds() != 0 {
		frameDelay = rt.Config.SilentInterval
	} else if rt.BaudRate <= 0 || rt.BaudRate > 19200 {
		frameDelay = 1750 * time.Microsecond
	} else {
		frameDelay = time.Duration(35000000/rt.BaudRate) * time.Microsecond
	}
	return
}
