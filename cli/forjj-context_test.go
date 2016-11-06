package cli

import (
	"github.com/forj-oss/forjj-modules/cli/kingpinMock"
	"testing"
)

func check_object_exist(t *testing.T, c *ForjCli, o_name, o_key, flag, value string) {
	if _, found := c.values[o_name]; !found {
		t.Errorf("Expected object '%s' to exist in values. Not found.", o_name)
		return
	}
	if _, found := c.values[o_name].records[o_key]; !found {
		t.Errorf("Expected object '%s', record '%s' to exist in values. Not found.", o_name, o_key)
		return
	}
	if v, found := c.values[o_name].records[o_key].attrs[flag]; !found {
		t.Errorf("Expected record '%s-%s' to have '%s = %s' in values. Not found.",
			o_name, o_key, flag, value)
		return
	} else {
		if v != value {
			t.Errorf("Expected key value '%s-%s-%s' to be set to '%s'. Got '%s'",
				o_name, o_key, flag, value, v)
		}
	}
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
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).err)
		return
	}

	app.NewContext().SetContext(update, test).SetContextValue(flag, flag_value)
	// --- Run the test ---
	cmd, err := c.LoadContext([]string{})

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

	if c.NewObject(test, test_help, false).AddKey(String, flag, flag_help).DefineActions(update) == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).err)
		return
	}

	context := app.NewContext().SetContext(update, test).SetContextValue(flag, flag_value)

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
		reposlist          = "repos-list"
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
	context = app.NewContext().SetContext(create, repos).SetContextValue(reposlist, reposlist_value)
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
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).err)
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
	check_object_exist(t, c, test, flag_value, flag, flag_value)
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
		DefineActions(create).
		AddFlag(flag, nil).

		// create list
		CreateList("to_update", ",", "#w").
		Field(1, flag).
		// <app> create tests "flag_key"
		AddActions(create) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(test).Error())
	}

	// <app> update --tests "flag_key"
	c.AddFlagFromObjectListAction(test, "to_update", update)

	context := app.NewContext().SetContext(update).SetContextValue(tests, flag_value)

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
	check_object_exist(t, c, test, flag_value, flag, flag_value)
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
		flag_value1 = "flag value"
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
		DefineActions(create).
		AddFlag(flag, nil).

		// create list
		CreateList("to_update", ",", "#w").
		Field(1, flag).
		// <app> create tests "flag_key"
		AddActions(create) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(test).Error())
	}

	// <app> update --tests "flag_key"
	c.AddFlagFromObjectListAction(test, "to_update", update)

	context := app.NewContext().SetContext(update, tests).SetContextValue(flag, flag_value1+","+flag_value2)

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
	check_object_exist(t, c, test, flag_value1, flag, flag_value1)
	check_object_exist(t, c, test, flag_value2, flag, flag_value2)
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
		flag_value1      = "flag value"
		flag_value2      = "other"
		myapp            = "app"
		apps             = "apps"
		myapp_help       = "app help"
		instance         = "intance"
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
		DefineActions(create).
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
		DefineActions(create).
		AddFlag(flag, nil).

		// create list
		CreateList("to_update", ",", "#w(:#w(:#w)?)?").
		Field(1, instance).
		// <app> create apps <data>
		AddActions(create) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(test).Error())
	}

	// <app> update --tests <data>
	c.AddFlagFromObjectListAction(test, "to_update", update)
	// <app> update --apps <data>
	c.AddFlagFromObjectListAction(myapp, "to_update", update)

	context := app.NewContext().SetContext(update).
		SetContextValue(tests, flag_value1).
		SetContextValue(apps, "type:driver:name")

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
	check_object_exist(t, c, test, flag_value1, flag, flag_value1)
	check_object_exist(t, c, myapp, instance, instance, "name")
	check_object_exist(t, c, myapp, instance, driver_type, "type")
	check_object_exist(t, c, myapp, instance, driver, "driver")
}

// TestForjCli_loadListData_contextObjectData :
// TODO: check if <app> create test --flag "flag value" --flag2 "value"
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
		DefineActions(create).
		AddFlag(flag, nil).
		AddFlag(flag2, nil) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(test).Error())
	}

	// <app> update --tests "flag_key"
	c.AddFlagFromObjectListAction(test, "to_update", update)

	context := app.NewContext().SetContext(update).SetContextValue(flag, flag_value1)

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
	check_object_exist(t, c, test, flag_value1, flag, flag_value1)
}

// TestForjCli_loadListData_contextMultipleObjectsListAndData :
// TODO: check if <app> create --tests "name1,name2" --name1-flag "value" --name2-flag "value2" --apps "test:blabla"
// => creates 1 object 'test' record with key and all data set.
func TestForjCli_loadListData_contextMultipleObjectsListAndData(t *testing.T) {
	t.Log("Expect ForjCli_loadListData() to create object list instances.")

	// --- Setting test context ---
	const (
		test        = "test"
		tests       = "tests"
		test_help   = "test help"
		flag        = "flag"
		flag_help   = "flag help"
		flag_value1 = "flag value"
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
		DefineActions(create).
		AddFlag(flag, nil).

		// create list
		CreateList("to_update", ",", "#w").
		Field(1, flag).
		// <app> create tests "flag_key"
		AddActions(create) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(test).Error())
	}

	// <app> update --tests "flag_key"
	c.AddFlagFromObjectListAction(test, "to_update", update)
	context := app.NewContext().SetContext(update).SetContextValue(flag, flag_value1)

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
	check_object_exist(t, c, test, flag_value1, flag, flag_value1)
}
