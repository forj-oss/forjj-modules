package kingpinCli

import (
	"github.com/alecthomas/kingpin"
	"github.com/forj-oss/forjj-modules/cli/interface"
)

type KFlagClause interface {
	GetFlag() *kingpin.FlagClause
}

type FlagClause struct {
	flag          *kingpin.FlagClause
	default_value *string
}

func (f *FlagClause) String() *string {
	return f.flag.String()
}

func (f *FlagClause) Bool() *bool {
	return f.flag.Bool()
}

func (f *FlagClause) Required() clier.FlagClauser {
	f.flag.Required()
	return f
}

func (f *FlagClause) Short(p1 rune) clier.FlagClauser {
	f.flag.Short(p1)
	return f
}

func (f *FlagClause) Hidden() clier.FlagClauser {
	f.flag.Hidden()
	return f
}

func (f *FlagClause) Default(p1 string) clier.FlagClauser {
	if f.default_value == nil {
		f.default_value = new(string)
	}
	f.flag.Default(p1)
	*f.default_value = p1
	return f
}

func (f *FlagClause) getDefaults() *string {
	return f.default_value
}

func (f *FlagClause) hasDefaults() bool {
	if f.default_value == nil {
		return false
	}
	return true
}

func (f *FlagClause) Envar(p1 string) clier.FlagClauser {
	f.flag.Envar(p1)
	return f
}

func (f *FlagClause) SetValue(p1 clier.Valuer) clier.FlagClauser {
	f.flag.SetValue(p1)
	return f
}

func (f *FlagClause) GetFlag() *kingpin.FlagClause {
	return f.flag
}
