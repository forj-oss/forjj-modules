package cli

import (
	"fmt"
	"github.com/forj-oss/forjj-modules/trace"
	"github.com/forj-oss/forjj-modules/cli/interface"
)

// ForjActionRef To define an action reference
type ForjAction struct {
	help          string                      // String which will 'printf' the object name as %s
	name          string                      // Action Name
	cmd           clier.CmdClauser            // Action used at action level
	params        map[string]ForjParam        // Collection of Arguments/Flags
	internal_only bool                        // True if this action cannot be enhanced by plugins
	to_refresh    map[string]*ForjContextTime // List of Object to refresh with context flags
}

// ForjContextTime. Structure used at context time to add more flags from Objectlist instances.
type ForjContextTime struct {
	objects_list *ForjObjectList   // List of Object list flags added - Used to add detailed flags
	action       *ForjObjectAction // Action to refresh with ObjectList detailed flags.
}

// AddFlagFromObjectListAction add one ObjectList action to the selected action.
func (c *ForjCli) AddFlagFromObjectListAction(obj_name, obj_list, obj_action string) *ForjCli {
	o_object, o_object_list, o_action, err := c.getObjectListAction(obj_list, obj_action)

	if err != nil {
		gotrace.Trace("Unable to find object '%s' action '%s'. Adding flags into selected actions ignored.",
			obj_list, obj_action)
		return c
	}

	for _, action := range c.sel_actions {
		d_flag := new(ForjFlagList)
		help := fmt.Sprintf("%s one or more %s", obj_action, o_object.desc)
		d_flag.set_cmd(action.cmd, String, obj_name, help, nil)
		action.params[obj_name+"s"] = d_flag

		// Need to add all others object fields not managed by the list, but At context time.
		action.to_refresh[obj_name] = &ForjContextTime{o_object_list, o_action}
	}
	return c
}

// AddFlagFromObjectListActions add one ObjectList action to the selected action.
func (c *ForjCli) AddFlagFromObjectListActions(obj_name, obj_list string, obj_actions ...string) *ForjCli {
	for _, obj_action := range obj_actions {
		o_object, o_object_list, o_action, _ := c.getObjectListAction(obj_list, obj_action)

		for _, action := range c.sel_actions {
			new_obj_name := action.name + "-" + obj_name
			d_flag := new(ForjFlagList)
			help := fmt.Sprintf("%s one or more %s", obj_action, o_object.desc)
			d_flag.set_cmd(action.cmd, String, new_obj_name, help, nil)
			action.params[new_obj_name+"s"] = d_flag

			// Need to add all others object fields not managed by the list, but At context time.
			action.to_refresh[obj_name] = &ForjContextTime{o_object_list, o_action}
		}
	}
	return c
}

// AddFlagsFromObjectAction create all flags defined on an object action
func (c *ForjCli) AddFlagsFromObjectAction(obj_name, obj_action string) *ForjCli {
	_, o_action, _ := c.getObjectAction(obj_name, obj_action)
	for _, action := range c.sel_actions {
		for _, param := range o_action.params {
			var fc ForjParamCopier
			fc = param

			d_flag := fc.CopyToFlag(action.cmd)
			action.params[obj_name+"s"] = d_flag
		}
	}
	return c
}

// AddArg Add an arg on selected actions
func (c *ForjCli) AddArg(value_type, name, help string, options *ForjOpts) *ForjCli {
	for _, action := range c.sel_actions {
		a := new(ForjArg)
		a.arg = c.App.Arg(name, help)
		a.set_options(options)

		switch value_type {
		case String:
			a.argv = a.arg.String()
		case Bool:
			a.argv = a.arg.Bool()
		}
		action.params[name] = a
	}
	return c
}

// AddFlag Add an flag on selected actions
func (c *ForjCli) AddFlag(value_type, name, help string, options *ForjOpts) *ForjCli {
	for _, action := range c.sel_actions {
		f := new(ForjFlag)
		f.flag = c.App.Flag(name, help)
		f.set_options(options)

		switch value_type {
		case String:
			f.flagv = f.flag.String()
		case Bool:
			f.flagv = f.flag.Bool()
		}
		action.params[name] = f
	}
	return c
}

// NewActions create the list of referenced valid actions supported. kingpin layer created.
// It add them to the kingpin application layer.
//
// name     : Name of the action to add
// help     : Generic help to add to the action.
// for_forjj: True if the action is protected against plugins features.
func (c *ForjCli) NewActions(name, act_help, compose_help string, for_forjj bool) (r *ForjAction) {
	r = new(ForjAction)
	r.cmd = c.App.Command(name, act_help)
	r.help = compose_help
	r.internal_only = for_forjj
	r.params = make(map[string]ForjParam)
	r.name = name
	c.actions[name] = r
	return
}

// OnActions Do a selection of action to apply more functionality
func (c *ForjCli) OnActions(actions ...string) *ForjCli {
	for _, action := range actions {
		if v, err := c.getAction(action); err == nil {
			c.sel_actions[action] = v
		}
	}
	return c
}
