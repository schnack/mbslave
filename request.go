package mbslave

type Request interface {
	GetSlaveId() uint8
	GetFunction() uint8
	GetAddress() uint16
	GetQuantity() uint16
	GetCountByte() uint8
	GetData() []byte
	GetCrc() uint16
	Validate() error
	GetADU() []byte
}
