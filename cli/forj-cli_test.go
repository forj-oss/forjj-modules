package cli

import (
	"github.com/forj-oss/forjj-modules/cli/kingpinMock"
	"reflect"
	"testing"
)

const (
	create   = "create"
	update   = "update"
	maintain = "maintain"
)
const create_help = "create-help"

const workspace = "workspace"

var app = kingpinMock.New("Application")

const (
	w_f  = `([a-z]+[a-z0-9_-]*)`
	ft_f = `([A-Za-z0-9_ !:/.-]+)`
)

func mustPanic(t *testing.T, f func()) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("Panic expected: No panic returned.")
		}
	}()

	f()
}

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
	c.AddFieldListCapture("w", w_f)
	c.AddFieldListCapture("ft", ft_f)

	if v, found := c.filters["w"]; !found || v != w_f {
		t.Fail()
	}
	if v, found := c.filters["ft"]; !found || v != ft_f {
		t.Fail()
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

	if action.FullCommand() != create {
		t.Errorf("Expected Command name to be '%s'. Got '%s'", create, action.FullCommand())
	}
}

func TestForjCli_NewObject(t *testing.T) {
	t.Log("Expect NewObject('workspace', 'forjj workspace', true) to create a new object at App level.")

	const workspace_help = "workspace help"

	c := NewForjCli(app)
	o := c.NewObject(workspace, workspace_help, true)

	ot := reflect.TypeOf(o).String()
	if ot != "*cli.ForjObject" {
		t.Errorf("Expected to get ForjObject type. Got: %s", ot)
	}
	of, found := c.objects[workspace]

	if of != o {
		t.Error("Expected to get the object created registered. Is not.")
	}

	if !found {
		t.Errorf("Expected %s registered in the App layer as new object. Not found.", workspace)
	}

	if !o.internal {
		t.Error("Expect to be an internal object. Is not")
	}
	if o.name != workspace {
		t.Errorf("Expect object name to be '%s'. Got '%s'", workspace, o.name)
	}

	if o.desc != workspace_help {
		t.Errorf("Expect object help to be '%s'. Got %s", workspace_help, o.desc)
	}

	o = c.NewObject(workspace, workspace_help, true)
	if len(c.objects) > 1 {
		t.Errorf("Expect to have only one workspace object. Got %d", len(c.objects))
	}
}

func TestForjObject_AddField(t *testing.T) {
	t.Log("Expect AddField(cli.String, 'docker-exe-path', docker_exe_path_help) to a field to workspace object.")

	const docker = "docker-exe-path"
	const docker_help = "docker-exe-path-help"

	c := NewForjCli(kingpinMock.New("Application"))
	o := c.NewObject(workspace, "", true)

	or := o.AddField(String, docker, docker_help)

	if or != o {
		t.Error("Expected to get the object 'object' updated. Is not.")
	}

	f, found := o.fields[docker]
	if !found {
		t.Errorf("Expected %s registered in the object as new field. Not found.", docker)
	}

	if f.name != workspace+"_"+docker {
		t.Errorf("Expect field name to be '%s'. Got %s", workspace+"_"+docker, f.name)
	}

	if f.help != docker_help {
		t.Errorf("Expect field help to be '%s'. Got %s", docker_help, f.help)
	}

	if f.value_type != String {
		t.Errorf("Expect field type to be '%s'. Got %s", String, f.value_type)
	}

	or = o.AddField(String, docker, "blabla")

	if or != o {
		t.Error("Expected to get the object 'object' updated. Is not.")
	}

	if len(o.fields) > 1 {
		t.Errorf("Expected to have a unique field '%s'. Got %d", docker, len(o.fields))
	}

	f, found = o.fields[docker]

	if f.help != docker_help {
		t.Errorf("Expect field help to stay at '%s'. Got %s", docker_help, f.help)
	}
}

func TestForjObject_DefineActions(t *testing.T) {
	t.Log("Expect DefineActions('create') to add the action as Command, with a ForjObjectAction.")

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	o := c.NewObject(workspace, "", true)

	or := o.DefineActions(create)

	if or != o {
		t.Error("Expected to get the object 'object' updated. Is not.")
	}

	f, found := o.actions[create]
	if found {
		t.Errorf("Expected %s registered in the object as inexistent. Found it.", create)
	}

	c.NewActions(create, create_help, "create %s", true)
	o.DefineActions(create)

	f, found = o.actions[create]
	if !found {
		t.Errorf("Expected %s registered in the object actions. Not found.", create)
	}

	if f.action == nil {
		t.Errorf("Expected action to refer to global action '%s'. Got nil", create)
	}

	if f.action.name != create {
		t.Errorf("Expected action name to refer to global action '%s'. Got %s", create, f.action.name)
	}

	var cmd *kingpinMock.CmdClause
	cmd = app.GetCommand(create, workspace)
	if cmd == nil {
		t.Errorf("Expected Command '%s' to be created in kingpin. Not found.", workspace)
	}

	if cmd.FullCommand() != workspace {
		t.Errorf("Expected Command to be '%s' in kingpin. Got '%s'", workspace, cmd.FullCommand())
	}
}

func TestForjObject_OnActions(t *testing.T) {
	t.Log("Expect OnAction() to select wanted action.")
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)
	c.NewActions(maintain, "", "maintain %s", false)
	o := c.NewObject(workspace, "", true).
		DefineActions(create, update).
		OnActions(create)

	if len(o.actions) != 2 {
		t.Errorf("Expected 2 actions in object '%s'. Got '%d'", workspace, len(o.actions))
	}

	if len(o.sel_actions) != 1 {
		t.Errorf("Expected 1 selected action. Got '%d'", len(o.sel_actions))
	}

	a, found := o.sel_actions[create]

	if !found {
		t.Errorf("expected '%s' selected. Got nothing", create)
	}

	if a.action.name != create {
		t.Errorf("expected '%s' selected. Got '%s'", create, a.action.name)
	}

	o.OnActions(update)
	a, found = o.sel_actions[update]

	if !found {
		t.Errorf("expected '%s' selected. Got nothing", update)
	}

	if a.action.name != update {
		t.Errorf("expected '%s' selected. Got '%s'", update, a.action.name)
	}

	o.OnActions()
	if len(o.sel_actions) != 2 {
		t.Errorf("Expected 2 selected action. Got '%d'", len(o.sel_actions))
	}
}

func TestForjObject_AddFlag(t *testing.T) {
	t.Log("Expect AddFlag() to be added to selected actions.")
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)
	c.NewActions(maintain, "", "maintain %s", false)
	o := c.NewObject(workspace, "", true).
		DefineActions(create, update).
		OnActions(create)

	const Path = "path"
	or := o.AddFlag(Path, nil)

	if or != o {
		t.Error("Expected to get the object updated. Is not.")
	}

	f := app.GetFlag(create, workspace, Path)

	if f == nil {
		t.Errorf("Expected flag '%s' to be added to kingpin '%s' command. Not found.", Path, workspace)
		return
	}

	if f.GetName() != Path {
		t.Errorf("Expected flag name to be '%s'. Got '%s'", Path, f.GetName())
	}
}
