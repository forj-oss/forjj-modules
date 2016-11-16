package cli

import (
	"fmt"
	"github.com/forj-oss/forjj-modules/cli/interface"
	"github.com/forj-oss/forjj-modules/trace"
)

type ForjCliContext struct {
	action *ForjAction     // Can be only one action
	object *ForjObject     // Can be only one object at a time. Ex: forj add repo
	list   *ForjObjectList // Can be only one list at a time.
	// forjj add apps ...
}

func (c *ForjCli) LoadContext(args []string, context interface{}) ([]clier.CmdClauser, error) {
	return c.loadContext(args, context)
}

// LoadContext gets data from context and store it in internal object model (ForjValue)
//
//
func (c *ForjCli) loadContext(args []string, context interface{}) (cmds []clier.CmdClauser, err error) {

	var cli_context clier.ParseContexter

	if v, err := c.App.ParseContext(args); err != nil {
		return cmds, err
	} else {
		cli_context = v
	}

	cmds = cli_context.SelectedCommands()
	if len(cmds) == 0 {
		err = c.contextHook(context)
		return
	}

	// Determine selected Action/object/object list.
	c.identifyObjects(cmds[len(cmds)-1])

	// Load object list instances
	c.loadListData(nil, cli_context, cmds[len(cmds)-1])

	// Load anything that could be required from any existing flags setup.
	// Ex: app driver - app object hook. - Add new flags/args/objects
	//     Settings of Defaults, flags attributes - Application hook. - Update existing flags.
	if err = c.contextHook(context); err != nil {
		return
	}

	// Define instance flags for each list.
	c.addInstanceFlags()

	return
}

// ContextHook
// Load anything that could be required from any existing flags setup.
// Ex: app driver - app object hook. - Add new flags/args/objects
//     Settings of Defaults, flags attributes - Application hook. - Update existing flags.
func (c *ForjCli) contextHook(context interface{}) error {
	for _, object := range c.objects {
		if object.context_hook == nil {
			continue
		}

		if err := object.context_hook(object, c, context); err != nil {
			c.err = err
			return nil
		}
	}

	if c.context_hook == nil {
		return nil
	}

	if err := c.context_hook(c, context); err != nil {
		return err
	}
	return nil
}

func (c *ForjCli) LoadListData(more_flags func(*ForjCli), context clier.ParseContexter, Cmd clier.CmdClauser) error {
	return c.loadListData(more_flags, context, Cmd)
}

// check List flag and start creating object instance.
func (c *ForjCli) loadListData(more_flags func(*ForjCli), context clier.ParseContexter, Cmd clier.CmdClauser) error {
	// check if the ObjectList is found.
	// Ex: forjj create repos <list>
	if c.context.list != nil {
		gotrace.Trace("Loading Data list from an Object list.")
		l := c.context.list

		key_name := l.obj.getKeyName()
		// loop on list data to create object records.
		for _, attrs := range l.list {
			// Get the list element key
			key_value := attrs.Data[key_name]
			if key_value == "" {
				return fmt.Errorf("Invalid key value for object list '%s-%s'. a key cannot be empty.",
					l.obj.name, l.name)
			}

			data := c.setObjectAttributes(c.context.action.name, l.obj.name, key_value)
			for key, value := range attrs.Data {
				field := l.obj.fields[key]
				if err := data.set(field.value_type, key, value); err != nil {
					return err
				}
				data.attrs[key] = value
			}
		}
		return nil
	}

	// Check if the Object is found
	// Ex: forjj create repo <list> # with any additional object lists flags.
	if c.context.object != nil {
		o := c.context.object
		gotrace.Trace("Loading Data list from the object '%s'.", o.name)
		var key_value string

		key_name := o.getKeyName()
		param := o.actions[c.context.action.name].params[key_name]
		if param == nil {
			return fmt.Errorf("Unable to find key '%s' in object action '%s-%s' parameters.",
				key_name, o.name, c.context.action.name)
		}
		if v, found := c.getContextValue(context, param.(forjParam)); !found {
			return fmt.Errorf("Unable to find key '%s' value from action '%s' parameters. "+
				"Missing OnActions().AddFlag(%s)?", key_name, c.context.action.name, key_name)
		} else {
			key_value = v
		}
		if key_value == "" {
			return fmt.Errorf("Invalid key value for object '%s'. a key cannot be empty.", o.name)
		}
		gotrace.Trace("New object record identified by key '%s' (%s).", key_value, o.getKeyName())

		// Search for object list flags
		for _, param := range o.actions[c.context.action.name].params {
			switch param.(type) {
			case *ForjFlagList:
				fl := param.(*ForjFlagList)
				key_name := fl.obj.obj.getKeyName()
				for _, list_data := range fl.obj.list {
					key_value := list_data.Data[key_name]
					data := c.setObjectAttributes(c.context.action.name, fl.obj.obj.name, key_value)
					if data == nil {
						return c.err
					}
					for key, attr := range list_data.Data {
						field := fl.obj.obj.fields[key]
						if err := data.set(field.value_type, key, attr); err != nil {
							return err
						}
					}
				}
			case *ForjArgList:
				al := param.(*ForjArgList)
				key_name := al.obj.obj.getKeyName()
				for _, list_data := range al.obj.list {
					key_value := list_data.Data[key_name]
					data := c.setObjectAttributes(c.context.action.name, al.obj.obj.name, key_value)
					for key, attr := range list_data.Data {
						field := al.obj.obj.fields[key]
						if err := data.set(field.value_type, key, attr); err != nil {
							return err
						}
					}
				}
			}
		}

		// get or create a record and populate it with all flags/args
		data := c.setObjectAttributes(c.context.action.name, o.name, key_value)
		for field_name, field := range o.fields {
			param := o.actions[c.context.action.name].params[field_name]
			if v, found := c.getContextValue(context, param.(forjParam)); found {
				if err := data.set(field.value_type, field_name, v); err != nil {
					return err
				}

			}
		}
		return nil
	}

	// Parse flags to determine if there is another objects list
	gotrace.Trace("Loading Data list from an action flag/arg.")
	for _, param := range c.context.action.params {
		switch param.(type) {
		case *ForjFlagList:
			fl := param.(*ForjFlagList)
			key_name := fl.obj.obj.getKeyName()
			for _, list_data := range fl.obj.list {
				key_value := list_data.Data[key_name]
				data := c.setObjectAttributes(c.context.action.name, fl.obj.obj.name, key_value)
				if data == nil {
					return c.err
				}
				for key, attr := range list_data.Data {
					data.attrs[key] = attr
				}
			}
		case *ForjArgList:
			al := param.(*ForjArgList)
			key_name := al.obj.obj.getKeyName()
			for _, list_data := range al.obj.list {
				key_value := list_data.Data[key_name]
				data := c.setObjectAttributes(c.context.action.name, al.obj.obj.name, key_value)
				for key, attr := range list_data.Data {
					data.attrs[key] = attr
				}
			}
		}
	}
	return nil
}

func (c *ForjCli) getContextValue(context clier.ParseContexter, param forjParam) (string, bool) {
	switch param.(type) {
	case *ForjArg:
		a := param.(*ForjArg)
		return context.GetArgValue(a.arg)
	case *ForjFlag:
		f := param.(*ForjFlag)
		return context.GetFlagValue(f.flag)
	}
	return "", false
}

func (c *ForjCli) addInstanceFlags() {
	for _, l := range c.list {
		if _, found := c.values[l.obj.name]; !found {
			continue
		}
		r := c.values[l.obj.name]
		if len(r.records) == 0 {
			continue
		}
		for instance_name := range r.records {
			for field_name, field := range l.obj.fields {
				found := false
				// Do not include fields defined by the list.
				for _, fname := range l.fields_name {
					if fname == field_name {
						found = true
						break
					}
				}
				if found {
					continue
				}
				// Add instance flags to `<app> <action> <object>s --...`
				flag_name := instance_name + "-" + field_name
				for _, action := range l.actions {
					// Do not recreate if already exist.
					if _, found := action.params[flag_name]; found {
						continue
					}

					f := new(ForjFlag)
					p := ForjParam(f)
					p.set_cmd(action.cmd, field.value_type, flag_name, field.help+" for "+instance_name, nil)
					f.list = l
					f.instance_name = instance_name
					f.field_name = field_name
					action.params[flag_name] = p
				}

				// Add instance flags to any object list flags added to actions or other objects.
				// defined by Add*Flag*FromObjectListAction* like functions
				for _, flag_list := range l.flags_list {
					// Do not recreate if already exist.
					if _, found := flag_list.params[flag_name]; found {
						continue
					}

					switch {
					case flag_list.action != nil:
						f := new(ForjFlag)
						f.list = l
						f.instance_name = instance_name
						f.field_name = field_name
						p := ForjParam(f)
						p.set_cmd(flag_list.action.cmd, field.value_type, flag_name, field.help+" for "+instance_name, nil)
						flag_list.action.params[flag_name] = p
						flag_list.params[flag_name] = p
					case flag_list.objectAction != nil:
						f := new(ForjFlag)
						f.list = l
						f.instance_name = instance_name
						f.field_name = field_name
						p := ForjParam(f)
						p.set_cmd(flag_list.objectAction.cmd, field.value_type, flag_name, field.help+" for "+instance_name, nil)
						flag_list.objectAction.params[flag_name] = p
						flag_list.params[flag_name] = p
					}
				}
			}
		}
	}
}

// loadObjectData is executed at final Parse task
// It loads Object data from any other object/instance flags
// and update the cli object data fields list
func (c *ForjCli) loadObjectData() {
	var params map[string]ForjParam
	switch {
	case c.context.list != nil: // <app> <action> <object>s
		l := c.context.list
		params = l.actions[c.context.action.name].params
	case c.context.object != nil: // <app> <action> <object>
		o := c.context.object
		params = o.actions[c.context.action.name].params
	case c.context.action != nil: // <app> <action>
		a := c.context.action
		params = a.params
	}
	for _, param := range params {
		if p, ok := param.(forjParamObject); ok {
			p.UpdateObject()
		}
	}
}

func (c *ForjCli) identifyObjects(cmd clier.CmdClauser) {
	c.context.action = nil
	c.context.object = nil
	c.context.list = nil
	// Identify in Actions, in Objects, then in ObjectList
	for _, action := range c.actions {
		if action.cmd == cmd {
			// ex: forjj =>create<=
			c.context.action = action
			return
		}
	}

	for _, object := range c.objects {
		for _, action := range object.actions {
			if action.cmd == cmd {
				// ex: forjj add =>repo<=
				c.context.object = object
				c.context.action = action.action
				return
			}
		}
	}

	for _, list := range c.list {
		for _, action := range list.actions {
			if action.cmd == cmd {
				// ex: forjj add =>repos<=
				c.context.action = action.action
				c.context.object = list.obj
				c.context.list = list
			}
		}
	}
}

// LoadValuesFrom load most of flags/arguments found in the cli context in values, like kingpin.execute do.
func (c *ForjCli) LoadValuesFrom(context clier.ParseContexter) {
	c.loadListValuesFrom(context)
	c.loadObjectValuesFrom(context)
	c.loadActionValuesFrom(context)
	c.loadAppValuesFrom(context)
}

func (c *ForjCli) loadListValuesFrom(context clier.ParseContexter) {
	if c.context.list == nil {
		return
	}
	for _, action := range c.context.list.actions {
		if action.action == c.context.action {
			for _, param := range action.params {
				param.loadFrom(context)
			}
		}
	}
}

func (c *ForjCli) loadObjectValuesFrom(context clier.ParseContexter) {
	if c.context.object == nil {
		return
	}
	for _, action := range c.context.object.actions {
		if action.action == c.context.action {
			for _, param := range action.params {
				param.loadFrom(context)
			}
		}
	}
}

func (c *ForjCli) loadActionValuesFrom(context clier.ParseContexter) {
	if c.context.action == nil {
		return
	}
	for _, param := range c.context.action.params {
		param.loadFrom(context)
	}
}

func (c *ForjCli) loadAppValuesFrom(context clier.ParseContexter) {
	for _, flag := range c.flags {
		flag.loadFrom(context)
	}
}

// GetStringValueAddr : Provide the Address to the flag/argument value if is a *string.
// It returns nil if the flag/arg was not found or is not a string.
func (c *ForjCli) GetStringValueAddr(name string) *string {
	if v, found := c.flags[name]; found {
		if is_string(v.flagv) {
			return v.flagv.(*string)
		}
		return nil
	}
	return nil
}
