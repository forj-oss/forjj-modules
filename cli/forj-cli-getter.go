package cli

import "github.com/forj-oss/forjj-modules/trace"

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
func (c *ForjCli) GetAppBoolValue(paramValue string) bool {
	if v, found := c.flags[paramValue]; found {
		return to_bool(v.flagv)
	}
	return false
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
func (c *ForjCli) GetAppStringValue(paramValue string) string {
	if v, found := c.flags[paramValue]; found {
		return to_string(v.flagv)
	}
	return ""
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
func (c *ForjCli) GetBoolValue(object, key, param_name string) (bool, bool) {
	if v, found := c.getValue(object, key, param_name); found {
		return to_bool(v), true
	}
	return false, false
}

// GetStringValue : Get a String of the parameter from cli.
//
// Get data from object defined.
// if object == "application", it will get data from the Application layer
func (c *ForjCli) GetStringValue(object, key, param_name string) (string, bool) {
	if v, found := c.getValue(param_name); found {
		return to_string(v), true
	}
	return "", false
}

// IsObjectList returns
// - true if the context is a list and is that object.
// - true if the action has a ObjectList
func (c *ForjCli) IsObjectList(object, key, obj_name string) bool {
	if c.context.list != nil {
		return true
	}
	// Search in flags if the object list has been added.

	return false
}

// LoadCli Same as LoadContext. But in final stage
//
// Load all cli data to internal object representative
func (c *ForjCli) LoadCli() error {

}
