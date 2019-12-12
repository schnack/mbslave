package mbslave

import "time"

func RtuFrameDelay(baudRate int) (frameDelay time.Duration) {
	if baudRate <= 0 || baudRate > 19200 {
		frameDelay = 1750 * time.Microsecond
	} else {
		frameDelay = time.Duration(35000000/baudRate) * time.Microsecond
	}
	return
}

func ExceptionFunction(function uint8) uint8 {
	return function | 1<<7
}
