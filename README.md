# mbslave

## USAGE

    package main
    
    import (
    	"github.com/goburrow/serial"
    	"github.com/schnack/mbslave/mbslave"
    	"github.com/sirupsen/logrus"
    )
    
    func main(){
    	logrus.SetLevel(logrus.DebugLevel)
    	logrus.Fatal(mbslave.NewRtuServer(serial.Config{
    		Address:  "/dev/ttyUSB1",
    		BaudRate: 9600,
    		DataBits: 8,
    		StopBits: 1,
    		Parity:   "N",
    		Timeout:  1 * time.Hour,
    	}, mbslave.NewDefaultDataModel(0xb1)).Listen())
    
    }