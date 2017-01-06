package cli

import (
	"github.com/alecthomas/kingpin"
)

// ForjFlagList defines the flag list structure for each object actions
type ForjFlagList struct {
	flag           *kingpin.FlagClause    // Flag clause.
	detailed_flags []*kingpin.FlagClause  // Additional flags prefixed by the list key.
	obj            *ForjObjectList        // Object list
	plugins        []string               // List of plugins that use this flag.
	actions        map[string]*ForjAction // List of actions where this flag could be requested.
	key            string                 // Prefix key name for detailed_flags
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

func (f *ForjFlagList) loadFrom(context *kingpin.ParseContext) {
	for _, element := range context.Elements {
		if flag, ok := element.Clause.(*kingpin.FlagClause); ok && flag == f.flag {
			f.obj.Set(*element.Value)
			f.obj.found = true
		}
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
