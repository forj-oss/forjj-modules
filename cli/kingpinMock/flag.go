package kingpinMock

import (
	"github.com/forj-oss/forjj-modules/cli/interface"
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
}

func NewFlag(name, help string) (f *FlagClause) {
	f = new(FlagClause)
	f.name = name
	f.help = help
	return f
}

func (f *FlagClause) GetHelp() string {
	return f.help
}

func (f *FlagClause) GetName() string {
	return f.name
}

func (f *FlagClause) String() *string {
	f.vtype = StringType
	return new(string)
}

func (f *FlagClause) Bool() *bool {
	f.vtype = BoolType
	return new(bool)
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
