package cli

import (
	"github.com/forj-oss/forjj-modules/cli/kingpinMock"
	"testing"
)

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
	if _, found := c.values[test]; !found {
		t.Errorf("Expected object '%s' to exist in values. Not found.", test)
		return
	}
	if _, found := c.values[test].records[flag_value]; !found {
		t.Errorf("Expected object '%s', record '%s' to exist in values. Not found.", test, flag_value)
		return
	}
	if v, found := c.values[test].records[flag_value].attrs[flag]; !found {
		t.Errorf("Expected record '%s-%s' to have '%s = %s' in values. Not found.", test, flag_value, flag, flag_value)
		return
	} else {
		if v != flag_value {
			t.Errorf("Expected key value '%s-%s-%s' to be set to '%s'. Got '%s'", test, flag_value, flag, flag_value, v)
		}
	}
}
