package cli

const internal_app = "_app"

// loadAppData Load application context data flag in forj_values
//
func (c *ForjCli)loadAppData() {
	action_name := "none"
	if c.cli_context.action != nil {
		action_name = c.cli_context.action.name
	}
	obj_data := c.setObjectAttributes(action_name, internal_app, c.App.Name())

	for name, flag := range c.flags {
		v, _ := flag.GetContextValue(c.cli_context.context)

		obj_data.set(flag.Type(), name, v)
	}

	if c.cli_context.action == nil {
		return
	}
	for name, param := range c.cli_context.action.params {
		if param.isListRelated() || param.isObjectRelated() || param.IsList() {
			continue
		}
		v, _ := param.GetContextValue(c.cli_context.context)

		obj_data.set(param.Type(), name, v)
	}
}

