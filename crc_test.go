package mbslave

import (
	"github.com/schnack/gotest"
	"testing"
)

func TestCalcCRC(t *testing.T) {
	if err := gotest.Expect(CalcCRC([]byte{0x02, 0x07})).Eq(uint16(0x1241)); err != nil {
		t.Error(err)
	}
}
