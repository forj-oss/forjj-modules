package cli

import (
	"fmt"
	"github.com/kr/text"
	"github.com/forj-oss/forjj-modules/trace"
	"regexp"
	"strings"
)

// ForjObjectList defines a list to apply as new Cmd or as flags/args to any other objects/actions
//
// Ex: forjj create =>repos list --<instance>-flags value ...<=
// identified by actions map[string]*ForjObjectAction.
//
// or with ForjObjectListFlags
//
// Ex: forjj create infra =>--repos list --<instance>-flags value ...<=
// identified by flags_list
type ForjObjectList struct {
	c               *ForjCli                        // Reference to the cli object
	name            string                          // List name
	help            string                          // help
	obj             *ForjObject                     // Object attached
	sep             string                          // List separator
	max_fields      uint                            // Number of captured fields defined by the RegExp.
	ext_regexp      *regexp.Regexp                  // Capturing Regexp
	fields_name     map[uint]string                 // Data fields extraction
	actions_related map[string]*ForjObjectAction    // Possible actions for this list
	actions         map[string]*ForjObjectAction    // Collection of actions per objects where flags are added.
	context         []ForjListData                  // Data list collected from the list of flags found in the cli context.
	list            []ForjListData                  // Data list collected from the list of flags found in the cli.
	data            []ForjData                      // Objects list generated from data list collected.
	found           bool                            // True if the list flag was provided.
	key_name        string                          // List key name to use for any detailed flags.
	valid_handler   func(*ForjListData) error       // Handler to validate data collected and correct if needed.
	flags_list      map[string]*ForjObjectListFlags // list of ObjectList flags added
}

type ForjObjectListFlags struct {
	objList *ForjObjectList // Object list referenced
	// If action and objectAction are nil, then the object list flag is attached to the App
	action        *ForjAction          // Otherwise this action is set where flags are added
	objectAction  *ForjObjectAction    // otherwise this ObjectAction is set where flags are added
	params        map[string]ForjParam // List of flags added.
	multi_actions bool                 // true if added for multiple object list actions
}

type ForjListData struct {
	Data map[string]string
}

func (o *ForjObject) getKeyName() string {
	for field_name, field := range o.fields {
		if field.key {
			return field_name
		}
	}
	return ""
}

func (l *ForjObjectList) getParamListObjectName() string {
	return l.obj.name + "s"
}

// AddActions Add the list actions.
// Ex: forjj add repos <blabla>.
//
// `add` must be an existing action
// `repos` is a new object action with a cmd attached.
// `<blabla>` is a arg string that must be attached to the new object action.
//
// The new object action is called with the referenced action to add action to.
// The Cmd is called `object_name` + `s`
//
// The new Argument is called as `object_name` + `s`.
// In Cmd, it is called `object_name` + `s`
//
// kingpin: The function creates a new command and an attached argument. The argument is managed by ForjObjectList.
//
// It returns the base object.
// The list key value will be used at context time to add contexted flag prefixed by the key value.
func (l *ForjObjectList) AddActions(actions ...string) *ForjObjectList {
	if l == nil {
		return l
	}

	for _, action := range actions {
		if v, found := l.actions_related[action]; found {
			// Create a new Command with 's' at the end.
			object_name := l.getParamListObjectName()
			list_action := newForjObjectAction(v.action, object_name, l.help)
			l.actions[action] = list_action

			// Create a new Argument of the object as list (the 's' is added automatically to the argument name)
			arg_list := new(ForjArgList)
			arg_list.obj = l
			arg_list.set_cmd(list_action.cmd, List, object_name, l.help, Opts().Required())
			list_action.params[object_name] = arg_list
		}
	}
	return l
}

// Field add a Map RegExp result to an object field parameter.
//
// - index is the parenthesis capture regexp index
//
// - field_name must be declared in the object list of fields.
//
// Return nil if there is any issue. Otherwise, returns the list object.
func (l *ForjObjectList) Field(index uint, field_name string) *ForjObjectList {
	if l == nil {
		return nil
	}
	if index < 1 {
		l.obj.err = fmt.Errorf("Index < 1 are invalid. Must start at %d. Ignored.", 1)
		return nil
	}
	if index > l.max_fields-1 {
		l.obj.err = fmt.Errorf("Cannot define field at position %d. Regexp '%s' has Max %d capture elements [0..%d]. Ignored.",
			index, l.ext_regexp.String(), l.max_fields, l.max_fields-1)
		return nil
	}
	if _, found := l.obj.fields[field_name]; !found {
		l.obj.err = fmt.Errorf("Cannot define field if the object '%s' has no field '%s'. Ignored.", l.obj.name, field_name)
		return nil
	}

	// Update the list of actions where this field is requested.
	// Final, we got a list of actions where all fields are requested.
	if l.inter_actions_list(l.obj.get_actions_list_from(field_name)) == nil {
		l.obj.err = fmt.Errorf("Adding field '%s' has reduced list of valid action to none. "+
			"Mainly because NO actions has all previous fields and '%s' at the same time. \n%s ",
			field_name, field_name, l.obj.err)
		return nil
	}

	l.fields_name[index] = field_name
	return l
}

// Set function for kingpin.Value interface
func (l *ForjObjectList) Set(value string) error {
	list := Split(" *"+l.sep+" *", value, l.sep)
	gotrace.Trace("Interpret list: %d records identified.", len(list))
	for i, v := range list {
		if err := l.add(v); err != nil {
			return fmt.Errorf("At index %d: %s", i, err)
		}
	}
	return nil
}

// Called by Set to add a new element in the list.
func (l *ForjObjectList) add(value string) error {
	res := l.ext_regexp.FindStringSubmatch(value)
	if res == nil {
		return fmt.Errorf("The string portion '%s' is an invalid %s description. It must respect regular expression '%s'.",
			value, l.obj.name, l.ext_regexp.String())
	}

	dd := ForjListData{make(map[string]string)}

	for index, field_name := range l.fields_name {
		dd.Data[field_name] = res[index]
	}

	if l.valid_handler != nil {
		if err := l.valid_handler(&dd); err != nil {
			return err
		}
	}
	if l.c.parse {
		l.list = append(l.list, dd)
	} else {
		l.context = append(l.context, dd)
	}

	gotrace.Trace("'%s'(%s) added '%s'", l.obj.name, l.name, value)
	return nil
}

// FIXME: kingpin is having trouble in the context case, where several --<object>s set, with some flags in between, is ignoring seconds and next --apps flags values. Workaround is to have them all followed or use the --apps APP[,APP ...] format.

// Inform kingpin that flag is cumulative.
func (d *ForjObjectList) IsCumulative() bool {
	return true
}

// String : Set function for kingpin.Value interface
func (d *ForjObjectList) String() (ret string) {
	if d == nil {
		return ""
	}
	ret += fmt.Sprintf("object referenced: %p (%s)\n", d.obj, d.obj.name)
	ret += fmt.Sprintf("Reg Ref : %s %s ...\n", d.ext_regexp.String(), d.sep)
	ret += fmt.Sprintf("Fields extracted : %d (Index Max : 0..%d)\n", len(d.fields_name), d.max_fields-1)
	for index, name := range d.fields_name {
		ret += text.Indent(fmt.Sprintf("%d: %s\n", index, name), "  ")
	}
	ret += fmt.Sprintf("key name: %s\n", d.key_name)
	ret += "context data list:\n"
	if len(d.context) > 0 {
		list := make([]string, 0, len(d.context))
		for _, v := range d.context {
			for key, value := range v.Data {
				list = append(list, key+"='"+value+"'")
			}
		}
		ret += text.Indent(strings.Join(list, ", ")+"\n", "  ")
	} else {
		ret += text.Indent("-- empty --\n", "  ")
	}
	ret += "data list:\n"
	if len(d.list) > 0 {
		list := make([]string, 0, len(d.list))
		for _, v := range d.list {
			for key, value := range v.Data {
				list = append(list, key+"='"+value+"'")
			}
		}
		ret += text.Indent(strings.Join(list, ", ")+"\n", "  ")
	} else {
		ret += text.Indent("-- empty --\n", "  ")
	}
	ret += "actions:\n"
	for key, action := range d.actions {
		ret += text.Indent(key+":\n", "  ")
		ret += text.Indent(action.String(), "    ")
	}
	return
}

// get_actions_list_from returns the list of actions which defines the 'field_name' parameter.
func (o *ForjObject) get_actions_list_from(field_name string) (res map[string]*ForjObjectAction) {
	if o == nil {
		return nil
	}

	res = make(map[string]*ForjObjectAction)
	for action, action_data := range o.actions {
		if _, found := action_data.params[field_name]; found {
			res[action] = action_data
			gotrace.Trace("field '%s' found in action '%s'", field_name, action)
		} else {
			gotrace.Trace("field '%s' NOT found in action '%s'", field_name, action)
		}
	}
	return
}

// Do intersection between actions_related and a filtered list of actions.
// A warning is given if the actions_related become empty.
// This means, the list of fields to extract are not all found in at least one action.
func (l *ForjObjectList) inter_actions_list(filtered_list map[string]*ForjObjectAction) *ForjObjectList {
	if l == nil {
		return nil
	}

	if len(l.actions_related) == 0 {
		return nil
	}
	for action := range l.actions_related {
		if _, found := filtered_list[action]; !found {
			delete(l.actions_related, action)
			gotrace.Trace("action '%s' eliminated.", action)
		}
	}
	if len(l.actions_related) == 0 {
		l.obj.err = fmt.Errorf("Warning! List '%s' can not be applied to any object actions. ", l.name)
		return nil
	}
	return l
}

func (l *ForjObjectList) AddValidateHandler(valid_handler func(*ForjListData) error) *ForjObjectList {
	if l == nil {
		return nil
	}
	l.valid_handler = valid_handler
	return l
}

// Search for a flag/Arg from the list or additional param (object/list)
func (l *ForjObjectList) search_object_param(action, object, key, param_name string) (p ForjParam) {
	for _, param := range l.actions[action].params {
		if fl, pi, pn := param.fromList(); fl == nil {
			if fl.obj.name != object || pi != key || pn != param_name {
				continue
			}
			return param
		} else {
			if l.obj.name != object {
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
