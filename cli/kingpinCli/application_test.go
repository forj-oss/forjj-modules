package kingpinCli

import "testing"

func TestNilDetection(t *testing.T) {
	var v *Application

	t.Log("without new application, expect app.IsNil() to be true")
	if !v.IsNil() {
		t.Error("fail: ", v)
	}
	v = new(Application)
	t.Log("with a new application, expect app.IsNil() to be false")
	if v.IsNil() {
		t.Error("fail ", v)
	}
}
