package cli

import (
	"github.com/alecthomas/kingpin"
)

func (c *ForjCli) SetCmdContext(context *kingpin.ParseContext) (cmd *kingpin.CmdClause) {
	cmd = context.SelectedCommand
	if cmd == nil {
		return
	}

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
	return
}

// LoadValuesFrom load most of flags/arguments found in the cli context in values, like kingpin.execute do.
func (c *ForjCli) LoadValuesFrom(context *kingpin.ParseContext) {
	c.loadListValuesFrom(context)
	c.loadObjectValuesFrom(context)
	c.loadActionValuesFrom(context)
	c.loadAppValuesFrom(context)
}

func (c *ForjCli) loadListValuesFrom(context *kingpin.ParseContext) {
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

func (c *ForjCli) loadObjectValuesFrom(context *kingpin.ParseContext) {
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

func (c *ForjCli) loadActionValuesFrom(context *kingpin.ParseContext) {
	if c.context.action == nil {
		return
	}
	for _, param := range c.context.action.params {
		param.loadFrom(context)
	}
}

func (c *ForjCli) loadAppValuesFrom(context *kingpin.ParseContext) {
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
