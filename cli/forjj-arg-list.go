package cli

import (
	"fmt"
	"github.com/kr/text"
	"github.com/forj-oss/forjj-modules/cli/interface"
)

// ForjArgList defines the flag list structure for each object actions
type ForjArgList struct {
	name           string                 // Arg list name
	help           string                 // help used for kingpin arg
	value_type     string                 // arg type
	arg            clier.ArgClauser       // Arg clause.
	detailed_flags []clier.FlagClauser    // Additional flags prefixed by the list key.
	obj            *ForjObjectList        // Object list
	plugins        []string               // List of plugins that use this flag.
	actions        map[string]*ForjAction // List of actions where this flag could be requested.
	key            string                 // Prefix key name for detailed_flags
}

func (a *ForjArgList) loadFrom(context clier.ParseContexter) {
	if v, found := context.GetArgValue(a.arg); found {
		a.obj.Set(v)
		a.obj.found = true
	} else {
		a.obj.found = false
	}
	return
}

// set_cmd do set the flag (Param)
// name: name
// help: help
// options: Collection of options. Support required, default, hidden, envar
// actions: List of actions to attach.
//
// It sets the kingpin flag.
// ex:
// forjj create --apps ...
// or
// forjj update infra --apps ...
func (f *ForjArgList) set_cmd(cmd clier.CmdClauser, paramIntType, name, help string, options *ForjOpts) {
	f.name = name
	f.help = help
	f.value_type = paramIntType
	f.arg = cmd.Arg(f.obj.obj.name+"s", help)

	f.set_options(options)

	f.arg.SetValue(f.obj)
}

// TODO: To apply to a new Arg interface.

func (f *ForjArgList) set_options(options *ForjOpts) {
	if options == nil {
		return
	}

	if v, ok := options.opts["required"]; ok && to_bool(v) {
		f.arg.Required()
	}

	if v, ok := options.opts["default"]; ok {
		f.arg.Default(to_string(v))
	}

	if v, ok := options.opts["envar"]; ok {
		f.arg.Envar(to_string(v))
	}
}

func (f *ForjArgList) GetBoolValue() bool {
	return false
}

func (f *ForjArgList) IsList() bool {
	return true
}

func (f *ForjArgList) GetStringValue() string {
	return ""
}

func (f *ForjArgList) GetListValues() []ForjListData {
	return f.obj.list
}

func (f *ForjArgList) GetValue() interface{} {
	return nil
}

func (f *ForjArgList) IsFound() bool {
	return f.obj.found
}

func (f *ForjArgList) Default(value string) (ret ForjParam) {
	if f.arg == nil {
		return
	}
	f.arg.Default(value)
	ret = f
	return
}

func (a *ForjArgList) String() (ret string) {
	ret = fmt.Sprintf("Arg list (%p)\n", a)
	ret += text.Indent(fmt.Sprintf("name : %s\n", a.name), "  ")
	ret += text.Indent(fmt.Sprintf("Object list ref: %p (%s)\n", a.obj, a.obj.name), "  ")
	return
}

func (a *ForjArgList) CopyToFlag(cmd clier.CmdClauser) *ForjFlag {
	return nil
}

func (a *ForjArgList) CopyToArg(cmd clier.CmdClauser) *ForjArg {
	return nil
}
