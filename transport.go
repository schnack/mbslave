package mbslave

type Transport interface {
	Listen() error
}
