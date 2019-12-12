package mbslave

import (
	"bytes"
	"github.com/goburrow/serial"
	"sync"
	"time"
)

type RtuTransport struct {
	Config     *serial.Config
	handler    func(request Request) Response
	FrameDelay time.Duration
	Port       serial.Port
}

func NewRtuTransport(config *serial.Config, handler func(request Request) Response) Transport {
	return &RtuTransport{
		Config:     config,
		handler:    handler,
		FrameDelay: RtuFrameDelay(config.BaudRate),
		Port:       serial.New(),
	}
}

func (rt *RtuTransport) Listen() (exitError error) {
	if err := rt.Port.Open(rt.Config); err != nil {
		return err
	}
	defer rt.Port.Close()

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
			case <-time.After(rt.FrameDelay):
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
	c := make(chan error)

	go func() {
		defer wg.Done()
		b := make([]byte, 1)
		n, err := port.Read(b)
		if n != 0 {
			mu.Lock()
			data.Write(b)
			mu.Unlock()
		}

		c <- err
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
	if adu == nil {
		return nil
	}

	response := rt.handler(NewRtuRequest(adu))

	if response == nil {
		return nil
	}

	if adu, err := response.GetADU(); err == nil {
		if _, err := rt.Port.Write(adu); err != nil {
			return err
		}
		time.Sleep(rt.FrameDelay)
	}
	return nil
}
