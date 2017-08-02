package cli

import "testing"

func TestOpts(t *testing.T) {
	t.Log("Expect Opts() to create a new Options object.")

	// --- Setting test context ---

	// --- Run the test ---
	opts := Opts()

	// --- Start testing ---
	if opts == nil {
		t.Error("Expected to get an allocated opts. But got Nil.")
		return
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid allocated opts. But got opts Nil.")
	}
}

func TestForjOpts_Default(t *testing.T) {
	t.Log("Expect Default() to add a default value in Options object.")
	// --- Setting test context ---
	opts := Opts()
	if opts == nil {
		t.Error("Expected to get an allocated opts. But got Nil.")
		return
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}

	// --- Run the test ---
	ret := opts.Default("test")

	// --- Start testing ---
	if ret == nil {
		t.Error("Expected to get the opts object at return. Got nil.")
	} else {
		if ret != opts {
			t.Error("Expected to get the opts object at return. Got a different one.")
		}
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}
	if v, found := opts.opts["default"] ; !found {
		t.Errorf("Expected to get '%s' key. Not found", "default")
	} else {
		if v != "test" {
			t.Errorf("Expected to get '%s' value '%s'. Got '%s'", "default", "test", v)
		}
	}
}

func TestForjOpts_Envar(t *testing.T) {
	t.Log("Expect Envar() to add envar in Options object.")
	// --- Setting test context ---
	opts := Opts()
	if opts == nil {
		t.Error("Expected to get an allocated opts. But got Nil.")
		return
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}

	// --- Run the test ---
	ret := opts.Envar("TEST")

	// --- Start testing ---
	if ret == nil {
		t.Error("Expected to get the opts object at return. Got nil.")
	} else {
		if ret != opts {
			t.Error("Expected to get the opts object at return. Got a different one.")
		}
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}
	if v, found := opts.opts["envar"] ; !found {
		t.Errorf("Expected to get '%s' key. Not found", "envar")
	} else {
		if v != "TEST" {
			t.Errorf("Expected to get '%s' value '%s'. Got '%s'", "envar", "TEST", v)
		}
	}
}

func TestForjOpts_NoEnvar(t *testing.T) {
	t.Log("Expect NoEnvar() to remove envar from Options object.")
	// --- Setting test context ---
	opts := Opts()
	if opts == nil {
		t.Error("Expected to get an allocated opts. But got Nil.")
		return
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}

	opts.Envar("TEST")
	// --- Run the test ---
	ret := opts.NoEnvar()

	// --- Start testing ---
	if ret == nil {
		t.Error("Expected to get the opts object at return. Got nil.")
	} else {
		if ret != opts {
			t.Error("Expected to get the opts object at return. Got a different one.")
		}
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}
	if _, found := opts.opts["envar"] ; found {
		t.Errorf("Expected to NOT get '%s' key. found", "envar")
	}
}

func TestForjOpts_NoDefault(t *testing.T) {
	t.Log("Expect NoDefault() to remove default value from Options object.")
	// --- Setting test context ---
	opts := Opts()
	if opts == nil {
		t.Error("Expected to get an allocated opts. But got Nil.")
		return
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}

	opts.Default("test")
	// --- Run the test ---
	ret := opts.NoDefault()

	// --- Start testing ---
	if ret == nil {
		t.Error("Expected to get the opts object at return. Got nil.")
	} else {
		if ret != opts {
			t.Error("Expected to get the opts object at return. Got a different one.")
		}
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}
	if _, found := opts.opts["default"] ; found {
		t.Errorf("Expected to NOT get '%s' key. found", "default")
	}
}

func TestForjOpts_Required(t *testing.T) {
		t.Log("Expect Required() to add required flag in Options object.")
	// --- Setting test context ---
	opts := Opts()
	if opts == nil {
		t.Error("Expected to get an allocated opts. But got Nil.")
		return
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}

	// --- Run the test ---
	ret := opts.Required()

	// --- Start testing ---
	if ret == nil {
		t.Error("Expected to get the opts object at return. Got nil.")
	} else {
		if ret != opts {
			t.Error("Expected to get the opts object at return. Got a different one.")
		}
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}
	if v, found := opts.opts["required"] ; !found {
		t.Errorf("Expected to get '%s' key. Not found", "required")
	} else {
		if b, ok := v.(bool) ; !ok {
			t.Errorf("Expected to get '%s' value '%s'. Got '%t'", "required", "TEST", v)
		} else {
			if ! b {
				t.Errorf("Expected to get '%s' value '%s'. Got '%t'", "required", "TEST", b)
			}
		}
	}
}

func TestForjOpts_NotRequired(t *testing.T) {
		t.Log("Expect NotRequired() to remove required flag from Options object.")
	// --- Setting test context ---
	opts := Opts()
	if opts == nil {
		t.Error("Expected to get an allocated opts. But got Nil.")
		return
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}

	opts.Required()
	// --- Run the test ---
	ret := opts.NotRequired()

	// --- Start testing ---
	if ret == nil {
		t.Error("Expected to get the opts object at return. Got nil.")
	} else {
		if ret != opts {
			t.Error("Expected to get the opts object at return. Got a different one.")
		}
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}
	if _, found := opts.opts["required"] ; found {
		t.Errorf("Expected to NOT get '%s' key. found", "required")
	}
}

func TestForjOpts_Short(t *testing.T) {
	t.Log("Expect Short() to add envar in Options object.")
	// --- Setting test context ---
	opts := Opts()
	if opts == nil {
		t.Error("Expected to get an allocated opts. But got Nil.")
		return
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}

	// --- Run the test ---
	ret := opts.Short('U')

	// --- Start testing ---
	if ret == nil {
		t.Error("Expected to get the opts object at return. Got nil.")
	} else {
		if ret != opts {
			t.Error("Expected to get the opts object at return. Got a different one.")
		}
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}
	if v, found := opts.opts["short"] ; !found {
		t.Errorf("Expected to get '%s' key. Not found", "short")
	} else {
		if r, ok := v.(byte) ; !ok {
			t.Errorf("Expected to get '%s' value to be a byte type. Is not.", "short")
		} else {
			if r != 'U' {
				t.Errorf("Expected to get '%s' value '%#v'. Got '%#v'", "short", 'U', v)
			}
		}
	}
}

func TestForjOpts_NoShort(t *testing.T) {
	t.Log("Expect NoShort() to remove short value from Options object.")
	// --- Setting test context ---
	opts := Opts()
	if opts == nil {
		t.Error("Expected to get an allocated opts. But got Nil.")
		return
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}

	opts.Default("test")
	// --- Run the test ---
	ret := opts.NoDefault()

	// --- Start testing ---
	if ret == nil {
		t.Error("Expected to get the opts object at return. Got nil.")
	} else {
		if ret != opts {
			t.Error("Expected to get the opts object at return. Got a different one.")
		}
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}
	if _, found := opts.opts["short"] ; found {
		t.Errorf("Expected to NOT get '%s' key. found", "short")
	}
}

func TestForjOpts_MergeWith(t *testing.T) {
	t.Log("Expect MergeWith() to merge list of options in Options object.")
	// --- Setting test context ---
	opts := Opts()
	if opts == nil {
		t.Error("Expected to get an allocated opts. But got Nil.")
		return
	}
	if opts.opts == nil {
		t.Error("Expected to get a valid opts. But got opts Nil.")
		return
	}
	opts2 := Opts().Required()

	// --- Run the test ---
	opts.MergeWith(opts2)

	// --- Start testing ---
	if n := len(opts2.opts) ; n != 1 {
		t.Errorf("Expected to get 1 options. Got %d options", n)
	}
	if v, found := opts.opts["required"] ; !found {
		t.Errorf("Expected to get '%s' key. Not found", "required")
	} else {
		if b, ok := v.(bool) ; !ok {
			t.Errorf("Expected to get '%s' value '%s'. Got '%t'", "required", "TEST", v)
		} else {
			if ! b {
				t.Errorf("Expected to get '%s' value '%s'. Got '%t'", "required", "TEST", b)
			}
		}
	}
}
