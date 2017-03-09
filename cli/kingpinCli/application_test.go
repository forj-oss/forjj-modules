package kingpinCli

import (
	"testing"
	"github.com/alecthomas/kingpin"
)

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

func TestNew(t *testing.T) {
	const app_name="myapp"
	t.Log("Expect New to create the application interface with kingpin.")
	app := New(new(kingpin.Application), app_name)

	if app.app == nil {
		t.Error("Expected kingpin object to exist. Got Nil.")
	}
	if app.name != app_name {
		t.Errorf("Expected kingpinCli to have application name stored. Got '%s'", app.name)
	}
}

func TestApplication_Name(t *testing.T) {
	const app_name="myapp"

	t.Log("Expect Name() to return the kingpinCli application name.")
	app := New(new(kingpin.Application), app_name)

	if v := app.Name() ; v != app_name {
		t.Errorf("Expected app to be named '%s'. Got '%s'", app_name, v)
	}
}
