package kingpinCli

import (
	"github.com/alecthomas/kingpin"
	"github.com/forj-oss/forjj-modules/cli/interface"
)

type ParseContext struct {
	context *kingpin.ParseContext
}

type ParseContexter interface {
	GetContext() *ParseContext
}

func (p *ParseContext) GetArgValue(a clier.ArgClauser) (string, bool) {
	arg := a.(KArgClause).GetArg()
	for _, element := range p.context.Elements {
		if a, ok := element.Clause.(*kingpin.ArgClause); ok && a == arg {
			return *element.Value, true
		}
	}
	return "", false
}

func (p *ParseContext) GetFlagValue(f clier.FlagClauser) (string, bool) {
	flag := f.(KFlagClause).GetFlag()
	for _, element := range p.context.Elements {
		if f, ok := element.Clause.(*kingpin.FlagClause); ok && f == flag {
			return *element.Value, true
		}
	}
	return "", false
}

func (p *ParseContext) GetValue(i interface{}) (string, bool) {
	switch i.(type) {
	case clier.ArgClauser:
		return p.GetArgValue(i.(clier.ArgClauser))
	case clier.FlagClauser:
		return p.GetFlagValue(i.(clier.FlagClauser))
	}
	return "", false
}

func (p *ParseContext) SelectedCommands() (res []clier.CmdClauser) {
	cmds, _ := p.context.SelectedCmds()
	res = make([]clier.CmdClauser, 0, len(cmds))
	for _, cmd := range cmds {
		res = append(res, &CmdClause{cmd})
	}
	return
}

func (p *ParseContext) GetParam(param_name string) (ret interface{}, err string) {

	return
}

// Used by local unit test.

type ParseContextTester interface {
	GetContext() *ParseContext
}
