package mbslave

import (
	"github.com/schnack/gotest"
	"math"
	"testing"
)

func TestNewDefaultDataModel(t *testing.T) {
	if err := gotest.Expect(NewDefaultDataModel(0x01).SlaveId).Eq(uint8(0x01)); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_Init(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)

	if err := gotest.Expect(len(ddm.discreteInputs)).Eq(math.MaxUint16); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(len(ddm.coils)).Eq(math.MaxUint16); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(len(ddm.inputRegisters)).Eq(math.MaxUint16); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(len(ddm.holdingRegisters)).Eq(math.MaxUint16); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(ddm.function[FuncReadCoils]).NotNil(); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(ddm.function[FuncReadDiscreteInputs]).NotNil(); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(ddm.function[FuncReadHoldingRegisters]).NotNil(); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(ddm.function[FuncReadInputRegisters]).NotNil(); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(ddm.function[FuncWriteSingleCoil]).NotNil(); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(ddm.function[FuncWriteSingleRegister]).NotNil(); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(ddm.function[FuncWriteMultipleCoils]).NotNil(); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(ddm.function[FuncWriteMultipleRegisters]).NotNil(); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_ReadCoils(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	ddm.coils[0] = true
	ddm.coils[2] = true
	request := NewRtuRequest([]byte{0x01, 0x01, 0x00, 0x00, 0x00, 0x08, 0x3d, 0xcc})

	if err := gotest.Expect(request.Parse()).NotError(); err != nil {
		t.Error(err)
	}
	response := NewRtuResponse(request)

	ddm.ReadCoils(request, response)

	if err := gotest.Expect(response.GetData()).Eq([]byte{0x05}); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_ReadDiscreteInputs(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	ddm.discreteInputs[0] = true
	ddm.discreteInputs[1] = true
	request := NewRtuRequest([]byte{0x01, 0x02, 0x00, 0x00, 0x00, 0x08, 0x79, 0xcc})

	if err := gotest.Expect(request.Parse()).NotError(); err != nil {
		t.Error(err)
	}
	response := NewRtuResponse(request)

	ddm.ReadDiscreteInputs(request, response)

	if err := gotest.Expect(response.GetData()).Eq([]byte{0x03}); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_ReadHoldingRegisters(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	ddm.holdingRegisters[0] = 0x0001
	ddm.holdingRegisters[1] = 0x0002
	request := NewRtuRequest([]byte{0x01, 0x03, 0x00, 0x00, 0x00, 0x02, 0xc4, 0x0b})

	if err := gotest.Expect(request.Parse()).NotError(); err != nil {
		t.Error(err)
	}
	response := NewRtuResponse(request)

	ddm.ReadHoldingRegisters(request, response)

	if err := gotest.Expect(response.GetData()).Eq([]byte{0x00, 0x01, 0x00, 0x02}); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_ReadInputRegisters(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	ddm.inputRegisters[0] = 0x0003
	ddm.inputRegisters[1] = 0x0004
	request := NewRtuRequest([]byte{0x01, 0x04, 0x00, 0x00, 0x00, 0x02, 0x71, 0xcb})

	if err := gotest.Expect(request.Parse()).NotError(); err != nil {
		t.Error(err)
	}
	response := NewRtuResponse(request)

	ddm.ReadInputRegisters(request, response)

	if err := gotest.Expect(response.GetData()).Eq([]byte{0x00, 0x03, 0x00, 0x04}); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_WriteSingleCoil(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	request := NewRtuRequest([]byte{0x01, 0x05, 0x00, 0x00, 0xff, 0x00, 0x8c, 0x3a})

	if err := gotest.Expect(request.Parse()).NotError(); err != nil {
		t.Error(err)
	}
	response := NewRtuResponse(request)

	ddm.WriteSingleCoil(request, response)

	if err := gotest.Expect(ddm.coils[0]).Eq(true); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_WriteSingleRegister(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	request := NewRtuRequest([]byte{0x01, 0x06, 0x00, 0x00, 0x01, 0x02, 0x09, 0x9b})

	if err := gotest.Expect(request.Parse()).NotError(); err != nil {
		t.Error(err)
	}
	response := NewRtuResponse(request)

	ddm.WriteSingleRegister(request, response)

	if err := gotest.Expect(ddm.holdingRegisters[0]).Eq(uint16(258)); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_WriteMultipleCoils(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	request := NewRtuRequest([]byte{0x01, 0x0f, 0x00, 0x00, 0x00, 0x02, 0x01, 0x03, 0x9e, 0x96})

	if err := gotest.Expect(request.Parse()).NotError(); err != nil {
		t.Error(err)
	}
	response := NewRtuResponse(request)

	ddm.WriteMultipleCoils(request, response)

	if err := gotest.Expect(ddm.coils[0]).Eq(true); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(ddm.coils[1]).Eq(true); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_WriteMultipleRegisters(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	request := NewRtuRequest([]byte{0x01, 0x10, 0x00, 0x00, 0x00, 0x02, 0x04, 0x00, 0x01, 0x00, 0x02, 0x23, 0xae})

	if err := gotest.Expect(request.Parse()).NotError(); err != nil {
		t.Error(err)
	}
	response := NewRtuResponse(request)

	ddm.WriteMultipleRegisters(request, response)

	if err := gotest.Expect(ddm.holdingRegisters[0]).Eq(uint16(1)); err != nil {
		t.Error(err)
	}
	if err := gotest.Expect(ddm.holdingRegisters[1]).Eq(uint16(2)); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_Handler(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	request := NewRtuRequest([]byte{0xff, 0x05, 0x00, 0x00, 0xff, 0x00, 0x99, 0xe4})

	if err := gotest.Expect(request.Parse()).NotError(); err != nil {
		t.Error(err)
	}
	response := NewRtuResponse(request)

	ddm.Handler(request, response)

	if err := gotest.Expect(ddm.coils[0]).Eq(true); err != nil {
		t.Error(err)
	}

	_, err := response.GetADU()
	if err := gotest.Expect(err).Error("not answer"); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_HandlerError(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	request := NewRtuRequest([]byte{0x01, 0x00, 0x00, 0x00, 0xff, 0x00, 0x40, 0x3a})

	if err := gotest.Expect(request.Parse()).NotError(); err != nil {
		t.Error(err)
	}
	response := NewRtuResponse(request)

	ddm.Handler(request, response)

	if err := gotest.Expect(response.GetError()).Eq(ErrorFunction); err != nil {
		t.Error(err)
	}
}
