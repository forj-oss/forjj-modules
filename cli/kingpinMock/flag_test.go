package kingpinMock

import (
	"reflect"
	"testing"
)

func TestNewFlag(t *testing.T) {
	t.Log("NewArg creates a new Flag with name and help")
	a := NewFlag("test", "help")
	if a.name != "test" {
		t.Errorf("name expected to be 'test'. Got %s", a.name)
	}
	if a.help != "help" {
		t.Errorf("help expected to be 'help'. Got %s", a.help)
	}
}

func TestFlagClause_Bool(t *testing.T) {
	t.Log("Setting bool type")
	a := NewFlag("test", "help")

	if a.GetType() != "any" {
		t.Errorf("Expected arg to be initialized as any. Got %s", a.GetType())
	}

	b := a.Bool()

	bt := reflect.TypeOf(b).String()
	if bt != "*bool" {
		t.Errorf("Expected returned value type to be *bool. Got: %s", bt)
	}

	if a.GetType() != "bool" {
		t.Errorf("Expected arg to be set as bool. Got %s", a.GetType())
	}
}

func TestFlagClause_String(t *testing.T) {
	t.Log("Setting string type")
	a := NewFlag("test", "help")

	if a.GetType() != "any" {
		t.Errorf("Expected arg to be initialized as any. Got %s", a.GetType())
	}

	b := a.String()

	bt := reflect.TypeOf(b).String()
	if bt != "*string" {
		t.Errorf("Expected returned value type to be *string. Got: %s", bt)
	}

	if a.GetType() != "string" {
		t.Errorf("Expected arg to be set as string. Got : %s", a.GetType())
	}

}

func TestFlagClause_Default(t *testing.T) {
	value := "default"
	function := "Default"
	t.Logf("Running %s(\"%s\")", function, value)
	a := NewFlag("test", "help")

	if len(a.vdefault) > 0 {
		t.Errorf("Expected %s() to not be set. Got '%s'", function, a.vdefault)
	}

	b := a.Default(value)

	if a != b {
		t.Fail()
	}

	if a.vdefault[0] != value {
		t.Errorf("Expected %s() to be set to '%s'. Got '%s'", function, value, a.vdefault[0])
	}
}

func TestFlagClause_Envar(t *testing.T) {
	value := "ARG"
	function := "Envar"
	t.Logf("Running %s(\"%s\")", function, value)
	a := NewFlag("test", "help")
	b := a.Envar("ARG")

	if a != b {
		t.Fail()
	}

	if a.envar != value {
		t.Errorf("Expected %s() to be set to '%s'. Got '%s'", function, value, a.vdefault[0])
	}

}

func TestFlagClause_Required(t *testing.T) {
	value := true
	function := "Required"
	t.Logf("Running %s(\"%s\")", function, value)

	a := NewFlag("test", "help")
	b := a.Required()

	if a != b {
		t.Fail()
	}

	if a.required != value {
		t.Errorf("Expected %s() to be true. Got '%s'", function, a.required)
	}

}

func TestFlagClause_Hidden(t *testing.T) {
	value := true
	function := "Hidden"
	t.Logf("Running %s(\"%s\")", function, value)

	a := NewFlag("test", "help")
	b := a.Hidden()

	if a != b {
		t.Fail()
	}

	if a.hidden != value {
		t.Errorf("Expected %s() to be true. Got '%s'", function, a.hidden)
	}

}

func TestFlagClause_Short(t *testing.T) {
	value := 'S'
	function := "Hidden"
	t.Logf("Running %s(\"%s\")", function, value)

	a := NewFlag("test", "help")
	b := a.Short(value)

	if a != b {
		t.Fail()
	}

	if a.short != value {
		t.Errorf("Expected %s() to be '%s'. Got '%s'", function, value, a.hidden)
	}

}
