package kingpinCli

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/forj-oss/forjj-modules/cli"
)

type KFlagClause interface {
	GetFlag() *kingpin.FlagClause
}

type FlagClause struct {
	flag          *kingpin.FlagClause
	default_value *string
}

func (a *FlagClause) Stringer() string {
	ret := fmt.Sprintf("FlagClause (%p):\n", a)
	ret += fmt.Sprintf("  name: '%s'\n", a.flag.Model().Name)
	if a.default_value == nil {
		ret += fmt.Sprint("  vdefault: nil\n")
	} else {
		ret += fmt.Sprintf("  vdefault: '%s' (%p)\n", *a.default_value, a.default_value)
	}
	return ret
}

func (f *FlagClause) String() *string {
	return f.flag.String()
}

func (f *FlagClause) Bool() *bool {
	return f.flag.Bool()
}

func (f *FlagClause) Required() cli.FlagClauser {
	f.flag.Required()
	return f
}

func (f *FlagClause) Short(p1 rune) cli.FlagClauser {
	f.flag.Short(p1)
	return f
}

func (f *FlagClause) Hidden() cli.FlagClauser {
	f.flag.Hidden()
	return f
}

func (f *FlagClause) Default(p1 string) cli.FlagClauser {
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

func (f *FlagClause) Envar(p1 string) cli.FlagClauser {
	f.flag.Envar(p1)
	return f
}

func (f *FlagClause) SetValue(p1 cli.Valuer) cli.FlagClauser {
	f.flag.SetValue(p1)
	return f
}

func (f *FlagClause) GetFlag() *kingpin.FlagClause {
	return f.flag
}

// NewFlag creates a generic FlagClause from kingpin.FlagClause
func NewFlag(flag *kingpin.FlagClause)(f* FlagClause) {
	f = new(FlagClause)
	f.flag = flag
	return
}