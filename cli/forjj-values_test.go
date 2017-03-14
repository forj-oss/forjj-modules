package cli

import (
	"testing"
	"sort"
)

func TestForjData_Keys(t *testing.T) {
	t.Log("Expect ForjData_Keys() to return list of keys.")

	// --- Setting test context ---
	anStr := "string2"
	r := ForjData{
		attrs: map[string]interface{} {
			"test1": "string",
			"test2": &anStr,
		},
	}
	// --- Run the test ---
	keys := r.Keys()
	// --- Start testing ---
	if keys == nil {
		t.Error("Expected keys to exist. Got nil")
		return
	}
	if v := len(keys) ; v != 2 {
		t.Errorf("Expected Keys() to return %d keys. Got %d keys", 2, v)
	}
	sort.Strings(keys)
	if v := keys[0] ; v != "test1" {
		t.Errorf("Expected to have '%s' key in list. Got '%s'", "test1", v)
	}
	if v := keys[1] ; v != "test2" {
		t.Errorf("Expected to have '%s' key in list. Got '%s'", "test2", v)
	}
}
