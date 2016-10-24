package kingpinMock

import (
	"github.com/forj-oss/forjj-modules/cli/interface"
)

type Application struct {
	flags map[string]*FlagClause
	args  map[string]*ArgClause
	cmds  map[string]*CmdClause
}

func New(_ string) *Application {
	a := new(Application)
	a.flags = make(map[string]*FlagClause)
	a.args = make(map[string]*ArgClause)
	a.cmds = make(map[string]*CmdClause)
	return a
}

func (a *Application) IsNil() bool {
	if a == nil {
		return true
	}
	return false
}

func (a *Application) Arg(p1, p2 string) clier.ArgClauser {
	arg := NewArg(p1, p2)
	a.args[p1] = arg
	return arg
}

func (a *Application) Flag(p1, p2 string) clier.FlagClauser {
	flag := NewFlag(p1, p2)
	a.flags[p1] = flag
	return flag
}

func (a *Application) Command(p1, p2 string) clier.CmdClauser {
	cmd := NewCmd(p1, p2)
	a.cmds[p1] = cmd
	return cmd
}

func (a *Application) GetCommand(p1 ...string) *CmdClause {
	cmd_len := len(p1)
	if cmd_len == 0 {
		return nil
	}

	if cmd_len == 1 {
		return a.cmds[p1[0]]
	}

	cmd := a.cmds[p1[0]]
	if cmd == nil {
		return nil
	}

	for _, value := range p1[1 : cmd_len-1] {
		cmd, _ = cmd.cmds[value]
		if cmd == nil {
			return nil
		}
	}

	cmd, _ = cmd.cmds[p1[cmd_len-1]]
	return cmd
}

func (a *Application) GetFlag(p1 ...string) *FlagClause {
	flag_len := len(p1)
	var cmd *CmdClause

	if flag_len == 0 {
		return nil
	}

	if flag_len == 1 {
		flag, _ := a.flags[p1[0]]
		return flag
	}

	cmd = a.cmds[p1[0]]
	if cmd == nil {
		return nil
	}

	for _, value := range p1[1 : flag_len-1] {
		cmd, _ = cmd.cmds[value]
		if cmd == nil {
			return nil
		}
	}

	flag, _ := cmd.flags[p1[flag_len-1]]
	return flag
}

func (a *Application) GetArg(p1 ...string) *ArgClause {
	flag_len := len(p1)
	var cmd *CmdClause

	if flag_len == 0 {
		return nil
	}

	if flag_len == 1 {
		flag, _ := a.args[p1[0]]
		return flag
	}

	cmd = a.cmds[p1[0]]
	if cmd == nil {
		return nil
	}

	for _, value := range p1[1 : flag_len-1] {
		cmd, _ = cmd.cmds[value]
		if cmd == nil {
			return nil
		}
	}

	arg, _ := cmd.args[p1[flag_len-1]]
	return arg
}

func (*Application) GetContext(_ []string) (clier.ParseContexter, error) {
	return &ParseContext{}, nil
}
