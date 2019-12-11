package mbslave

type Transport interface {
	Listen() error
	HandlerFunc(func(Request) Response)
}
