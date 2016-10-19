package kingpinCli

import (
	"github.com/alecthomas/kingpin"
	"github.com/forj-oss/forjj-modules/cli/interface"
)

type CmdClause struct {
	cmd *kingpin.CmdClause
}

func (c *CmdClause) Command(p1, p2 string) clier.CmdClauser {
	return &CmdClause{c.cmd.Command(p1, p2)}
}

func (c *CmdClause) Flag(p1, p2 string) clier.FlagClauser {
	return new(FlagClause)
}

func (c *CmdClause) Arg(p1, p2 string) clier.ArgClauser {
	return new(ArgClause)
}
