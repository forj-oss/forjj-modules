package cli

import (
	"forjj-modules/cli/kingpinMock"
	"testing"
)

func TestNewField(t *testing.T) {
	t.Log("Expect NewField() to create a new Field object.")

	// --- Setting test context ---
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	o := c.NewObject("test", "", "")
	opts := Opts()

	// --- Run the test ---
	res := NewField(o, String, "field", "field help", ".*", opts)

	// --- Start testing ---
	if res == nil {
		t.Error("Expected Field object to be created. Got nil")
		return
	}
	if res.obj != o {
		t.Error("Expected NewField to create a field which refer to the object. Got different one.")
	}
	if res.options != opts {
		t.Error("Expected NewField to create a field which refer to options. Got different one.")
	}
	if res.name != "field" {
		t.Errorf("Expected NewField to create a field named '%s'. Got '%s'.", "field", res.name)
	}
	if res.help != "field help" {
		t.Errorf("Expected NewField to create a field help to '%s'. Got '%s'.", "field help", res.help)
	}
	if res.regexp != ".*" {
		t.Errorf("Expected NewField to create a field regexp with '%s' string. Got '%s'.", "%s", res.regexp)
	}
	if res.value_type != String {
		t.Errorf("Expected NewField to create a field type '%s'. Got '%s'.", String, res.value_type)
	}
}
