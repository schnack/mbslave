package mbslave

import (
	"github.com/schnack/gotest"
	"math"
	"sync"
	"testing"
)

func TestNewDefaultDataModel(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)

	if err := gotest.Expect(ddm.SlaveId).Eq(uint8(0x01)); err != nil {
		t.Error(err)
	}

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

	if err := gotest.Expect(len(ddm.callbackDiscreteInputs)).Eq(math.MaxUint16); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(len(ddm.callbackCoils)).Eq(math.MaxUint16); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(len(ddm.callbackInputRegisters)).Eq(math.MaxUint16); err != nil {
		t.Error(err)
	}

	if err := gotest.Expect(len(ddm.callbackHoldingRegisters)).Eq(math.MaxUint16); err != nil {
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

func TestDefaultDataModel_SetCallbackDiscreteInputs(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	ddm.SetCallbackDiscreteInputs(0x0001, func(e Event, a uint16, v bool) {})
	if err := gotest.Expect(ddm.callbackDiscreteInputs[1]).NotNil(); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_SetCallbackCoils(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	ddm.SetCallbackCoils(0x0001, func(e Event, a uint16, v bool) {})
	if err := gotest.Expect(ddm.callbackCoils[1]).NotNil(); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_SetCallbackInputRegisters(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	ddm.SetCallbackInputRegisters(0x0001, func(e Event, a uint16, v uint16) {})
	if err := gotest.Expect(ddm.callbackInputRegisters[1]).NotNil(); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_SetCallbackHoldingRegisters(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	ddm.SetCallbackHoldingRegisters(0x0001, func(e Event, a uint16, v uint16) {})
	if err := gotest.Expect(ddm.callbackHoldingRegisters[1]).NotNil(); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_LengthDiscreteInputs(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	if err := gotest.Expect(ddm.LengthDiscreteInputs()).Eq(math.MaxUint16); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_LengthCoils(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	if err := gotest.Expect(ddm.LengthCoils()).Eq(math.MaxUint16); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_LengthInputRegisters(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	if err := gotest.Expect(ddm.LengthInputRegisters()).Eq(math.MaxUint16); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_LengthHoldingRegisters(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	if err := gotest.Expect(ddm.LengthHoldingRegisters()).Eq(math.MaxUint16); err != nil {
		t.Error(err)
	}
}

func TestDefaultDataModel_SetDiscreteInputs(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	ddm.SetCallbackDiscreteInputs(0x0001, func(e Event, a uint16, v bool) {
		if err := gotest.Expect(e).Eq(EventWrite); err != nil {
			t.Error(err)
		}
		if err := gotest.Expect(a).Eq(uint16(0x0001)); err != nil {
			t.Error(err)
		}
		if err := gotest.Expect(v).True(); err != nil {
			t.Error(err)
		}
		wg.Done()
	})

	if err := gotest.Expect(ddm.SetDiscreteInputs(0x0001, true)).NotError(); err != nil {
		t.Error()
	}

	wg.Wait()
}

func TestDefaultDataModel_SetCoils(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	ddm.SetCallbackCoils(0x0001, func(e Event, a uint16, v bool) {
		if err := gotest.Expect(e).Eq(EventWrite); err != nil {
			t.Error(err)
		}
		if err := gotest.Expect(a).Eq(uint16(0x0001)); err != nil {
			t.Error(err)
		}
		if err := gotest.Expect(v).True(); err != nil {
			t.Error(err)
		}
		wg.Done()
	})

	if err := gotest.Expect(ddm.SetCoils(0x0001, true)).NotError(); err != nil {
		t.Error()
	}

	wg.Wait()
}

func TestDefaultDataModel_SetHoldingRegisters(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	ddm.SetCallbackHoldingRegisters(0x0001, func(e Event, a uint16, v uint16) {
		if err := gotest.Expect(e).Eq(EventWrite); err != nil {
			t.Error(err)
		}
		if err := gotest.Expect(a).Eq(uint16(0x0001)); err != nil {
			t.Error(err)
		}
		if err := gotest.Expect(v).Eq(uint16(0xff00)); err != nil {
			t.Error(err)
		}
		wg.Done()
	})

	if err := gotest.Expect(ddm.SetHoldingRegisters(0x0001, 0xff00)).NotError(); err != nil {
		t.Error()
	}

	wg.Wait()
}

func TestDefaultDataModel_SetInputRegisters(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	ddm.SetCallbackInputRegisters(0x0001, func(e Event, a uint16, v uint16) {
		if err := gotest.Expect(e).Eq(EventWrite); err != nil {
			t.Error(err)
		}
		if err := gotest.Expect(a).Eq(uint16(0x0001)); err != nil {
			t.Error(err)
		}
		if err := gotest.Expect(v).Eq(uint16(0xff00)); err != nil {
			t.Error(err)
		}
		wg.Done()
	})

	if err := gotest.Expect(ddm.SetInputRegisters(0x0001, 0xff00)).NotError(); err != nil {
		t.Error()
	}

	wg.Wait()
}

func TestDefaultDataModel_GetDiscreteInputs(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	ddm.SetCallbackDiscreteInputs(0x0001, func(e Event, a uint16, v bool) {
		if err := gotest.Expect(e).Eq(EventRead); err != nil {
			t.Error(err)
		}
		if err := gotest.Expect(a).Eq(uint16(0x0001)); err != nil {
			t.Error(err)
		}
		if err := gotest.Expect(v).False(); err != nil {
			t.Error(err)
		}
		wg.Done()
	})

	if err := gotest.Expect(ddm.GetDiscreteInputs(0x0001)).False(); err != nil {
		t.Error()
	}

	wg.Wait()
}

func TestDefaultDataModel_GetCoils(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	ddm.SetCallbackCoils(0x0001, func(e Event, a uint16, v bool) {
		if err := gotest.Expect(e).Eq(EventRead); err != nil {
			t.Error(err)
		}
		if err := gotest.Expect(a).Eq(uint16(0x0001)); err != nil {
			t.Error(err)
		}
		if err := gotest.Expect(v).False(); err != nil {
			t.Error(err)
		}
		wg.Done()
	})

	if err := gotest.Expect(ddm.GetCoils(0x0001)).False(); err != nil {
		t.Error()
	}

	wg.Wait()
}

func TestDefaultDataModel_GetHoldingRegisters(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	ddm.SetCallbackHoldingRegisters(0x0001, func(e Event, a uint16, v uint16) {
		if err := gotest.Expect(e).Eq(EventRead); err != nil {
			t.Error(err)
		}
		if err := gotest.Expect(a).Eq(uint16(0x0001)); err != nil {
			t.Error(err)
		}
		if err := gotest.Expect(v).Eq(uint16(0)); err != nil {
			t.Error(err)
		}
		wg.Done()
	})

	if err := gotest.Expect(ddm.GetHoldingRegisters(0x0001)).Eq(uint16(0)); err != nil {
		t.Error()
	}

	wg.Wait()
}

func TestDefaultDataModel_GetInputRegisters(t *testing.T) {
	ddm := NewDefaultDataModel(0x01)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	ddm.SetCallbackInputRegisters(0x0001, func(e Event, a uint16, v uint16) {
		if err := gotest.Expect(e).Eq(EventRead); err != nil {
			t.Error(err)
		}
		if err := gotest.Expect(a).Eq(uint16(0x0001)); err != nil {
			t.Error(err)
		}
		if err := gotest.Expect(v).Eq(uint16(0)); err != nil {
			t.Error(err)
		}
		wg.Done()
	})

	if err := gotest.Expect(ddm.GetInputRegisters(0x0001)).Eq(uint16(0)); err != nil {
		t.Error()
	}

	wg.Wait()
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
