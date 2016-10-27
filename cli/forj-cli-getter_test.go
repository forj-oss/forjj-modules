package cli

import (
	"github.com/forj-oss/forjj-modules/cli/kingpinMock"
	"testing"
)

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
		DefineActions(update).
		AddFlag(key, Opts().Required()).
		AddFlag(flag, nil)
	app.NewContext().SetContext(update, test).
		SetContextValue(key, key_value).
		SetContextValue(flag, flag_value)

	c.LoadContext([]string{update, test, "--" + flag, flag_value})
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
		DefineActions(update).
		AddArg(key, Opts().Required()).
		AddFlag(flag, nil)
	app.NewContext().SetContext(update, test).
		SetContextValue(key, key_value).
		SetContextValue(flag, flag_value)

	c.LoadContext([]string{update, test, "--" + flag, flag_value})

	// --- Run the test ---
	ret, found := c.GetBoolValue(test, key_value, flag)

	// --- Start testing ---
	if !found {
		t.Error("Expected GetStringValue() to find the value. Not found")
	}
	if ret != true {
		t.Errorf("Expected GetStringValue() to return '%t'. Got '%t'", true, ret)
	}
}

func TestForjCli_GetAppBoolValue(t *testing.T) {
	t.Log("Expect ForjCli_GetAppBoolValue to .")

	// --- Setting test context ---

	// --- Run the test ---

	// --- Start testing ---
}
