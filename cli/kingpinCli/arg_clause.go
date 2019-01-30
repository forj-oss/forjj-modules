package kingpinCli

import (
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/forj-oss/forjj-modules/cli"
)

type KArgClause interface {
	GetArg() *kingpin.ArgClause
}

type ArgClause struct {
	arg           *kingpin.ArgClause
	default_value *string
}

func (a *ArgClause) Stringer() string {
	ret := fmt.Sprintf("ArgClause (%p):\n", a)
	ret += fmt.Sprintf("  name: '%s'\n", a.arg.Model().Name)
	if a.default_value == nil {
		ret += fmt.Sprint("  vdefault: nil\n")
	} else {
		ret += fmt.Sprintf("  vdefault: '%s' (%p)\n", *a.default_value, a.default_value)
	}
	return ret
}

func (a *ArgClause) String() *string {
	return a.arg.String()
}

func (a *ArgClause) Bool() *bool {
	return a.arg.Bool()
}

func (a *ArgClause) Required() cli.ArgClauser {
	a.arg.Required()
	return a
}

func (a *ArgClause) Default(p1 string) cli.ArgClauser {
	if a.default_value == nil {
		a.default_value = new(string)
	}
	a.arg.Default(p1)
	*a.default_value = p1
	return a
}

func (f *ArgClause) getDefaults() *string {
	return f.default_value
}

func (f *ArgClause) hasDefaults() bool {
	return (f.default_value != nil)
}

func (a *ArgClause) Envar(p1 string) cli.ArgClauser {
	a.arg.Envar(p1)
	return a
}

func (a *ArgClause) SetValue(p1 cli.Valuer) cli.ArgClauser {
	a.arg.SetValue(p1)
	return a
}

func (a *ArgClause) GetArg() *kingpin.ArgClause {
	return a.arg
}

// NewArg creates a generic ArgClause from kingpin.ArgClause
func NewArg(arg *kingpin.ArgClause) (f *ArgClause) {
	f = new(ArgClause)
	f.arg = arg
	return
}
