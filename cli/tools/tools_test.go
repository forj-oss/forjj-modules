package tools

import "testing"

func TestToBoolWithAddr(t *testing.T) {
	t.Log("Expect ToBoolWithAddr to return appropriate values and type")

	// --------------- Set Context

	var val1 interface{}
	b := true
	val1 = b

	// --------------- running test
	res, err := ToBoolWithAddr(val1)
	// --------------- Testing
	if err != nil {
		t.Errorf("Expected ToBoolWithAddr(&%t) to return a valid value. But got an error '%s'", b, err)
	} else {
		if v, ok := res.(bool) ; !ok {
			t.Errorf("Expected ToBoolWithAddr(&%t) to get a bool type. Is not.", b)
		} else {
			if !v {
				t.Errorf("Expected to return 'true'. Got %t", v)
			}
		}
	}

	// --------------- Set Context

	val1 = &b

	// --------------- running test
	res, err = ToBoolWithAddr(val1)
	// --------------- Testing
	if err != nil {
		t.Errorf("Expected ToBoolWithAddr(&%t) to return a valid value. But got an error '%s'", b, err)
	} else {
		if v, ok := res.(*bool) ; !ok {
			t.Errorf("Expected ToBoolWithAddr(&%t) to get a *bool type. Is not.", b)
		} else {
			if !*v {
				t.Errorf("Expected ToBoolWithAddr(&%t) to return 'true'. Got %t", b, v)
			}
		}
	}

	s := "false"
	val1 = s

	// --------------- running test
	res, err = ToBoolWithAddr(val1)
	// --------------- Testing
	if err != nil {
		t.Errorf("Expected ToBoolWithAddr('%s') to return a valid value. But got an error '%s'", s, err)
	} else {
		if v, ok := res.(bool) ; !ok {
			t.Errorf("Expected ToBoolWithAddr('%s') to get a bool type. Is not.", s)
		} else {
			if v {
				t.Errorf("Expected ToBoolWithAddr('%s') to return 'false'. Got %t", s, v)
			}
		}
	}

	// --------------- Set Context

	val1 = &s

	// --------------- running test
	res, err = ToBoolWithAddr(val1)
	// --------------- Testing
	if err != nil {
		t.Errorf("Expected ToBoolWithAddr('%s') to return a valid value. But got an error '%s'", s, err)
	} else {
		if v, ok := res.(*bool); !ok {
			t.Errorf("Expected ToBoolWithAddr('%s') to get a *bool type. Is not.", s)
		} else {
			if *v {
				t.Errorf("Expected ToBoolWithAddr('%s') to return 'true'. Got %t", s, v)
			}
		}
	}
	// --------------- Set Context

	s = "true"
	val1 = &s

	// --------------- running test
	res, err = ToBoolWithAddr(val1)
	// --------------- Testing
	if err != nil {
		t.Errorf("Expected ToBoolWithAddr(&'%s') to return a valid value. But got an error '%s'", s, err)
	} else {
		if v, ok := res.(*bool); !ok {
			t.Errorf("Expected ToBoolWithAddr(&'%s') to get a *bool type. Is not.", s)
		} else {
			if !*v {
				t.Errorf("Expected ToBoolWithAddr(&'%s') to return 'false'. Got %t", s, v)
			}
		}
	}
}

func TestToBool(t *testing.T) {
	t.Log("Expect ToBool to return appropriate value")
	// --------------- Set Context

	var val1 interface{}
	b := true
	val1 = b

	// --------------- running test
	res, err := ToBool(val1)
	// --------------- Testing
	if err != nil {
		t.Errorf("Expected ToBool(%s) to return a valid value. But got an error '%s'", "true", err)
	} else {
		if !res {
			t.Errorf("Expected ToBool(%s) to return 'true'. Got %t", "true", res)
		}
	}

	// --------------- Set Context

	val1 = &b

	// --------------- running test
	res, err = ToBool(val1)
	// --------------- Testing
	if err != nil {
		t.Errorf("Expected ToBool(&%s) to return a valid value. But got an error '%s'", "true", err)
	} else {
		if !res {
			t.Errorf("Expected ToBool(&%s) to return 'true'. Got %t", "true", res)
		}
	}

	s := "false"
	val1 = s

	// --------------- running test
	res, err = ToBool(val1)
	// --------------- Testing
	if err != nil {
		t.Errorf("Expected ToBool('%s') to return a valid value. But got an error '%s'", s, err)
	} else {
		if res {
			t.Errorf("Expected ToBool('%s') to return 'false'. Got %t", s, res)
		}
	}

	// --------------- Set Context

	val1 = &s

	// --------------- running test
	res, err = ToBool(val1)
	// --------------- Testing
	if err != nil {
		t.Errorf("Expected ToBool(&'%s') to return a valid value. But got an error '%s'", s, err)
	} else {
		if res {
			t.Errorf("Expected ToBool(&'%s') to return 'true'. Got %t", s, res)
		}
	}
	// --------------- Set Context

	s = "true"
	val1 = &s

	// --------------- running test
	res, err = ToBool(val1)
	// --------------- Testing
	if err != nil {
		t.Errorf("Expected ToBool(&'%s') to return a valid value. But got an error '%s'", s, err)
	} else {
		if !res {
			t.Errorf("Expected ToBool(&'%s') to return 'false'. Got %t", s, res)
		}
	}
}
