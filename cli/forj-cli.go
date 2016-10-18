package cli

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/forj-oss/goforjj"
	"github.com/forj-oss/forjj-modules/trace"
	"log"
	"strings"
)

// ForjCli is the Core cli for forjj command.
type ForjCli struct {
	App     *kingpin.Application       // Kingpin Application object
	flags   map[string]*ForjFlag       // Collection of Objects at Application level
	objects map[string]*ForjObject     // Collection of Objects that forjj will manage.
	actions map[string]*ForjAction     // Collection recognized actions
	list    map[string]*ForjObjectList // Collection of object list
	context ForjCliContext             // Context from cli parsing
}

type ForjCliContext struct {
	action *ForjAction
	object *ForjObject
	list   *ForjObjectList
}

// ForjActionRef To define an action reference
type ForjAction struct {
	help          string               // String which will 'printf' the object name as %s
	name          string               // Action Name
	cmd           *kingpin.CmdClause   // Action used at action level
	params        map[string]ForjParam // Collection of Arguments/Flags
	values        []ForjValues         //
	internal_only bool                 // True if this action cannot be enhanced by plugins
}

// ForjObject defines the Object structure
type ForjObject struct {
	name     string                       // name of the action to add for objects
	help     string                       // Generic Action help string.
	actions  map[string]*ForjObjectAction // Collection of actions per objects where flags are added.
	list     *ForjObjectList              // List configured for this object.
	internal bool                         // true if the object is forjj internal
}

// ForjObjectAction defines the action structure for each object
type ForjObjectAction struct {
	cmd     *kingpin.CmdClause   // Object
	action  *ForjAction          // Action name and help
	plugins []string             // Plugins implementing this object action.
	params  map[string]*ForjFlag // Collection of flags
}

type ForjParam interface {
	set_cmd(*kingpin.CmdClause, string, string, string, *ForjOpts)
	loadFrom(*kingpin.ParseContext)
}

/*type ForjList interface {
    SetListFlag()
}*/

type ForjValues struct {
}

// ForjParams type
const (
	// Arg : To set a ForjParam as Argument.
	Arg = "arg"
	// Flag : To set a ForjParam as Flag.
	Flag = "flag"
)

// List of ForjParams internal types
const (
	// String - Define the Param data type to string.
	String = "string"
	// Bool - Define the Param data type to bool.
	Bool = "bool"
	List = "list"
)

// NewForjCli : Initialize a new ForjCli object
//
// panic if app is nil.
func NewForjCli(app *kingpin.Application) (c *ForjCli) {
	if app == nil {
		panic("kingpin.Application cannot be nil.")
	}
	c = new(ForjCli)
	c.objects = make(map[string]*ForjObject)
	c.actions = make(map[string]*ForjAction)
	c.flags = make(map[string]*ForjFlag)
	c.App = app
	return
}

// AddFlag create a flag object at the application layer.
func (c *ForjCli) AddFlag(paramIntType, name, help string, options *ForjOpts) {
	f := new(ForjFlag)
	f.flag = c.App.Flag(name, help)
	f.set_options(options)

	switch paramIntType {
	case String:
		f.flagv = f.flag.String()
	case Bool:
		f.flagv = f.flag.Bool()
	}
	c.flags[name] = f
}

// AddActions create the list of referenced valid actions supported. kingpin layer created.
// It add them to the kingpin application layer.
//
// name     : Name of the action to add
// help     : Generic help to add to the action.
// for_forjj: True if the action is protected against plugins features.
func (c *ForjCli) AddActions(name, act_help, compose_help string, for_forjj bool) (r *ForjAction) {
	r = new(ForjAction)
	r.cmd = c.App.Command(name, act_help)
	r.help = compose_help
	r.internal_only = for_forjj
	c.actions[name] = r
	return
}

// AddObjectActions add a new object and the list of actions.
// It creates the ForjAction object for each action/object couple, to attach the object to kingpin object layer.
func (c *ForjCli) AddObjectActions(name, desc string, internal bool, actions ...string) {
	o := c.newForjObject(name, desc, internal)
	for _, action := range actions {
		if ar, found := c.actions[action]; found {
			o.actions[action] = newForjObjectAction(ar, name, fmt.Sprintf(ar.help, desc))
		} else {
			log.Printf("unknown action '%s'. Ignored.", action)
		}
	}
}

// AddObjectActionsParam get the Object and add several actions
// The command line is :
// forjj <action> <object>
func (c *ForjCli) AddObjectActionsParam(pType, pIntType, obj, name, desc string, options *ForjOpts, actions ...string) (err error) {
	var o *ForjObject
	if v, found := c.objects[obj]; !found {
		return fmt.Errorf("Unknown object '%s'. It must be created before with 'AddObjectActions'.", obj)
	} else {
		o = v
	}

	actionsInError := make([]string, 0, 2)

	for _, action := range actions {
		param := c.newParam(pType, name)

		var oa *ForjObjectAction

		if v, found := o.actions[action]; !found {
			actionsInError = append(actionsInError, obj+"/"+action)
			continue
		} else {
			oa = v
		}

		param.set_cmd(oa.cmd, pIntType, name, desc, options)
	}
	if len(actionsInError) > 0 {
		err = fmt.Errorf("Object/Actions '%s' are invalid. Argument '%s' ignored.", strings.Join(actionsInError, "', '"), name)
	}
	return

}

// AddActionsParam add a ForjParam to several actions. It creates the action layer of cmd in kingpin.
//
// name: name
// help: help
// options: Collection of options. See set().
// actions: List of actions to attach.
func (c *ForjCli) AddActionsParam(pType, pIntType, name, help string, options *ForjOpts, actions ...string) (err error) {
	actionsInError := make([]string, 0, 2)

	for _, action := range actions {
		param := c.newParam(pType, name)

		var act *ForjAction

		if v, found := c.actions[action]; found {
			act = v
		} else {
			actionsInError = append(actionsInError, action)
		}

		param.set_cmd(act.cmd, pIntType, name, help, options)
	}
	if len(actionsInError) > 0 {
		err = fmt.Errorf("Actions '%s' are invalid. Argument '%s' ignored.", strings.Join(actionsInError, "', '"), name)
	}
	return
}

// newParam create a ForjFlag or ForjArg defined by `paramType`
func (c *ForjCli) newParam(paramType, name string) ForjParam {
	switch paramType {
	case Arg:
		return new(ForjArg)
	case Flag:
		return new(ForjFlag)
	case List:
		l := new(ForjFlagList)
		if v, found := c.list[name]; found {
			l.obj = v
		} else {
			gotrace.Trace("Unable to find '%s' list.", name)
		}
		return l
	}
	return nil
}

// Create the ForjAction object to attach to the object parent.
func newForjObjectAction(ar *ForjAction, name, desc string) *ForjObjectAction {
	a := new(ForjObjectAction)
	a.action = ar
	a.cmd = ar.cmd.Command(name, fmt.Sprintf(ar.help, desc))
	a.params = make(map[string]*ForjFlag)
	a.plugins = make([]string, 0, 5)
	return a
}

func (c *ForjCli) newForjObject(object_name, description string, internal bool) (o *ForjObject) {
	o = new(ForjObject)
	o.actions = make(map[string]*ForjObjectAction)
	o.help = description
	o.internal = internal
	c.objects[object_name] = o
	return
}

// LoadPluginData: Load Plugin Definition in cli.
func (c *ForjCli) LoadPluginData(data *goforjj.YamlPluginComm) error {
	return nil
}
