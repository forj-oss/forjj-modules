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

	c.NewObject(test, test_help, false).AddField(String, flag, flag_help).DefineActions(update)

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
