package cli

import "github.com/forj-oss/forjj-modules/trace"

type ForjValues struct {
	name            string               // Name of the flag.
	app             *ForjFlag            // Flag defined in the App layer.
	actions         map[string]ForjParam // ForjParams defined at the action layer.
	objects_actions map[string]ForjParam // ForjParams defined at the action/object layer.
	object_lists    map[string]ForjParam // Flag found from one or more Objects list.
}

func (v *ForjValues) AtApp(flag *ForjFlag) {
	if v.app != nil && v.app != flag {
		gotrace.Trace("Probably a bug: a different flag object '%s' has been registered at Application layer.", v.name)
	}
	v.app = flag
}

func (v *ForjValues) AtAction(action string, p ForjParam) {
	if d, found := v.actions[action]; found && d != p {
		gotrace.Trace("Probably a bug: a different flag object '%s' has been registered at Action layer '%s'.", v.name, action)
	}
	v.actions[action] = p
}

func (v *ForjValues) AtObjectAction(objAction string, p ForjParam) {
	if d, found := v.objects_actions[objAction]; found && d != p {
		gotrace.Trace("Probably a bug: a different flag object '%s' has been registered at Object action layer '%s'.", v.name, objAction)
	}
	v.objects_actions[objAction] = p
}

func (v *ForjValues) AtObjectListAction(objList string, p ForjParam) {
	if d, found := v.object_lists[objList]; found && d != p {
		gotrace.Trace("Probably a bug: a different flag object '%s' has been registered at Object list action layer '%s'.", v.name, objList)
	}
	v.object_lists[objList] = p
}

// GetFrom, get the param value from the defined context.
// If no context exists, it tries to get from App Flag layer
// It search in action_object and then action.
// If the value context is a list, it moves to get it from the App layer directly.
func (v *ForjValues) GetFrom(cli *ForjCli) (ret interface{}, found bool) {
	defer func() {
		if !found {
			ret = v.app.GetValue()
			found = v.app.IsFound()
		}
	}()

	if cli.context.list != nil {
		// Context is on an object list. GetFrom can't work.
		return
	}

	if cli.context.object != nil {
		// Context is on an object action.
		action_object := cli.context.action.name + "_" + cli.context.object.name
		if value, ok := v.objects_actions[action_object]; ok {
			ret = value.GetValue()
			found = value.IsFound()
			return
		}
	}

	if cli.context.action != nil {
		if value, ok := v.actions[cli.context.action.name]; ok {
			ret = value.GetValue()
			found = value.IsFound()
			return
		}
	}
	return
}
