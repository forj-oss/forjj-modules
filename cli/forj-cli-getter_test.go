package cli

import (
	"github.com/forj-oss/forjj-modules/cli/interface"
	"forjj-modules/cli/kingpinMock"
	"testing"
)

func TestForjCli_Parse(t *testing.T) {
	t.Log("Expect ForjCli_Parse() to parse and update objects data automatically.")

	// --- Setting test context ---
	const (
		test       = "test"
		tests      = "tests"
		test_help  = "test help"
		flag       = "flag"
		flag_help  = "flag help"
		flag_value = "true"
		key        = "key"
		key_help   = "key help"
		key_value  = "name1,name2"
		name1      = "name1"
		name2      = "name2"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)

	c.AddFieldListCapture("w", w_f)

	c.NewObject(test, test_help, "").
		AddKey(String, key, key_help, "#w", nil).
		AddField(String, flag, flag_help, "", nil).
		DefineActions(update).OnActions().
		AddArg(key, Opts().Required()).
		AddFlag(flag, nil).
		CreateList("to_update", ",", "key", test_help).
		AddActions(update)

	context := []string{
		"cmd:" + update, "cmd:" + tests,
		tests, key_value,
		name1 + "-" + flag, flag_value,
	}

	if c.Error() != nil {
		t.Errorf("Expected context to work. Got '%s'", c.Error())
	}

	if o := c.GetObject(test); o == nil {
		t.Errorf("Expected context to work. Unable to find '%s' object", test)
		return
	} else {
		if o.Error() != nil {
			t.Errorf("Expected context to work. Got '%s'", o.Error())
		}
	}
	// --- Run the test ---
	cmd, err := c.Parse(context, nil)

	// --- Start testing ---
	if cmd != "update tests" {
		t.Errorf("Expected Parse() to return '%s'. Got '%s'", "update tests", cmd)
	}
	if err != nil {
		t.Errorf("Expected Parse() to work successfully. Got '%s'", err)
	}

	// test in cli
	if _, found := c.values[test]; !found {
		t.Errorf("Expected '%s' object to exist. Not found.", test)
		return
	}

	v := c.values[test]
	if _, found := v.records[name1]; !found {
		t.Errorf("Expected '%s' object to have '%s' as record. Not found.", test, name1)
	}

	r := v.records[name1]
	if value, found := r.attrs[flag]; !found {
		t.Errorf("Expected '%s' object to have '%s' '%s' as record field. Not found.",
			test, name1, flag)
	} else {
		if value != flag_value {
			t.Errorf("Expected '%s' '%s' '%s' = '%s'. Got '%s' ",
				test, name1, flag, flag_value, value)
		}
	}

	// Updating context
	if c.OnActions(create).AddActionFlagFromObjectListAction(test, "to_update", update) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.Error())
	}

	context = []string{
		"cmd:" + create,
		tests, key_value,
		name2 + "-" + flag, flag_value,
	}

	// --- Run the test ---
	cmd, err = c.Parse(context, nil)

	// --- Start testing ---
	if cmd != "create" {
		t.Errorf("Expected Parse() to return '%s'. Got '%s'", "create", cmd)
	}
	if err != nil {
		t.Errorf("Expected Parse() to work successfully. Got '%s'", err)
	}
	// test in cli
	if _, found := c.values[test]; !found {
		t.Errorf("Expected '%s' object to exist. Not found.", test)
		return
	}

	v = c.values[test]
	if _, found := v.records[name2]; !found {
		t.Errorf("Expected '%s' object to have '%s' as record. Not found.", test, name2)
	}

	r = v.records[name2]
	if value, found := r.attrs[flag]; !found {
		t.Errorf("Expected '%s' object to have '%s' '%s' as record field. Not found.",
			test, name2, flag)
	} else {
		if value != flag_value {
			t.Errorf("Expected '%s' '%s' '%s' = '%s'. Got '%s' ",
				test, name2, flag, flag_value, value)
		}
	}
}

func TestForjCli_GetStringValue(t *testing.T) {
	t.Log("Expect GetStringValue() to be get the Command flag value as string.")

	const (
		test       = "test"
		test_help  = "test help"
		key        = "key"
		key_help   = "key help"
		key_value  = "key-value"
		flag       = "flag"
		flag_help  = "flag help"
		flag_value = "flag value"
	)
	// --- Setting test context ---
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)

	c.NewObject(test, test_help, "").
		AddKey(String, key, key_help, "", nil).
		AddField(String, flag, flag_help, "", nil).
		DefineActions(update).OnActions().
		AddFlag(key, Opts().Required()).
		AddFlag(flag, nil)
	context := []string{"cmd:" + update, "cmd:" + test, key, key_value, flag, flag_value}

	_, err := c.Parse(context, nil)
	if err != nil {
		t.Errorf("Expected Parse() to work successfully. Got '%s'", err)
	}

	// --- Run the test ---
	ret, found, isDefault, err := c.GetStringValue(test, key_value, flag)

	// --- Start testing ---
	if !found {
		t.Error("Expected GetStringValue() to find the value. Not found")
	}
	if isDefault {
		t.Error("Expected GetStringValue() to find a real value. Not default. Found default one.")
	}
	if ret != flag_value {
		t.Errorf("Expected GetStringValue() to return '%s'. Got '%s'", flag_value, ret)
	}
}

func TestForjCli_GetBoolValue(t *testing.T) {
	t.Log("Expect ForjCli_GetBoolValue() to be get the Command flag value as bool.")
	const (
		test       = "test"
		test_help  = "test help"
		flag       = "flag"
		flag_help  = "flag help"
		flag_value = "true"
		key        = "key"
		key_help   = "key help"
		key_value  = "key-value"
	)
	// --- Setting test context ---
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)

	c.NewObject(test, test_help, "").
		AddKey(String, key, key_help, "", nil).
		AddField(Bool, flag, flag_help, "", nil).
		DefineActions(update).OnActions().
		AddArg(key, Opts().Required()).
		AddFlag(flag, nil)

	context := []string{"cmd:" + update, "cmd:" + test, key, key_value, flag, flag_value}

	if _, err := c.Parse(context, nil); err != nil {
		t.Errorf("Expected Parse() to work successfully. Got '%s'", err)
	}

	// --- Run the test ---
	ret, found, err := c.GetBoolValue(test, key_value, flag)

	// --- Start testing ---
	if !found {
		t.Errorf("Expected GetStringValue() to find the value. Not found. %s", err)
	}
	if ret != true {
		t.Errorf("Expected GetBoolValue() to return '%t'. Got '%t'", true, ret)
	}
}

func TestForjCli_GetStringValue_FromObjectListContext(t *testing.T) {
	t.Log("Expect GetStringValue() to get flag value from object list flag (context or final).")

	const (
		test             = "test"
		test_help        = "test help"
		key              = "key"
		key_help         = "key help"
		key_value        = "key-value"
		flag             = "flag"
		flag_help        = "flag help"
		flag_value       = "flag value"
		myapp            = "app"
		apps             = "apps"
		app_help         = "app help"
		instance         = "instance"
		instance_help    = "instance help"
		driver           = "driver"
		driver_help      = "driver help"
		driver_type      = "driver_type"
		driver_type_help = "driver_type help"
		flag2            = "flag2"
		flag2_help       = "flag2 help"
		flag2_value      = "flag2 value"
		myinstance       = "myapp"
	)
	// --- Setting test context ---
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)
	c.AddFieldListCapture("w", w_f)

	if c.NewObject(test, test_help, "").
		AddKey(String, key, key_help, "", nil).
		AddField(String, flag, flag_help, "", nil).
		DefineActions(update).OnActions().
		AddFlag(key, Opts().Required()).
		AddFlag(flag, nil) == nil {
		t.Error(c.GetObject(test).Error())
	}

	if c.NewObject(myapp, app_help, "").
		AddKey(String, instance, instance_help, "#w", nil).
		AddField(String, driver, driver_help, "#w", nil).
		AddField(String, driver_type, driver_type_help, "#w", nil).
		AddField(String, flag2, flag2_help, "", nil).
		ParseHook(func(_ *ForjObject, c *ForjCli, _ interface{}) (err error, updated bool) {
		ret, found, _, err := c.GetStringValue(myapp, myinstance, flag2)
		if found {
			t.Error("Expected GetStringValue() to NOT find the context value. Got one.")
		}
		if ret != "" {
			t.Errorf("Expected GetStringValue() to return '' from context. Got '%s'", ret)
		}

		ret, found, _, err = c.GetStringValue(test, key_value, flag)
		if !found {
			t.Errorf("Expected GetStringValue() to find the context value. Got none. %s", err)
		}
		if ret != flag_value {
			t.Errorf("Expected GetStringValue() to return '%s' from context. Got '%s'", flag_value, ret)
		}

		ret, found, _, err = c.GetStringValue(test, "", flag)
		if !found {
			t.Errorf("Expected GetStringValue() to find the context value. Got none. %s", err)
		}
		if ret != flag_value {
			t.Errorf("Expected GetStringValue() to return '%s' from context. Got '%s'", flag_value, ret)
		}
		return nil, false
	}).DefineActions(create).OnActions().
		AddFlag(driver_type, nil).
		AddFlag(driver, nil).
		AddFlag(instance, Opts().Required()).
		AddFlag(flag2, nil).
		CreateList("to_create", ",", "driver_type:driver[:instance]", app_help).
		AddValidateHandler(func(l *ForjListData) (err error) {
		if v, found := l.Data[instance]; !found || v == "" {
			l.Data[instance] = l.Data[driver]
		}
		return nil
	}) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(myapp).Error())
	}

	c.GetObject(test).AddFlagFromObjectListAction(myapp, "to_create", create)

	context := []string{"cmd:" + update, "cmd:" + test, key, key_value, flag, flag_value,
		apps, "mytype:mydriver", "mydriver-flag2", flag2_value}

	if _, err := c.Parse(context, nil); err != nil {
		t.Errorf("Expected Parse() to work successfully. Got '%s'", err)
	}

	// --- Run the test ---
	ret, found, _, err := c.GetStringValue(myapp, "mydriver", flag2)

	// --- Start testing ---
	if !found {
		t.Errorf("Expected GetStringValue() to find the value. Not found. %s", err)
	}
	if ret != flag2_value {
		t.Errorf("Expected GetStringValue() to return '%s'. Got '%s'", flag2_value, ret)
	}
}

func TestForjCli_GetBoolValue_FromObjectListContext(t *testing.T) {
	t.Log("Expect GetBoolValue() to get flag value from object list flag (context or final).")

	const (
		test             = "test"
		test_help        = "test help"
		key              = "key"
		key_help         = "key help"
		key_value        = "key-value"
		flag             = "flag"
		flag_help        = "flag help"
		flag_value       = "true"
		myapp            = "app"
		apps             = "apps"
		app_help         = "app help"
		instance         = "instance"
		instance_help    = "instance help"
		driver           = "driver"
		driver_help      = "driver help"
		driver_type      = "driver_type"
		driver_type_help = "driver_type help"
		flag2            = "flag2"
		flag2_help       = "flag2 help"
		flag2_value      = "true"
		myinstance       = "myapp"
	)
	// --- Setting test context ---
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)
	c.AddFieldListCapture("w", w_f)

	if c.NewObject(test, test_help, "").
		AddKey(String, key, key_help, "", nil).
		AddField(Bool, flag, flag_help, "", nil).
		DefineActions(update).OnActions().
		AddFlag(key, Opts().Required()).
		AddFlag(flag, nil) == nil {
		t.Error(c.GetObject(test).Error())
	}

	if c.NewObject(myapp, app_help, "").
		AddKey(String, instance, instance_help, "#w", nil).
		AddField(String, driver, driver_help, "#w", nil).
		AddField(String, driver_type, driver_type_help, "#w", nil).
		AddField(Bool, flag2, flag2_help, "", nil).
		ParseHook(func(_ *ForjObject, c *ForjCli, _ interface{}) (err error, updated bool) {
		ret, found, err := c.GetBoolValue(myapp, myinstance, flag2)
		if found {
			t.Error("Expected GetStringValue() to NOT find the context value. Got one.")
		}
		if ret {
			t.Error("Expected GetStringValue() to return 'false' from context. Got 'true'")
		}

		ret, found, err = c.GetBoolValue(test, key_value, flag)
		if !found {
			t.Errorf("Expected GetStringValue() to find the context value. Got none. %s", err)
		}
		if !ret {
			t.Error("Expected GetStringValue() to return 'true' from context. Got 'false'")
		}

		ret, found, err = c.GetBoolValue(test, "", flag)
		if !found {
			t.Errorf("Expected GetStringValue() to find the context value. Got none. %s", err)
		}
		if !ret {
			t.Error("Expected GetStringValue() to return 'true' from context. Got 'false'")
		}
		return nil, false
	}).DefineActions(create).OnActions().
		AddFlag(driver_type, nil).
		AddFlag(driver, nil).
		AddFlag(instance, Opts().Required()).
		AddFlag(flag2, nil).
		CreateList("to_create", ",", "driver_type:driver[:instance]", app_help).
		AddValidateHandler(func(l *ForjListData) (err error) {
		if v, found := l.Data[instance]; !found || v == "" {
			l.Data[instance] = l.Data[driver]
		}
		return nil
	}) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(myapp).Error())
	}

	c.GetObject(test).AddFlagFromObjectListAction(myapp, "to_create", create)

	context := []string{"cmd:" + update, "cmd:" + test, key, key_value, flag, flag_value,
		apps, "mytype:mydriver", "mydriver-flag2", flag2_value}

	if _, err := c.Parse(context, nil); err != nil {
		t.Errorf("Expected Parse() to work successfully. Got '%s'", err)
	}

	// --- Run the test ---
	ret, found, err := c.GetBoolValue(myapp, "mydriver", flag2)

	// --- Start testing ---
	if !found {
		t.Errorf("Expected GetStringValue() to find the value. Not found. %s", err)
	}
	if !ret {
		t.Error("Expected GetStringValue() to return 'true'. Got 'false'")
	}
}

func TestForjCli_GetAppBoolValue(t *testing.T) {
	t.Log("Expect ForjCli_GetAppBoolValue to .")

	// --- Setting test context ---
	const (
		flag       = "flag"
		flag_help  = "flag help"
		flag_value = "true"
	)
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.AddAppFlag(Bool, flag, flag_help, nil)

	context := app.NewContext()
	c.cli_context.context = clier.ParseContexter(context)
	// --- Run the test ---
	v, err := c.GetAppBoolValue(flag)
	// --- Start testing ---
	if v != false {
		t.Error("Expected GetAppBoolValue() to return false. Got true")
	}
	if err != nil {
		t.Errorf("Expected have no error. Got %s.", err)
	}

	// --- Updating test context --- parse time
	context.SetContextAppValue(flag, flag_value)

	// --- Run the test ---
	v, err = c.GetAppBoolValue(flag)

	// --- Start testing ---
	if v != true {
		t.Error("Expected GetAppBoolValue() to return true. Got false.")
	}
	if err != nil {
		t.Errorf("Expected no error. Got one. %s", err)
	}

	// --- Update test context ---
	c.parse = true // Parse time over.

	// --- Run the test ---
	v, err = c.GetAppBoolValue(flag)

	// --- Start testing ---
	if v != false {
		t.Error("Expected GetAppBoolValue() to return false. Got true")
	}
	if err != nil {
		t.Errorf("Expected no error. Got one. %s", err)
	}

	// --- Update test context ---
	context.SetParsedAppValue(flag, flag_value)

	// --- Run the test ---
	v, err = c.GetAppBoolValue(flag)
	// --- Start testing ---
	if v != true {
		t.Error("Expected GetAppBoolValue() to return true. Got false")
	}
	if err != nil {
		t.Errorf("Expected no error. Got one. %s", err)
	}

}

func TestForjCli_GetAppStringValue(t *testing.T) {
	t.Log("Expect ForjCli_GetAppStringValue() to get App string value.")

	// --- Setting test context ---
	const (
		flag       = "flag"
		flag_help  = "flag help"
		flag_value = "true"
	)
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.AddAppFlag(String, flag, flag_help, nil)

	context := app.NewContext()
	c.cli_context.context = clier.ParseContexter(context)
	// --- Run the test ---
	v, err := c.GetAppStringValue(flag)
	// --- Start testing ---
	if v != false {
		t.Error("Expected GetAppBoolValue() to return false. Got true")
	}
	if err != nil {
		t.Errorf("Expected have no error. Got %s.", err)
	}

	c.cli_context.context = clier.ParseContexter(app.NewContext().SetContextAppValue(flag, flag_value))
	// --- Run the test ---
	v, err = c.GetAppStringValue(flag)
	// --- Start testing ---
	if v != flag_value {
		t.Errorf("Expected GetAppBoolValue() top return '%s'. Got '%s'", flag_value, v)
	}
	if err != nil {
		t.Errorf("Expected no error. Got one. %s", err)
	}
}

func TestForjCli_GetActionStringValue(t *testing.T) {
	t.Log("Expect ForjCli_GetActionStringValue() to get action flag string value.")

	// --- Setting test context ---
	const (
		flag = "flag"
		flag_value = "true"
		no_maintain_f = "no-maintain"
	)
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, "", "", false)
	c.OnActions(create).AddFlag(String, no_maintain_f, "", nil)

	if ctx, err := app.NewContext().
		SetContext(create).
		SetContextValue(no_maintain_f, flag_value) ; err != nil {
		t.Error("Unable to set context.")
	} else {
		c.cli_context.context = clier.ParseContexter(ctx)
	}

	// --- Run the test ---
	v, err := c.GetActionStringValue(create, no_maintain_f)
	// --- Start testing ---
	if v != flag_value {
		t.Errorf("Expected GetActionStringValue() top return '%s'. Got '%s'", flag_value, v)
	}
	if err != nil {
		t.Errorf("Expected no error. Got one. %s", err)
	}
}

func TestForjCli_GetActionBoolValue(t *testing.T) {
	t.Log("Expect ForjCli_GetActionBoolValue() to get action flag bool value.")

	// --- Setting test context ---
	const (
		flag = "flag"
		flag_value = true
		no_maintain_f = "no-maintain"
	)
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, "", "", false)
	c.OnActions(create).AddFlag(Bool, no_maintain_f, "", nil)

	if ctx, err := app.NewContext().
		SetContext(create).
		SetContextValue(no_maintain_f, "true"); err != nil {
		t.Error("Unable to set context.")
	} else {
		c.cli_context.context = clier.ParseContexter(ctx)
	}
	// --- Run the test ---
	v, err := c.GetActionBoolValue(create, no_maintain_f)
	// --- Start testing ---
	if v != flag_value {
		t.Errorf("Expected GetActionBoolValue() top return '%s'. Got '%s'", flag_value, v)
	}
	if err != nil {
		t.Errorf("Expected no error. Got one. %s", err)
	}
}
