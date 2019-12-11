package mbslave

type Response interface {
	GetSlaveId() uint8
	GetFunction() uint8
	GetAddress() uint16
	GetError() uint8
	GetData() []byte
	GetADU() ([]byte, error)
}
