package cli

import (
	"github.com/alecthomas/kingpin"
)

// ForjFlagList defines the flag list structure for each object actions
type ForjFlagList struct {
	flag    *kingpin.FlagClause    // Flag clause.
	obj     *ForjObjectList        // Object list
	plugins []string               // List of plugins that use this flag.
	actions map[string]*ForjAction // List of actions where this flag could be requested.
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
func (f *ForjFlagList) set_cmd(cmd *kingpin.CmdClause, paramIntType, name, help string, options *ForjOpts) {
	f.flag = cmd.Flag(f.obj.name+"s", help)

	f.set_options(options)

	f.flag.SetValue(f.obj)
}

// TODO: To apply to a new Arg interface.

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

	if v, ok := options.opts["envar"]; ok {
		f.flag.Envar(to_string(v))
	}
}
