package cli

import (
	"github.com/alecthomas/kingpin"
)

// ForjArg defines the flag structure for each object actions
type ForjArg struct {
	arg     *kingpin.ArgClause     // Arg clause.
	argv    interface{}            // Arg value.
	found   bool                   // True if the flag was used.
	plugins []string               // List of plugins that use this flag.
	actions map[string]*ForjAction // List of actions where this flag could be requested.
}

// Part of ForjParam interface

// set the Argument (Param)
// name: name
// help: help
// options: Collection of options. Support required, default.
// actions: List of actions to attach.
func (a *ForjArg) set_cmd(cmd *kingpin.CmdClause, paramIntType, name, help string, options *ForjOpts) {
	a.arg = cmd.Arg(name, help)
	a.set_options(options)

	switch paramIntType {
	case String:
		a.argv = a.arg.String()
	case Bool:
		a.argv = a.arg.Bool()
	}
}

func (a *ForjArg) loadFrom(context *kingpin.ParseContext) {
	for _, element := range context.Elements {
		if arg, ok := element.Clause.(*kingpin.ArgClause); ok && arg == a.arg {
			copyValue(a.argv, element.Value)
			a.found = true
		}
	}
	return
}

// TODO: To apply to a new arg interface.

func (a *ForjArg) set_options(options *ForjOpts) {
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

func (a *ForjArg) GetListValues() []ForjData {
	return []ForjData{}
}

func (f *ForjArg) GetValue() interface{} {
	return f.argv
}

func (a *ForjArg) IsFound() bool {
	return a.found
}
