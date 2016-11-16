package cli

import (
	"github.com/forj-oss/forjj-modules/cli/interface"
	"github.com/forj-oss/forjj-modules/cli/kingpinMock"
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

	c.NewObject(test, test_help, false).
		AddKey(String, key, key_help).
		AddField(String, flag, flag_help).
		DefineActions(update).OnActions().
		AddArg(key, Opts().Required()).
		AddFlag(flag, nil).
		CreateList("to_update", ",", "#w").
		Field(1, key).
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
	if c.AddActionFlagFromObjectListAction(create, test, "to_update", update) == nil {
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

	c.NewObject(test, test_help, false).
		AddKey(String, key, key_help).
		AddField(String, flag, flag_help).
		DefineActions(update).OnActions().
		AddFlag(key, Opts().Required()).
		AddFlag(flag, nil)
	context := []string{"cmd:" + update, "cmd:" + test, key, key_value, flag, flag_value}

	_, err := c.Parse(context, nil)
	if err != nil {
		t.Errorf("Expected Parse() to work successfully. Got '%s'", err)
	}

	// --- Run the test ---
	ret, found := c.GetStringValue(test, key_value, flag)

	// --- Start testing ---
	if !found {
		t.Error("Expected GetStringValue() to find the value. Not found")
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

	c.NewObject(test, test_help, false).
		AddKey(String, key, key_help).
		AddField(Bool, flag, flag_help).
		DefineActions(update).OnActions().
		AddArg(key, Opts().Required()).
		AddFlag(flag, nil)

	context := []string{"cmd:" + update, "cmd:" + test, key, key_value, flag, flag_value}

	_, err := c.Parse(context, nil)
	if err != nil {
		t.Errorf("Expected Parse() to work successfully. Got '%s'", err)
	}

	// --- Run the test ---
	ret, found := c.GetBoolValue(test, key_value, flag)

	// --- Start testing ---
	if !found {
		t.Error("Expected GetStringValue() to find the value. Not found")
	}
	if ret != true {
		t.Errorf("Expected GetBoolValue() to return '%t'. Got '%t'", true, ret)
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
	v := c.GetAppBoolValue(flag)
	// --- Start testing ---
	if v != false {
		t.Error("Expected GetAppBoolValue() to return false. Got true")
	}

	// --- Updating test context --- parse time
	context.SetContextAppValue(flag, flag_value)
	// --- Run the test ---
	v = c.GetAppBoolValue(flag)

	// --- Start testing ---
	if v != true {
		t.Error("Expected GetAppBoolValue() to return true. Got false")
	}

	// --- Update test context ---
	c.parse = true // Parse time over.

	// --- Run the test ---
	v = c.GetAppBoolValue(flag)

	// --- Start testing ---
	if v != false {
		t.Error("Expected GetAppBoolValue() to return false. Got true")
	}

	// --- Update test context ---
	context.SetParsedAppValue(flag, flag_value)

	// --- Run the test ---
	v = c.GetAppBoolValue(flag)
	// --- Start testing ---
	if v != true {
		t.Error("Expected GetAppBoolValue() to return true. Got false")
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

	c.cli_context.context = clier.ParseContexter(app.NewContext().SetContextAppValue(flag, flag_value))
	// --- Run the test ---
	v := c.GetAppStringValue(flag)
	// --- Start testing ---
	if v != flag_value {
		t.Error("Expected GetAppBoolValue() top return true. Got false")
	}
}
