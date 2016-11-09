package cli

import (
	"fmt"
	"github.com/kr/text"
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

func (a *ForjAction) String() string {
	ret := fmt.Sprintf("Action (%p):\n", a)
	ret += fmt.Sprintf("  name: '%s'\n", a.name)
	ret += fmt.Sprintf("  help: '%s'\n", a.help)
	ret += fmt.Sprintf("  internal_only: '%b'\n", a.internal_only)
	ret += fmt.Sprintf("  cmd: '%p'\n", a.cmd)
	ret += fmt.Sprintf("  params: %d\n", len(a.params))
	for key, param := range a.params {
		ret += fmt.Sprintf("    %s:\n", key)
		ret += text.Indent(param.String(), "      ")
	}
	return ret
}

// ForjContextTime. Structure used at context time to add more flags from Objectlist instances.
type ForjContextTime struct {
	objects_list *ForjObjectList   // List of Object list flags added - Used to add detailed flags
	action       *ForjObjectAction // Action to refresh with ObjectList detailed flags.
}

// AddActionFlagFromObjectListAction add one ObjectList action to the selected action.
//
// Ex:<app> update --tests "flag_key"
func (c *ForjCli) AddActionFlagFromObjectListAction(action_name, obj_name, obj_list, obj_action string) *ForjCli {
	o_object, o_object_list, o_action, err := c.getObjectListAction(obj_name+"_"+obj_list, obj_action)

	var action *ForjAction

	if a, found := c.actions[action_name]; !found {
		c.err = fmt.Errorf("Unable to find action '%s'. Adding object list action %s '%s-%s' as flag ignored.",
			action_name, obj_action, obj_name, obj_list)
		return nil
	} else {
		action = a
	}

	if err != nil {
		c.err = fmt.Errorf("Unable to find object '%s' action '%s'. %s. Adding flags into selected actions ignored.",
			obj_name+"_"+obj_list, obj_action, err)
		return nil
	}

	d_flag := new(ForjFlagList)

	new_object_name := obj_name + "s"
	d_flag.obj = o_object_list

	help := fmt.Sprintf("%s one or more %s", obj_action, o_object.desc)
	d_flag.set_cmd(action.cmd, String, new_object_name, help, nil)
	action.params[new_object_name] = d_flag

	// Need to add all others object fields not managed by the list, but At context time.
	action.to_refresh[obj_name] = &ForjContextTime{o_object_list, o_action}
	return c
}

// AddActionFlagsFromObjectListActions add one ObjectList action to the selected action.
// Ex: <app> update --add-tests "flag_key" --remove-tests "test,test2"
func (c *ForjCli) AddActionFlagsFromObjectListActions(action_name, obj_name, obj_list string, obj_actions ...string) *ForjCli {
	for _, obj_action := range obj_actions {
		o_object, o_object_list, o_action, err := c.getObjectListAction(obj_name+"_"+obj_list, obj_action)

		var action *ForjAction

		if a, found := c.actions[action_name]; !found {
			c.err = fmt.Errorf("Unable to find action '%s'. Adding object list action %s '%s-%s' as flag ignored.",
				action_name, obj_action, obj_name, obj_list)
			return nil
		} else {
			action = a
		}

		if err != nil {
			c.err = fmt.Errorf("Unable to find object '%s' action '%s'. %s. Adding flags into selected actions ignored.",
				obj_name+"_"+obj_list, obj_action, err)
			return nil
		}

		new_obj_name := obj_action + "-" + obj_name + "s"
		d_flag := new(ForjFlagList)
		d_flag.obj = o_object_list
		help := fmt.Sprintf("%s one or more %s", obj_action, o_object.desc)
		d_flag.set_cmd(action.cmd, String, new_obj_name, help, nil)
		action.params[new_obj_name] = d_flag

		// Need to add all others object fields not managed by the list, but At context time.
		action.to_refresh[obj_name] = &ForjContextTime{o_object_list, o_action}
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
			action.params[d_flag.name] = d_flag
		}
	}
	return c
}

// AddArg Add an arg on selected actions
func (c *ForjCli) AddArg(value_type, name, help string, options *ForjOpts) *ForjCli {
	return c.addFlag(func() ForjParam {
		return new(ForjArg)
	}, value_type, name, help, options)
}

// AddFlag Add an flag on selected actions
func (c *ForjCli) AddFlag(value_type, name, help string, options *ForjOpts) *ForjCli {
	return c.addFlag(func() ForjParam {
		return new(ForjFlag)
	}, value_type, name, help, options)
}

func (c *ForjCli) addFlag(newParam func() ForjParam, value_type, name, help string, options *ForjOpts) *ForjCli {
	if c == nil {
		return nil
	}
	for _, action := range c.sel_actions {
		p := newParam()

		p.set_cmd(action.cmd, value_type, name, help, options)

		action.params[name] = p
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
	r.to_refresh = make(map[string]*ForjContextTime)
	r.name = name
	c.actions[name] = r
	return
}

// OnActions Do a selection of action to apply more functionality
func (c *ForjCli) OnActions(actions ...string) *ForjCli {
	if len(actions) == 0 {
		c.sel_actions = c.actions
		return c
	}

	for _, action := range actions {
		if v, err := c.getAction(action); err == nil {
			c.sel_actions[action] = v
		}
	}
	return c
}
