package mbslave

import (
	"github.com/schnack/gotest"
	"testing"
)

func TestExceptionFunction(t *testing.T) {
	if err := gotest.Expect(ExceptionFunction(FuncReadCoils)).Eq(uint8(0x81)); err != nil {
		t.Error(err)
	}
}
