package mbslave

type Response interface {
	GetSlaveId() uint8
	GetFunction() uint8
	GetAddress() uint16
	GetError() uint8
	GetData() []byte
	GetADU() ([]byte, error)
	SetError(errCode uint8)
	SetRead(data []byte)
	SetSingleWrite(address uint16, data []byte)
	SetMultiWrite(address uint16, countReg uint16)
	Unanswered(on bool)
}
