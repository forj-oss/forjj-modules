package kingpinMock

import (
	"reflect"
	"testing"
)

func TestNilDetection(t *testing.T) {
	var v *Application

	t.Log("without application, expected app.IsNil() to be true")
	if !v.IsNil() {
		t.Error("fail: ", v)
	}
	v = new(Application)
	t.Log("with a new application, expected app.IsNil() to be false")
	if v.IsNil() {
		t.Error("fail: ", v)
	}
}

func TestAppArg(t *testing.T) {
	t.Log("New(Application).Arg() creates a new Arg with name and help")
	a := New("TesApplication")

	arg := a.Arg("arg1", "help")

	if arg == nil {
		t.Error("Expect having a new Arg object. Got nil.")
	}

	if reflect.TypeOf(arg).String() != "*kingpinMock.ArgClause" {
		t.Errorf("Expect having a new Cmd object type. Got %s", reflect.TypeOf(arg))
	}
	if _, found := a.args["arg1"]; !found {
		t.Error("Expect having a new Arg in App layer. Not found.")
	}

}

func TestAppFlag(t *testing.T) {
	t.Log("New(Application).Flag() creates a new Flag with name and help")
	a := New("TesApplication")

	flag := a.Flag("flag1", "help")

	if flag == nil {
		t.Error("Expect having a new Arg object. Got nil.")
	}
	if reflect.TypeOf(flag).String() != "*kingpinMock.FlagClause" {
		t.Errorf("Expect having a new Cmd object type. Got %s", reflect.TypeOf(flag))
	}
	if _, found := a.flags["flag1"]; !found {
		t.Error("Expect having a new Arg in App layer. Not found.")
	}

}

func TestApplication_Command(t *testing.T) {
	t.Log("New(Application).Command() creates a new Command with name and help")
	a := New("TesApplication")

	cmd := a.Command("cmd1", "help")

	if cmd == nil {
		t.Error("Expect having a new Arg object. Got nil.")
	}

	if reflect.TypeOf(cmd).String() != "*kingpinMock.CmdClause" {
		t.Errorf("Expect having a new Cmd object type. Got %s", reflect.TypeOf(cmd))
	}

	if _, found := a.cmds["cmd1"]; !found {
		t.Error("Expect having a new Command in App layer. Not found.")
	}
}
