package cli

import (
	"fmt"
	"github.com/kr/text"
	"github.com/forj-oss/forjj-modules/cli/interface"
	"github.com/forj-oss/forjj-modules/trace"
	"log"
	"regexp"
	"strings"
)

const no_fields = "none"

// ForjObject defines the Object structure
type ForjObject struct {
	cli          *ForjCli                                       // Reference to the parent
	name         string                                         // name of the action to add for objects
	desc         string                                         // Object description string.
	actions      map[string]*ForjObjectAction                   // Collection of actions per objects where flags are added.
	list         map[string]*ForjObjectList                     // List configured for this object.
	internal     bool                                           // true if the object is forjj internal
	sel_actions  map[string]*ForjObjectAction                   // Select several actions to apply for AddParam
	fields       map[string]*ForjField                          // List of fields of this object
	instances    map[string]*ForjObjectInstance                 // Instance detected at Context time.
	err          error                                          // Last error found.
	context_hook func(*ForjObject, *ForjCli, interface{}) error // Parse hook related to this object. Can use cli to create more.
}

func (o *ForjObject) Error() error {
	if o == nil {
		return nil
	}

	return o.err
}

func (o *ForjObject) String() string {
	ret := fmt.Sprintf("Object (%p):\n", o)
	ret += fmt.Sprintf("  cli: %p\n", o.cli)
	ret += fmt.Sprintf("  name: '%s'\n", o.name)
	ret += fmt.Sprintf("  desc: '%s'\n", o.desc)
	ret += fmt.Sprint("  object actions: \n")

	for key, action := range o.actions {
		ret += fmt.Sprintf("    %s: \n", key)
		ret += text.Indent(action.String(), "      ")
	}

	ret += fmt.Sprintf("  internal: '%s'\n", o.internal)
	ret += fmt.Sprint("  fields:\n")
	for key, field := range o.fields {
		ret += fmt.Sprintf("    %s: \n", key)
		ret += text.Indent(field.String(), "      ")
	}
	ret += fmt.Sprint("  instances:\n")
	for key, instance := range o.instances {
		ret += fmt.Sprintf("    %s: \n", key)
		ret += text.Indent(instance.String(), "      ")
	}
	return ret

}

type ForjField struct {
	name       string // name
	help       string // help
	value_type string // Expected value type
	key        bool   // true if this field is a key for list.

	found   bool     // True if the flag was used.
	plugins []string // List of plugins that use this flag.
}

func (f *ForjField) String() string {
	ret := fmt.Sprintf("Field (%p):\n", f)
	ret += fmt.Sprintf("  name: '%s'\n", f.name)
	ret += fmt.Sprintf("  help: '%s'\n", f.help)
	ret += fmt.Sprintf("  value_type: '%s'\n", f.value_type)
	ret += fmt.Sprintf("  found: '%s'\n", f.found)
	return ret
}

// ForjObjectAction defines the action structure for each object
type ForjObjectAction struct {
	name    string               // object action name (formatted as <action>_<object>)
	cmd     clier.CmdClauser     // Object
	action  *ForjAction          // Parent Action name and help
	plugins []string             // Plugins implementing this object action.
	params  map[string]ForjParam // Collection of flags
}

func (a *ForjObjectAction) String() string {
	ret := fmt.Sprintf("Object Action (%p):\n", a)
	ret += fmt.Sprintf("  name: '%s'\n", a.name)
	ret += fmt.Sprintf("  cmd: '%p'\n", a.cmd)
	ret += fmt.Sprint("  params:\n")
	for key, param := range a.params {
		ret += fmt.Sprintf("    %s: \n", key)
		ret += text.Indent(param.String(), "      ")
	}
	ret += fmt.Sprint("  action attached:\n")
	ret += text.Indent(a.action.String(), "      ")
	return ret
}

type ForjObjectInstance struct {
	name              string // Instance name
	additional_fields map[string]*ForjField
}

func (i *ForjObjectInstance) String() string {
	ret := fmt.Sprintf("Object Instance (%p):\n", i)
	ret += fmt.Sprintf("  name: '%s'\n", i.name)
	ret += fmt.Sprint("  fields (map):\n")
	for key, field := range i.additional_fields {
		ret += fmt.Sprintf("    %s: \n", key)
		ret += text.Indent(field.String(), "      ")
	}
	return ret
}

// ---------------------

// NewObjectActions add a new object and the list of actions.
// It creates the ForjAction object for each action/object couple, to attach the object to kingpin object layer.
func (c *ForjCli) NewObject(name, desc string, internal bool) *ForjObject {
	return c.newForjObject(name, desc, internal)
}

func (c *ForjCli) newForjObject(object_name, description string, internal bool) (o *ForjObject) {
	o = new(ForjObject)
	o.actions = make(map[string]*ForjObjectAction)
	o.sel_actions = make(map[string]*ForjObjectAction)
	o.instances = make(map[string]*ForjObjectInstance)
	o.fields = make(map[string]*ForjField)
	o.list = make(map[string]*ForjObjectList)
	o.desc = description
	o.internal = internal
	o.name = object_name
	c.objects[object_name] = o
	o.cli = c
	return
}

// OnActions select several actions from ObjectActions. If list is empty, used the declared object actions.
func (o *ForjObject) OnActions(list ...string) *ForjObject {
	if o == nil {
		return nil
	}
	actions := make([]string, 0, len(o.actions))
	if len(list) == 0 {
		for action_name := range o.actions {
			actions = append(actions, action_name)
		}
	} else {
		actions = list
	}

	// Should reset the map.
	o.sel_actions = make(map[string]*ForjObjectAction)

	for _, name := range actions {
		if action, found := o.actions[name]; found {
			o.sel_actions[name] = action
		}
	}
	return o
}

func (o *ForjObject) ParseHook(context_hook func(*ForjObject, *ForjCli, interface{}) error) *ForjObject {
	if o == nil {
		return nil
	}
	o.context_hook = context_hook
	return o
}

// AddFlag add a flag on the selected list of actions (OnActions)
func (o *ForjObject) AddFlag(name string, options *ForjOpts) *ForjObject {
	if o == nil {
		return nil
	}

	return o.addFlag(func() ForjParam {
		return new(ForjFlag)
	}, name, options)
}

func (o *ForjObject) AddArg(name string, options *ForjOpts) *ForjObject {
	if o == nil {
		return nil
	}
	return o.addFlag(func() ForjParam {
		return new(ForjArg)
	}, name, options)
}

func (o *ForjObject) addFlag(newParam func() ForjParam, name string, options *ForjOpts) *ForjObject {
	if o == nil {
		return nil
	}
	var field *ForjField

	if v, found := o.fields[name]; !found {
		o.err = fmt.Errorf("Unable to find '%s' field in Object '%s'.", name, o.name)
		return nil
	} else {
		field = v
	}

	for _, action := range o.sel_actions {
		p := newParam()

		p.set_cmd(action.cmd, field.value_type, name, field.help, options)

		action.params[name] = p
	}

	return o
}

// DefineActions add a new object and the list of actions.
// It creates the ForjAction object for each action/object couple, to attach the object to kingpin object layer.
func (o *ForjObject) DefineActions(actions ...string) *ForjObject {
	if o == nil {
		return nil
	}

	key_field_found := false
	for _, field := range o.fields {
		if field.key {
			key_field_found = true
			break
		}
	}

	if !key_field_found {
		o.err = fmt.Errorf("Missing key in the object '%s'", o.name)
		return nil
	}

	for _, action := range actions {
		if ar, found := o.cli.actions[action]; found {
			o.actions[action] = newForjObjectAction(ar, o.name, o.desc)
		} else {
			log.Printf("unknown action '%s'. Ignored.", action)
		}
	}
	return o
}

// NoFields add a Key field to the object.
func (o *ForjObject) NoFields() *ForjObject {
	if o == nil {
		return nil
	}

	if len(o.fields) > 0 {
		o.err = fmt.Errorf("The object '%s' cannot be defined no fields if at least field has been added", o.name)
		return nil
	}

	if o.AddField(String, no_fields, "help") == nil {
		return nil
	}

	field := o.fields[no_fields]
	field.key = true
	return o
}

func (o *ForjObject) keyName() string {
	if o == nil {
		return ""
	}
	for field_name, field := range o.fields {
		if field.key {
			return field_name
		}
	}
	return ""
}

// AddKey add a Key field to the object.
func (o *ForjObject) AddKey(pIntType, name, help string) *ForjObject {
	if o == nil {
		return nil
	}

	for field_name, field := range o.fields {
		if field.key {
			o.err = fmt.Errorf("One key already exist in the object '%s', called '%s'", o.name, field_name)
			return nil
		}
	}

	if o.AddField(pIntType, name, help) == nil {
		return nil
	}

	field := o.fields[name]
	field.key = true
	return o
}

// AddField add a field to the object.
func (o *ForjObject) AddField(pIntType, name, help string) *ForjObject {
	if o == nil {
		return nil
	}

	if _, found := o.fields[no_fields]; found {
		o.err = fmt.Errorf("Unable to Add field on a Fake Object.")
	}

	if _, found := o.fields[name]; found {
		gotrace.Trace("Field %s already added in %s. Ignored.", name, o.name)
		return o
	}

	o.fields[name] = &ForjField{
		name:       o.name + "_" + name,
		help:       help,
		value_type: pIntType,
	}
	return o
}

// CreateList create a new list. It returns the ForjObjectList to set it and configure actions
func (o *ForjObject) CreateList(name, list_sep, ext_regexp string) *ForjObjectList {
	if o == nil {
		return nil
	}
	ext_regexp = o.cli.buildCapture(ext_regexp)
	l := new(ForjObjectList)
	if r, err := regexp.Compile(ext_regexp); err != nil {
		o.err = fmt.Errorf("%s_%s not created: Regexp error found: %s", o, name, err)
		return nil
	} else {
		l.ext_regexp = r
		parentheses_reg, _ := regexp.Compile(`\(`)
		l.max_fields = uint(len(parentheses_reg.FindAllString(strings.Replace(ext_regexp, `\(`, "", -1), -1)) + 1)
		gotrace.Trace("Found '%d' group in '%s'", l.max_fields-1, ext_regexp)
	}

	l.fields_name = make(map[uint]string)
	l.name = name
	l.obj = o
	l.obj.list[name] = l
	l.sep = list_sep
	l.key_name = o.keyName()
	l.actions_related = o.actions
	l.actions = make(map[string]*ForjObjectAction)
	l.list = make([]ForjListData, 0, 5)
	l.data = make([]ForjData, 0, 5)
	l.flags_list = make(map[string]*ForjObjectListFlags)
	l.c = o.cli
	o.cli.list[o.name+"_"+name] = l
	return l
}

// AddFlagFromObjectListAction add flag on the select object selected action (OnActions) from object list actions
// identified by obj_name, obj_list, []obj_actions. The flag will be named as --<obj_action>-<obj_name>s
//
// - obj_name, obj_list, obj_action identify the list and action to add as flag
//
// - action where flags will be created.
//
// ex: forjj create workspace --repos ...
//
// At context time this object list will create more detailed flags.
//
// return nil if the obj_list is not found. Otherwise, return the object updated.
func (o *ForjObject) AddFlagFromObjectListAction(obj_name, obj_list, obj_action string) *ForjObject {
	if o == nil {
		return nil
	}

	if obj_name == o.name {
		o.err = fmt.Errorf("Unable to add '%s' object list action flag on itself.", obj_name)
		return nil
	}

	o_object, o_object_list, o_action, err := o.cli.getObjectListAction(obj_name+"_"+obj_list, obj_action)

	if err != nil {
		o.err = fmt.Errorf("Unable to find Object/Object list/action '%s/%s/%s'", obj_name, obj_list, obj_action)
		return nil
	}

	for _, action := range o.sel_actions {
		d_flag := new(ForjFlagList)
		new_object_name := obj_name + "s"

		d_flag.obj = o_object_list
		help := fmt.Sprintf("%s one or more %s", obj_action, o_object.desc)
		d_flag.set_cmd(action.cmd, String, new_object_name, help, nil)
		action.params[new_object_name] = d_flag

		// Need to add all others object fields not managed by the list, but At context time.
		action.action.to_refresh[obj_name] = &ForjContextTime{o_object_list, o_action}

		// Add reference to the Object list for context instance flags creation.
		flags_ref := new(ForjObjectListFlags)
		flags_ref.params = make(map[string]ForjParam)
		flags_ref.multi_actions = false
		flags_ref.objList = o_object_list
		flags_ref.objectAction = action
		o_object_list.flags_list[o.name+" --"+new_object_name] = flags_ref
	}
	return o
}

// AddFlagsFromObjectListActions add flags on the select object selected action (OnActions) from object list actions
// identified by obj_name, obj_list, []obj_actions. The flag will be named as --<obj_action>-<obj_name>s
//
// - obj_name, obj_list, obj_action identify the list and action to add as flags
//
// - action where flags will be created.
//
// ex: forjj create --add-repos ...
//
// At context time this object list will create more detailed flags.
//
// return nil if the obj_list is not found. Otherwise, return the object updated.
func (o *ForjObject) AddFlagsFromObjectListActions(obj_name, obj_list string, obj_actions ...string) *ForjObject {
	if o == nil {
		return nil
	}

	if obj_name == o.name {
		o.err = fmt.Errorf("Unable to add '%s' object list actions flags on itself.", obj_name)
		return nil
	}

	for _, obj_action := range obj_actions {
		o_object, o_object_list, o_action, err := o.cli.getObjectListAction(obj_name+"_"+obj_list, obj_action)

		if err != nil {
			o.err = fmt.Errorf("Unable to find object '%s' action '%s'. Adding flags into selected actions of object '%s' ignored.",
				obj_list, obj_action, o.name)
			return nil
		}

		for _, action := range o.sel_actions {

			new_object_name := obj_action + "-" + obj_name + "s"

			d_flag := new(ForjFlagList)
			d_flag.obj = o_object_list
			help := fmt.Sprintf("%s one or more %s", obj_action, o_object.desc)
			d_flag.set_cmd(action.cmd, String, new_object_name, help, nil)
			action.params[new_object_name] = d_flag

			// Need to add all others object fields not managed by the list, but At context time.
			action.action.to_refresh[obj_name] = &ForjContextTime{o_object_list, o_action}

			// Add reference to the Object list for context instance flags creation.
			flags_ref := new(ForjObjectListFlags)
			flags_ref.params = make(map[string]ForjParam)
			flags_ref.multi_actions = true
			flags_ref.objList = o_object_list
			flags_ref.objectAction = action
			o_object_list.flags_list[action.action.name+" "+o.name+" --"+new_object_name] = flags_ref
		}

	}
	return o
}

func (o *ForjObject) AddFlagsFromObjectAction(obj_name, obj_action string) *ForjObject {
	if o == nil {
		return nil
	}

	if obj_name == o.name {
		o.err = fmt.Errorf("Unable to add '%s' object action flags on itself.", obj_name)
		return nil
	}

	_, o_action, _ := o.cli.getObjectAction(obj_name, obj_action)
	for _, action := range o.sel_actions {
		for param_name, param := range o_action.params {
			var fc ForjParamCopier
			fc = param

			d_flag := fc.CopyToFlag(action.cmd)
			action.params[param_name] = d_flag
		}
	}

	return o
}

// Search for a flag/Arg from the list or additional param (object/list)
func (o *ForjObject) search_object_param(action, object, key, param_name string) (p ForjParam) {
	for _, param := range o.actions[action].params {
		if fl, pi, pn := param.fromList(); fl == nil {
			if o.name != object || pi != key || pn != param_name {
				continue
			}
			return param
		} else {
			if o.name != object {
				continue
			}
			name := param.Name()
			if name == key+"-"+param_name {
				return param
			}
			if name == action+"-"+key+"-"+param_name {
				return param
			}
		}
	}
	return p
}
