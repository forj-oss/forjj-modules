package cli

import (
	"github.com/forj-oss/forjj-modules/cli/interface"
)

// ForjFlag defines the flag structure for each object actions
type ForjFlag struct {
	name       string                 // flag name
	help       string                 // help used for kingpin flag
	value_type string                 // flag type
	options    ForjOpts               // Options
	flag       clier.FlagClauser      // Flag clause.
	flagv      interface{}            // Flag value.
	found      bool                   // True if the flag was used.
	plugins    []string               // List of plugins that use this flag.
	actions    map[string]*ForjAction // List of actions where this flag could be requested.
}

// set the Argument (Param)
// name: name
// help: help
// options: Collection of options. Support required, default, hidden, envar
// actions: List of actions to attach.
func (f *ForjFlag) set_cmd(cmd clier.CmdClauser, paramIntType, name, help string, options *ForjOpts) {
	f.flag = cmd.Flag(name, help)
	f.name = name
	f.help = help
	f.value_type = paramIntType
	if options != nil {
		f.options = *options
	}
	f.set_options(options)

	switch paramIntType {
	case String:
		f.flagv = f.flag.String()
	case Bool:
		f.flagv = f.flag.Bool()
	}
}

func (f *ForjFlag) loadFrom(context clier.ParseContexter) {
	if v, found := context.GetFlagValue(f.flag); found {
		copyValue(f.flagv, v)
		f.found = true
	} else {
		f.found = false
	}
	return
}

// TODO: To apply to a new flag interface.

func (f *ForjFlag) set_options(options *ForjOpts) {
	if options == nil {
		options = &f.options
	}

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

	if v, ok := options.opts["short"]; ok && is_byte(v) {
		f.flag.Short(to_byte(v))
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

func (f *ForjFlag) String() string {
	return f.name
}

func (f *ForjFlag) CopyToFlag(cmd clier.CmdClauser) *ForjFlag {
	flag := new(ForjFlag)
	flag.set_cmd(cmd, f.value_type, f.name, f.help, &f.options)
	return flag
}

func (f *ForjFlag) CopyToArg(cmd clier.CmdClauser) *ForjArg {
	arg := new(ForjArg)
	arg.set_cmd(cmd, f.value_type, f.name, f.help, &f.options)
	return arg
}
