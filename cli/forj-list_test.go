package cli

import (
	"fmt"
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

	l := o.CreateList("to_create", ",", "[blabla]]", repo_help)
	if l != nil {
		t.Error("Expected CreateList() to return nil if the regexp is failing. But got one list.")
	}

	l = o.CreateList("to_create", ",", "[[blabla]", repo_help)
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
	o := c.NewObject(repo, repo_help, true).
		AddKey(String, "name", "name help", "#w", nil).
		AddField(String, "name2", "name2 help", "#ft", nil)

	l := o.CreateList("to_create", ",", "name", repo_help)
	if l == nil {
		t.Errorf("Expected list to be created. Got nil. %s", o.Error())
		return
	}
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

	l = o.CreateList("another_list", ",", "name[:name2]", repo_help)
	if l == nil {
		t.Errorf("Expected list to be created. Got nil. %s", o.Error())
		return
	}
	if l.name != "another_list" {
		t.Errorf("Expected list name to be '%s'. Got '%s'", "another_list", l.name)
	}
	expected_reg = "(" + w_f + ")(:(" + ft_f + "))?"
	if l.ext_regexp.String() != expected_reg {
		t.Errorf("Expected list regexp to be '%s'. Got '%s'", expected_reg, l.ext_regexp)
	}

	l = o.CreateList("another_list2", ",", "[name/]name[:name2[:name]]", repo_help)
	if l == nil {
		t.Errorf("Expected list to be created. Got nil. %s", o.Error())
		return
	}
	if l.name != "another_list2" {
		t.Errorf("Expected list name to be '%s'. Got '%s'", "another_list2", l.name)
	}
	expected_reg = "((" + w_f + ")/)?(" + w_f + ")(:(" + ft_f + ")(:(" + w_f + "))?)?"
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
		AddKey(String, "name", "help", "#w", nil).
		DefineActions(create).
		OnActions().
		AddFlag("name", nil)
	if o == nil {
		t.Errorf("Expected context failure. %s", c.GetObject(repo).Error())
		return
	}

	l := o.CreateList("to_create", ",", "name[:instance]", repo_help)
	if l != nil {
		t.Error("Expected CreateList() to return the object list. Got one. ")
		return
	}
	if _, found := o.list["to_create"]; found {
		t.Errorf("Expected to not have list '%s' created. But got it.", "to_create")
	}

	o.AddField(String, "instance", "instance help", "ft", nil).
		OnActions().
		AddFlag("instance", nil)

	l = o.CreateList("to_create", ",", "name[:instance]", repo_help)
	if l == nil {
		t.Errorf("Expected CreateList() to return the object list. Got Nil. %s", o.Error())
		return
	}

	field, found := l.fields_name[1]
	if !found {
		t.Errorf("Expected Field '%s' to be added. Not found.", "name")
	}
	if field != "name" {
		t.Errorf("Expected new field to be named '%s'. Got '%s'.", "name", field)
	}

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
		maintain_help = "maintain help"
	)

	c := NewForjCli(app)

	c.NewActions(create, create_help, "create %s", false)
	c.NewActions(update, update_help, "update %s", false)
	c.NewActions(maintain, maintain_help, "maintain %s", false)

	c.AddFieldListCapture("w", w_f)
	c.AddFieldListCapture("ft", ft_f)

	o := c.NewObject(repo, repo_help, true).
		AddKey(String, "name", "help", "#w", nil).
		AddField(String, "instance", "instance help", "#ft", nil).
		DefineActions(create, update, maintain).
		OnActions(create).AddFlag("name", nil).AddFlag("instance", nil).
		OnActions(update).AddFlag("name", nil)
	if o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(repo).Error())
		return
	}

	l := o.CreateList("to_create", ",", "name[:instance]", repo_help)
	// --- Check internal actions_related list --- Must decrease because create repo has 2 flags,
	// while update repo has only one flag.
	// If we create a list with name AND Instance, only 'create repos' can be used.
	if l == nil {
		t.Error("Expected list to be created. But is not.")
		return
	}

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

	arg := app.GetArg(create, repos, repos)
	if arg == nil {
		t.Errorf("Expected '%s' to be created as Argument for Command '%s'. Got nil.", repos, repos)
		return
	}
	if arg.GetName() != repos {
		t.Errorf("Expected Argument '%s' to be called '%s'. But got '%s'", repos, repos, arg.GetName())
	}

	// --- Run another test on the same context ---
	l_ret = l.AddActions(update)
	// --- Start testing ---
	if l != l_ret {
		t.Error("Expected AddActions() to return the list object. Is not.")
	}
}

func TestForjObjectList_Set(t *testing.T) {
	t.Log("Expect ForjObjectList_Set() to create a data list from cli setup.")

	// --- Setting test context ---
	const (
		repo_help       = "repo help"
		repo            = "repo"
		repos           = "repos"
		maintain_help   = "maintain help"
		f_name          = "name"
		f_name_help     = "field name help"
		f_instance      = "instance"
		f_instance_help = "Field instance help"
	)

	c := NewForjCli(app)

	c.NewActions(create, create_help, "create %s", false)
	c.NewActions(update, update_help, "update %s", false)
	c.NewActions(maintain, maintain_help, "maintain %s", false)

	c.AddFieldListCapture("w", w_f)

	o := c.NewObject(repo, repo_help, true).
		AddKey(String, f_name, f_name_help, "#w", nil).
		AddField(String, f_instance, f_instance_help, "#w", nil).
		DefineActions(create).
		OnActions().AddFlag(f_name, nil).AddFlag(f_instance, nil)
	if o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(repo).Error())
		return
	}

	l := o.CreateList("to_create", ",", "name[:instance]", repo_help).
		AddActions(create)
	if l == nil {
		t.Errorf("Expected Context list declaration to work. %s", c.GetObject(repo).Error())
		return
	}
	// --- Run the test ---
	err := l.Set("blabla")
	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected Set() to work properly. Got '%s'", err)
	}
	if len(l.context) != 1 {
		t.Errorf("Expected to find at least one record. Got '%d' records.", len(l.context))
	}
	if v, found := l.context[0].Data[f_name]; !found {
		t.Errorf("Expected to find out '%s'. But got nothing.", f_name)
	} else {
		if v != "blabla" {
			t.Errorf("Expected to find out '%s' = '%s'. But got '%s'.", f_name, "blabla", v)
		}
	}
	if v, found := l.context[0].Data[f_instance]; found && v != "" {
		t.Errorf("Expected to not found any '%s'. But got one with '%s'.", f_instance, v)
	}
	// --- Run the test ---
	err = l.Set("value:instance")
	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected Set() to work properly. Got '%s'", err)
	}
	if len(l.context) != 2 {
		t.Errorf("Expected to find at least 2 records. Got '%d' records.", len(l.context))
	}
	if v, found := l.context[1].Data[f_name]; !found {
		t.Errorf("Expected to find out '%s'. But got nothing.", f_name)
	} else {
		if v != "value" {
			t.Errorf("Expected to find out '%s' = '%s'. But got '%s'.", f_name, "value", v)
		}
	}
	if v, found := l.context[1].Data[f_instance]; !found {
		t.Errorf("Expected to find out '%s'. But got nothing.", f_instance)
	} else {
		if v != "instance" {
			t.Errorf("Expected to find out '%s' = '%s'. But got '%s'.", f_instance, "instance", v)
		}
	}
	// --- Run the test ---
	err = l.Set("last,result:instance2")
	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected Set() to work properly. Got '%s'", err)
	}
	if len(l.context) != 4 {
		t.Errorf("Expected to find at least 4 records. Got '%d' records.", len(l.context))
	}
	if v, found := l.context[2].Data[f_name]; !found {
		t.Errorf("Expected to find out '%s'. But got nothing.", f_name)
	} else {
		if v != "last" {
			t.Errorf("Expected to find out '%s' = '%s'. But got '%s'.", f_name, "last", v)
		}
	}
	if v, found := l.context[2].Data[f_instance]; found && v != "" {
		t.Errorf("Expected to not found any '%s'. But got one with '%s'.", f_instance, v)
	}
	if v, found := l.context[3].Data[f_name]; !found {
		t.Errorf("Expected to find out '%s'. But got nothing.", f_name)
	} else {
		if v != "result" {
			t.Errorf("Expected to find out '%s' = '%s'. But got '%s'.", f_name, "result", v)
		}
	}
	if v, found := l.context[3].Data[f_instance]; !found {
		t.Errorf("Expected to find out '%s'. But got nothing.", f_instance)
	} else {
		if v != "instance2" {
			t.Errorf("Expected to find out '%s' = '%s'. But got '%s'.", f_instance, "instance2", v)
		}
	}
}

func TestForjObjectList_AddValidateHandler(t *testing.T) {
	t.Log("Expect AddActions() to add some action for the list.")
	// --- Setting test context ---
	const (
		repo_help       = "repo help"
		repo            = "repo"
		repos           = "repos"
		maintain_help   = "maintain help"
		f_name          = "name"
		f_name_help     = "field name help"
		f_instance      = "instance"
		f_instance_help = "Field instance help"
	)

	c := NewForjCli(app)

	c.NewActions(create, create_help, "create %s", false)
	c.NewActions(update, update_help, "update %s", false)
	c.NewActions(maintain, maintain_help, "maintain %s", false)

	c.AddFieldListCapture("w", w_f)

	o := c.NewObject(repo, repo_help, true).
		AddKey(String, f_name, f_name_help, "#w", nil).
		AddField(String, f_instance, f_instance_help, "#w", nil).
		DefineActions(create).
		OnActions().AddFlag(f_name, nil).AddFlag(f_instance, nil)
	if o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(repo).Error())
		return
	}

	l := o.CreateList("to_create", ",", "name[:instance]", repo_help).
		AddActions(create)
	if l == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(repo).Error())
		return
	}

	// --- Run the test ---
	valid_handler := func(d *ForjListData) error {
		var name string
		if v, found := d.Data[f_name]; !found || v == "" {
			return fmt.Errorf("Field '%s' is missing or empty.", f_name)
		} else {
			name = v
		}

		if v, found := d.Data["instance"]; !found || v == "" {
			d.Data["instance"] = name
		}
		return nil
	}
	l_ret := l.AddValidateHandler(valid_handler)

	// --- Start testing ---
	// Check in cli
	if l != l_ret {
		t.Error("Expected AddValidateHandler() to return the list object. Is not.")
	}
	if l.valid_handler == nil {
		t.Error("Expected AddValidateHandler() to store the handler. Is not.")
	}
	// --- Update test context ---
	err := l.Set("last,result:instance2,")
	// --- Start testing ---
	if err == nil {
		t.Error("Expected Set() to return an error. Got none.")
	}
	if len(l.context) != 2 {
		t.Errorf("Expected Set to save 2 records. Got %d", len(l.context))
		return
	}
	if v, found := l.context[0].Data[f_name]; !found {
		t.Errorf("Expected Set to add a record with '%s' field. Not found.", f_name)
	} else {
		if v != "last" {
			t.Errorf("Expected set to have field '%s' = '%s' for record %d. Got '%s'", f_name, "last", 0, v)
		}
	}
	if v, found := l.context[0].Data[f_instance]; !found {
		t.Errorf("Expected Set to add a record with '%s' field. Not found.", f_instance)
	} else {
		if v != "last" {
			t.Errorf("Expected set to have field '%s' = '%s' for record %d. Got '%s'", f_instance, "last", 0, v)
		}
	}
	if v, found := l.context[1].Data[f_name]; !found {
		t.Errorf("Expected Set to add a record with '%s' field. Not found.", f_name)
	} else {
		if v != "result" {
			t.Errorf("Expected set to have field '%s' = '%s' for record %d. Got '%s'", f_name, "result", 1, v)
		}
	}
	if v, found := l.context[1].Data[f_instance]; !found {
		t.Errorf("Expected Set to add a record with '%s' field. Not found.", f_instance)
	} else {
		if v != "instance2" {
			t.Errorf("Expected set to have field '%s' = '%s' for record %d. Got '%s'", f_instance, "instance2", 1, v)
		}
	}
}
