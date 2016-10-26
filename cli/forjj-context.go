package cli

import (
	"github.com/forj-oss/forjj-modules/cli/interface"
)

type ForjCliContext struct {
	action *ForjAction     // Can be only one action
	object *ForjObject     // Can be only one object at a time. Ex: forj add repo
	list   *ForjObjectList // Can be only one list at a time.
	// forjj add apps ...
}

func (c *ForjCli) LoadContext(args []string) (cmds []clier.CmdClauser, err error) {

	var context clier.ParseContexter

	if v, err := c.App.GetContext(args); err != nil {
		return cmds, err
	} else {
		context = v
	}

	cmds = context.SelectedCommands()
	if len(cmds) == 0 {
		return
	}

	c.context.action = nil
	c.context.object = nil
	c.context.list = nil
	// Identify in Actions, in Objects, then in ObjectList
	for _, action := range c.actions {
		if action.cmd == cmds[0] {
			// ex: forjj =>create<=
			c.context.action = action
			return
		}
	}

	for _, object := range c.objects {
		for _, action := range object.actions {
			if action.cmd == cmds[0] {
				// ex: forjj add repo
				c.context.object = object
				c.context.action = action.action
				return
			}
		}
	}

	for _, list := range c.list {
		for _, action := range list.actions {
			if action.cmd == cmds[0] {
				// ex: forjj add repos
				c.context.action = action.action
				c.context.object = list.obj
				c.context.list = list
			}
		}
	}
	return
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
