package kingpinCli

import (
	"github.com/alecthomas/kingpin"
	"github.com/forj-oss/forjj-modules/cli"
)

type KCmdClause interface {
	GetCmd() *kingpin.CmdClause
}

type CmdClause struct {
	cmd *kingpin.CmdClause
}

func (c *CmdClause) Command(p1, p2 string) cli.CmdClauser {
	return &CmdClause{c.cmd.Command(p1, p2)}
}

func (c *CmdClause) Flag(p1, p2 string) cli.FlagClauser {
	return NewFlag(c.cmd.Flag(p1, p2))
}

func (c *CmdClause) Arg(p1, p2 string) cli.ArgClauser {
	return NewArg(c.cmd.Arg(p1, p2))
}

func (c *CmdClause) FullCommand() string {
	return c.cmd.FullCommand()
}

func (c *CmdClause) GetCmd() *kingpin.CmdClause {
	return c.cmd
}

func (c *CmdClause) IsEqualTo(c_ref cli.CmdClauser) bool {
	return (c.cmd == c_ref.(*CmdClause).cmd)
}
