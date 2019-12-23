package main

import (
	"github.com/schnack/mbslave"
	"github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"math"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "Jan _2 15:04:05.000",
	})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Fatal(mbslave.NewRtuServer(&mbslave.Config{
		Port:     "/dev/ttyUSB0",
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.TwoStopBits,
		//SilentInterval: 50 * time.Millisecond,

		SlaveId:              0xb1,
		SizeDiscreteInputs:   math.MaxUint16,
		SizeCoils:            math.MaxUint16,
		SizeInputRegisters:   math.MaxUint16,
		SizeHoldingRegisters: math.MaxUint16,
	}).Listen())
}
