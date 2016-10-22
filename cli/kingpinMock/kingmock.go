package kingpinMock

import (
	"github.com/forj-oss/forjj-modules/cli/interface"
)

const (
	NilType = 0
	StringType
	BoolType
)

// **************************************

type Application struct {
}

func NewMock(_ []string, _ string) *Application {
	a := new(Application)
	return a
}

func (a *Application) IsNil() bool {
	if a == nil {
		return true
	}
	return false
}

func (a *Application) Arg(_, _ string) clier.ArgClauser {
	ac := new(ArgClause)
	return ac
}

func (a *Application) Flag(_, _ string) clier.FlagClauser {
	f := new(FlagClause)
	return f
}

func (a *Application) Command(_, _ string) clier.CmdClauser {
	c := new(CmdClause)
	return c
}

// **************************************

type ParseContext struct {
}

func (*ParseContext) loadfrom(_ *ParseContext) {

}

// **************************************

type CmdClause struct {
}

func (c *CmdClause) Command(_, _ string) clier.CmdClauser {
	return c
}

func (c *CmdClause) Flag(_, _ string) clier.FlagClauser {
	return new(FlagClause)
}

func (c *CmdClause) Arg(_, _ string) clier.ArgClauser {
	return new(ArgClause)
}

// **************************************

type ArgClause struct {
	vtype int // Value type requested.
}

func (a *ArgClause) String() *string {
	a.vtype = StringType
	return new(string)
}

func (f *ArgClause) Bool() *bool {
	f.vtype = BoolType
	return new(bool)
}

func (f *ArgClause) Required() clier.ArgClauser {
	return f
}

func (f *ArgClause) Default(_ ...string) clier.ArgClauser {
	return f
}

func (f *ArgClause) Envar(_ string) clier.ArgClauser {
	return f
}

func (f *ArgClause) SetValue(_ clier.Valuer) clier.ArgClauser {
	return f
}

// **************************************

type FlagClause struct {
	vtype int // Value type requested.
}

func (f *FlagClause) String() *string {
	f.vtype = StringType
	return new(string)
}

func (f *FlagClause) Bool() *bool {
	f.vtype = BoolType
	return new(bool)
}

func (f *FlagClause) Required() clier.FlagClauser {
	return f
}

func (f *FlagClause) Short(_ byte) clier.FlagClauser {
	return f
}

func (f *FlagClause) Hidden() clier.FlagClauser {
	return f
}

func (f *FlagClause) Default(_ ...string) clier.FlagClauser {
	return f
}

func (f *FlagClause) Envar(_ string) clier.FlagClauser {
	return f
}

func (f *FlagClause) SetValue(_ clier.Valuer) clier.FlagClauser {
	return f
}
