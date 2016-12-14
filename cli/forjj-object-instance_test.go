package cli

import (
	"github.com/forj-oss/forjj-modules/cli/kingpinMock"
	"testing"
)

func TestNewObjectInstance(t *testing.T) {
	t.Log("Expect NewObjectInstance() to be initialized.")

	// --- Setting test context ---

	// --- Run the test ---
	v := NewObjectInstance("test")

	// --- Start testing ---
	if v == nil {
		t.Error("Expected to get a new ObjectInstance created. Got nil")
		return
	}
	if v.name != "test" {
		t.Errorf("Expected ObjectInstanace name to be '%s'. Got '%s'", "test", v.name)
	}
	if v.additional_fields == nil {
		t.Error("Expected ObjectInstanace fields to be initialized. Got nil")
	}
}

func TestForjObjectInstance_addField(t *testing.T) {
	t.Log("Expect ForjObjectInstance_addField() to add field to an ObjectInstance.")

	// --- Setting test context ---

	var oi *ForjObjectInstance
	// --- Run the test ---
	res := oi.addField(nil, "", "", "", "", nil)

	// --- Start testing ---
	if res != nil {
		t.Error("Expected ObjectInstance.addField() to not create or get an Object Instance. got one")
	}

	// --- Setting test context ---
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	o := c.NewObject("test", "", false)
	oi = NewObjectInstance("test")

	// --- Run the test ---
	res = oi.addField(o, String, "field", "", "", nil)
	// --- Start testing ---
	if res == nil {
		t.Errorf("Expected addField to not fail. got %s", o.err)
		return
	}
	if res != oi {
		t.Error("Expected addField to return the Object Instance called. Got a different one")
	}
	if oi.additional_fields == nil {
		t.Error("Expected ObjectInstanace fields to be initialized. Got nil")
	}
	if len(oi.additional_fields) != 1 {
		t.Errorf("Expected to get %d field in the list. got %d", 1, len(oi.additional_fields))
	}
	if v, found := oi.additional_fields["field"]; !found {
		t.Errorf("Expected field '%s' to be added. Got none.", "field")
	} else {
		if v.name != "field" {
			t.Errorf("Expected addField to add field name '%s'. Got '%s'", "field", v.name)
		}
	}
}

func TestForjObjectInstance_hasField(t *testing.T) {
	t.Log("Expect ForjObjectInstance_hasField() to check field existence in an ObjectInstance.")

	// --- Setting test context ---
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	o := c.NewObject("test", "", false)
	oi := NewObjectInstance("test").
		addField(o, String, "field", "", "", nil)

	// --- Run the test ---
	if oi.hasField("field1") {
		t.Error("Expected hasField to return false. Got true.")
	}
	// --- Start testing ---
	if !oi.hasField("field") {
		t.Error("Expected hasField to return true. Got false.")
	}
}
