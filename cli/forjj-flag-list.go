package cli

import (
	"fmt"
	"github.com/kr/text"
	"forjj-modules/cli/interface"
	"github.com/forj-oss/forjj-modules/trace"
)

// ForjFlagList defines the flag list structure for each object actions
type ForjFlagList struct {
	name       string            // flag list name
	help       string            // help used for kingpin flag
	value_type string            // flag type
	flag       clier.FlagClauser // Flag clause.
	obj        *ForjObjectList   // Object list
	plugins    []string          // List of plugins that use this flag.
	action     string            // Flag context - Action name.

	detailed       bool                // true to add detailed flags from context
	detailed_flags []clier.FlagClauser // Additional flags prefixed by the list key.
}

func (fl *ForjFlagList) Name() string {
	return fl.name
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
	f.flag = cmd.Flag(name, help)
	f.set_options(options)

	f.flag.SetValue(f.obj)
}

func (f *ForjFlagList) loadFrom(context clier.ParseContexter) {
	if v, found := context.GetFlagValue(f.flag); found {
		f.obj.Set(to_string(v))
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

func (a *ForjFlagList) isListRelated() bool {
	return false
}

func (a *ForjFlagList) isObjectRelated() bool {
	return false
}

func (f *ForjFlagList) IsFromObject(obj *ForjObject) bool {
	return (obj == f.obj.obj)
}

func (f *ForjFlagList) getObject() *ForjObject {
	return f.obj.obj
}

func (*ForjFlagList) fromList() (*ForjObjectList, string, string) {
	return nil, "", ""
}

func (f *ForjFlagList) GetBoolValue() bool {
	return false
}

func (f *ForjFlagList) GetStringValue() string {
	return ""
}

func (f *ForjFlagList) GetBoolAddr() *bool {
	return nil
}

func (f *ForjFlagList) GetStringAddr() *string {
	return nil
}

func (f *ForjFlagList) GetContextValue(context clier.ParseContexter) (interface{}, bool) {
	return context.GetFlagValue(f.flag)
}

func (f *ForjFlagList) GetListValues() []ForjData {
	return f.obj.data
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

func (f *ForjFlagList) String() (ret string) {
	ret = fmt.Sprintf("Flag list (%p)\n", f)
	ret += text.Indent(fmt.Sprintf("name : %s\n", f.name), "  ")
	ret += text.Indent(fmt.Sprintf("created in context : %s\n", f.action), "  ")
	ret += text.Indent(fmt.Sprintf("Object list ref: %p (%s)\n", f.obj, f.obj.name), "  ")
	return
}

func (a *ForjFlagList) Copier() (p ForjParamCopier) {
	p = ForjParamCopier(a)
	return
}

func (a *ForjFlagList) CopyToFlag(cmd clier.CmdClauser) *ForjFlag {
	return nil
}

func (a *ForjFlagList) CopyToArg(cmd clier.CmdClauser) *ForjArg {
	return nil
}

func (f *ForjFlagList) forjParam() forjParam {
	return nil
}

func (a *ForjFlagList) GetFlagClauser() clier.FlagClauser {
	return a.flag
}

func (*ForjFlagList) forjParamRelated() forjParamRelated {
	return nil
}

func (*ForjFlagList) getObjectAction() *ForjObjectAction {
	return nil
}

// --------------------------------
// forjParamRelatedSetter Interface - Not Defined for a list

func (a *ForjFlagList) forjParamRelatedSetter() (p forjParamRelatedSetter) {
	p = forjParamRelatedSetter(nil)
	return
}

// --------------------------------
// forjParamSetter Interface

func (a *ForjFlagList) forjParamSetter() forjParamSetter {
	return forjParamSetter(a)
}

func (a *ForjFlagList) forjParamList() forjParamList {
	return forjParamList(a)
}

func (f *ForjFlagList) createObjectDataFromParams(params map[string]ForjParam) error {
	// Initialize context list from context if context is set.
	if f.obj.c.cli_context.context != nil {
		if v, found := f.obj.c.cli_context.context.GetFlagValue(f.flag); found {
			gotrace.Trace("Initializing context list with '%s'", v)
			f.obj.Set(to_string(v))
		} else {
			return nil
		}
	}
	key_name := f.obj.obj.getKeyName()

	var lists_data []ForjListData
	if f.obj.c.parse {
		lists_data = f.obj.list
	} else {
		lists_data = f.obj.context
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

func (a *ForjFlagList) forjParamUpdater() forjParamUpdater {
	return forjParamUpdater(nil)
}

// getInstances return the list of key values found in ParseContext or not.
func (a *ForjFlagList) getInstances() (instances []string) {
	objList := a.obj
	var data_list []ForjListData
	if !a.obj.c.parse {
		data_list = objList.context
	} else {
		data_list = objList.list
	}
	instances = make([]string, 0, len(data_list))
	for _, element := range data_list {
		instances = append(instances, element.Data[objList.key_name])
	}
	return
}

func (f *ForjFlagList) Type() string {
	return f.value_type
}
