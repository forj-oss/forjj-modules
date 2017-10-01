package cli

import (
	"fmt"
	"github.com/forj-oss/forjj-modules/trace"
)

// IsAppValueFound return true if the parameter value is found on App Layer
func (c *ForjCli) IsAppValueFound(paramValue string) bool {
	if _, found := c.flags[paramValue]; found {
		return true
	}
	return false
}

func (c *ForjCli) GetObject(obj_name string) *ForjObject {
	if o, err := c.getObject(obj_name); err == nil {
		return o
	} else {
		gotrace.Trace("%s", err)
	}
	return nil
}

// GetAppFlag return the Application layer flag named paramValue.
func (c *ForjCli) GetAppFlag(paramValue string) *ForjFlag {
	if v, found := c.flags[paramValue]; found {
		return v
	}
	return nil
}

// GetAppBoolValue return a bool value of the parameter found at App layer.
// return false if
//
// - the parameter is not found
//
// - a different type
//
// - parameter value is false
//
// To check if the parameter exist, use IsAppValueFound.
func (c *ForjCli) GetAppBoolValue(paramValue string) (bool, error) {
	var f *ForjFlag

	if v, found := c.flags[paramValue]; found {
		f = v
	} else {
		return false, fmt.Errorf("Unable to find '%s' parameter from Application layer.", paramValue)
	}

	if c.parse {
		return to_bool(f.flagv), nil
	}

	// Get from Parse time
	if c.cli_context.context == nil {
		return false, fmt.Errorf("Unable to find '%s' parameter from Application layer context. Context nil.", paramValue)
	}

	var (
		value interface{}
		found bool
	)
	if value, found = c.cli_context.context.GetFlagValue(f.flag); found {
		return to_bool(value), nil
	}
	return false, fmt.Errorf("Unable to find '%s' parameter from Application layer context.", paramValue)
}

// GetAppStringValue return a string of the parameter at App layer.
// An empty string is returned if:
//
// - the parameter is not found
//
// - a different type
//
// - parameter value is ""
//
// To check if the parameter exist, use IsAppValueFound.
//
// Note that this function works during the parse context.
func (c *ForjCli) GetAppStringValue(paramValue string) (string, error) {
	var f *ForjFlag

	if v, found := c.flags[paramValue]; found {
		f = v
	} else {
		return "", fmt.Errorf("Unable to find '%s' parameter from Application layer.", paramValue)
	}

	if c.parse {
		return to_string(f.flagv), nil
	}
	// Get from Parse time
	if c.cli_context.context == nil {
		return "", fmt.Errorf("Unable to find '%s' parameter from Application layer context. Context nil.", paramValue)
	}
	if v, found := c.cli_context.context.GetFlagValue(f.flag); found {
		return to_string(v), nil
	}
	return "", fmt.Errorf("Unable to find '%s' parameter from Application layer context.", paramValue)
}

// GetActionStringValue return a string of the parameter (flag/arg) attached to an action.
// An empty string is returned if:
//
// - the parameter is not found
//
// - a different type
//
// - parameter value is ""
//
// To check if the parameter exist, use IsAppValueFound.
//
// Note that this function works during the parse context.
func (c *ForjCli) GetActionStringValue(action_name, paramValue string) (string, error) {
	var f ForjParam

	var action *ForjAction
	if a, found := c.actions[action_name] ; !found {
		return "", fmt.Errorf("'%s' action not found")
	} else {
		action = a
	}

	if v, found := action.params[paramValue]; found {
		f = v
	} else {
		return "", fmt.Errorf("Unable to find '%s' parameter from action '%s'.", paramValue, action_name)
	}

	if c.parse {
		return to_string(f.GetStringAddr()), nil
	}
	// Get from Parse time
	if c.cli_context.context == nil {
		return "", fmt.Errorf("Unable to find '%s' parameter from action '%s' context. Context nil.", paramValue, action_name)
	}
	if v, found := f.GetContextValue(c.cli_context.context); found {
		return to_string(v), nil
	}
	return "", fmt.Errorf("Unable to find '%s' parameter from action '%s' context.", paramValue, action_name)
}

// GetActionStringValue return a string of the parameter (flag/arg) attached to an action.
// An empty string is returned if:
//
// - the parameter is not found
//
// - a different type
//
// - parameter value is ""
//
// To check if the parameter exist, use IsAppValueFound.
//
// Note that this function works during the parse context.
func (c *ForjCli) GetActionBoolValue(action_name, paramValue string) (bool, error) {
	var f ForjParam

	var action *ForjAction
	if a, found := c.actions[action_name]; !found {
		return false, fmt.Errorf("'%s' action not found")
	} else {
		action = a
	}

	if v, found := action.params[paramValue]; found {
		f = v
	} else {
		return false, fmt.Errorf("Unable to find '%s' parameter from action '%s'.", paramValue, action_name)
	}

	if c.parse {
		return to_bool(f.GetBoolAddr()), nil
	}
	// Get from Parse time
	if c.cli_context.context == nil {
		return false, fmt.Errorf("Unable to find '%s' parameter from action '%s' context. Context nil.", paramValue, action_name)
	}
	if v, found := f.GetContextValue(c.cli_context.context); found {
		return to_bool(v), nil
	}
	return false, fmt.Errorf("Unable to find '%s' parameter from action '%s' context.", paramValue, action_name)
}

// IsParamFound. Search in defined parameter if it exists
func (c *ForjCli) IsParamFound(param_name string) (found bool) {
	_, found = c.values[param_name]
	return
}

// GetBoolValue : Get a Boolean of the parameter from cli.
//
// Get data from object defined.
// if object == "application", it will get data from the Application layer
func (c *ForjCli) GetBoolValue(object, key, param_name string) (bool, bool, error) {

	if v, found, err := c.getValue(object, key, param_name); found {
		return to_bool(v), true, nil
	} else {
		return false, false, err
	}
}

// GetStringValue : Get a String of the parameter from cli.
//
// Get data from object defined.
// if object == "application", it will get data from the Application layer
//
// returns:
// - value
// - found
// - default : If value is a pointer to a string, default is set to true.
// - error
//
func (c *ForjCli) GetStringValue(object, key, param_name string) (string, bool, bool, error) {
	if v, found, err := c.getValue(object, key, param_name); found {
		if _, ok := v.(*string); ok {
			return to_string(v), true, true, nil
		}
		return to_string(v), true, false, nil
	} else {
		return "", false, false, err
	}
}

// IsObjectList returns
// - true if the context is a list and is that object.
// - true if the action has a ObjectList
func (c *ForjCli) IsObjectList(object, key, obj_name string) bool {
	if c.cli_context.list != nil {
		return true
	}
	// Search in flags if the object list has been added.

	return false
}

// LoadCli Same as LoadContext. But in final stage
//
// Load all cli data to internal object representative
func (c *ForjCli) LoadCli() error {
	return nil
}

func (c *ForjCli) GetObjectValues(obj_name string) map[string]*ForjData {
	if v, found := c.values[obj_name]; found {
		return v.records
	}
	return make(map[string]*ForjData)
}
