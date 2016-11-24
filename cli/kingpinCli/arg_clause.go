package kingpinCli

import (
	"github.com/alecthomas/kingpin"
	"github.com/forj-oss/forjj-modules/cli/interface"
)

type KArgClause interface {
	GetArg() *kingpin.ArgClause
}

type ArgClause struct {
	arg            *kingpin.ArgClause
	default_values []string
}

func (a *ArgClause) String() *string {
	return a.arg.String()
}

func (a *ArgClause) Bool() *bool {
	return a.arg.Bool()
}

func (a *ArgClause) Required() clier.ArgClauser {
	a.arg.Required()
	return a
}

func (a *ArgClause) Default(p1 ...string) clier.ArgClauser {
	a.arg.Default(p1...)
	a.default_values = p1
	return a
}

func (f *ArgClause) getDefaults() []string {
	if f.default_values == nil {
		return []string{}
	}
	return f.default_values
}

func (f *ArgClause) hasDefaults() bool {
	if f.default_values == nil {
		return false
	}
	return true
}

func (a *ArgClause) Envar(p1 string) clier.ArgClauser {
	a.arg.Envar(p1)
	return a
}

func (a *ArgClause) SetValue(p1 clier.Valuer) clier.ArgClauser {
	a.arg.SetValue(p1)
	return a
}

func (a *ArgClause) GetArg() *kingpin.ArgClause {
	return a.arg
}
