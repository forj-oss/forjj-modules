package kingpinCli

import (
	"github.com/alecthomas/kingpin"
	"testing"
)

func TestParseContext_GetArgValue(t *testing.T) {
	t.Log("Expect ParseContext_GetArgValue() to return the appropriate context value.")

	const (
		test       = "test"
		test_help  = "test help"
		test_value = "test value"
		myapp      = "myapp"
	)
	// --- Setting test context ---
	a := New(kingpin.New(test, test_help), myapp)
	c1 := a.Command("add", "")
	c2 := c1.Command("test", "")

	a3 := c2.Arg("arg", test_help)
	a3.String()

	// --- Run the test ---
	c, err := a.ParseContext([]string{"add", test, test_value})
	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected GetContext to not return an error. Got %s", err)
		return
	}

	v, found := c.GetArgValue(a3)
	if !found {
		t.Errorf("Expected GetArgValue() to get '%s' arg. Not found", test)
	}
	if v != test_value {
		t.Errorf("Expected GetArgValue() to get and return '%s'. But Got '%s'", test_value, v)
	}
}

func TestParseContext_GetFlagValue(t *testing.T) {
	t.Log("Expect ParseContext_GetFlagValue() to return the appropriate context value.")

	const (
		test       = "test"
		test_help  = "test help"
		test_value = "test value"
		myapp      = "myapp"
	)
	// --- Setting test context ---
	a := New(kingpin.New(test, test_help), myapp)
	c1 := a.Command("add", "")
	c2 := c1.Command("test", "")

	f3 := c2.Flag("flag", test_help)
	f3.String()

	// --- Run the test ---
	c, err := a.ParseContext([]string{"add", test, "--flag", test_value})
	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected GetContext to not return an error. Got %s", err)
		return
	}

	v, found := c.GetFlagValue(f3)
	if !found {
		t.Errorf("Expected GetFlagValue() to get '%s' arg. Not found", test)
	}
	if v != test_value {
		t.Errorf("Expected GetFlagValue() to get and return '%s'. But Got '%s'", test_value, v)
	}
}
