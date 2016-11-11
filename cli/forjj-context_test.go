package cli

import (
	"fmt"
	"github.com/forj-oss/forjj-modules/cli/kingpinMock"
	"reflect"
	"testing"
)

func check_object_exist(c *ForjCli, o_name, o_key, flag, value string) error {
	if _, found := c.values[o_name]; !found {
		return fmt.Errorf("Expected object '%s' to exist in values. Not found.", o_name)
	}
	if _, found := c.values[o_name].records[o_key]; !found {
		return fmt.Errorf("Expected object '%s', record '%s' to exist in values. Not found.", o_name, o_key)
	}
	if v, found := c.values[o_name].records[o_key].attrs[flag]; !found {
		return fmt.Errorf("Expected record '%s-%s' to have '%s = %s' in values. Not found.",
			o_name, o_key, flag, value)
	} else {
		if v != value {
			return fmt.Errorf("Expected key value '%s-%s-%s' to be set to '%s'. Got '%s'",
				o_name, o_key, flag, value, v)
		}
	}
	return nil
}

func TestForjCli_LoadContext(t *testing.T) {
	t.Log("Expect LoadContext() to report the context with values.")

	// --- Setting test context ---
	const (
		test       = "test"
		test_help  = "test help"
		flag       = "flag"
		flag_help  = "flag help"
		flag_value = "flag value"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, update_help, "create %s", true)

	if c.NewObject(test, test_help, false).AddKey(String, flag, flag_help).DefineActions(update) == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).Error())
		return
	}

	app.NewContext().SetContext(update, test).SetContextValue(flag, flag_value)
	// --- Run the test ---
	cmd, err := c.loadContext([]string{}, nil)

	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected LoadContext() to not fail. Got '%s'", err)
	}
	if cmd == nil {
		t.Error("Expected LoadContext() to return the last context command. Got none.")
	}
	if len(cmd) != 2 {
		t.Errorf("Expected to have '%d' context commands. Got '%d'", 2, len(cmd))
	}
}

func TestForjCli_identifyObjects(t *testing.T) {
	t.Log("Expect ForjCli_identifyObjects() to identify and store context reference to action, object and object list.")

	// --- Setting test context ---
	const (
		test       = "test"
		test_help  = "test help"
		flag       = "flag"
		flag_help  = "flag help"
		flag_value = "flag value"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, update_help, "create %s", true)

	if c.NewObject(test, test_help, false).
		AddKey(String, flag, flag_help).
		DefineActions(update).OnActions().
		AddFlag(flag, nil) == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).Error())
		return
	}

	context := app.NewContext()
	if context.SetContext(update, test) == nil {
		t.Error("Expected context with SetContext() to set context. But fails.")
	}
	if context.SetContextValue(flag, flag_value) == nil {
		t.Error("Expected context with SetContextValue() to set values. But fails.")
	}

	if _, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	}

	cmds := context.SelectedCommands()
	if len(cmds) != 2 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 2, len(cmds))
		return
	}

	// --- Run the test ---
	c.identifyObjects(cmds[len(cmds)-1])

	// --- Start testing ---
	if c.context.action == nil {
		t.Error("Expected action to be identified. Got nil.")
		return
	}
	if c.context.action != c.objects[test].actions[update].action {
		t.Errorf("Expected Action to be '%s'. Got '%s.", update, c.context.action.name)
	}
	if c.context.object == nil {
		t.Error("Expected object to be identified. Got nil.")
		return
	}
	if c.context.object != c.objects[test] {
		t.Errorf("Expected Object to be '%s'. Got '%s.", test, c.context.object.name)
	}
	if c.context.list != nil {
		t.Errorf("Expected object to be nil. Got '%s'.", c.context.list.name)
		return
	}

	// ------------------------------------------------------------------
	// --- Updating test context ---
	const (
		flag2       = "flag2"
		flag2_help  = "flag2 help"
		flag2_value = "flag2 value"
	)
	c.OnActions(create).AddFlag(String, flag2, flag2_help, nil)

	context = app.NewContext().SetContext(create).SetContextValue(flag2, flag2_value)
	if _, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	}

	cmds = context.SelectedCommands()
	if len(cmds) != 1 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 1, len(cmds))
		return
	}

	// --- Run the test ---
	c.identifyObjects(cmds[len(cmds)-1])

	// --- Start testing ---
	if c.context.action == nil {
		t.Error("Expected action to be identified. Got nil.")
		return
	}
	if c.context.action != c.actions[create] {
		t.Errorf("Expected Action to be '%s'. Got '%s.", create, c.context.action.name)
	}
	if c.context.object != nil {
		t.Errorf("Expected object to be nil. Got '%s'.", c.context.object.name)
		return
	}
	if c.context.list != nil {
		t.Errorf("Expected object to be nil. Got '%s'.", c.context.list.name)
		return
	}

	// ------------------------------------------------------------------
	// --- Updating test context ---
	const (
		repo               = "repo"
		repos              = "repos"
		reposlist_value    = "myinstance:myname,otherinstance"
		repo_help          = "repo help"
		reponame           = "name"
		reponame_help      = "repo name help"
		repo_instance      = "repo_instance"
		repo_instance_help = "repo instance help"
	)

	c.AddFieldListCapture("w", w_f)

	o := c.NewObject(repo, repo_help, false).
		AddKey(String, repo_instance, repo_instance_help).
		AddField(String, reponame, reponame_help).
		DefineActions(create).OnActions().
		AddFlag(repo_instance, nil).
		AddFlag(reponame, nil).
		CreateList("list", ",", "#w(:#w)?").
		Field(1, repo_instance).Field(3, reponame).
		AddActions(create)

	if o == nil {
		t.Errorf("Expected context failed to work with error:\n%s", c.GetObject(repo).Error())
		return
	}
	context = app.NewContext().SetContext(create, repos).SetContextValue(repos, reposlist_value)
	if _, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	}

	cmds = context.SelectedCommands()
	if len(cmds) != 2 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 1, len(cmds))
		return
	}

	// --- Run the test ---
	c.identifyObjects(cmds[len(cmds)-1])

	// --- Start testing ---
	if c.context.action == nil {
		t.Error("Expected action to be identified. Got nil.")
		return
	}
	if v, found := c.objects[repo].actions[create]; found {
		if c.context.action != v.action {
			t.Errorf("Expected Action to be '%s'. Got '%s.", create, c.context.action.name)
		}
	} else {
		t.Errorf("Expected Action '%s' to exist in Object '%s'. Got Nil.", create, repo)
	}

	if c.context.object == nil {
		t.Error("Expected object to be set. Got Nil.")
		return
	}
	if c.context.object != c.objects[repo] {
		t.Errorf("Expected Object to be '%s'. Got '%s.", repo, c.context.object.name)
	}
	if c.context.list == nil {
		t.Error("Expected object to be set. Got Nil.")
		return
	}
	if c.context.list != c.objects[repo].list["list"] {
		t.Errorf("Expected Object to be '%s'. Got '%s.", repo, c.context.object.name)
	}
}

// TestForjCli_loadListData_contextObject :
// check if <app> update test --flag "flag value"
// => creates an unique object 'test' record with key and data set.
func TestForjCli_loadListData_contextObject(t *testing.T) {
	t.Log("Expect ForjCli_loadListData() to create object list instances.")

	// --- Setting test context ---
	const (
		test       = "test"
		test_help  = "test help"
		flag       = "flag"
		flag_help  = "flag help"
		flag_value = "flag value"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, update_help, "update %s", true)

	if c.NewObject(test, test_help, false).
		AddKey(String, flag, flag_help).
		DefineActions(update).
		OnActions().
		AddFlag(flag, nil) == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).Error())
		return
	}

	context := app.NewContext().SetContext(update, test).SetContextValue(flag, flag_value)

	if _, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	}

	cmds := context.SelectedCommands()
	if len(cmds) == 0 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 2, len(cmds))
		return
	}
	// Ensure objects are identified properly.
	c.identifyObjects(cmds[len(cmds)-1])

	// --- Run the test ---
	err := c.loadListData(nil, context, cmds[len(cmds)-1])

	// --- Start testing ---
	// check in cli.
	if err != nil {
		t.Errorf("Expected loadListData to return successfully. But got an error. %s", err)
		return
	}
	if err := check_object_exist(c, test, flag_value, flag, flag_value); err != nil {
		t.Errorf("%s", err)
	}
}

// TestForjCli_loadListData_contextAction :
// check if <app> update --tests "flag_key"
// => creates an unique object 'test' record with key and data set.
func TestForjCli_loadListData_contextAction(t *testing.T) {
	t.Log("Expect ForjCli_loadListData() to create object list instances.")

	// --- Setting test context ---
	const (
		test       = "test"
		tests      = "tests"
		test_help  = "test help"
		flag       = "flag"
		flag_help  = "flag help"
		flag_value = "flag_key"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)

	c.AddFieldListCapture("w", w_f)

	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, update_help, "update %s", true)

	if c.NewObject(test, "test object help", false).
		AddKey(String, flag, flag_help).
		// <app> create test --flag <data>
		DefineActions(update).OnActions().
		AddFlag(flag, nil).

		// create list
		CreateList("to_update", ",", "#w").
		Field(1, flag).
		// <app> create tests "flag_key"
		AddActions(update) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(test).Error())
	}

	// <app> update --tests "flag_key"
	if c.AddActionFlagFromObjectListAction(create, test, "to_update", update) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.Error())
	}

	context := app.NewContext().SetContext(create).SetContextValue(tests, flag_value)

	if _, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	}

	cmds := context.SelectedCommands()
	if len(cmds) != 1 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 1, len(cmds))
		return
	}
	// Ensure objects are identified properly.
	c.identifyObjects(cmds[len(cmds)-1])

	// --- Run the test ---
	err := c.loadListData(nil, context, cmds[len(cmds)-1])

	// --- Start testing ---
	// check in cli.
	if err != nil {
		t.Errorf("Expected loadListData to return successfully. But got an error. %s", err)
		return
	}
	if err := check_object_exist(c, test, flag_value, flag, flag_value); err != nil {
		t.Errorf("%s", err)
	}
}

// TestForjCli_loadListData_contextObjectList:
// check if <app> update tests "flag value,other"
// => creates 2 objects 'test' records with key and data set.
func TestForjCli_loadListData_contextObjectList(t *testing.T) {
	t.Log("Expect ForjCli_loadListData() to create object list instances.")

	// --- Setting test context ---
	const (
		test        = "test"
		tests       = "tests"
		test_help   = "test help"
		flag        = "flag"
		flag_help   = "flag help"
		flag_value1 = "flag_value"
		flag_value2 = "other"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)

	c.AddFieldListCapture("w", w_f)

	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, update_help, "update %s", true)

	if c.NewObject(test, test_help, false).
		AddKey(String, flag, flag_help).
		// <app> create test --flag <data>
		DefineActions(create).OnActions().
		AddFlag(flag, nil).

		// create list
		CreateList("to_update", ",", "#w").
		Field(1, flag).
		// <app> create tests "flag_key"
		AddActions(create) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(test).Error())
	}

	context := app.NewContext().SetContext(create, tests).SetContextValue(tests, flag_value1+","+flag_value2)

	if _, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	}

	cmds := context.SelectedCommands()
	if len(cmds) == 0 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 2, len(cmds))
		return
	}
	// Ensure objects are identified properly.
	c.identifyObjects(cmds[len(cmds)-1])

	// --- Run the test ---
	err := c.loadListData(nil, context, cmds[len(cmds)-1])

	// --- Start testing ---
	// check in cli.
	if err != nil {
		t.Errorf("Expected loadListData to return successfully. But got an error. %s", err)
		return
	}
	if err := check_object_exist(c, test, flag_value1, flag, flag_value1); err != nil {
		t.Errorf("%s", err)
	}
	if err := check_object_exist(c, test, flag_value2, flag, flag_value2); err != nil {
		t.Errorf("%s", err)
	}
}

// TestForjCli_loadListData_contextMultipleObjectList :
// check if <app> update --tests "flag value, other" --apps "type:driver:name"
// => creates 2 different object 'test' and 'app' records with key and data set.
func TestForjCli_loadListData_contextMultipleObjectList(t *testing.T) {
	t.Log("Expect ForjCli_loadListData() to create object list instances.")

	// --- Setting test context ---
	const (
		test             = "test"
		tests            = "tests"
		test_help        = "test help"
		flag             = "flag"
		flag_help        = "flag help"
		flag_value1      = "flag-value"
		flag_value2      = "other"
		myapp            = "app"
		apps             = "apps"
		myapp_help       = "app help"
		instance         = "instance"
		instance_help    = "instance_help"
		driver_type      = "type"
		driver_type_help = "type help"
		driver           = "driver"
		driver_help      = "driver help"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)

	c.AddFieldListCapture("w", w_f)

	// <app> create
	c.NewActions(create, create_help, "create %s", true)
	// <app> update
	c.NewActions(update, update_help, "update %s", true)

	if c.NewObject(test, test_help, false).
		AddKey(String, flag, flag_help).
		// <app> create test --flag <data>
		DefineActions(create).OnActions().
		AddFlag(flag, nil).

		// create list
		CreateList("to_update", ",", "#w").
		Field(1, flag).
		// <app> create tests <data>
		AddActions(create) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(test).Error())
	}

	if c.NewObject(myapp, myapp_help, false).
		AddKey(String, instance, instance_help).
		AddField(String, driver_type, driver_type_help).
		AddField(String, driver, driver_help).
		// <app> create app --instance <instance1> --type <type> --driver <driver>
		DefineActions(create).OnActions().
		AddFlag(instance, nil).
		AddFlag(driver_type, nil).
		AddFlag(driver, nil).

		// create list
		CreateList("to_update", ",", "#w(:#w(:#w)?)?").
		Field(1, instance).
		Field(3, driver).
		Field(5, driver_type).
		// <app> create apps <data>
		AddActions(create) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(myapp).Error())
	}

	// <app> update --tests <data>
	c.AddActionFlagFromObjectListAction(update, test, "to_update", create)
	// <app> update --apps <data>
	c.AddActionFlagFromObjectListAction(update, myapp, "to_update", create)

	context := app.NewContext().SetContext(update)
	if context.SetContextValue(tests, flag_value1) == nil {
		t.Errorf("Expected context to work. Unable to add '%s' context value.", tests)
	}
	if context.SetContextValue(apps, "type:driver:name") == nil {
		t.Errorf("Expected context to work. Unable to add '%s' context value.", apps)
	}

	if _, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	}

	cmds := context.SelectedCommands()
	if len(cmds) == 0 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 2, len(cmds))
		return
	}
	// Ensure objects are identified properly.
	c.identifyObjects(cmds[len(cmds)-1])

	// --- Run the test ---
	err := c.loadListData(nil, context, cmds[len(cmds)-1])

	// --- Start testing ---
	// check in cli.
	if err != nil {
		t.Errorf("Expected loadListData to return successfully. But got an error. %s", err)
		return
	}
	if err := check_object_exist(c, test, flag_value1, flag, flag_value1); err != nil {
		t.Errorf("%s", err)
	}
	if err := check_object_exist(c, myapp, "type", instance, "type"); err != nil {
		t.Errorf("%s", err)
	}
	if err := check_object_exist(c, myapp, "type", driver, "driver"); err != nil {
		t.Errorf("%s", err)
	}
	if err := check_object_exist(c, myapp, "type", driver_type, "name"); err != nil {
		t.Errorf("%s", err)
	}
}

// TestForjCli_loadListData_contextObjectData :
// check if <app> create test --flag "flag value" --flag2 "value"
// => creates 1 object 'test' record with key and all data set.
func TestForjCli_loadListData_contextObjectData(t *testing.T) {
	t.Log("Expect ForjCli_loadListData() to create object list instances.")

	// --- Setting test context ---
	const (
		test        = "test"
		tests       = "tests"
		test_help   = "test help"
		flag        = "flag"
		flag_help   = "flag help"
		flag_value1 = "flag value"
		flag2       = "flag2"
		flag2_help  = "flag2 help"
		flag_value2 = "other"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)

	c.AddFieldListCapture("w", w_f)

	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, update_help, "update %s", true)

	if c.NewObject(test, test_help, false).
		AddKey(String, flag, flag_help).
		AddField(String, flag2, flag2_help).
		// <app> create test --flag <data> --flag2 <data>
		DefineActions(create).OnActions().
		AddFlag(flag, nil).
		AddFlag(flag2, nil) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(test).Error())
	}

	// <app> update --tests "flag_key"
	c.AddActionFlagFromObjectListAction(update, test, "to_update", create)

	context := app.NewContext().SetContext(create, test).
		SetContextValue(flag, flag_value1).
		SetContextValue(flag2, flag_value2)

	if _, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	}

	cmds := context.SelectedCommands()
	if len(cmds) == 0 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 2, len(cmds))
		return
	}
	// Ensure objects are identified properly.
	c.identifyObjects(cmds[len(cmds)-1])

	// --- Run the test ---
	err := c.loadListData(nil, context, cmds[len(cmds)-1])

	// --- Start testing ---
	// check in cli.
	if err != nil {
		t.Errorf("Expected loadListData to return successfully. But got an error. %s", err)
		return
	}
	if err := check_object_exist(c, test, flag_value1, flag, flag_value1); err != nil {
		t.Errorf("%s", err)
	}
	if err := check_object_exist(c, test, flag_value1, flag2, flag_value2); err != nil {
		t.Errorf("%s", err)
	}
}

// TestForjCli_loadListData_contextMultipleObjectsListAndData :
// TODO: check if <app> update --tests "name1,name2" --name1-flag "value" --name2-flag "value2" --apps "test:blabla"
// => creates 1 object 'test' record with key and all data set.
func TestForjCli_addInstanceFlags(t *testing.T) {
	t.Log("Expect ForjCli_LoadContext_withMoreFlags() to create object list instances.")

	// --- Setting test context ---
	const (
		test             = "test"
		tests            = "tests"
		test_help        = "test help"
		flag             = "flag"
		flag_help        = "flag help"
		flag2            = "flag2"
		flag2_help       = "flag2 help"
		flag_value1      = "value"
		flag_value2      = "value2"
		myapp            = "app"
		apps             = "apps"
		myapp_help       = "app help"
		instance         = "instance"
		instance_help    = "instance_help"
		driver_type      = "type"
		driver_type_help = "type help"
		driver           = "driver"
		driver_help      = "driver help"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)

	c.AddFieldListCapture("w", w_f)

	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, update_help, "update %s", true)

	if c.NewObject(test, test_help, false).
		AddKey(String, flag, flag_help).
		AddField(String, flag2, flag2_help).
		// <app> create test --flag <data>
		DefineActions(create).OnActions().
		AddFlag(flag, nil).
		AddFlag(flag2, nil).

		// create list
		CreateList("to_update", ",", "#w").
		Field(1, flag).
		// <app> create tests "flag_key"
		AddActions(create) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(test).Error())
	}

	if c.NewObject(myapp, myapp_help, false).
		AddKey(String, instance, instance_help).
		AddField(String, driver_type, driver_type_help).
		AddField(String, driver, driver_help).
		// <app> create test --flag <data>
		DefineActions(create).OnActions().
		AddFlag(instance, nil).AddFlag(driver_type, nil).AddFlag(driver, nil).

		// create list
		CreateList("to_update", ",", "#w(:#w(:#w)?)?").
		Field(1, driver_type).Field(3, driver).Field(5, instance).
		// <app> create tests "flag_key"
		AddActions(create) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(myapp).Error())
	}

	// <app> update --tests "flag_key"
	c.AddActionFlagFromObjectListAction(update, test, "to_update", create)
	// <app> update --apps "type:driver"
	c.AddActionFlagFromObjectListAction(update, myapp, "to_update", create)

	context := app.NewContext()
	if context.SetContext(update) == nil {
		t.Error("Expected SetContext() to work. It fails")
	}
	if context.SetContextValue(tests, "name1,name2") == nil {
		t.Error("Expected SetContextValue(tests) to work. It fails")
	}
	if context.SetContextValue(apps, "test:blabla:instance") == nil {
		t.Error("Expected SetContext(apps) to work. It fails")
	}

	if _, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	}

	cmds := context.SelectedCommands()
	if len(cmds) != 1 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 1, len(cmds))
		return
	}
	// Ensure objects are identified properly.
	c.identifyObjects(cmds[len(cmds)-1])

	if err := c.loadListData(nil, context, cmds[len(cmds)-1]); err != nil {
		t.Errorf("Expected loadListData() to work. got '%s'", err)
	}

	// --- Run the test ---
	c.addInstanceFlags()

	// --- Start testing ---
	// checking in cli
	// <app> create tests blabla --blabla-...
	if _, found := c.list[test+"_to_update"].actions[create].params["name1-flag"]; found {
		t.Errorf("Expected instance flag '%s' to NOT exist. But found it.", "name1-flag")
	}
	if v, found := c.list[test+"_to_update"].actions[create].params["name1-flag2"]; !found {
		t.Errorf("Expected instance flag '%s' to exist. But not found.", "name1-flag2")
	} else {
		if f, ok := v.(*ForjFlag); ok {
			if f.list == nil {
				t.Errorf("Expected '%s' to be attached to the list. Got Nil", "name1-flag2")
			}
			if f.instance_name != "name1" {
				t.Errorf("Expected '%s' to be attached to instance_name '%s'. got '%s'",
					"name1-flag2", "name1", f.instance_name)
			}
			if f.field_name != "flag2" {
				t.Errorf("Expected '%s' to be attached to instance_name '%s'. got '%s'",
					"name1-flag2", "flag2", f.field_name)
			}
		} else {
			t.Errorf("Expected '%s' to be '%s' type. got '%s'",
				"name1-flag2", "*ForjFlag", reflect.TypeOf(v))

		}
	}
	if _, found := c.list[test+"_to_update"].actions[create].params["name2-flag2"]; !found {
		t.Errorf("Expected instance flag '%s' to exist. But not found.", "name2-flag2")
	}
	if _, found := c.list[myapp+"_to_update"].actions[create].params["instance-"+instance]; found {
		t.Errorf("Expected instance flag '%s' to NOT exist. But found it.", "instance-"+instance)
	}
	if _, found := c.list[myapp+"_to_update"].actions[create].params["instance-"+driver]; found {
		t.Errorf("Expected instance flag '%s' to NOT exist. But found it.", "instance-"+driver)
	}
	if _, found := c.list[myapp+"_to_update"].actions[create].params["instance-"+driver_type]; found {
		t.Errorf("Expected instance flag '%s' to NOT exist. But found it.", "instance-"+driver_type)
	}

	// <add> update --tests blabla,bloblo --blabla-... --bloblo-... --apps blabla:blibli ...
	if _, found := c.actions[update].params["name1-flag2"]; !found {
		t.Errorf("Expected instance flag '%s' to exist. But not found.", "name1-flag2")
	}
	if _, found := c.actions[update].params["name2-flag2"]; !found {
		t.Errorf("Expected instance flag '%s' to exist. But not found.", "name2-flag2")
	}

	// checking in kingpin
	// <app> create tests blabla --blabla-...
	if app.GetFlag(create, tests, "name1-flag2") == nil {
		t.Errorf("Expected instance flag '%s' to exist in kingpin. But not found.", "name1-flag")
	}
	if app.GetFlag(create, tests, "instance-"+driver_type) != nil {
		t.Errorf("Expected instance flag '%s' to NOT exist in kingpin. But found it.", "instance-"+driver_type)
	}

	// <add> update --tests blabla,bloblo --blabla-... --bloblo-... --apps blabla:blibli ...
	if app.GetFlag(update, "name1-flag2") == nil {
		t.Errorf("Expected instance flag '%s' to exist in kingpin. But not found.", "name1-flag2")
	}
	if app.GetFlag(update, "name2-flag2") == nil {
		t.Errorf("Expected instance flag '%s' to exist in kingpin. But not found it.", "name2-flag2")
	}
	// At context time, instance created flags are not parsed. It will be at next Parse time.
}

func TestForjCli_contextHook(t *testing.T) {
	t.Log("Expect ForjCli_contextHook() to manipulate cli/objects.")

	// --- Setting test context ---
	const (
		test  = "test"
		test2 = "test2"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	// --- Run the test ---
	err := c.contextHook(nil)

	// --- Start testing ---
	if o := c.GetObject(test); o != nil {
		t.Errorf("Expected contextHook() to do nothing. But found the '%s' object.", test)
	}

	// Update the context
	c.ParseHook(func(c *ForjCli, _ interface{}) error {
		if c == nil {
			return nil
		}
		if c.GetObject(test) == nil {
			c.NewObject(test, "", false)
			return nil
		}
		return fmt.Errorf("Found object '%s'.", test)
	})

	// --- Run the test ---

	err = c.contextHook(nil)

	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected contextHook() to return no error. Got '%s'", err)
	}
	if o := c.GetObject(test); o == nil {
		t.Errorf("Expected contextHook() to create the ' %s' object. Not found.", test)
	}

	// --- Run another test ---
	err = c.contextHook(nil)

	// --- Start testing ---
	if err == nil {
		t.Error("Expected contextHook() to return an error. Got none")
	}
	if fmt.Sprintf("%s", err) != "Found object 'test'." {
		t.Errorf("Expected contextHook() to return a specific error. Got '%s'", err)
	}

	// Update the context
	c.ParseHook(nil).
		GetObject(test).ParseHook(func(o *ForjObject, c *ForjCli, _ interface{}) error {
		if c == nil {
			return nil
		}
		if c.GetObject(test2) == nil {
			c.NewObject(test2, "", false)
			o.AddKey(String, "flag_key", "flag help")
			return nil
		}
		return fmt.Errorf("Found object '%s'.", test2)
	})

}
