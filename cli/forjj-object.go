package cli

import (
	"fmt"
	"github.com/forj-oss/forjj-modules/trace"
	"log"
	"regexp"
	"github.com/forj-oss/forjj-modules/cli/interface"
)

// ForjObject defines the Object structure
type ForjObject struct {
	cli         *ForjCli                       // Reference to the parent
	name        string                         // name of the action to add for objects
	desc        string                         // Object description string.
	actions     map[string]*ForjObjectAction   // Collection of actions per objects where flags are added.
	list        *ForjObjectList                // List configured for this object.
	internal    bool                           // true if the object is forjj internal
	sel_actions map[string]*ForjObjectAction   // Select several actions to apply for AddParam
	fields      map[string]*ForjField          // List of fields of this object
	instances   map[string]*ForjObjectInstance // Instance detected at Context time.
}

type ForjField struct {
	name       string // name
	help       string // help
	value_type string // Expected value type

	found   bool     // True if the flag was used.
	plugins []string // List of plugins that use this flag.
}

// ForjObjectAction defines the action structure for each object
type ForjObjectAction struct {
	cmd     clier.CmdClauser     // Object
	action  *ForjAction          // Action name and help
	plugins []string             // Plugins implementing this object action.
	params  map[string]ForjParam // Collection of flags
}

type ForjObjectInstance struct {
	name              string // Instance name
	additional_fields map[string]*ForjField
}

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
	o.desc = description
	o.internal = internal
	c.objects[object_name] = o
	o.cli = c
	return
}

// OnActions select several actions from ObjectActions. If list is empty, used the declared object actions.
func (o *ForjObject) OnActions(list ...string) *ForjObject {
	actions := make([]string, 0, len(o.actions))
	if len(actions) == 0 {
		for action_name := range o.actions {
			actions = append(actions, action_name)
		}
	} else {
		actions = list
	}

	for action_to_delete := range o.sel_actions {
		delete(o.sel_actions, action_to_delete)
	}

	for _, name := range actions {
		if action, found := o.actions[name]; found {
			o.sel_actions[name] = action
		}
	}
	return o
}

// AddFlag add a flag on the selected list of actions (OnActions)
func (o *ForjObject) AddFlag(name string, options *ForjOpts) *ForjObject {
	return o.addFlag(func() ForjParam {
		return new(ForjFlag)
	}, name, options)
}

func (o *ForjObject) AddArg(name string, options *ForjOpts) *ForjObject {
	return o.addFlag(func() ForjParam {
		return new(ForjArg)
	}, name, options)
}

func (o *ForjObject) addFlag(newParam func() ForjParam, name string, options *ForjOpts) *ForjObject {
	var field *ForjField

	if v, found := o.fields[name]; !found {
		gotrace.Trace("Unable to find '%s' field in Object '%s'.", name, o.name)
		return o
	} else {
		field = v
	}

	for _, action := range o.sel_actions {
		p := newParam()

		p.set_cmd(action.cmd, field.value_type, name, field.help, options)
	}

	return o
}

// DefineActions add a new object and the list of actions.
// It creates the ForjAction object for each action/object couple, to attach the object to kingpin object layer.
func (o *ForjObject) DefineActions(actions ...string) *ForjObject {
	for _, action := range actions {
		if ar, found := o.cli.actions[action]; found {
			o.actions[action] = newForjObjectAction(ar, o.name, fmt.Sprintf(ar.help, o.desc))
		} else {
			log.Printf("unknown action '%s'. Ignored.", action)
		}
	}
	return o
}

// AddField add a field to the object.
func (o *ForjObject) AddField(pIntType, name, help string) *ForjObject {
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
func (o *ForjObject) CreateList(name, list_sep, ext_regexp, key_name string) *ForjObjectList {
	ext_regexp = o.cli.buildCapture(ext_regexp)
	l := new(ForjObjectList)
	if r, err := regexp.Compile(ext_regexp); err != nil {
		gotrace.Trace("%s_%s not created: Regexp error found: %s", o, name, err)
	} else {
		l.ext_regexp = r
	}

	l.fields_name = make(map[uint]string)
	l.name = name
	l.obj = o
	l.obj.list = l
	l.sep = list_sep
	l.key_name = key_name
	l.actions_related = o.actions
	l.list = make([]ForjData, 0, 5)
	l.c = o.cli
	o.cli.list[o.name+"_"+name] = l
	return l
}

// AddFlagFromObjectListAction add an object list action.
//
// - obj_name, obj_list, obj_action identify the list and action to add as flags
//
// - prefix true if the action is added to flags created.
//
// - detailed true if detailed object fields will be added and prefixed by the object list key value.
//
// - action where flags will be created.
//
// ex: forjj create --repos ...
//
// At context time this object list will create more detailed flags.
func (o *ForjObject) AddFlagFromObjectListAction(obj_name, obj_list, obj_action string) *ForjObject {
	o_object, o_object_list, o_action, _ := o.cli.getObjectListAction(obj_list, obj_action)

	for _, action := range o.sel_actions {
		d_flag := new(ForjFlagList)
		help := fmt.Sprintf("%s one or more %s", obj_action, o_object.desc)
		d_flag.set_cmd(action.cmd, String, obj_name, help, nil)
		action.params[obj_name+"s"] = d_flag

		// Need to add all others object fields not managed by the list, but At context time.
		action.action.to_refresh[obj_name] = &ForjContextTime{o_object_list, o_action}
	}
	return o
}

// AddFlagFromObjectListActions add an object list action.
//
// - obj_name, obj_list, obj_action identify the list and action to add as flags
//
// - prefix true if the action is added to flags created.
//
// - detailed true if detailed object fields will be added and prefixed by the object list key value.
//
// - action where flags will be created.
//
// ex: forjj create --repos ...
//
// At context time this object list will create more detailed flags.
func (o *ForjObject) AddFlagFromObjectListActions(obj_name, obj_list string, obj_actions ...string) *ForjObject {
	for _, obj_action := range obj_actions {
		o_object, o_object_list, o_action, err := o.cli.getObjectListAction(obj_list, obj_action)

		if err != nil {
			gotrace.Trace("Unable to find object '%s' action '%s'. Adding flags into selected actions of object '%s' ignored.",
				obj_list, obj_action, o.name)
			return o
		}

		for _, action := range o.sel_actions {

			new_object_name := obj_action + "-" + obj_name

			d_flag := new(ForjFlagList)
			help := fmt.Sprintf("%s one or more %s", obj_action, o_object.desc)
			d_flag.set_cmd(action.cmd, String, new_object_name, help, nil)
			action.params[new_object_name+"s"] = d_flag

			// Need to add all others object fields not managed by the list, but At context time.
			action.action.to_refresh[obj_name] = &ForjContextTime{o_object_list, o_action}
		}

	}
	return o
}

func (o *ForjObject) AddFlagsFromObjectAction(obj_name, obj_action string) *ForjObject {
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
