package cli

import (
	"github.com/alecthomas/kingpin"
)

// ForjArgList defines the flag list structure for each object actions
type ForjArgList struct {
	flag    *kingpin.FlagClause    // Flag clause.
	obj     *ForjObjectList        // Object list
	plugins []string               // List of plugins that use this flag.
	actions map[string]*ForjAction // List of actions where this flag could be requested.
}

func (f *ForjArgList) loadFrom(context *kingpin.ParseContext) {
	for _, element := range context.Elements {
		if flag, ok := element.Clause.(*kingpin.FlagClause); ok && flag == f.flag {
			f.obj.Set(*element.Value)
			f.obj.found = true
		}
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
func (f *ForjArgList) set_cmd(cmd *kingpin.CmdClause, paramIntType, name, help string, options *ForjOpts) {
	f.flag = cmd.Flag(f.obj.name+"s", help)

	f.set_options(options)

	f.flag.SetValue(f.obj)
}

// TODO: To apply to a new Arg interface.

func (f *ForjArgList) set_options(options *ForjOpts) {
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

func (f *ForjArgList) GetBoolValue() bool {
	return false
}

func (f *ForjArgList) GetStringValue() string {
	return ""
}

func (f *ForjArgList) GetListValues() []ForjData {
	return f.obj.list
}

func (f *ForjArgList) GetValue() interface{} {
	return nil
}

func (f *ForjArgList) IsFound() bool {
	return f.obj.found
}
