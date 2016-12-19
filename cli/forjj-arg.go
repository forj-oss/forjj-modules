package cli

import (
	"fmt"
	"github.com/forj-oss/forjj-modules/cli/interface"
	"github.com/forj-oss/forjj-modules/trace"
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
	obj_act    *ForjObjectAction      // Set if the flag has been created by an object field. list must be nil.
	obj        *ForjObject            // Set if the flag has been created by an object field. list must be nil.
	// The object instance name must be set to create the object data.
	list          *ForjObjectList // Set if the flag has been created by a list
	instance_name string          // List related: Instance name where this flag is attached.
	field_name    string          // List related: Field name where this flag is attached
	data          *ForjData       // Data set from this flag.
}

func (a *ForjArg) Name() string {
	return a.name
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
	if options != nil {
		if a.options != nil {
			a.options.MergeWith(options)
		} else {
			a.options = options
		}

	}
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

func (a *ForjArg) GetContextValue(context clier.ParseContexter) (interface{}, bool) {
	return context.GetArgValue(a.arg)
}

func (f *ForjArg) IsList() bool {
	return false
}

func (a *ForjArg) isListRelated() bool {
	return (a.list != nil)
}

func (a *ForjArg) isObjectRelated() bool {
	return (a.obj != nil || a.obj_act != nil)
}

func (f *ForjArg) fromList() (*ForjObjectList, string, string) {
	return f.list, f.instance_name, f.field_name
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
	ret := fmt.Sprintf("Arg (%p) - %s", a, a.name)
	if a.data != nil {
		ret += fmt.Sprintf(" (data attached : %p - %d attributes)\n", a.data, len(a.data.attrs))
	} else {
		ret += fmt.Sprint(" (NO data attached )\n")
	}
	return ret
}

// ForjParamCopier interface

func (a *ForjArg) Copier() (p ForjParamCopier) {
	p = ForjParamCopier(a)
	return
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

func (a *ForjArg) UpdateObject() error {
	if a.instance_name == "" || a.field_name == "" {
		if a.field_name != "" {
			gotrace.Trace("Possible issue: Flag field '%s' were created without an instance name attached.", a.field_name)
		}
		return nil
	}

	if a.list != nil {
		return a.updateObject(a.list.obj.cli, a.list.obj.name)
	}

	return a.updateObject(a.obj.cli, a.obj.name)
}

func (a *ForjArg) updateObject(c *ForjCli, object_name string) error {
	var value interface{}

	switch a.argv.(type) {
	case *string:
		value = *a.argv.(*string)
	case *bool:
		value = *a.argv.(*bool)
	default:
		return fmt.Errorf("Unable to convert flagv to object attribute value.")
	}
	c.SetValue(object_name, a.instance_name, a.value_type, a.field_name, value)
	return nil

}

func (a *ForjArg) forjParam() (p forjParam) {
	p = forjParam(a)
	return
}

// ParamListRelated Interface

func (a *ForjArg) forjParamRelated() (p forjParamRelated) {
	p = forjParamRelated(a)
	return
}

func (a *ForjArg) getFieldName() string {
	return a.field_name
}

func (a *ForjArg) getInstanceName() string {
	return a.instance_name
}

func (a *ForjArg) getObjectList() *ForjObjectList {
	return a.list
}

func (a *ForjArg) getObjectAction() *ForjObjectAction {
	return a.obj_act
}

func (a *ForjArg) getObject() *ForjObject {
	return a.obj
}

// --------------------------------
// forjParamRelatedSetter Interface

func (a *ForjArg) forjParamRelatedSetter() (p forjParamRelatedSetter) {
	p = forjParamRelatedSetter(a)
	return
}

func (a *ForjArg) setList(ol *ForjObjectList, instance, field string) {
	a.list = ol
	a.field_name = field
	a.instance_name = instance
}

func (a *ForjArg) setObjectAction(oa *ForjObjectAction, field string) {
	a.obj_act = oa
	a.obj = oa.obj
	a.field_name = field
}

func (a *ForjArg) setObjectField(o *ForjObject, field string) {
	a.obj = o
	a.field_name = field
}

func (a *ForjArg) setObjectInstance(instance string) {
	if a.obj == nil {
		return
	}
	a.instance_name = instance
}

// --------------------------------
// forjParamSetter Interface

func (a *ForjArg) forjParamSetter() forjParamSetter {
	return forjParamSetter(a)
}

func (a *ForjArg) createObjectDataFromParams(params map[string]ForjParam) error {
	if a.obj == nil {
		// Not an object flag.
		return nil
	}
	if err := a.obj.createObjectDataFromParams(params); err != nil {
		return fmt.Errorf("Unable to update Object '%s' from context. %s", a.obj.name, err)
	}
	return nil
}

// --------------------------------
// forjParamDataUpdater Interface

func (a *ForjArg) forjParamUpdater() forjParamUpdater {
	return forjParamUpdater(a)
}

// updateContextData do the context data update as soon as some flag options (default) has been updated/set.
func (a *ForjArg) updateContextData() {
	if a.data == nil {
		return
	}
	if a.obj == nil && a.list == nil {
		return
	}
	if a.obj.cli.cli_context.context == nil || a.obj.cli.parse {
		return
	}
	ctxt := a.obj.cli.cli_context.context
	if v, found := a.GetContextValue(ctxt); found {
		a.data.set(a.value_type, a.field_name, v)
	}
}

func (f *ForjArg) set_ref(data *ForjData) {
	f.data = data
}
