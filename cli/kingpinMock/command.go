package kingpinMock

import "github.com/forj-oss/forjj-modules/cli/interface"

type CmdClause struct {
	command string
	help    string
	flags   map[string]*FlagClause
	args    map[string]*ArgClause
	cmds    map[string]*CmdClause
}

func NewCmd(name, help string) (cmd *CmdClause) {
	cmd = new(CmdClause)
	cmd.flags = make(map[string]*FlagClause)
	cmd.args = make(map[string]*ArgClause)
	cmd.cmds = make(map[string]*CmdClause)
	cmd.command = name
	cmd.help = help
	return cmd
}

func (f *CmdClause) IsHelp(help string) bool {
	return (f.help == help)
}

func (c *CmdClause) Command(p1, p2 string) clier.CmdClauser {
	cmd := NewCmd(p1, p2)
	c.cmds[p1] = NewCmd(p1, p2)
	return cmd
}

func (c *CmdClause) IsCommand(p1 string, p2 string) bool {
	return (c.command == p1 && c.help == p2)
}

func (c *CmdClause) Flag(p1, p2 string) clier.FlagClauser {
	flag := NewFlag(p1, p2)
	c.flags[p1] = flag
	return flag
}

func (c *CmdClause) Arg(p1, p2 string) clier.ArgClauser {
	arg := NewArg(p1, p2)
	c.args[p1] = arg
	return arg
}

func (c *CmdClause) FullCommand() string {
	return c.command
}
