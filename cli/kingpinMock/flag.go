package kingpinMock

import (
	"fmt"
	"github.com/kr/text"
	"github.com/forj-oss/forjj-modules/cli/interface"
	"github.com/forj-oss/forjj-modules/trace"
	"reflect"
)

type FlagClause struct {
	vtype     int // Value type requested.
	name      string
	help      string
	vdefault  []string
	hidden    bool
	required  bool
	short     rune
	envar     string
	context   string // Context value
	value     interface{}
	set_value ClauseList
}

func (a *FlagClause) Stringer() string {
	ret := fmt.Sprintf("Flag (%p):\n", a)
	ret += fmt.Sprintf("  name: '%s'\n", a.name)
	ret += fmt.Sprintf("  help: '%s'\n", a.help)
	ret += fmt.Sprintf("  vtype: '%d'\n", a.vtype)
	ret += fmt.Sprintf("  required: '%t'\n", a.required)
	ret += fmt.Sprintf("  vdefault: '%s'\n", a.vdefault)
	ret += fmt.Sprintf("  envar: '%s'\n", a.envar)
	ret += fmt.Sprintf("  hidden: '%t'\n", a.hidden)
	ret += fmt.Sprintf("  short: '%b'\n", a.short)
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

func NewFlag(name, help string) (f *FlagClause) {
	f = new(FlagClause)
	f.name = name
	f.help = help
	gotrace.Trace("Flag created : (%p)%#v", f, f)
	return f
}

func (f *FlagClause) GetHelp() string {
	return f.help
}

func (f *FlagClause) GetName() string {
	return f.name
}

func (f *FlagClause) String() (ret *string) {
	if f.vtype != NilType && f.vtype != StringType {
		return nil
	}

	if f.vtype == NilType {
		f.vtype = StringType
		ret = new(string)
		f.value = ret
	} else {
		ret = f.value.(*string)
	}
	return
}

func (f *FlagClause) Bool() (ret *bool) {
	if f.vtype != NilType && f.vtype != BoolType {
		return nil
	}
	if f.vtype == NilType {
		f.vtype = BoolType
		ret = new(bool)
		f.value = ret
	} else {
		ret = f.value.(*bool)
	}
	return
}

func (a *FlagClause) GetType() string {
	switch {
	case a.vtype == StringType:
		return "string"
	case a.vtype == BoolType:
		return "bool"
	}
	return "any"
}

func (f *FlagClause) Required() clier.FlagClauser {
	f.required = true
	return f
}

func (f *FlagClause) IsRequired() bool {
	return f.required
}

func (f *FlagClause) Short(p1 rune) clier.FlagClauser {
	f.short = p1
	return f
}

func (f *FlagClause) IsShort(p1 rune) bool {
	return (f.short == p1)
}

func (f *FlagClause) Hidden() clier.FlagClauser {
	f.hidden = true
	return f
}

func (f *FlagClause) IsHidden() bool {
	return f.hidden
}

func (f *FlagClause) Default(p1 ...string) clier.FlagClauser {
	f.vdefault = p1
	return f
}

func (f *FlagClause) IsDefault(p1 ...string) bool {
	return reflect.DeepEqual(f.vdefault, p1)
}

func (f *FlagClause) Envar(p1 string) clier.FlagClauser {
	f.envar = p1
	return f
}

func (f *FlagClause) IsEnvar(p1 string) bool {
	return (f.envar == p1)
}

func (f *FlagClause) SetValue(v clier.Valuer) clier.FlagClauser {
	f.set_value = v.(ClauseList)
	return f
}

func (f *FlagClause) IsSetValue(_ clier.Valuer) (ret bool) {
	if f.set_value == nil {
		return
	}
	ret = true
	return
}

// Context interface

func (f *FlagClause) SetContextValue(s string) (*FlagClause, error) {
	f.context = s
	return f, nil
}

func (f *FlagClause) GetContextValue() string {
	return f.context
}

func (f *FlagClause) update_data() {
	switch f.value.(type) {
	case *string:
		s := f.value.(*string)
		*s = f.context
	case *bool:
		b := f.value.(*bool)
		if f.context == "true" {
			*b = true
		}
	}
}
