package kingpinMock

import (
	"fmt"
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
	set_value string
	short     rune
	envar     string
	context   string // Context value
	value     interface{}
}

func (a *FlagClause) Stringer() string {
	ret := fmt.Sprintf("Flag (%p):\n", a)
	ret += fmt.Sprintf("  name: '%s'\n", a.name)
	ret += fmt.Sprintf("  help: '%s'\n", a.help)
	ret += fmt.Sprintf("  vtype: '%d'\n", a.vtype)
	ret += fmt.Sprintf("  required: '%t'\n", a.required)
	ret += fmt.Sprintf("  vdefault: '%s'\n", a.vdefault)
	ret += fmt.Sprintf("  envar: '%s'\n", a.envar)
	ret += fmt.Sprintf("  set_value: '%s'\n", a.set_value)
	ret += fmt.Sprintf("  hidden: '%t'\n", a.hidden)
	ret += fmt.Sprintf("  short: '%b'\n", a.short)
	ret += fmt.Sprintf("  value: '%s'", a.value)
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

func (f *FlagClause) SetValue(_ clier.Valuer) clier.FlagClauser {
	return f
}

func (f *FlagClause) IsSetValue(_ clier.Valuer) bool {
	return false
}

// Context interface

func (f *FlagClause) SetContextValue(s string) *FlagClause {
	f.context = s
	return f
}

func (f *FlagClause) GetContextValue() string {
	return f.context
}
