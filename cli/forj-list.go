package cli

import (
	"fmt"
	"github.com/forj-oss/forjj-modules/trace"
	"regexp"
	"strings"
)

type ForjObjectList struct {
	c               *ForjCli                     // Reference to the cli object
	name            string                       // List name
	obj             *ForjObject                  // Object attached
	sep             string                       // List separator
	ext_regexp      *regexp.Regexp               // Capturing Regexp
	fields_name     map[uint]string              // Data fields extraction
	actions_related map[string]*ForjObjectAction // Possible actions for this list
	actions         map[string]*ForjObjectAction // Collection of actions per objects where flags are added.
	list            []ForjData                   // Data collected from the list.
	found           bool                         // True if the list flag was provided.
	key_name        string                       // List key name to use for any detailed flags.
}

type ForjData struct {
	data map[string]string
}

// CreateObjectList create an object list description
func (c *ForjCli) CreateObjectList(obj, name, list_sep, ext_regexp, key_name string) *ForjObjectList {
	if v, found := c.list[obj+"_"+name]; found {
		gotrace.Trace("%s_%s already exist. Not updated.", obj, name)
		return v
	}

	var o *ForjObject

	if v, found := c.objects[obj]; found {
		o = v
	} else {
		gotrace.Trace("%s_%s not created: Object '%s' not found: %s", obj, name, obj)
	}

	l := new(ForjObjectList)
	if r, err := regexp.Compile(ext_regexp); err != nil {
		gotrace.Trace("%s_%s not created: Regexp error found: %s", obj, name, err)
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
	l.c = c
	c.list[obj+"_"+name] = l
	return l
}

func (l *ForjObjectList) AddActions(actions ...string) {
	for _, action := range actions {
		if v, found := l.actions_related[action]; found {
			l.actions[action] = v
			for name, param := range v.params {
				l.c.addTracked(name).AtObjectListAction(l.name, param)
			}
		}
	}
}

// Field add a Map RegExp result to an object field parameter.
func (l *ForjObjectList) Field(index uint, field_name string) *ForjObjectList {
	l.inter_actions_list(l.obj.get_actions_list_from(field_name))
	l.fields_name[index] = field_name
	return l
}

// Set function for kingpin.Value interface
func (l *ForjObjectList) Set(value string) error {
	for _, v := range Split(" *"+l.sep+" *", value, l.sep) {
		if err := l.Add(v); err != nil {
			return err
		}
	}
	return nil
}

func (l *ForjObjectList) Add(value string) error {
	res := l.ext_regexp.FindStringSubmatch(value)
	if res == nil {
		return fmt.Errorf("%s is an invalid application driver. APP must be formated as '<type>:<DriverName>[:<InstanceName>]' all lower case. if <Name> is missed, <Name> will be set to <app>", value)
	}

	dd := ForjData{make(map[string]string)}

	for index, field_name := range l.fields_name {
		dd.data[field_name] = res[index]
	}

	l.list = append(l.list, dd)
	gotrace.Trace("'%s'(%s) added '%s'", l.obj.name, l.name, value)
	return nil
}

// FIXME: kingpin is having trouble in the context case, where several --<object>s set, with some flags in between, is ignoring seconds and next --apps flags values. Workaround is to have them all followed or use the --apps APP[,APP ...] format.

func (d *ForjObjectList) IsCumulative() bool {
	return true
}

// String : Set function for kingpin.Value interface
func (d *ForjObjectList) String() string {
	list := make([]string, 0, 2)

	for _, v := range d.list {
		for key, value := range v.data {
			list = append(list, key+"='"+value+"'")
		}
	}
	return strings.Join(list, ", ")
}

// get_actions_list_from returns the list of actions which defines the 'field_name' parameter.
func (o *ForjObject) get_actions_list_from(field_name string) (res map[string]*ForjObjectAction) {
	res = make(map[string]*ForjObjectAction)
	for action, action_data := range o.actions {
		if _, found := action_data.params[field_name]; found {
			res[action] = action_data
			gotrace.Trace("field '%s' found in action '%s'", field_name, action)
		}
	}
	return
}

// Do intersection between actions_related and a filtered list of actions.
// A warning is given if the actions_related become empty.
// This means, the list of fields to extract are not all found in at least one action.
func (l *ForjObjectList) inter_actions_list(filtered_list map[string]*ForjObjectAction) {
	if len(l.actions_related) == 0 {
		return
	}
	for action := range l.actions_related {
		if _, found := filtered_list[action]; !found {
			delete(l.actions_related, action)
			gotrace.Trace("action '%s' eliminated.", action)
		}
	}
	if len(l.actions_related) == 0 {
		gotrace.Trace("Warning! List '%s' can not be applied to any object actions.", l.name)
	}
}
