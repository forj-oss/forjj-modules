package cli

import (
	"github.com/forj-oss/forjj-modules/cli/kingpinMock"
	"reflect"
	"testing"
)

// -------------------------------
type ForjParamTester interface {
	GetFlag() *ForjFlag
	GetArg() *ForjArg
}

// -------------------------------

var app = kingpinMock.New("Application")

const (
	create   = "create"
	update   = "update"
	maintain = "maintain"
)
const (
	create_help = "create-help"
	update_help = "update-help"
)

const (
	workspace  = "workspace"
	infra      = "infra"
	infra_help = "infra help"
)

const (
	w_f  = `[a-z]+[a-z0-9_-]*`
	ft_f = `[A-Za-z0-9_ !:/.-]+`
)

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
	if len(c.objects) != 1 {
		t.Errorf("Expect to have only one workspace object. Got %d", len(c.objects))
	}

	o = c.NewObject(infra, infra_help, true)
	if len(c.objects) != 2 {
		t.Errorf("Expect to have only one workspace object. Got %d", len(c.objects))
	}
}

func TestForjCli_GetObject(t *testing.T) {
	t.Log("Expect NewObject('workspace', 'forjj workspace', true) to create a new object at App level.")

	const workspace_help = "workspace help"

	c := NewForjCli(app)
	o := c.NewObject(workspace, workspace_help, true)

	o_found := c.GetObject(workspace)
	if o_found != o {
		t.Error("Expected any created object to be found and returned. Is not.")
	}
}

func TestForjObject_AddKey(t *testing.T) {
	t.Log("Expect ForjObject_AddKey() to add a new field key in the object.")

	// --- Setting test context ---
	const (
		docker      = "docker-exe-path"
		docker_help = "docker-exe-path-help"
	)
	c := NewForjCli(kingpinMock.New("Application"))
	o := c.NewObject(workspace, "", true)

	// --- Run the test ---
	or := o.AddKey(String, docker, docker_help)

	// --- Start testing ---
	if or != o {
		t.Error("Expected to get the object 'object' updated. Is not.")
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

func TestForjObject_NoFields(t *testing.T) {
	t.Log("Expect ForjObject_NoFields() to register the object with no fields.")

	// --- Setting test context ---
	c := NewForjCli(kingpinMock.New("Application"))
	o := c.NewObject(workspace, "", true)

	// --- Run the test ---
	o = o.NoFields()

	// --- Start testing ---
	if o == nil {
		t.Error("Expected NoFields() to fails. but it works.")
	}
	if v, found := o.fields[no_fields]; !found {
		t.Error("Expected NoFields() to create 'no_field' record. Got nothing.")
	} else {
		if !v.key {
			t.Error("Expected NoFields() to create 'no_field' record as key. Is is not")
		}
	}

	// --- Setting test context ---
	c = NewForjCli(kingpinMock.New("Application"))
	o = c.NewObject(workspace, "", true).AddKey(String, "test", "help")

	// --- Run the test ---
	o = o.NoFields()

	// --- Start testing ---
	if o != nil {
		t.Errorf("Expected NoFields() to work. But it fails. %s", c.GetObject(workspace).err)
	}

	// --- Setting test context ---
	c = NewForjCli(kingpinMock.New("Application"))
	o = c.NewObject(workspace, "", true)

	// --- Run the test ---
	o = o.NoFields().AddKey(String, "test", "help")

	// --- Start testing ---
	if o != nil {
		t.Errorf("Expected NoFields() to work. But it fails. %s", c.GetObject(workspace).err)
	}
}

func TestForjObject_DefineActions(t *testing.T) {
	t.Log("Expect DefineActions('create') adding an action to fail if no action gets created from app.")

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	o := c.NewObject(workspace, "", true)
	or := o.DefineActions(create)
	if or != nil {
		t.Error("Expected DefineActions() to fail. Got one.")
	}

	o.AddKey(String, "test", "test help")
	or = o.DefineActions(create)
	if or != o {
		t.Error("Expected to get the object 'object' updated. Is not.")
	}

	_, found := o.actions[create]
	if found {
		t.Errorf("Expected %s registered in the object as inexistent. Found it.", create)
	}
}

func TestForjObject_DefineActions2(t *testing.T) {
	t.Log("Expect actions to be added to the object.")
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	o := c.NewObject(workspace, "", true).AddKey(String, "test", "test help")
	if o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).err)
		return
	}

	c.NewActions(create, create_help, "create %s", true)
	if o.DefineActions(create) == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).err)
		return
	}

	f, found := o.actions[create]
	if !found {
		t.Errorf("Expected %s registered in the object actions. Not found.", create)
	}
	if f.action == nil {
		t.Errorf("Expected action to refer to global action '%s'. Got nil", create)
	}
	if f.action.name != create {
		t.Errorf("Expected action name to refer to global action '%s'. Got %s", create, f.action.name)
	}

	cmd := app.GetCommand(create, workspace)
	if cmd == nil {
		t.Errorf("Expected Command '%s' to be created in kingpin. Not found.", workspace)
	}
	if cmd.FullCommand() != workspace {
		t.Errorf("Expected Command to be '%s' in kingpin. Got '%s'", workspace, cmd.FullCommand())
	}

}

func TestForjObject_DefineActions3(t *testing.T) {
	t.Log("Expect double actions to be added to the object.")
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, update_help, "update %s", false)
	o := c.NewObject(workspace, "", true).AddKey(String, "test", "test help").
		DefineActions(create, update)

	if o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).err)
		return
	}

	// Check in cli
	cmd := app.GetCommand(create, workspace)
	cmd_in_cli := o.actions[create].cmd
	if cmd_in_cli == nil {
		t.Errorf("Expected Command '%s' to be found in kingpin. Is nil.", update)
	}
	if cmd != cmd_in_cli {
		t.Errorf("Expected Command '%s' to be found identical in kingpin. Is not.", create)
	}
	if c.objects[workspace].actions[create].cmd.FullCommand() != workspace {
		t.Errorf("Expected Command '%s' associated to object '%s' to be named '%s'. Got '%s'",
			create, workspace, workspace, c.objects[workspace].actions[create].action.cmd.FullCommand())
	}

	// Check in kingpin
	cmd = app.GetCommand(create)
	if cmd == nil {
		t.Errorf("Expected Command '%s' to exist. Not found.", create)
	}
	if cmd.FullCommand() != create {
		t.Errorf("Expected '%s' has an command named '%s'", create, create)
	}

	cmd = app.GetCommand(create, workspace)
	if cmd == nil {
		t.Errorf("Expected Command '%s' to be created under '%s'. Not found.", workspace, create)
	}
	if cmd.FullCommand() != workspace {
		t.Errorf("Expected '%s/%s' has an command named '%s'", create, workspace, workspace)
	}

	cmd = app.GetCommand(update)
	if cmd == nil {
		t.Errorf("Expected Command '%s' to exist. Not found.", update)
	}
	if cmd.FullCommand() != update {
		t.Errorf("Expected '%s' has an command named '%s'", update, update)
	}

	cmd = app.GetCommand(update, workspace)
	if cmd == nil {
		t.Errorf("Expected Command '%s' to be created under '%s'. Not found.", workspace, update)
	}
	if cmd.FullCommand() != workspace {
		t.Errorf("Expected '%s/%s' has an command named '%s'", update, workspace, workspace)
	}
}

func TestForjObject_OnActions(t *testing.T) {
	t.Log("Expect OnAction() to select wanted action.")
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)
	c.NewActions(maintain, "", "maintain %s", false)
	o := c.NewObject(workspace, "", true).AddKey(String, "test", "test help").
		DefineActions(create, update).
		OnActions(create)

	if o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).err)
		return
	}
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

	const Path = "path"

	o := c.NewObject(workspace, "", true).
		AddKey(String, Path, "path help").
		DefineActions(create, update).
		OnActions(create)
	if o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).err)
		return
	}

	or := o.AddFlag(Path, nil)
	if or != o {
		t.Error("Expected to get the object updated. Is not.")
	}

	f := app.GetFlag(create, workspace, Path)
	if f == nil {
		t.Errorf("Expected flag '%s' to be added to kingpin '%s' command. Got '%s'.",
			Path, workspace, app.ListOf(create, workspace))
		return
	}
	if f.GetName() != Path {
		t.Errorf("Expected flag name to be '%s'. Got '%s'", Path, f.GetName())
	}
}

func TestForjObject_AddFlagsFromObjectAction(t *testing.T) {
	t.Log("Expect AddFlagsFromObjectAction() to be added to selected actions.")

	// --- Set context ---
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

	infra_obj := c.NewObject(infra, "", true).NoFields().
		DefineActions(update).
		OnActions()

	// --- Running the test ---
	o := infra_obj.AddFlagsFromObjectAction(workspace, update)

	// --- Start Testing ---
	if o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).err)
		return
	}
	if o != infra_obj {
		t.Error("Expected to get the object updated. Is not.")
	}

	// Checking in cli
	expected_name := test
	param, found := o.actions[update].params[expected_name]
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
	f := app.GetFlag(update, infra, test)
	if f == nil {
		t.Error("Expected to get flags from workspace added to another object action. Not found.")
		return
	}

	if f.GetName() != test {
		t.Errorf("Expected to get '%s' as flag name. Got '%s'", expected_name, f.GetName())
	}
}

func TestForjObject_AddFlagsFromObjectListActions(t *testing.T) {
	t.Log("Expect AddFlagFromObjectListActions() to be added to an object action as Flag.")

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

	infra_obj := c.NewObject(infra, "", true).NoFields().
		DefineActions(update).
		OnActions()

	if infra_obj == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).err)
		return
	}

	// Checking in cli
	o := infra_obj.AddFlagsFromObjectListActions(workspace, "to_create", update)
	if o != infra_obj {
		t.Error("Expected to get the object updated. Is not.")
	}

	expected_name := update + "-" + workspace + "s"
	if _, found := c.objects[infra].actions[update].params[expected_name]; !found {
		t.Errorf("Expected to get a new Flag '%s'related to the Objectlist added. Not found.", expected_name)
	}

	// Checking in kingpin
	flag := app.GetFlag(update, infra, expected_name)
	if flag == nil {
		t.Errorf("Expected to get a Flag in kingpin called '%s'. Got '%s'",
			update+"-"+workspace+"s", app.ListOf(update, infra))
	}
}
