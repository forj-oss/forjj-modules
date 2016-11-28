package cli

import (
	"fmt"
	"github.com/forj-oss/forjj-modules/cli/interface"
)

// ForjFlag defines the flag structure for each object actions
type ForjFlag struct {
	name       string                 // flag name
	help       string                 // help used for kingpin flag
	value_type string                 // flag type
	options    *ForjOpts              // Options
	flag       clier.FlagClauser      // Flag clause.
	flagv      interface{}            // Flag value.
	found      bool                   // True if the flag was used.
	plugins    []string               // List of plugins that use this flag.
	actions    map[string]*ForjAction // List of actions where this flag could be requested.
	obj        *ForjObjectAction      // Set if the flag has been created by an object field. list must be nil.
	// The object instance name must be set to create the object data.
	list          *ForjObjectList // Set if the flag has been created by a list. obj must be nil.
	instance_name string          // List/object related: Instance name where this flag is attached.
	field_name    string          // List/object related: Field name where this flag is attached
}

func (f *ForjFlag) Name() string {
	return f.name
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
		if f.options != nil {
			f.options.MergeWith(options)
		} else {
			f.options = options
		}

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
		options = f.options
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

func (f *ForjFlag) GetContextValue(context clier.ParseContexter) (interface{}, bool) {
	return context.GetFlagValue(f.flag)
}

func (f *ForjFlag) IsList() bool {
	return false
}

func (f *ForjFlag) isListRelated() bool {
	return (f.list != nil)
}

func (f *ForjFlag) isObjectRelated() bool {
	return (f.obj != nil)
}

func (f *ForjFlag) fromList() (*ForjObjectList, string, string) {
	return f.list, f.instance_name, f.field_name
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
	return fmt.Sprintf("Flag (%p) - %s \n", f, f.name)
}

// ----------------------------
// ForjParamCopier interface

func (a *ForjFlag) Copier() (p ForjParamCopier) {
	p = ForjParamCopier(a)
	return
}

func (f *ForjFlag) CopyToFlag(cmd clier.CmdClauser) *ForjFlag {
	flag := new(ForjFlag)
	flag.set_cmd(cmd, f.value_type, f.name, f.help, f.options)
	return flag
}

func (f *ForjFlag) CopyToArg(cmd clier.CmdClauser) *ForjArg {
	arg := new(ForjArg)
	arg.set_cmd(cmd, f.value_type, f.name, f.help, f.options)
	return arg
}

func (*ForjFlag) GetArg() *ForjArg {
	return nil
}

func (f *ForjFlag) GetFlag() *ForjFlag {
	return f
}

func (f *ForjFlag) UpdateObject() error {
	if f.list == nil {
		return nil
	}
	if f.instance_name == "" || f.field_name == "" {
		return nil
	}

	var value interface{}

	switch f.flagv.(type) {
	case *string:
		value = *f.flagv.(*string)
	case *bool:
		value = *f.flagv.(*bool)
	default:
		return fmt.Errorf("Unable to convert flagv to object attribute value.")
	}
	c := f.list.obj.cli
	c.values[f.list.obj.name].records[f.instance_name].attrs[f.field_name] = value
	return nil
}

func (f *ForjFlag) forjParam() (p forjParam) {
	p = forjParam(f)
	return
}

// ----------------------------
// ParamListRelated Interface

func (a *ForjFlag) forjParamRelated() (p forjParamRelated) {
	p = forjParamRelated(a)
	return
}

func (a *ForjFlag) getFieldName() string {
	return a.field_name
}

func (a *ForjFlag) getInstanceName() string {
	return a.instance_name
}

func (a *ForjFlag) getObjectList() *ForjObjectList {
	return a.list
}

func (a *ForjFlag) getObjectAction() *ForjObjectAction {
	return a.obj
}

// ----------------------------
// forjParamRelatedSetter Interface

func (a *ForjFlag) forjParamRelatedSetter() (p forjParamRelatedSetter) {
	p = forjParamRelatedSetter(a)
	return
}

func (a *ForjFlag) setList(ol *ForjObjectList, instance, field string) {
	a.list = ol
	a.field_name = field
	a.instance_name = instance
}

func (a *ForjFlag) setObject(oa *ForjObjectAction, field string) {
	a.obj = oa
	a.field_name = field
}

func (a *ForjFlag) setObjectInstance(instance string) {
	if a.obj == nil {
		return
	}
	a.instance_name = instance
}

// --------------------------------
// forjParamSetter Interface

func (a *ForjFlag) forjParamSetter() forjParamSetter {
	return forjParamSetter(a)
}

func (f *ForjFlag) createObjectDataFromParams(params map[string]ForjParam) error {
	if f.obj == nil {
		// Not an object flag.
		return nil
	}
	if err := f.obj.obj.createObjectDataFromParams(params); err != nil {
		return fmt.Errorf("Unable to update Object '%s' from context. %s", f.obj.obj.name, err)
	}
	return nil
}
