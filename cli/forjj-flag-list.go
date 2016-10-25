package cli

import (
	"github.com/forj-oss/forjj-modules/cli/interface"
)

// ForjFlagList defines the flag list structure for each object actions
type ForjFlagList struct {
	name       string                 // flag list name
	help       string                 // help used for kingpin flag
	value_type string                 // flag type
	flag       clier.FlagClauser      // Flag clause.
	obj        *ForjObjectList        // Object list
	plugins    []string               // List of plugins that use this flag.
	actions    map[string]*ForjAction // List of actions where this flag could be requested.

	detailed       bool                // true to add detailed flags from context
	detailed_flags []clier.FlagClauser // Additional flags prefixed by the list key.
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
func (f *ForjFlagList) set_cmd(cmd clier.CmdClauser, paramIntType, name, help string, options *ForjOpts) {
	f.name = name
	f.help = help
	f.value_type = paramIntType
	f.flag = cmd.Flag(name+"s", help)
	f.set_options(options)

	f.flag.SetValue(f.obj)
}

func (f *ForjFlagList) loadFrom(context clier.ParseContexter) {
	if v, found := context.GetFlagValue(f.flag); found {
		f.obj.Set(v)
		f.obj.found = true
	} else {
		f.obj.found = false
	}
	return
}

// TODO: To apply to a new flag interface.

func (f *ForjFlagList) set_options(options *ForjOpts) {
	if options == nil {
		return
	}

	if v, ok := options.opts["required"]; ok && to_bool(v) {
		f.flag.Required()
	}

	if v, ok := options.opts["default"]; ok {
		f.flag.Default(to_string(v))
	}

	if v, ok := options.opts["hidden"]; ok && to_bool(v) {
		f.flag.Hidden()
	}

	if v, ok := options.opts["envar"]; ok {
		f.flag.Envar(to_string(v))
	}

	if v, ok := options.opts["short"]; ok && is_rune(v) {
		f.flag.Short(to_rune(v))
	}
}

func (f *ForjFlagList) IsList() bool {
	return true
}

func (f *ForjFlagList) GetBoolValue() bool {
	return false
}

func (f *ForjFlagList) GetStringValue() string {
	return ""
}

func (f *ForjFlagList) GetListValues() []ForjData {
	return f.obj.list
}

func (f *ForjFlagList) GetValue() interface{} {
	return nil
}

func (f *ForjFlagList) IsFound() bool {
	return f.obj.found
}

func (f *ForjFlagList) Default(value string) ForjParam {
	if f.flag == nil {
		return nil
	}
	f.flag.Default(value)
	return f
}

func (f *ForjFlagList) String() string {
	return f.name
}

func (a *ForjFlagList) CopyToFlag(cmd clier.CmdClauser) *ForjFlag {
	return nil
}

func (a *ForjFlagList) CopyToArg(cmd clier.CmdClauser) *ForjArg {
	return nil
}
