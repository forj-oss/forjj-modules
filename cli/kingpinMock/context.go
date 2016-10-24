package kingpinMock

import "github.com/forj-oss/forjj-modules/cli/interface"

type ParseContext struct {
}

func (*ParseContext) GetFlagValue(_ clier.FlagClauser) (string, bool) {
	return "", false
}

func (*ParseContext) GetArgValue(_ clier.ArgClauser) (string, bool) {
	return "", false
}

func (*ParseContext) SelectedCommands() (res []clier.CmdClauser) {
	res = make([]clier.CmdClauser, 0)
	return
}
