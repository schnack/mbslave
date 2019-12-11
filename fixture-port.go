package mbslave

import (
	"bytes"
	"github.com/goburrow/serial"
)

type FixturePort struct {
	Config    *serial.Config
	ReadBuff  *bytes.Buffer
	WriteBuff *bytes.Buffer
	Closed    bool
	OpenError error
}

func NewFixturePort(read []byte, closed bool, err error) *FixturePort {
	return &FixturePort{
		Config:    nil,
		ReadBuff:  bytes.NewBuffer(read),
		WriteBuff: new(bytes.Buffer),
		Closed:    closed,
		OpenError: err,
	}
}

func (f *FixturePort) Read(p []byte) (n int, err error) {
	return f.ReadBuff.Read(p)
}

func (f *FixturePort) Write(p []byte) (n int, err error) {
	return f.WriteBuff.Write(p)
}

func (f *FixturePort) Close() error {
	f.Closed = true
	return nil
}

func (f *FixturePort) Open(c *serial.Config) error {
	f.Config = c
	return f.OpenError
}
