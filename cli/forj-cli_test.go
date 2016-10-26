package cli

import (
	"github.com/forj-oss/forjj-modules/cli/kingpinMock"
	"testing"
)

// --- Global Test definition ---
type ParamTester interface {
	GetFlag() *ForjFlag
	GetArg() *ForjArg
}

func mustPanic(t *testing.T, f func()) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("Panic expected: No panic returned.")
		}
	}()

	f()
}

// ------------------------------

func TestNewForjCli(t *testing.T) {
	var app_nil *kingpinMock.Application

	t.Log("Expect an exception if the App is nil.")
	mustPanic(t, func() {
		NewForjCli(app_nil)
	})

	t.Log("Expect application to be registered.")
	c := NewForjCli(app)
	if c.App != app {
		t.Fail()
	}
}

func TestForjCli_AddFieldListCapture(t *testing.T) {
	t.Log("Expect AddFieldListCapture to add capture list.")
	c := NewForjCli(app)
	err := c.AddFieldListCapture("w", w_f)
	if err != nil {
		t.Errorf("Expected AddFieldListCapture() to return No error. Got %s", err)
	}
	if v, found := c.filters["w"]; !found || v != "("+w_f+")" {
		t.Errorf("Expected Capture to be registered as '%s'. Got '%s'", "("+w_f+")", v)
	}

	err = c.AddFieldListCapture("ft", ft_f)
	if err != nil {
		t.Error("Expected AddFieldListCapture() to return an error. Got none.")
	}
	if v, found := c.filters["ft"]; !found || v != "("+ft_f+")" {
		t.Errorf("Expected Capture to be registered as '%s'. Got '%s'", "("+ft_f+")", v)
	}

	const (
		test_f = ":(.*)"
		test   = "test"
	)
	c.AddFieldListCapture(test, test_f)
	if v, found := c.filters[test]; !found || v != test_f {
		t.Errorf("Expected Capture to be registered as '%s'. Got '%s'", test_f, v)
	}

	err = c.AddFieldListCapture(test, test_f)
	if err == nil {
		t.Errorf("Expected AddFieldListCapture() to return no error. Got %s.", err)
	}

	err = c.AddFieldListCapture("test2", `:\(.*`)
	if err != nil {
		t.Errorf("Expected AddFieldListCapture() to return no error. Got %s.", err)
	}

	err = c.AddFieldListCapture("test3", "(.*")
	if err == nil {
		t.Errorf("Expected AddFieldListCapture() to return an error. Got none.")
	}
}

func TestForjCli_AddAppFlag(t *testing.T) {
	t.Log("Expect AddAppFlag to create a Flag at App level.")

	c := NewForjCli(app)
	c.AddAppFlag(String, "test1", "test_help", nil)

	if app.GetFlag("test1").GetHelp() != "test_help" {
		t.Fail()
	}
}

func TestForjCli_NewActions(t *testing.T) {
	t.Log("Expect NewActions('create', 'direct create help', 'create %s', true) to create a new action at App level.")

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)

	a, found := c.actions[create]
	if !found {
		t.Errorf("Expected %s registered in the App layer. Not found.", create)
	}
	if !a.internal_only {
		t.Error("Expected to be an internal action. Is not")
	}
	if a.name != create {
		t.Errorf("Expected action name to be '%s'. Got '%s'", create, a.name)
	}
	if a.help != "create %s" {
		t.Errorf("Expected action help to be '%s'. Got '%s'", "create %s", a.help)
	}

	action := app.GetCommand(create)
	if action == nil {
		t.Error("Expected Command created in kingpin. Not found")
	}
	if a.cmd != action {
		t.Errorf("Expected Action '%s' created in kingpin to be identical to action.cmd. Is not.", create)
	}
	if action.FullCommand() != create {
		t.Errorf("Expected Command name to be '%s'. Got '%s'", create, action.FullCommand())
	}

	c.NewActions(update, update_help, "update %s", true)
	if len(c.actions) != 2 {
		t.Errorf("Expected 2 Actions in cli. Got '%d'", len(c.actions))
	}

	a, found = c.actions[create]
	if !found {
		t.Errorf("Expected %s registered in the App layer. Not found.", create)
	}

	action = app.GetCommand(create)
	if action == nil {
		t.Error("Expected Command created in kingpin. Not found")
	}
	if a.cmd != action {
		t.Errorf("Expected Action '%s' created in kingpin to be identical to action.cmd. Is not.", create)
	}
	if action.FullCommand() != create {
		t.Errorf("Expected Command name to be '%s'. Got '%s'", create, action.FullCommand())
	}

	a, found = c.actions[update]
	if !found {
		t.Errorf("Expected %s registered in the App layer. Not found.", update)
	}

	action = app.GetCommand(update)
	if action == nil {
		t.Error("Expected Command created in kingpin. Not found")
	}
	if a.cmd != action {
		t.Errorf("Expected Action '%s' created in kingpin to be identical to action.cmd. Is not.", update)
	}
	if action.FullCommand() != update {
		t.Errorf("Expected Command name to be '%s'. Got '%s'", update, action.FullCommand())
	}
}

func TestForjCli_OnActions(t *testing.T) {
	t.Log("Expect OnActions() to be added to selected actions at App layer.")

	// --- Setting test context ---
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)

	// --- Run the test ---
	c_ret := c.OnActions(create)

	// --- Start testing ---
	if s := "OnAction() to return cli object"; c_ret != c {
		t.Errorf("Expected %s. Is not.", s)
	}
	if s, o := "selected actions map", c.sel_actions; o == nil {
		t.Errorf("Expected %s to be initialized. Got Nil.", s)
	}
	if s, i, count := "number of actions selected", 1, len(c.sel_actions); count != i {
		t.Errorf("Expected %s to be '%d'. Got '%d'", s, i, count)
	}
	if _, found := c.sel_actions[create]; !found {
		t.Errorf("Expected to get '%s' as selected actions. Got '%s'", create, c.sel_actions)
	}

	// --- Run the test ---
	c_ret = c.OnActions()

	// --- Start testing ---
	if s := "OnAction() to return cli object"; c_ret != c {
		t.Errorf("Expected %s. Is not.", s)
	}
	if s, i, count := "number of actions selected to be ", 2, len(c.sel_actions); count != i {
		t.Errorf("Expected %s '%d'. Got '%d'", s, i, count)
	}

	// --- Run the test ---
	c.OnActions(create, update)

	// --- Start testing ---
	if s, i, count := "number of actions selected to be ", 2, len(c.sel_actions); count != i {
		t.Errorf("Expected %s '%d'. Got '%d'", s, i, count)
	}
}

func TestForjCli_GetStringValue(t *testing.T) {
	t.Log("Expect GetStringValue() to be get the Command flag value as string.")

	const (
		test       = "test"
		test_help  = "test help"
		flag       = "flag"
		flag_help  = "flag help"
		flag_value = "flag value"
	)
	// --- Setting test context ---
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)

	c.NewObject(test, test_help, false).AddField(String, flag, flag_help).DefineActions(update)
	app.NewContext().SetContext(update, test).SetContextValue(flag, flag_value)

	c.LoadContext([]string{update, test, "--" + flag, flag_value})
	// --- Run the test ---
	ret, found := c.GetStringValue(flag)

	// --- Start testing ---
	if !found {
		t.Error("Expected GetStringValue() to find the value. Not found")
	}
	if ret != flag_value {
		t.Errorf("Expected GetStringValue() to return '%s'. Got '%s'", flag_value, ret)
	}
}

func TestForjCli_GetBoolValue(t *testing.T) {
	t.Log("Expect ForjCli_GetBoolValue() to .")

	// --- Setting test context ---

	// --- Run the test ---

	// --- Start testing ---
}

func TestForjCli_GetAppBoolValue(t *testing.T) {
	t.Log("Expect ForjCli_GetAppBoolValue to .")

	// --- Setting test context ---

	// --- Run the test ---

	// --- Start testing ---
}
