package cli

import (
	"github.com/alecthomas/kingpin"
)

// ForjArgList defines the flag list structure for each object actions
type ForjArgList struct {
	arg            *kingpin.ArgClause     // Arg clause.
	detailed_flags []*kingpin.FlagClause  // Additional flags prefixed by the list key.
	obj            *ForjObjectList        // Object list
	plugins        []string               // List of plugins that use this flag.
	actions        map[string]*ForjAction // List of actions where this flag could be requested.
	key            string                 // Prefix key name for detailed_flags
}

func (f *ForjArgList) loadFrom(context *kingpin.ParseContext) {
	for _, element := range context.Elements {
		if arg, ok := element.Clause.(*kingpin.ArgClause); ok && arg == f.arg {
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
	f.arg = cmd.Arg(f.obj.name+"s", help)

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

func (f *ForjArgList) GetListValues() []ForjData {
	return f.obj.list
}

func (f *ForjArgList) GetValue() interface{} {
	return nil
}

func (f *ForjArgList) IsFound() bool {
	return f.obj.found
}

func (f *ForjArgList) Default(value string) ForjParam {
	if f.arg == nil {
		return nil
	}
	f.arg.Default(value)
	return f
}
