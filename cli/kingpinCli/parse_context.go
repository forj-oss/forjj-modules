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

// GetArgValue get value from cli, or if missing, ENV or if missing, defaults
func (p *ParseContext) GetArgValue(a clier.ArgClauser) (interface{}, bool) {
	karg := a.(KArgClause).GetArg()
	argClause := a.(*ArgClause)
	for _, element := range p.context.Elements {
		if a, ok := element.Clause.(*kingpin.ArgClause); ok && a == karg {
			return *element.Value, true
		}
	}
	if karg.HasEnvarValue() {
		return karg.GetEnvarValue(), true
	}
	if argClause.hasDefaults() {
		return argClause.getDefaults(), true
	}
	return nil, false
}

// GetFlagValue get value from cli, or if missing, ENV or if missing, defaults
func (p *ParseContext) GetFlagValue(f clier.FlagClauser) (interface{}, bool) {
	kflag := f.(KFlagClause).GetFlag()
	flagClause := f.(*FlagClause)
	for _, element := range p.context.Elements {
		if f, ok := element.Clause.(*kingpin.FlagClause); ok && f == kflag {
			return *element.Value, true
		}
	}
	if kflag.HasEnvarValue() {
		return kflag.GetEnvarValue(), true
	}
	if flagClause.hasDefaults() {
		return flagClause.getDefaults(), true
	}
	return nil, false
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
