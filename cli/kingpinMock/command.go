package kingpinMock

import (
	"fmt"
	"github.com/kr/text"
	"github.com/forj-oss/forjj-modules/cli/clier"
	"github.com/forj-oss/forjj-modules/trace"
)

type CmdClause struct {
	command string
	help    string
	flags   map[string]*FlagClause
	args    map[string]*ArgClause
	cmds    map[string]*CmdClause
}

func (f *CmdClause) String() string {
	ret := fmt.Sprintf("Cmd (%p):\n", f)
	ret += fmt.Sprintf("  command: '%s'\n", f.command)
	ret += fmt.Sprintf("  help: '%s'\n", f.help)
	ret += fmt.Sprint("  Cmds (map):\n")
	for key, cmd := range f.cmds {
		ret += fmt.Sprintf("    key: %s : \n", key)
		ret += text.Indent(cmd.String(), "      ")
	}
	ret += fmt.Sprint("   Args (map):\n")
	for key, arg := range f.args {
		ret += fmt.Sprintf("    key: %s : \n", key)
		ret += text.Indent(arg.Stringer(), "      ")
	}
	ret += fmt.Sprint("   Flags (map):\n")
	for key, flag := range f.flags {
		ret += fmt.Sprintf("    key: %s : \n", key)
		ret += text.Indent(flag.Stringer(), "      ")
	}
	return ret

}

func NewCmd(name, help string) (cmd *CmdClause) {
	cmd = new(CmdClause)
	cmd.flags = make(map[string]*FlagClause)
	cmd.args = make(map[string]*ArgClause)
	cmd.cmds = make(map[string]*CmdClause)
	cmd.command = name
	cmd.help = help
	gotrace.Trace("Cmd created : %#v", cmd)
	return cmd
}

func (f *CmdClause) IsHelp(help string) bool {
	return (f.help == help)
}

func (c *CmdClause) Command(p1, p2 string) clier.CmdClauser {
	gotrace.Trace("Parent CMD : (%p)%#v", c, c)
	cmd := NewCmd(p1, p2)
	c.cmds[p1] = cmd
	gotrace.Trace("Parent CMD: (%p)%#v, Child Cmd: (%p)%#v (Key: %s)", c, c, c.cmds[p1], c.cmds[p1], p1)
	return cmd
}

func (c *CmdClause) IsCommand(p1 string, p2 string) bool {
	return (c.command == p1 && c.help == p2)
}

func (c *CmdClause) Flag(p1, p2 string) clier.FlagClauser {
	flag := NewFlag(p1, p2)
	c.flags[p1] = flag
	gotrace.Trace("Parent CMD: (%p)%#v, Child Flag: (%p)%#v (Key: %s)", c, c, c.flags[p1], c.flags[p1], p1)
	return flag
}

func (c *CmdClause) Arg(p1, p2 string) clier.ArgClauser {
	arg := NewArg(p1, p2)
	c.args[p1] = arg
	gotrace.Trace("Parent CMD: (%p)%#v, Child Arg: (%p)%#v (Key: %s)", c, c, c.args[p1], c.args[p1], p1)
	return arg
}

func (c *CmdClause) FullCommand() string {
	return c.command
}

func (c *CmdClause) IsEqualTo(c_ref clier.CmdClauser) bool {
	return (c == c_ref.(*CmdClause))
}
