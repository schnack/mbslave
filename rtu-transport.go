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
	serial.Mode
	handler func(request Request, response Response)
	Port    serial.Port
	Log     logrus.FieldLogger
}

func NewRtuTransport(config serial.Mode) *RtuTransport {
	return &RtuTransport{
		Mode: config,
		Log:  logrus.StandardLogger(),
	}
}

func (rt *RtuTransport) SetHandler(f func(request Request, response Response)) {
	rt.handler = f
}

func (rt *RtuTransport) Listen() (exitError error) {
	var err error
	if rt.Port, err = serial.Open("/dev/ttyUSB0", &rt.Mode); err != nil {
		return err
	}
	defer rt.Port.Close()
	rt.Log.Debugf("start listing %s %d %d %d %s ", "com3", rt.BaudRate, rt.DataBits, rt.StopBits, rt.Parity)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		// Буфер фрейма
		buff := new(bytes.Buffer)
		var muBuff sync.Mutex

		for {
			// Ожидаем поток чтения
			var wgRead sync.WaitGroup
			wgRead.Add(1)

			select {
			// Читаем  по 1 символу
			case final := <-rt.read(rt.Port, buff, muBuff, &wgRead):
				// Если ошибка значит порт закрыт
				if final != nil {
					exitError = final
					// Обрабатываем финальный пакет удобно для тестов
					_ = rt.newFrame(buff, muBuff)
					return
				}

				//Если будут проблемы с чтением wgRead.Wait()  fmt.Printf("%p on\n", &wgRead)

			// Ждем окончание ADU и парсим его
			case <-time.After(rt.rtuFrameDelay()):
				if err := rt.newFrame(buff, muBuff); err != nil {
					exitError = err
					return
				}
				// Ждем начало следующей ADU
				wgRead.Wait()
			}
		}
	}()
	wg.Wait()
	return
}

// читаем по байту для отслеживания таймингов modbus
func (*RtuTransport) read(port serial.Port, data *bytes.Buffer, mu sync.Mutex, wg *sync.WaitGroup) <-chan error {
	c := make(chan error, 1)
	go func() {
		defer wg.Done()
		b := make([]byte, 255)
		n, err := port.Read(b)
		// TODO Первый байт получили засикаем время фрейма и возвращаем фрейм в канал
		if n != 0 {
			mu.Lock()
			data.Write(b[:n])
			mu.Unlock()
			c <- err
		} else {
			c <- fmt.Errorf("unable to read data from serial port")
		}
	}()
	return c
}

// getFrame - синхронизирует буфер
func (*RtuTransport) getFrame(buff *bytes.Buffer, mu sync.Mutex) []byte {
	mu.Lock()
	defer mu.Unlock()
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

func (rt *RtuTransport) rtuFrameDelay() (frameDelay time.Duration) {
	if rt.BaudRate <= 0 || rt.BaudRate > 19200 {
		frameDelay = 1750 * time.Microsecond
	} else {
		frameDelay = time.Duration(35000000/rt.BaudRate) * time.Microsecond
	}
	return
}
