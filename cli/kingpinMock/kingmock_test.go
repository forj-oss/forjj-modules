package kingpinMock

import "testing"

func TestNilDetection(t *testing.T) {
	var v *Application

	if !v.IsNil() {
		t.Error("Expected IsNil to be true ", v)
	}
	v = new(Application)
	if v.IsNil() {
		t.Error("Expected IsNil to be false ", v)
	}
}
