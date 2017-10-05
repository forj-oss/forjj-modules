package kingpinCli

import (
	"github.com/alecthomas/kingpin"
	"forjj-modules/cli/interface"
)

type Application struct {
	app *kingpin.Application
	name string
}

func New(app *kingpin.Application, name string) *Application {
	return &Application{app: app, name: name}
}

func (a *Application) IsNil() bool {
	if a == nil {
		return true
	}
	return false
}

func (a *Application) Arg(p1, p2 string) clier.ArgClauser {
	return &ArgClause{arg: a.app.Arg(p1, p2)}
}

func (a *Application) Flag(p1, p2 string) clier.FlagClauser {
	return &FlagClause{flag: a.app.Flag(p1, p2)}
}

func (a *Application) Command(p1, p2 string) clier.CmdClauser {
	return &CmdClause{a.app.Command(p1, p2)}
}

func (a *Application) ParseContext(args []string) (p clier.ParseContexter, err error) {
	context := new(ParseContext)
	context.context, err = a.app.ParseContext(args)
	p = context
	return
}

func (a *Application) Parse(args []string) (cmd string, err error) {
	return a.app.Parse(args)
}

func (a *Application) Name() string {
	if a == nil {
		return ""
	}
	return a.name
}
