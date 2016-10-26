package kingpinMock

import (
	"testing"
)

func TestParseContext_SetContext(t *testing.T) {
	const (
		test      = "test"
		test_help = "test_help"
		add       = "add"
		add_help  = "add help"
	)

	t.Log("NewContext set The command context")
	a := New("TesApplication")
	a.Command(add, add_help).Command(test, test_help)

	p := a.NewContext().SetContext(add, test)
	if p == nil {
		t.Error("Expected to get a nil (Error). Got one")
	}

	if len(a.context.cmds) != 2 {
		t.Errorf("Expected to have 2 commands in context. Got '%d'", len(a.context.cmds))
		return
	}
	if a.context.cmds[0] != a.cmds[add] {
		t.Errorf("Expected '%s' command context added. But not the same one.", add)
	}
	if a.context.cmds[1] != a.cmds[add].cmds[test] {
		t.Errorf("Expected '%s' command context added. But not the same one.", test)
	}
}

func TestParseContext_SetContextValue(t *testing.T) {
	const (
		test       = "test"
		test_help  = "test_help"
		add        = "add"
		add_help   = "add help"
		test_value = "a value test"
	)

	t.Log("SetContextValue set the command context value")
	a := New("TesApplication")
	a.Command(add, add_help).Command(test, test_help).
		Flag("flag", "flag help").String()
	a.NewContext().SetContext(add, test).
		SetContextValue("flag", test_value)

	flag := a.context.cmds[len(a.context.cmds)-1].flags["flag"]
	if flag == nil {
		t.Error("Expected flag to exist. Failed.")
		return
	}
	if *flag.String() != test_value {
		t.Errorf("Expected to get '%s' as value. Got %s", test_value, flag.value)
	}
}

func TestParseContext_SetContextAppValue(t *testing.T) {
	const (
		test_value = "a value test"
	)

	t.Log("SetContextAppValue() set the APP context value")
	a := New("TesApplication")
	a.Flag("flag", "flag help").String()
	a.NewContext().
		SetContextAppValue("flag", test_value)

	flag := a.flags["flag"]
	if flag == nil {
		t.Error("Expected flag to exist. Failed.")
		return
	}
	if *flag.String() != test_value {
		t.Errorf("Expected to get '%s' as value. Got %s", test_value, flag.value)
	}

}
