package kingpinMock

import (
	"fmt"
	"github.com/kr/text"
	"github.com/forj-oss/forjj-modules/cli/interface"
	"github.com/forj-oss/forjj-modules/trace"
	"strings"
)

type ParseContext struct {
	cmds     []*CmdClause
	app      *Application
	Elements []interface{}
}

func (p *ParseContext) String() (ret string) {
	ret = "Cmds:\n"
	list := make([]string, 0, len(p.cmds))
	for _, cmd := range p.cmds {
		list = append(list, fmt.Sprintf("%s(%p)", cmd.command, cmd))
	}
	ret += text.Indent(strings.Join(list, " ")+"\n", "  ")
	ret += "Elements:\n"
	for _, element := range p.Elements {
		switch element.(type) {
		case *FlagClause:
			f := element.(*FlagClause)
			ret += text.Indent(fmt.Sprintf("FlagClause %s (%p)\n", f.name, f), "  ")
		case *ArgClause:
			a := element.(*ArgClause)
			ret += text.Indent(fmt.Sprintf("ArgClause %s (%p)\n", a.name, a), "  ")
		case *CmdClause:
			c := element.(*CmdClause)
			ret += text.Indent(fmt.Sprintf("CmdClause %s (%p)\n", c.command, c), "  ")
		default:
			ret += fmt.Sprintf("Unknown type:\n%s", element)
		}
	}
	return
}

type ParseContextTester interface {
	GetContext() *ParseContext
}

// Following functions are implemented by clier.ParseContexter

// IMPORTANT! The Mock GetFlagValue && GetArgValue are not configured to get from process environment

func (p *ParseContext) GetFlagValue(f clier.FlagClauser) (string, bool) {
	var flag *FlagClause

	if v, ok := f.(*FlagClause); !ok {
		return "", false
	} else {
		flag = v
	}

	for _, element := range p.app.context.Elements {
		if f, ok := element.(*FlagClause); ok && f == flag {
			return f.context, true
		}
	}
	if flag.hasDefaults() {
		return flag.getDefaults()[0], true
	}
	return "", false
}

func (p *ParseContext) GetArgValue(a clier.ArgClauser) (string, bool) {
	var arg *ArgClause

	if v, ok := a.(*ArgClause); !ok {
		return "", false
	} else {
		arg = v
	}

	for _, element := range p.app.context.Elements {
		if v, ok := element.(*ArgClause); ok && a == arg {
			return v.context, true
		}
	}
	if arg.hasDefaults() {
		return arg.getDefaults()[0], true
	}
	return "", false
}

func (p *ParseContext) GetParam(name string) (interface{}, string) {
	for _, element := range p.app.context.Elements {
		switch element.(type) {
		case *ArgClause:
			a := element.(*ArgClause)
			if a.name == name {
				return a, "*ArgClause"
			}
		case *FlagClause:
			f := element.(*FlagClause)
			if f.name == name {
				return f, "*FlagClause"
			}
		}
	}
	return nil, ""
}

func (p *ParseContext) SelectedCommands() (res []clier.CmdClauser) {
	if p == nil {
		return
	}
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
	a.context.cmds = make([]*CmdClause, 0)
	a.context.Elements = make([]interface{}, 0)
	return a.context
}

func (p *ParseContext) SetContext(p1 ...string) *ParseContext {
	p.cmds = make([]*CmdClause, 0, len(p1))
	p.Elements = make([]interface{}, 0, len(p1))
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
	p.Elements = append(p.Elements, cmd)

	if len(p1) == 1 {
		return p
	}

	for _, cmd_name := range p1[1:] {
		if v, found := cmd.cmds[cmd_name]; !found {
			gotrace.Trace("Unable to find %s Command from Application layer.", cmd)
			return nil
		} else {
			p.cmds = append(p.cmds, v)
			p.Elements = append(p.Elements, v)
		}
	}
	return p
}

func (p *ParseContext) SetContextValue(name string, value string) (*ParseContext, error) {
	return p.setValue(true, false, name, value)
}

func (p *ParseContext) SetCliValue(name string, value string) (*ParseContext, error) {
	return p.setValue(false, true, name, value)
}

func (p *ParseContext) SetValue(name string, value string) (*ParseContext, error) {
	return p.setValue(true, true, name, value)
}

func (p *ParseContext) setValue(context, cli bool, name string, value string) (*ParseContext, error) {
	if p == nil {
		return nil, nil
	}

	// App
	if len(p.cmds) == 0 {
		if v, found := p.app.flags[name]; found {
			if context {
				if _, err := v.SetContextValue(value); err != nil {
					return nil, err
				}
			}
			if cli {
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
			p.Elements = append(p.Elements, v)
			return p, nil
		}

		// Args
		if v, found := p.app.args[name]; found {
			if context {
				if _, err := v.SetContextValue(value); err != nil {
					return nil, err
				}
			}
			if cli {
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
			p.Elements = append(p.Elements, v)
			return p, nil
		}

		return nil, nil
	}

	cmd := p.cmds[len(p.cmds)-1]

	// Flags
	if v, found := cmd.flags[name]; found {
		if context {
			if _, err := v.SetContextValue(value); err != nil {
				return nil, err
			}
		}
		if cli {
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
		p.Elements = append(p.Elements, v)
		return p, nil
	}

	// Args
	if v, found := cmd.args[name]; found {
		if context {
			if _, err := v.SetContextValue(value); err != nil {
				return nil, err
			}
		}
		if cli {
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
		p.Elements = append(p.Elements, v)
		return p, nil
	}
	return nil, nil
}

func (p *ParseContext) SetContextAppValue(name string, value string) *ParseContext {
	if v, found := p.app.flags[name]; found {
		v.context = value
		p.Elements = append(p.Elements, v)
	}
	return p
}

func (p *ParseContext) SetParsedAppValue(name string, value string) *ParseContext {
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

func (p *ParseContext) GetContext() *ParseContext {
	return p
}
