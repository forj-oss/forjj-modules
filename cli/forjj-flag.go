package cli

import (
	"github.com/alecthomas/kingpin"
)

// ForjFlag defines the flag structure for each object actions
type ForjFlag struct {
	flag    *kingpin.FlagClause    // Flag clause.
	flagv   interface{}            // Flag value.
	found   bool                   // True if the flag was used.
	plugins []string               // List of plugins that use this flag.
	actions map[string]*ForjAction // List of actions where this flag could be requested.
}

// set the Argument (Param)
// name: name
// help: help
// options: Collection of options. Support required, default, hidden, envar
// actions: List of actions to attach.
func (f *ForjFlag) set_cmd(cmd *kingpin.CmdClause, paramIntType, name, help string, options *ForjOpts) {
	f.flag = cmd.Flag(name, help)

	f.set_options(options)

	switch paramIntType {
	case String:
		f.flagv = f.flag.String()
	case Bool:
		f.flagv = f.flag.Bool()
	}
}

func (f *ForjFlag) loadFrom(context *kingpin.ParseContext) {
	for _, element := range context.Elements {
		if flag, ok := element.Clause.(*kingpin.FlagClause); ok && flag == f.flag {
			f.found = true
			copyValue(f.flagv, element.Value)
		}
	}
	return
}

// TODO: To apply to a new flag interface.

func (f *ForjFlag) set_options(options *ForjOpts) {
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

func (f *ForjFlag) GetBoolValue() bool {
	return to_bool(f.flagv)
}

func (f *ForjFlag) GetStringValue() string {
	return to_string(f.flagv)
}

func (f *ForjFlag) IsList() bool {
	return false
}

func (f *ForjFlag) GetListValues() []ForjData {
	return []ForjData{}
}

func (f *ForjFlag) GetValue() interface{} {
	return f.flagv
}

func (f *ForjFlag) IsFound() bool {
	return f.found
}

func (f *ForjFlag) Default(value string) ForjParam {
	if f.flag == nil {
		return nil
	}
	f.flag.Default(value)
	return f
}
