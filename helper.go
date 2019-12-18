package mbslave

func ExceptionFunction(function uint8) uint8 {
	return function | 1<<7
}
