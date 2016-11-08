package cli

import (
	"github.com/forj-oss/forjj-modules/cli/kingpinMock"
	"testing"
)

func TestForjCli_AddFlagsFromObjectAction(t *testing.T) {
	t.Log("Expect AddFlagsFromObjectAction() to be added to selected actions at app layer.")

	// --- Setting test context ---
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)

	if o := c.NewObject(workspace, "", true).
		AddKey(String, "test", "test help").
		DefineActions(update).
		OnActions(update).
		AddFlag("test", nil); o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).err)
		return
	}

	const test = "test"

	// Checking if test flag exist.
	f := app.GetFlag(update, workspace, test)
	if f == nil {
		t.Errorf("Expected flag '%s' to be added to kingpin '%s' command. Got '%s'.",
			test, workspace, app.ListOf(update, workspace))
		return
	}
	if f.GetName() != test {
		t.Errorf("Expected flag name to be '%s'. Got '%s'", test, f.GetName())
	}

	// Checking if create action can get test flag from workspace.
	c.OnActions(create)

	// --- Run the test ---
	c_ret := c.AddFlagsFromObjectAction(workspace, update)

	// --- Start testing ---
	if c_ret != c {
		t.Error("Expected to get the object updated. Is not.")
	}

	// Checking in cli
	param, found := c.actions[create].params[test]
	if !found {
		t.Errorf("Expected flag '%s' added as in object action.params", test)
		return
	}

	f_cli := param.(forjParam).GetFlag()
	if f_cli == nil {
		t.Errorf("Expected to get a Flag from the object action '%s-%s'. Not found or is not a flag.",
			workspace, update)
	}

	// Checking in kingpin
	f = app.GetFlag(create, test)
	if f == nil {
		t.Error("Expected to get flags from workspace added to another object action. Not found.")
		return
	}
	if f.GetName() != test {
		t.Errorf("Expected to get '%s' as flag name. Got '%s'", test, f.GetName())
	}
}

func TestForjCli_AddFlagsFromObjectListActions(t *testing.T) {
	t.Log("Expect AddFlagsFromObjectListActions() to be added to an object action as Flag.")

	// --- Setting test context ---
	const test = "test"

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)
	c.AddFieldListCapture("w", w_f)

	if o := c.NewObject(workspace, "", true).
		AddKey(String, test, "test help").
		DefineActions(update).
		OnActions(update).
		AddFlag(test, nil).
		CreateList("to_create", ",", "#w").
		Field(1, test).
		AddActions(update); o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).err)
		return
	}

	// --- Run the test ---
	// ex : <app> create --create-workspaces "work1,work2"
	c_ret := c.AddActionFlagsFromObjectListActions(create, workspace, "to_create", update)

	// --- Start testing ---
	if c_ret != c {
		t.Errorf("Expected to get the cli object. But got an error: '%s'.", c.Error())
	}

	// Checking in cli
	expected_name := create + "-" + workspace + "s"
	if _, found := c.actions[create].params[expected_name]; !found {
		t.Errorf("Expected to get a new Flag '%s' related to the Objectlist added. Not found.", expected_name)
	}

	// Checking in kingpin
	flag := app.GetFlag(update, expected_name)
	if flag == nil {
		t.Errorf("Expected to get a Flag in kingpin called '%s'. Got '%s'",
			update+"-"+workspace+"s", app.ListOf(update))
	}
}

func TestForjCli_AddFlagFromObjectListActions(t *testing.T) {
	t.Log("Expect AddFlagsFromObjectListActions() to be added to an object action as Flag.")

	// --- Setting test context ---
	const test = "test"

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)
	c.AddFieldListCapture("w", w_f)

	if o := c.NewObject(workspace, "", true).
		AddKey(String, test, "test help").
		DefineActions(update).
		OnActions(update).
		AddFlag(test, nil).
		CreateList("to_create", ",", "#w").
		Field(1, test).
		AddActions(update); o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).err)
		return
	}

	// --- Run the test ---
	// ex : <app> create --workspaces "work1,work2"
	c_ret := c.AddActionFlagFromObjectListAction(create, workspace, "to_create", update)

	// --- Start testing ---
	if c_ret != c {
		t.Errorf("Expected to get the cli object. But got an error: '%s'.", c.Error())
	}

	// Checking in cli
	expected_name := workspace + "s"
	if _, found := c.actions[create].params[expected_name]; !found {
		t.Errorf("Expected to get a new Flag '%s' related to the Objectlist added. Not found.", expected_name)
	}

	// Checking in kingpin
	flag := app.GetFlag(create, expected_name)
	if flag == nil {
		t.Errorf("Expected to get a Flag in kingpin called '%s'. Got '%s'",
			workspace+"s", app.ListOf(update))
	}
}

func TestForjCli_AddFlag(t *testing.T) {
	t.Log("Expect AddFlag() to be added to an object action as Flag.")

	// --- Setting test context ---
	const (
		test      = "test"
		test_help = "test help"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.AddFieldListCapture("w", w_f)

	c.OnActions(create)

	// --- Run the test ---
	c_ret := c.AddFlag(String, test, test_help, nil)

	// --- Start testing ---
	if c_ret == nil {
		t.Error("Expected AddFlag() to not return Nil. Got Nil.")
		return
	}
	if c_ret != c {
		t.Error("Expected AddFlag() to return the cli object. But returned something else.")
	}

	p, found := c.actions[create].params[test]
	if !found {
		t.Errorf("Expected to create the '%s' flag in '%s'. Got nothing.", test, create)
		return
	}

	f := p.(ParamTester).GetFlag()
	if f == nil {
		t.Errorf("Expected '%s' to exist as Flag. Not found as Flag.", test)
	}
	if s, d, v := "name", f.name, test; d != v {
		t.Errorf("Expected %s to be '%s'. Got '%s'", s, v, d)
	}
	if s, d, v := "help", f.help, test_help; d != v {
		t.Errorf("Expected %s to be '%s'. Got '%s'", s, v, d)
	}

	// Test on kingpin
	flag := app.GetFlag(create, test)
	if flag == nil {
		t.Errorf("Expected flag '%s' be created. Not found", test)
	}
	if flag != f.flag {
		t.Error("Expected kingpin flag created is referenced in cli.")
	}

}

func TestForjCli_AddArg(t *testing.T) {
	t.Log("Expect AddArg() to be added to an object action as Flag.")

	// --- Setting test context ---
	const (
		test      = "test"
		test_help = "test help"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.AddFieldListCapture("w", w_f)

	c.OnActions(create)

	// --- Run the test ---
	c_ret := c.AddArg(String, test, test_help, nil)

	// --- Start testing ---
	if c_ret == nil {
		t.Error("Expected AddFlag() to not return Nil. Got Nil.")
		return
	}
	if c_ret != c {
		t.Error("Expected AddFlag() to return the cli object. But returned something else.")
	}

	p, found := c.actions[create].params[test]
	if !found {
		t.Errorf("Expected to create the '%s' flag in '%s'. Got nothing.", test, create)
		return
	}

	f := p.(ParamTester).GetArg()
	if f == nil {
		t.Errorf("Expected '%s' to exist as Arg. Not found as Arg.", test)
	}
	if s, d, v := "name", f.name, test; d != v {
		t.Errorf("Expected %s to be '%s'. Got '%s'", s, v, d)
	}
	if s, d, v := "help", f.help, test_help; d != v {
		t.Errorf("Expected %s to be '%s'. Got '%s'", s, v, d)
	}

	// Test on kingpin
	arg := app.GetArg(create, test)
	if arg == nil {
		t.Errorf("Expected flag '%s' be created. Not found", test)
	}
	if arg != f.arg {
		t.Error("Expected kingpin flag created is referenced in cli.")
	}

}
