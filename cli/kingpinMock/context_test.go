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

func TestParseContext_SelectedCommands(t *testing.T) {
	const (
		test      = "test"
		test_help = "test_help"
		add       = "add"
		add_help  = "add help"
	)

	t.Log("SetContextValue set the command context value")
	a := New("TestApplication")
	a.Command(add, add_help).Command(test, test_help).
		Flag("flag", "flag help").String()
	p := a.NewContext().SetContext(add, test)

	ret := p.SelectedCommands()

	if len(ret) != 2 {
		t.Errorf("Expected SelectedCommands() to return '%d' commands. Got '%d'", 2, len(ret))
	}
	if ret[0].FullCommand() != add && ret[1].FullCommand() != test {
		t.Errorf("Expected SelectedCommands() to return %s and %s. Got %s and %s.",
			add, test, ret[0].FullCommand(), ret[1].FullCommand())
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
	a := New("TestApplication")
	a.Command(add, add_help).Command(test, test_help).
		Flag("flag", "flag help").String()
	a.NewContext().SetContext(add, test).
		SetContextValue("flag", test_value)

	flag := a.context.cmds[len(a.context.cmds)-1].flags["flag"]
	if flag == nil {
		t.Error("Expected flag to exist. Failed.")
		return
	}
	if *flag.String() != "" {
		t.Errorf("Expected to get '' as value. Got '%s'", flag.value)
	}
	if flag.GetContextValue() != test_value {
		t.Errorf("Expected to get '%s' as value. Got '%s'", test_value, flag.GetContextValue())
	}
}

func TestParseContext_SetValue(t *testing.T) {
	const (
		test       = "test"
		test_help  = "test_help"
		add        = "add"
		add_help   = "add help"
		test_value = "a value test"
	)

	t.Log("SetContextValue set the command context value")
	a := New("TestApplication")
	a.Command(add, add_help).Command(test, test_help).
		Flag("flag", "flag help").String()
	a.NewContext().SetContext(add, test).
		SetValue("flag", test_value)

	flag := a.context.cmds[len(a.context.cmds)-1].flags["flag"]
	if flag == nil {
		t.Error("Expected flag to exist. Failed.")
		return
	}
	if *flag.String() != test_value {
		t.Errorf("Expected to get '%s' as value. Got '%s'", test_value, flag.value)
	}
	if flag.GetContextValue() != test_value {
		t.Errorf("Expected to get '%s' as value. Got '%s'", test_value, flag.GetContextValue())
	}
}

func TestParseContext_SetCliValue(t *testing.T) {
	const (
		test       = "test"
		test_help  = "test_help"
		add        = "add"
		add_help   = "add help"
		test_value = "a value test"
	)

	t.Log("SetContextValue set the command context value")
	a := New("TestApplication")
	a.Command(add, add_help).Command(test, test_help).
		Flag("flag", "flag help").String()
	a.NewContext().SetContext(add, test).
		SetCliValue("flag", test_value)

	flag := a.context.cmds[len(a.context.cmds)-1].flags["flag"]
	if flag == nil {
		t.Error("Expected flag to exist. Failed.")
		return
	}
	if *flag.String() != test_value {
		t.Errorf("Expected to get '%s' as value. Got '%s'", test_value, flag.value)
	}
	if flag.GetContextValue() != "" {
		t.Errorf("Expected to get '' as value. Got '%s'", flag.GetContextValue())
	}
}

func TestParseContext_SetContextAppValue(t *testing.T) {
	const (
		test_value = "a value test"
	)

	t.Log("SetContextAppValue() set the APP context value")
	a := New("TestApplication")
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

func TestApplication_GetContext(t *testing.T) {
	t.Log("GetContext get The internal context object")
	a := New("TestApplication")
	p := a.NewContext()

	p_ret, err := a.ParseContext([]string{})
	if p_ret == nil {
		t.Error("Expected GetContext() to return a context object. Got nil.")
	}
	if err != nil {
		t.Errorf("Expected GetContext() to work without error. Got '%s'", err)
	}
	if p != a.context {
		t.Error("Expected GetContext() to return the internal context object. Got another one.")

	}
}

func TestParseContext_GetArgValue(t *testing.T) {
	t.Log("Expect ParseContext_GetArgValue() to .")

	const (
		test       = "test"
		test_help  = "test help"
		test_value = "test value"
	)
	// --- Setting test context ---
	a := New("TestApplication")
	arg := a.Arg(test, test_help)
	arg.String()
	a.NewContext().SetContextValue(test, test_value)

	// --- Run the test ---
	public_c, err := a.ParseContext([]string{})
	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected GetContext to not return an error. Got %s", err)
		return
	}

	c := public_c.(ParseContextTester).GetContext()
	v, found := c.GetArgValue(arg)
	if !found {
		t.Errorf("Expected GetArgValue() to get '%s' arg. Not found", test)
	}
	if v != test_value {
		t.Errorf("Expected GetArgValue() to get and return '%s'. But Got '%s'", test_value, v)
	}

	// ----------------------------
	// --- Setting test context ---
	a = New("TestApplication")
	cmd := a.Command(test, test_help)
	arg = cmd.Arg(test, test_help)
	arg.String()
	a.NewContext().SetContext(test).SetContextValue(test, test_value)

	// --- Run the test ---
	public_c, err = a.ParseContext([]string{})
	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected GetContext to not return an error. Got %s", err)
		return
	}

	c = public_c.(ParseContextTester).GetContext()
	v, found = c.GetArgValue(arg)
	if !found {
		t.Errorf("Expected GetArgValue() to get '%s' arg. Not found", test)
	}
	if v != test_value {
		t.Errorf("Expected GetArgValue() to get and return '%s'. But Got '%s'", test_value, v)
	}
}

func TestParseContext_GetFlagValue(t *testing.T) {
	t.Log("Expect ParseContext_GetFlagValue() to .")

	const (
		test       = "test"
		test_help  = "test help"
		test_value = "test value"
	)
	// --- Setting test context ---
	a := New("TestApplication")
	flag := a.Flag(test, test_help)
	flag.String()
	a.NewContext().SetContextValue(test, test_value)

	// --- Run the test ---
	public_c, err := a.ParseContext([]string{})
	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected GetContext to not return an error. Got %s", err)
		return
	}

	c := public_c.(ParseContextTester).GetContext()
	v, found := c.GetFlagValue(flag)
	if !found {
		t.Errorf("Expected GetFlagValue() to get '%s' arg. Not found", test)
	}
	if v != test_value {
		t.Errorf("Expected GetFlagValue() to get and return '%s'. But Got '%s'", test_value, v)
	}

	// ----------------------------
	// --- Setting test context ---
	a = New("TestApplication")
	cmd := a.Command(test, test_help)
	flag = cmd.Flag(test, test_help)
	flag.String()
	a.NewContext().SetContext(test).SetContextValue(test, test_value)

	// --- Run the test ---
	public_c, err = a.ParseContext([]string{})
	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected GetContext to not return an error. Got %s", err)
		return
	}

	c = public_c.(ParseContextTester).GetContext()
	v, found = c.GetFlagValue(flag)
	if !found {
		t.Errorf("Expected GetFlagValue() to get '%s' arg. Not found", test)
	}
	if v != test_value {
		t.Errorf("Expected GetFlaggValue() to get and return '%s'. But Got '%s'", test_value, v)
	}
}
