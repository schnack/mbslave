package mbslave

import (
	"github.com/schnack/gotest"
	"testing"
)

func TestRtuFrameDelay(t *testing.T) {
	if err := gotest.Expect(RtuFrameDelay(9600).String()).Eq("3.645ms"); err != nil {
		t.Error(err)
	}
}

func TestExceptionFunction(t *testing.T) {
	if err := gotest.Expect(ExceptionFunction(FuncReadCoils)).Eq(uint8(0x81)); err != nil {
		t.Error(err)
	}
}
