package kingpinMock

import (
	"fmt"
	"github.com/kr/text"
	"github.com/forj-oss/forjj-modules/cli/interface"
	"github.com/forj-oss/forjj-modules/trace"
	"log"
	"strings"
)

type Application struct {
	flags   map[string]*FlagClause
	args    map[string]*ArgClause
	cmds    map[string]*CmdClause
	context *ParseContext
}

type ClauseList interface {
	Set(string) error
	IsCumulative() bool
	String() string
}

func (a *Application) String() string {
	ret := fmt.Sprintf("Application (%p):\n", a)
	ret += fmt.Sprint("  Cmds (map):\n")
	for key, cmd := range a.cmds {
		ret += fmt.Sprintf("    %s: \n", key)
		ret += text.Indent(cmd.String(), "      ")
	}
	ret += fmt.Sprint("  Args (map):\n")
	for key, arg := range a.args {
		ret += fmt.Sprintf("    %s: \n", key)
		ret += text.Indent(arg.Stringer(), "      ")
	}
	ret += fmt.Sprint("  Flags (map):\n")
	for key, flag := range a.flags {
		ret += fmt.Sprintf("    %s: \n", key)
		ret += text.Indent(flag.Stringer(), "      ")
	}
	if a.context != nil {
		ret += fmt.Sprint("  Context:\n")
		ret += text.Indent(a.context.String(), "    ")
	}
	return ret
}

func New(_ string) *Application {
	a := new(Application)
	a.flags = make(map[string]*FlagClause)
	a.args = make(map[string]*ArgClause)
	a.cmds = make(map[string]*CmdClause)
	gotrace.Trace("Application created : %#v", a)
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
	gotrace.Trace("Parent App: (%p)%#v, Child Cmd: (%p)%#v (Key: %s)", a, a, a.cmds[p1], a.cmds[p1], p1)
	return arg
}

func (a *Application) Flag(p1, p2 string) clier.FlagClauser {
	flag := NewFlag(p1, p2)
	a.flags[p1] = flag
	gotrace.Trace("Parent App: (%p)%#v, Child Flag: (%p)%#v (Key: %s)", a, a, a.flags[p1], a.flags[p1], p1)
	return flag
}

func (a *Application) Command(p1, p2 string) clier.CmdClauser {
	cmd := NewCmd(p1, p2)
	a.cmds[p1] = cmd
	gotrace.Trace("Parent App: (%p)%#v, Child Arg: (%p)%#v (Key: %s)", a, a, a.args[p1], a.args[p1], p1)
	return cmd
}

func (a *Application) ListOf(p1 ...string) (ret []string) {
	cmd := a.GetCommand(p1...)

	if cmd == nil {
		ret = []string{"nil"}
	}

	ret = make([]string, 0, len(cmd.args)+len(cmd.flags)+len(cmd.cmds))

	for key := range cmd.cmds {
		ret = append(ret, "Cmd:"+key)
	}
	for key := range cmd.args {
		ret = append(ret, "Arg:"+key)
	}
	for key := range cmd.flags {
		ret = append(ret, "Flag:"+key)
	}

	return
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

	if flag_len == 0 {
		return nil
	}

	if flag_len == 1 {
		flag, _ := a.flags[p1[0]]
		return flag
	}

	cmd := a.GetCommand(p1[0 : flag_len-1]...)
	if cmd == nil {
		return nil
	}

	flag, _ := cmd.flags[p1[flag_len-1]]
	return flag
}

func (a *Application) GetArg(p1 ...string) *ArgClause {
	flag_len := len(p1)

	if flag_len == 0 {
		return nil
	}

	if flag_len == 1 {
		flag, _ := a.args[p1[0]]
		return flag
	}

	cmd := a.GetCommand(p1[0 : flag_len-1]...)
	if cmd == nil {
		return nil
	}

	arg, _ := cmd.args[p1[flag_len-1]]
	return arg
}

func (a *Application) ParseContext(args []string) (clier.ParseContexter, error) {
	if a.context == nil {
		a.NewContext()
	}
	if len(args) == 0 {
		return a.context, nil
	}

	cmds := make([]string, 0)
	for _, arg := range args {
		if strings.Contains(arg, "cmd:") {
			cmds = append(cmds, arg[4:])
			continue
		}
		break
	}
	if a.context.SetContext(cmds...) == nil {
		return nil, fmt.Errorf("Issue to set Commands context '%s'", strings.Join(cmds, " "))
	}

	var flag_name string
	for _, arg := range args {
		if strings.Contains(arg, "cmd:") {
			continue
		}
		if flag_name == "" {
			flag_name = arg
		} else {
			if a.context.SetContextValue(flag_name, arg) == nil {
				log.Printf("Unable to add flag/arg '%s' value. Not found. Ignored.", flag_name)
			}
			flag_name = ""
		}
	}
	return a.context, nil
}

func (a *Application) Parse(args []string) (string, error) {
	_, err := a.ParseContext(args)

	list := make([]string, 0, len(a.context.cmds))
	for _, cmd := range a.context.cmds {
		list = append(list, cmd.command)
	}
	for _, element := range a.context.Elements {
		switch element.(type) {
		case *FlagClause:
			f := element.(*FlagClause)
			f.update_data()
		case *ArgClause:
			a := element.(*ArgClause)
			a.update_data()
		}
	}
	return strings.Join(list, " "), err
}
