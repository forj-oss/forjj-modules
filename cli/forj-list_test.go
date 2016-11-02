package cli

import (
	"reflect"
	"testing"
)

func TestForjObject_CreateList2(t *testing.T) {
	t.Log("Expect CreateList to return nil if regexp is wrong.")

	const (
		repo_help = "repo help"
		repo      = "repo"
	)

	c := NewForjCli(app)
	c.AddFieldListCapture("w", w_f)
	o := c.NewObject(repo, repo_help, true)

	l := o.CreateList("to_create", ",", "(#w")
	if l != nil {
		t.Error("Expected CreateList() to return nil if the regexp is failing. But got one list.")
	}
}

func TestForjObject_CreateList(t *testing.T) {
	t.Log("Expect CreateList to create a new object List at App level.")

	const (
		repo_help = "repo help"
		repo      = "repo"
	)

	c := NewForjCli(app)
	c.AddFieldListCapture("w", w_f)
	c.AddFieldListCapture("ft", ft_f)
	o := c.NewObject(repo, repo_help, true).AddKey(String, "name", "name help")

	l := o.CreateList("to_create", ",", "#w")
	if l.name != "to_create" {
		t.Errorf("Expected list name to be '%s'. Got '%s'", "to_create", l.name)
	}
	expected_reg := "(" + w_f + ")"
	if l.ext_regexp.String() != "("+w_f+")" {
		t.Errorf("Expected list regexp to be '%s'. Got '%s'", expected_reg, l.ext_regexp)
	}
	if l.sep != "," {
		t.Errorf("Expected list separator to be '%s'. Got '%s'.", ",", l.sep)
	}
	if l.key_name != "name" {
		t.Errorf("Expected list key name to be '%s'. Got '%s'.", "name", l.key_name)
	}
	if _, found := o.list["to_create"]; !found {
		t.Errorf("Expected list '%s' not found in object", "to_create")
	}
	if _, found := c.list[repo+"_to_create"]; !found {
		t.Errorf("Expected list '%s' not found in cli", repo+"_to_create")
	}

	l = o.CreateList("another_list", ",", "#w(:#ft)?")
	if l.name != "another_list" {
		t.Errorf("Expected list name to be '%s'. Got '%s'", "to_create", l.name)
	}
	expected_reg = "(" + w_f + ")(:(" + ft_f + "))?"
	if l.ext_regexp.String() != expected_reg {
		t.Errorf("Expected list regexp to be '%s'. Got '%s'", expected_reg, l.ext_regexp)
	}
}

func TestForjObjectList_Field(t *testing.T) {
	t.Log("Expect Field() to add a new extract field in a list.")

	const (
		repo_help = "repo help"
		repo      = "repo"
	)

	c := NewForjCli(app)
	c.NewActions(create, create_help, "%s", false)
	c.AddFieldListCapture("w", w_f)
	c.AddFieldListCapture("ft", ft_f)
	o := c.NewObject(repo, repo_help, true).
		AddKey(String, "name", "help").
		DefineActions(create).
		OnActions().
		AddFlag("name", nil)
	if o == nil {
		t.Errorf("Expected context failure. %s", c.GetObject(repo).Error())
		return
	}

	l := o.CreateList("to_create", ",", "#w(:#ft)?")
	if l == nil {
		t.Error("Expected CreateList() to return the object list. Got nil.")
	}

	l_ret := l.Field(1, "name")
	if l_ret != l {
		t.Error("Expected Field() to return the list object. Is not.")
	}

	field, found := l.fields_name[1]
	if !found {
		t.Errorf("Expected Field '%s' to be added. Not found.", "name")
	}
	if field != "name" {
		t.Errorf("Expected new field to be named '%s'. Got '%s'.", "name", field)
	}

	l_ret = l.Field(3, "instance")
	if l_ret != nil {
		t.Error("Expected Field() to return nil if an error is found. Got something else.")
	}

	field, found = l.fields_name[3]
	if found {
		t.Errorf("Expected Field '%s' to NOT be added, because object has no field '%s'. Got it.",
			"instance", "instance")
	}

	o.AddField(String, "instance", "instance help").
		OnActions().
		AddFlag("instance", nil)

	l.Field(3, "instance")
	field, found = l.fields_name[3]
	if !found {
		t.Errorf("Expected Field '%s' to be added. Not found.", "instance")
	}
	if field != "instance" {
		t.Errorf("Expected new field to be named '%s'. Got '%s'.", "instance", field)
	}

}

func TestForjObjectList_AddActions(t *testing.T) {
	t.Log("Expect AddActions() to add some action for the list.")
	// --- Setting test context ---
	const (
		repo_help     = "repo help"
		repo          = "repo"
		repos         = "repos"
		reposlist     = "repos-list"
		maintain_help = "maintain help"
	)

	c := NewForjCli(app)

	c.NewActions(create, create_help, "create %s", false)
	c.NewActions(update, update_help, "update %s", false)
	c.NewActions(maintain, maintain_help, "maintain %s", false)

	c.AddFieldListCapture("w", w_f)
	c.AddFieldListCapture("ft", ft_f)

	o := c.NewObject(repo, repo_help, true).
		AddKey(String, "name", "help").
		AddField(String, "instance", "instance help").
		DefineActions(create, update, maintain).
		OnActions(create).AddFlag("name", nil).AddFlag("instance", nil).
		OnActions(update).AddFlag("name", nil)
	if o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(repo).err)
		return
	}

	l := o.CreateList("to_create", ",", "#w(:#ft)?")
	// --- Check internal actions_related list --- Must decrease because create repo has 2 flags,
	// while update repo has only one flag.
	// If we create a list with name AND Instance, only 'create repos' can be used.
	if len(l.actions_related) != 3 {
		t.Errorf("Expected to have all object actions as possible. Got '%d'", len(l.actions_related))
	}

	l.Field(1, "name")
	if len(l.actions_related) != 2 {
		t.Errorf("Expected to have 2 object actions as possible. Got '%s'", len(l.actions_related))
	}

	l.Field(3, "instance")
	if len(l.actions_related) != 1 {
		t.Errorf("Expected to have 1 object action as possible. Got '%s'", len(l.actions_related))
	}
	if _, found := l.actions_related[create]; !found {
		t.Errorf("Expected '%s' as possible actions. Not found.", create)
	}

	// --- Run the test ---
	// This action should create a new command with one argument managed by the ObjectList object.
	l_ret := l.AddActions(create)

	// --- Start testing ---
	// Check in cli
	if l != l_ret {
		t.Error("Expected AddActions() to return the list object. Is not.")
	}
	if _, found := l.actions[create]; !found {
		t.Errorf("Expected '%s' as accepted actions. Not found.", create)
		return
	}
	if v, found := l.actions[create].params[repos]; !found {
		t.Errorf("Expected '%s' to have an argument representing the list of '%s'. Not found.", create, repos)
	} else {
		if r := reflect.TypeOf(v); r.String() != "*cli.ForjArgList" {
			t.Errorf("Expected '%s' to be an argument. But is '%s'.", repos, r.String())
		}
	}

	// check in kingpin
	if app.GetCommand(create, repos) == nil {
		t.Errorf("Expected '%s' to be created as command. Got nil.", repos)
		return
	}

	arg := app.GetArg(create, repos, reposlist)
	if arg == nil {
		t.Errorf("Expected '%s' to be created as Argument for Command '%s'. Got nil.", "repos-list", repos)
		return
	}
	if arg.GetName() != reposlist {
		t.Errorf("Expected Argument '%s' to be called '%s'. But got '%s'", repos, reposlist, arg.GetName())
	}

	// --- Run another test on the same context ---
	l_ret = l.AddActions(update)
	// --- Start testing ---
	if l != l_ret {
		t.Error("Expected AddActions() to return the list object. Is not.")
	}
}
