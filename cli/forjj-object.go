package cli

import (
	"bytes"
	"fmt"
	"github.com/kr/text"
	"github.com/forj-oss/forjj-modules/cli/interface"
	"github.com/forj-oss/forjj-modules/trace"
	"log"
	"regexp"
	"strings"
	"text/template"
	"unicode"
)

const no_fields = "none"

// ForjObject defines the Object structure
type ForjObject struct {
	cli           *ForjCli                                       // Reference to the parent
	name          string                                         // name of the action to add for objects
	desc          string                                         // Object description string.
	actions       map[string]*ForjObjectAction                   // Collection of actions per objects where object cmd flags are added.
	list          map[string]*ForjObjectList                     // List configured for this object.
	internal      bool                                           // true if the object is forjj internal
	sel_actions   map[string]*ForjObjectAction                   // Select several actions to apply for AddParam
	fields        map[string]*ForjField                          // List of fields of this object
	instances     map[string]*ForjObjectInstance                 // Instance detected at Context time.
	instance_name string                                         // Instance name for a uniq record.
	err           error                                          // Last error found.
	context_hook  func(*ForjObject, *ForjCli, interface{}) error // Parse hook related to this object. Can use cli to create more.
}

// createObjectDataFromParams creates object data from the given list of params
func (o *ForjObject) createObjectDataFromParams(params map[string]ForjParam) error {
	if o.instance_name != "" {
		return nil
	}
	if o.setInstanceNameFromParams(params) == "" && o.err != nil {
		return o.err
	}
	obj_data := o.cli.setObjectAttributes(o.cli.cli_context.action.name, o.name, o.instance_name)
	key_name := o.getKeyName()
	obj_data.set(o.fields[key_name].value_type, key_name, o.instance_name)
	for _, p := range params {
		if p.Name() == key_name {
			p.forjParamUpdater().set_ref(obj_data)
			// Found it
			continue
		}
		if !p.isObjectRelated() {
			continue
		}
		if p.forjParamRelated().getObjectAction().obj != o {
			continue
		}
		v, _ := p.GetContextValue(o.cli.cli_context.context)
		field_name := p.forjParamRelated().getFieldName()
		obj_data.set(o.fields[field_name].value_type, field_name, v)
		p.forjParamUpdater().set_ref(obj_data)

	}
	return nil
}

func (o *ForjObject) setInstanceNameFromParams(params map[string]ForjParam) string {
	if o.cli.cli_context.context == nil {
		o.err = fmt.Errorf("Internal error! Context object is missing")
		return ""
	}
	key_name := o.getKeyName()
	// Search for key value to create the object
	for _, p := range params {
		if p.Name() != key_name {
			// Found it
			continue
		}
		if v, found := p.GetContextValue(o.cli.cli_context.context); !p.IsList() && found {
			o.instance_name = to_string(v)
		}
	}
	return ""
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
	ret += fmt.Sprintf("  fields: %d\n", len(o.fields))
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

	found     bool                 // True if the flag was used.
	plugins   []string             // List of plugins that use this flag.
	inActions map[string]ForjParam // Collection of flags linked to Main actions. From
	// AddActionFlagsFromObjectAction
	regexp string // Regexp to validate input.
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
//
// Ex: forjj create =>repo --flags value ...<=
//   where repo is a cmd and params store all object flags/args
type ForjObjectAction struct {
	name    string               // object action name (formatted as <action>_<object>)
	cmd     clier.CmdClauser     // Object
	action  *ForjAction          // Parent Action name and help
	plugins []string             // Plugins implementing this object action.
	params  map[string]ForjParam // Collection of flags
	obj     *ForjObject          // Object referenced
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

	return o.addParam(func() ForjParam {
		return new(ForjFlag)
	}, name, options)
}

// SetParamOptions update flag/arg options anywhere param_name has been defined, except flag/arg list.
//
func (o *ForjObject) SetParamOptions(param_name string, options *ForjOpts) {
	for _, action := range o.actions {
		if p, found := action.params[param_name]; found {
			p.set_options(options)
			p.forjParamUpdater().updateContextData()
		}
	}
	for _, list := range o.list {
		for _, flag_list := range list.flags_list {
			for _, param := range flag_list.params {
				if param.forjParamRelated().getFieldName() == param_name {
					param.set_options(options)
					param.forjParamUpdater().updateContextData()
				}
			}
		}
	}
	for _, param := range o.fields[param_name].inActions {
		param.set_options(options)
		param.forjParamUpdater().updateContextData()
	}
}

func (o *ForjObject) AddArg(name string, options *ForjOpts) *ForjObject {
	if o == nil {
		return nil
	}
	return o.addParam(func() ForjParam {
		return new(ForjArg)
	}, name, options)
}

func (o *ForjObject) addParam(newParam func() ForjParam, name string, options *ForjOpts) *ForjObject {
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
		p.forjParamRelatedSetter().setObject(action, field.name)

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
			o.actions[action] = newForjObjectAction(ar, o, o.name, o.desc)
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

	if o.AddField(String, no_fields, "help", "") == nil {
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
func (o *ForjObject) AddKey(pIntType, name, help, re string) *ForjObject {
	if o == nil {
		return nil
	}

	for field_name, field := range o.fields {
		if field.key {
			o.err = fmt.Errorf("One key already exist in the object '%s', called '%s'", o.name, field_name)
			return nil
		}
	}

	if o.AddField(pIntType, name, help, re) == nil {
		return nil
	}

	field := o.fields[name]
	field.key = true
	return o
}

// AddField add a field to the object.
func (o *ForjObject) AddField(pIntType, name, help, re string) *ForjObject {
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

	if re == "" {
		gotrace.Trace("Warning. Field '%s' was configured with a regexp. Defaulting to '.*'", name)
		re = ".*"
	}
	o.fields[name] = &ForjField{
		name:       name,
		help:       help,
		value_type: pIntType,
		inActions:  make(map[string]ForjParam),
		regexp:     re,
	}
	return o
}

// buildListRegExp convert a human readable to Regexp
// [] are considered as optional and replaced by ()?
// Any word string are identified as a field are replaced by the template object field associated RegExp.
// Ex: name
// field index is calculated from a field detected and the list of []
func (o *ForjObject) buildListRegExp(sample string, l *ForjObjectList) (ret string, err error) {
	ret = sample
	l.sample = sample

	if !hasValidSquare(sample) {
		err = fmt.Errorf("Invalid syntax. square delimiter error in %s.", sample)
		return
	}

	// identify fields and their position
	fs := splitSepAndFields(sample, func(c rune) bool {
		return c == '['
	}, field_detect)
	l.max_fields = uint(len(fs)) + 1 // Number of Regexp matches.
	for i, value := range fs {

		if value != "[" {
			if _, found := o.fields[value]; found {
				if l.field(uint(i+1), value) == nil {
					return "", o.Error()
				}
			} else {
				return "", fmt.Errorf("'%s' is not a valid object field.", value)
			}
		}
	}

	sample, err = buildFromSepAndFields(sample, sep_detect, field_detect, regexpTmplReplacer)

	var t *template.Template
	t, err = template.New("regexp").Parse(sample)
	if err != nil {
		gotrace.Trace("'%s' is not a valid regexp template. %s. Ignored.", err)
		return
	}
	fields_data := make(map[string]string)
	for key, field := range o.fields {
		fields_data[key] = field.regexp
	}

	buf := bytes.NewBufferString("")
	err = t.Execute(buf, fields_data)
	if err != nil {
		gotrace.Trace("Unable to set regexp correctly. %s. Ignored.", err)
		return
	}
	ret = buf.String()

	return
}

//
func regexpTmplReplacer(s string) (string, error) {
	switch s {
	case "[":
		return "(", nil
	case "]":
		return ")?", nil
	default:
		return "{{ ." + s + " }}", nil
	}
}

func sep_detect(c rune) bool {
	return c == '[' || c == ']'
}

func field_detect(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsNumber(c) || c == '_' || c == '-'
}

func hasValidSquare(sample string) (isValid bool) {
	f := func(c rune) bool {
		return c != '[' && c != ']'
	}
	clean_sample := strings.Replace(strings.Replace(sample, `\[`, "", -1), `\]`, "", -1)
	a := splitSep(clean_sample, f)
	i := 0
	for _, value := range a {
		switch {
		case i < 0:
			return
		case value == "[":
			i++
			continue
		case value == "]":
			i--
		}
	}
	if i != 0 {
		return
	}
	isValid = true
	return
}

func splitSep(s string, f func(rune) bool) []string {
	n := 0
	for _, rune := range s {
		if f(rune) {
			n++
		}
	}

	// Now create them.
	a := make([]string, n)
	na := 0
	for i, rune := range s {
		if f(rune) {
			a[na] = s[i:i]
			na++
		}
	}
	return a
}

func buildFromSepAndFields(s string, sep, field func(rune) bool, replacer func(s string) (string, error)) (string, error) {
	result := ""
	fieldStart := -1    // Set to -1 when looking for start of field.
	nonFieldStart := -1 // Set to -1 when looking for start of non field or sep.
	escaped := false
	len_s := len(s)
	for i, rune := range s {
		isSep := sep(rune)
		isField := field(rune)

		if !isSep && !isField {
			if nonFieldStart == -1 {
				nonFieldStart = i
			}
			if fieldStart >= 0 {
				if res, err := replacer(s[fieldStart:i]); err != nil {
					return "", err
				} else {
					result += res
				}
				fieldStart = -1
			}
			continue
		}

		// Make \ as escape for mainly [] and \. In any other case, a single \ is ignored.
		if rune == '\\' {
			if !escaped {
				if fieldStart >= 0 {
					if res, err := replacer(s[fieldStart:i]); err != nil {
						return "", err
					} else {
						result += res
					}
					fieldStart = -1
				}
				if nonFieldStart >= 0 {
					if res, err := replacer(s[nonFieldStart:i]); err != nil {
						return "", err
					} else {
						result += res
					}
					nonFieldStart = -1
				}
			}
			escaped = !escaped
			continue
		}

		if isSep && !escaped {
			// Is a recognized separator
			if fieldStart >= 0 {
				if res, err := replacer(s[fieldStart:i]); err != nil {
					return "", err
				} else {
					result += res
				}
				fieldStart = -1
			}
			if nonFieldStart >= 0 {
				result += s[nonFieldStart:i]
				nonFieldStart = -1
			}

			sep_found := ""
			if len_s == i+1 {
				sep_found = s[i:]
			} else {
				sep_found = s[i : i+1]
			}
			if res, err := replacer(sep_found); err != nil {
				return "", err
			} else {
				result += res
			}
			continue
		}

		if isField {
			// is a Field
			if fieldStart == -1 {
				fieldStart = i
			}
			if nonFieldStart >= 0 {
				result += s[nonFieldStart:i]
				nonFieldStart = -1
			}
			continue
		}
		// Is not a field
		if nonFieldStart == -1 {
			nonFieldStart = i
		}
		if fieldStart >= 0 {
			if res, err := replacer(s[fieldStart:i]); err != nil {
				return "", err
			} else {
				result += res
			}
			fieldStart = -1
		}

	}

	if fieldStart >= 0 {
		// Last field might end at EOF.
		if res, err := replacer(s[fieldStart:]); err != nil {
			return "", err
		} else {
			result += res
		}
	}
	if nonFieldStart >= 0 {
		// Last non field might end at EOF.
		if res, err := replacer(s[nonFieldStart:]); err != nil {
			return "", err
		} else {
			result += res
		}
	}
	return result, nil
}

// splitSepAndFields return a slice of string identifying fields and separators.
// sep func(rune) bool return true if the rune is a single separator
// same for field. If sep and field return both true on a rune, sep is chosen
// If a Sep is prefixed by a single \, the sep won't be considered as a separator.
func splitSepAndFields(s string, sep, field func(rune) bool) []string {
	n := 0
	inField := false
	wasInField := false
	for _, rune := range s {
		if sep(rune) {
			n++
			wasInField = false
			continue
		}
		inField = field(rune)
		if inField && !wasInField {
			n++
		}
		wasInField = inField
	}

	// Now create them.
	a := make([]string, n)
	na := 0
	fieldStart := -1 // Set to -1 when looking for start of field.
	escaped := false
	len_s := len(s)
	for i, rune := range s {
		if rune == '\\' {
			escaped = !escaped
			continue
		}
		if sep(rune) && !escaped {
			if fieldStart >= 0 {
				a[na] = s[fieldStart:i]
				na++
				fieldStart = -1
			}
			if len_s == i+1 {
				a[na] = s[i:]
			} else {
				a[na] = s[i : i+1]
			}
			na++
			continue
		}
		if field(rune) {
			if fieldStart == -1 {
				fieldStart = i
			}
			continue
		}
		if fieldStart >= 0 {
			a[na] = s[fieldStart:i]
			na++
			fieldStart = -1
		}
	}
	if fieldStart >= 0 {
		// Last field might end at EOF.
		a[na] = s[fieldStart:]
	}
	return a
}

// CreateList create a new list. It returns the ForjObjectList to set it and configure actions
func (o *ForjObject) CreateList(name, list_sep, ext_regexp, help string) *ForjObjectList {
	if o == nil {
		return nil
	}

	l := new(ForjObjectList)
	l.obj = o
	l.fields_name = make(map[uint]string)
	l.name = name
	l.help = help
	l.sep = list_sep
	l.key_name = o.keyName()
	l.actions_related = make(map[string]*ForjObjectAction)
	for k, v := range o.actions {
		l.actions_related[k] = v
	}
	l.actions = make(map[string]*ForjObjectAction)
	l.list = make([]ForjListData, 0, 5)
	l.context = make([]ForjListData, 0, 5)
	l.data = make([]ForjData, 0, 5)
	l.flags_list = make(map[string]*ForjObjectListFlags)
	l.c = o.cli

	if r, err := o.buildListRegExp(ext_regexp, l); err != nil {
		o.err = err
		return nil
	} else {
		ext_regexp = r
	}

	ext_regexp = o.cli.buildCapture(ext_regexp)
	if r, err := regexp.Compile(ext_regexp); err != nil {
		o.err = fmt.Errorf("%s_%s not created: Regexp error found: %s", o, name, err)
		return nil
	} else {
		l.ext_regexp = r
		gotrace.Trace("Found '%d' group in '%s' (sample: %s)", l.max_fields-1, ext_regexp, l.sample)
	}

	// registering list
	l.obj.list[name] = l
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

	o_dest, o_action, _ := o.cli.getObjectAction(obj_name, obj_action)
	for _, action := range o.sel_actions {
		for fname := range o_dest.fields {
			if p, found := o_action.params[fname]; found {
				d_flag := p.Copier().CopyToFlag(action.cmd)
				d_flag.field_name = fname
				d_flag.obj = o_action
				action.params[fname] = d_flag
			}
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

func (o *ForjObject) Name() string {
	if o == nil {
		return ""
	}
	return o.name
}
