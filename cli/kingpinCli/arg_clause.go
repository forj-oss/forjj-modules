package kingpinCli

import (
	"github.com/alecthomas/kingpin"
	"github.com/forj-oss/forjj-modules/cli/interface"
)

type ArgClause struct {
	arg *kingpin.ArgClause
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
	return a
}

func (a *ArgClause) Envar(p1 string) clier.ArgClauser {
	a.arg.Envar(p1)
	return a
}

func (a *ArgClause) SetValue(p1 clier.Valuer) clier.ArgClauser {
	a.arg.SetValue(p1)
	return a
}
