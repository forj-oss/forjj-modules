package kingpinMock

import (
	"fmt"
	"github.com/kr/text"
	"forjj-modules/cli/interface"
	"forjj-modules/trace"
)

type ArgClause struct {
	vtype     int // Value type requested.
	name      string
	help      string
	required  bool
	vdefault  *string
	envar     string
	set_value ClauseList
	context   string // Context value
	value     interface{}
}

func (a *ArgClause) Stringer() string {
	ret := fmt.Sprintf("Arg (%p):\n", a)
	ret += fmt.Sprintf("  name: '%s'\n", a.name)
	ret += fmt.Sprintf("  help: '%s'\n", a.help)
	ret += fmt.Sprintf("  vtype: '%d'\n", a.vtype)
	ret += fmt.Sprintf("  required: '%s'\n", a.required)
	ret += fmt.Sprintf("  vdefault: '%s'\n", a.vdefault)
	ret += fmt.Sprintf("  envar: '%s'\n", a.envar)
	switch a.value.(type) {
	case *string:
		ret += fmt.Sprintf("  value: '%s' (string - %p)\n", *a.value.(*string), a.value)
	case *bool:
		ret += fmt.Sprintf("  value: '%t' (bool - %p)\n", *a.value.(*bool), a.value)
	}
	ret += fmt.Sprintf("  context value: '%s'\n", a.context)
	if a.set_value != nil {
		ret += "  set_value:\n"
		ret += text.Indent(a.set_value.String(), "    ")
	}
	return ret
}

func NewArg(name, help string) (f *ArgClause) {
	f = new(ArgClause)
	f.name = name
	f.help = help
	gotrace.Trace("Arg created : (%p)%#v", f, f)
	return f
}

func (a *ArgClause) GetHelp() string {
	return a.help
}

func (a *ArgClause) GetName() string {
	return a.name
}

func (a *ArgClause) String() (ret *string) {
	if a.vtype != NilType && a.vtype != StringType {
		return nil
	}

	if a.vtype == NilType {
		a.vtype = StringType
		ret = new(string)
		a.value = ret
	} else {
		ret = a.value.(*string)
	}
	return
}

func (a *ArgClause) GetType() string {
	switch {
	case a.vtype == BoolType:
		return "bool"
	case a.vtype == StringType:
		return "string"
	}
	return "any"
}

func (a *ArgClause) Bool() (ret *bool) {
	if a.vtype != NilType && a.vtype != BoolType {
		return nil
	}
	if a.vtype == NilType {
		a.vtype = BoolType
		ret = new(bool)
		a.value = ret
	} else {
		ret = a.value.(*bool)
	}
	return
}

func (f *ArgClause) IsBool() bool {
	return (f.vtype == BoolType)
}

func (f *ArgClause) Required() clier.ArgClauser {
	f.required = true
	return f
}

func (f *ArgClause) IsRequired() bool {
	return (f.required == true)
}

func (f *ArgClause) Default(p1 string) clier.ArgClauser {
	if f.vdefault == nil {
		f.vdefault = new(string)
	}
	*f.vdefault = p1
	return f
}

func (f *ArgClause) getDefaults() *string {
	return f.vdefault
}

func (f *ArgClause) hasDefaults() bool {
	return (f.vdefault != nil)
}

func (f *ArgClause) IsDefault(p1 string) bool {
	return (p1 == *f.vdefault)
}

func (f *ArgClause) Envar(p1 string) clier.ArgClauser {
	f.envar = p1
	return f
}

func (f *ArgClause) IsEnvar(p1 string) bool {
	return (f.envar == p1)
}

func (f *ArgClause) SetValue(v clier.Valuer) clier.ArgClauser {
	f.set_value = v.(ClauseList)
	return f
}

func (f *ArgClause) IsSetValue(_ clier.Valuer) bool {
	if f.set_value == nil {
		return false
	}
	return true
}

// Context interface

func (a *ArgClause) SetContextValue(s string) (*ArgClause, error) {
	a.context = s
	if a.set_value != nil {
		if err := a.set_value.Set(s); err != nil {
			return nil, err
		}
	}

	return a, nil
}

func (a *ArgClause) GetContextValue() string {
	return a.context
}

func (a *ArgClause) update_data() {
	switch a.value.(type) {
	case *string:
		s := a.value.(*string)
		*s = a.context
	case *bool:
		b := a.value.(*bool)
		if a.context == "true" {
			*b = true
		}
	}
}
