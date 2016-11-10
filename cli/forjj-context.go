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

// LoadContext gets data from context and store it in internal object model (ForjValue)
//
//
func (c *ForjCli) LoadContext(args []string) (cmds []clier.CmdClauser, err error) {

	var context clier.ParseContexter

	if v, err := c.App.ParseContext(args); err != nil {
		return cmds, err
	} else {
		context = v
	}

	cmds = context.SelectedCommands()
	if len(cmds) == 0 {
		return
	}

	// Determine selected Action/object/object list.
	c.identifyObjects(cmds[len(cmds)-1])

	// Load object list instances
	c.loadListData(nil, context, cmds[len(cmds)-1])

	// Then load additional flag/args from a function call - plugin
	// Then create additional flags/args
	//	c.loadAdditionalParams(cmds[len(cmds)-1], nil)

	// Then set Application flags/args values
	c.loadContextValue(cmds[0].FullCommand(), cmds[len(cmds)-1])
	return
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
				data.attrs[key] = value
			}
		}
		return nil
	}

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

			}
		}

		// get or create a record and populate it with all flags/args
		data := c.setObjectAttributes(c.context.action.name, o.name, key_value)
		for field_name := range o.fields {
			param := o.actions[c.context.action.name].params[field_name]
			if v, found := c.getContextValue(context, param.(forjParam)); found {
				data.attrs[field_name] = v
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

func (c *ForjCli) loadContextValue(action string, context clier.CmdClauser) {
	if c.context.action != nil {
		// Search all args/flags in the context
		/*		for param_name, param := range c.context.action.params {
				var value interface{}
				switch param.(type) {
				case *ForjArg: // Single value
					arg := param.(forjParam).GetArg()
					if v, found := c.context.GetArgValue(arg.arg); found {
						value = v
					}
				case *ForjFlag: // Single value
					flag := param.(forjParam).GetFlag()
					if v, found := c.context.GetFlagValue(flag.flag); found {
						value = v
					}
				case *ForjArgList: // Multiple Values
					arg := param.(forjParam).GetArg()
					if v, found := c.context.GetArgValue(arg.arg); found {
						value = v
					}
				case *ForjFlagList: // Multiple values
					flag := param.(forjParam).GetFlag()
					if v, found := c.context.GetFlagValue(flag.flag); found {
						value = v
					}
				}

			}*/
	}
}

// For Debug only.

func (c *ForjCli) IdentifyObjects(cmd clier.CmdClauser) {
	c.identifyObjects(cmd)
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
				// ex: forjj add repo
				c.context.object = object
				c.context.action = action.action
				return
			}
		}
	}

	for _, list := range c.list {
		for _, action := range list.actions {
			if action.cmd == cmd {
				// ex: forjj add repos
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
