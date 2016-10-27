package cli

import (
	"github.com/forj-oss/forjj-modules/cli/interface"
)

// ForjArg defines the flag structure for each object actions
type ForjArg struct {
	name       string                 // Argument name
	help       string                 // help used for kingpin arg
	value_type string                 // arg type
	options    *ForjOpts              // options used to create arg
	arg        clier.ArgClauser       // Arg clause.
	argv       interface{}            // Arg value.
	found      bool                   // True if the flag was used.
	plugins    []string               // List of plugins that use this flag.
	actions    map[string]*ForjAction // List of actions where this flag could be requested.
	list       *ForjObjectList        // Set if the flag has been created by a list
	objectData *ForjData              // Object instance Data where this flag will store data
}

// Part of ForjParam interface

// set the Argument (Param)
// name: name
// help: help
// options: Collection of options. Support required, default.
// actions: List of actions to attach.
func (a *ForjArg) set_cmd(cmd clier.CmdClauser, paramIntType, name, help string, options *ForjOpts) {
	a.arg = cmd.Arg(name, help)
	a.name = name
	a.help = help
	a.value_type = paramIntType
	a.options = options
	a.set_options(options)

	switch paramIntType {
	case String:
		a.argv = a.arg.String()
	case Bool:
		a.argv = a.arg.Bool()
	}
}

func (a *ForjArg) loadFrom(context clier.ParseContexter) {
	if v, found := context.GetArgValue(a.arg); found {
		copyValue(a.argv, v)
		a.found = true
	} else {
		a.found = false
	}
	return
}

// TODO: To apply to a new arg interface.

func (a *ForjArg) set_options(options *ForjOpts) {
	if options == nil {
		options = a.options
	}

	if options == nil {
		return
	}

	if v, ok := options.opts["required"]; ok && to_bool(v) {
		a.arg.Required()
	}

	if v, ok := options.opts["default"]; ok {
		a.arg.Default(to_string(v))
	}

	/*    if v, ok := options.opts["hidden"]; ok && to_bool(v) {
	          f.arg.Hidden()
	      }

	      if v, ok := options.opts["envar"]; ok {
	          f.arg.Envar(to_string(v))
	      } */
}

func (a *ForjArg) GetBoolValue() bool {
	return to_bool(a.argv)
}

func (a *ForjArg) GetStringValue() string {
	return to_string(a.argv)
}

func (f *ForjArg) IsList() bool {
	return false
}

func (a *ForjArg) GetListValues() []ForjData {
	return []ForjData{}
}

func (f *ForjArg) GetValue() interface{} {
	return f.argv
}

func (a *ForjArg) IsFound() bool {
	return a.found
}

func (a *ForjArg) Default(value string) ForjParam {
	if a.arg == nil {
		return nil
	}
	a.arg.Default(value)
	return a
}

func (a *ForjArg) String() string {
	return a.name
}

func (a *ForjArg) CopyToFlag(cmd clier.CmdClauser) *ForjFlag {
	flag := new(ForjFlag)
	flag.set_cmd(cmd, a.value_type, a.name, a.help, a.options)
	return flag
}

func (a *ForjArg) CopyToArg(cmd clier.CmdClauser) *ForjArg {
	arg := new(ForjArg)
	arg.set_cmd(cmd, a.value_type, a.name, a.help, a.options)
	return arg
}

func (a *ForjArg) GetArg() *ForjArg {
	return a
}

func (*ForjArg) GetFlag() *ForjFlag {
	return nil
}
