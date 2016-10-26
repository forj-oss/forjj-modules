package kingpinMock

import (
	"github.com/forj-oss/forjj-modules/cli/interface"
	"github.com/forj-oss/forjj-modules/trace"
)

type ParseContext struct {
	cmds []*CmdClause
	app  *Application
}

// Following functions are implemented by clier.ParseContexter
func (*ParseContext) GetFlagValue(_ clier.FlagClauser) (string, bool) {
	return "", false
}

func (*ParseContext) GetArgValue(_ clier.ArgClauser) (string, bool) {
	return "", false
}

func (p *ParseContext) SelectedCommands() (res []clier.CmdClauser) {
	res = make([]clier.CmdClauser, 0, len(p.cmds))
	for _, cmd := range p.cmds {
		res = append(res, cmd)
	}
	return
}

// Following functions are specific to the Mock
func (a *Application) NewContext() *ParseContext {
	a.context = new(ParseContext)
	a.context.app = a
	return a.context
}

func (p *ParseContext) SetContext(p1 ...string) *ParseContext {
	p.cmds = make([]*CmdClause, 0, len(p1))
	if len(p1) == 0 {
		return p
	}

	var cmd *CmdClause
	if v, found := p.app.cmds[p1[0]]; !found {
		gotrace.Trace("Unable to find %s Command from Application layer.", p1[0])
		return nil
	} else {
		cmd = v
	}
	p.cmds = append(p.cmds, cmd)

	if len(p1) == 1 {
		return p
	}

	for _, cmd_name := range p1[1:] {
		if v, found := cmd.cmds[cmd_name]; !found {
			gotrace.Trace("Unable to find %s Command from Application layer.", cmd)
			return nil
		} else {
			p.cmds = append(p.cmds, v)
		}
	}
	return p
}

func (p *ParseContext) SetContextValue(name string, value string) *ParseContext {
	if p == nil {
		return nil
	}
	if len(p.cmds) == 0 {
		return nil
	}

	cmd := p.cmds[len(p.cmds)-1]
	if v, found := cmd.flags[name]; found {
		switch v.value.(type) {
		case *string:
			*v.value.(*string) = value
		case *bool:
			if value == "true" {
				*v.value.(*bool) = true
			} else {
				*v.value.(*bool) = false
			}
		}
	}
	return p
}

func (p *ParseContext) SetContextAppValue(name string, value string) *ParseContext {
	if p == nil {
		return nil
	}

	if v, found := p.app.flags[name]; found {
		switch v.value.(type) {
		case *string:
			*v.value.(*string) = value
		case *bool:
			if value == "true" {
				*v.value.(*bool) = true
			} else {
				*v.value.(*bool) = false
			}
		}
	}
	return p
}
