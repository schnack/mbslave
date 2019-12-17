package mbslave

type Transport interface {
	Listen() error
	SetHandler(func(Request, Response))
}
