package cli

import (
	"fmt"
	"github.com/kr/text"
	"github.com/forj-oss/forjj-modules/cli/interface"
	"github.com/forj-oss/forjj-modules/trace"
)

// ForjArgList defines the flag list structure for each object actions
type ForjArgList struct {
	name           string              // Arg list name
	help           string              // help used for kingpin arg
	value_type     string              // arg type
	arg            clier.ArgClauser    // Arg clause.
	detailed_flags []clier.FlagClauser // Additional flags prefixed by the list key.
	obj            *ForjObjectList     // Object list
	plugins        []string            // List of plugins that use this flag.
	action         string              // Argument context - Action name.
	key            string              // Prefix key name for detailed_flags
}

func (a *ForjArgList) Name() string {
	return a.name
}

func (a *ForjArgList) loadFrom(context clier.ParseContexter) {
	if v, found := context.GetArgValue(a.arg); found {
		a.obj.Set(to_string(v))
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

func (f *ForjArgList) GetStringValue() string {
	return ""
}

func (f *ForjArgList) GetBoolAddr() *bool {
	return nil
}

func (f *ForjArgList) GetStringAddr() *string {
	return nil
}

func (f *ForjArgList) IsList() bool {
	return true
}

func (a *ForjArgList) isListRelated() bool {
	return false
}

func (a *ForjArgList) isObjectRelated() bool {
	return false
}

func (*ForjArgList) fromList() (*ForjObjectList, string, string) {
	return nil, "", ""
}

func (a *ForjArgList) GetContextValue(context clier.ParseContexter) (interface{}, bool) {
	return context.GetArgValue(a.arg)
}

func (f *ForjArgList) GetListValues() []ForjListData {
	if f.obj.c.parse {
		return f.obj.list
	} else {
		return f.obj.context
	}
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
	ret += text.Indent(fmt.Sprintf("created in context : %s\n", a.action), "  ")
	ret += text.Indent(fmt.Sprintf("Object list ref: %p (%s)\n", a.obj, a.obj.name), "  ")
	return
}

func (a *ForjArgList) Copier() (p ForjParamCopier) {
	p = ForjParamCopier(a)
	return
}

func (a *ForjArgList) CopyToFlag(cmd clier.CmdClauser) *ForjFlag {
	return nil
}

func (a *ForjArgList) CopyToArg(cmd clier.CmdClauser) *ForjArg {
	return nil
}

func (a *ForjArgList) forjParam() forjParam {
	return nil
}

func (a *ForjArgList) GetArgClauser() clier.ArgClauser {
	return a.arg
}

func (a *ForjArgList) forjParamRelated() forjParamRelated {
	return nil
}

func (a *ForjArgList) getObjectAction() *ForjObjectAction {
	return nil
}

// forjParamRelatedSetter Interface - Not Defined for a list

func (a *ForjArgList) forjParamRelatedSetter() (p forjParamRelatedSetter) {
	p = forjParamRelatedSetter(nil)
	return
}

// --------------------------------
// forjParamSetter Interface

func (a *ForjArgList) forjParamSetter() forjParamSetter {
	return forjParamSetter(a)
}

func (f *ForjArgList) createObjectDataFromParams(params map[string]ForjParam) error {
	// Initialize context list from context if context is set.
	if f.obj.c.cli_context.context != nil {
		if v, found := f.obj.c.cli_context.context.GetArgValue(f.arg); found {
			gotrace.Trace("Initializing context list with '%s'", v)
			f.obj.Set(to_string(v))
		} else {
			return nil
		}
	}
	key_name := f.obj.obj.getKeyName()

	var lists_data []ForjListData
	if f.obj.c.parse {
		lists_data = f.obj.context
	} else {
		lists_data = f.obj.list
	}

	for _, list_data := range lists_data {
		key_value := list_data.Data[key_name]
		data := f.obj.c.setObjectAttributes(f.action, f.obj.obj.name, key_value)
		if data == nil {
			return f.obj.c.err
		}
		for key, attr := range list_data.Data {
			data.attrs[key] = attr
		}
	}
	return nil
}

func (a *ForjArgList) forjParamUpdater() forjParamUpdater {
	return forjParamUpdater(nil)
}
