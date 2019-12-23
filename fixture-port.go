package mbslave

import (
	"bytes"
	"go.bug.st/serial"
)

var OpenSerialPort = func(config *Config) (serial.Port, error) {
	return serial.Open(config.Port, &serial.Mode{
		BaudRate: config.BaudRate,
		DataBits: config.DataBits,
		Parity:   config.Parity,
		StopBits: config.StopBits,
	})
}

var InoutSerialPort = inoutSerialPort{srcOpen: OpenSerialPort}

type inoutSerialPort struct {
	srcOpen func(config *Config) (serial.Port, error)

	Config map[string]*Config
	In     map[string]*bytes.Buffer
	Out    map[string]*bytes.Buffer
	Error  map[string]error
	Closed map[string]bool
}

func (i *inoutSerialPort) Load() {
	i.Config = make(map[string]*Config)
	i.In = make(map[string]*bytes.Buffer)
	i.Out = make(map[string]*bytes.Buffer)
	i.Error = make(map[string]error)
	i.Closed = make(map[string]bool)

	OpenSerialPort = func(config *Config) (serial.Port, error) {
		i.Config[config.Port] = config
		i.GetIn(config.Port)
		i.GetOut(config.Port)
		return &fixtureSerialPort{address: config.Port, fixture: i}, i.Error[config.Port]
	}
}

func (i *inoutSerialPort) Unload() {
	OpenSerialPort = i.srcOpen
}

func (i *inoutSerialPort) GetIn(address string) *bytes.Buffer {
	if _, ok := i.In[address]; !ok {
		i.In[address] = bytes.NewBuffer([]byte{})
	}
	return i.In[address]
}

func (i *inoutSerialPort) GetOut(address string) *bytes.Buffer {
	if _, ok := i.Out[address]; !ok {
		i.Out[address] = bytes.NewBuffer([]byte{})
	}
	return i.Out[address]
}

type fixtureSerialPort struct {
	address string
	fixture *inoutSerialPort
}

func (f *fixtureSerialPort) Read(p []byte) (n int, err error) {
	return f.fixture.GetOut(f.address).Read(p)
}

func (f *fixtureSerialPort) Write(p []byte) (n int, err error) {
	return f.fixture.GetIn(f.address).Write(p)
}

func (f *fixtureSerialPort) Close() error {
	f.fixture.Closed[f.address] = true
	return nil
}

func (f *fixtureSerialPort) SetMode(mode *serial.Mode) error {
	return nil
}

func (f *fixtureSerialPort) ResetInputBuffer() error {
	return nil
}

func (f *fixtureSerialPort) ResetOutputBuffer() error {
	return nil
}

func (f *fixtureSerialPort) SetDTR(dtr bool) error {
	return nil
}

func (f *fixtureSerialPort) SetRTS(rts bool) error {
	return nil
}

func (f *fixtureSerialPort) GetModemStatusBits() (*serial.ModemStatusBits, error) {
	return &serial.ModemStatusBits{}, nil
}
