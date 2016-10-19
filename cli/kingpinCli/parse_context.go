package kingpinCli

import "github.com/alecthomas/kingpin"

type ParseContext struct {
	context *kingpin.ParseContext
}

func GetContext(app *kingpin.Application, args []string) (p *ParseContext, err error) {
	p.context, err = app.ParseContext(args)
	return
}

func (p *ParseContext) GetArgValue(a *kingpin.ArgClause) (string, bool) {
	for _, element := range p.context.Elements {
		if arg, ok := element.Clause.(*kingpin.ArgClause); ok && arg == a {
			return *element.Value, true
		}
	}
	return "", false
}

func (p *ParseContext) GetFlagValue(f *kingpin.FlagClause) (string, bool) {
	for _, element := range p.context.Elements {
		if flag, ok := element.Clause.(*kingpin.FlagClause); ok && flag == f {
			return *element.Value, true
		}
	}
	return "", false
}

func (p *ParseContext) SelectedCommand() *kingpin.CmdClause {
	return p.context.SelectedCommand
}
