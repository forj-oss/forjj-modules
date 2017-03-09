package cli

import (
	"testing"
	"github.com/forj-oss/forjj-modules/cli/kingpinMock"
)

func TestForjCli_LoadAppData(t *testing.T) {
	const (
		c_test = "test"
		c_test2 = "test2"
		c_test_value = "test value"
		c_cmd = "cmd:"
		c_app = "internal_name"
	)

	// --- Setting test context ---
	t.Log("Expect App data to be collected to forjj_values")

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.AddAppFlag(String, c_test, "", nil)

	// Run the test
	if _, err := c.Parse([]string{c_test2, c_test_value}, nil) ; err != nil {
		t.Errorf("Expected to succeed. But it fails. %s", err)
	}

	// verify results
	v, found, deflt, err := c.GetStringValue(internal_app, c_app, c_test2)
	if err == nil {
		t.Error("Expected to fail. But it succeeds.")
	}
	if found {
		t.Errorf("Expected to not found any value. Found one. v = %#v", v)
	}
	if deflt {
		t.Errorf("Expected to not getting default value - Not set. Found '%s'", v)
	}
	if v != "" {
		t.Errorf("Expected to not get any value. Got '%s'", v)
	}

	// Run another test
	if _, err := c.Parse([]string{c_test, c_test_value}, nil) ; err != nil {
		t.Errorf("Expected to succeed. But it fails. %s", err)
	}

	// verify results
	v, found, deflt, err = c.GetStringValue(internal_app, c_app, c_test)
	if err != nil {
		t.Errorf("Expected to succeed. But it fails. %s", err)
	}
	if !found {
		t.Error("Expected to found a value. Found none.")
	}
	if deflt {
		t.Errorf("Expected to not getting default value - Not set. Found '%s'", v)
	}
	if v != c_test_value {
		t.Errorf("Expected to not get '%s'. Got '%s'", c_test_value, v)
	}
}
