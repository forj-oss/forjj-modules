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

// GetBoolValue : Get a Boolean of the parameter from context.
// If the parameter is not used in the context. Try to get it from App layer.
func (c *ForjCli) GetBoolValue(param_name string) bool {
	if v := c.getValue(param_name); v != nil {
		return to_bool(v)
	}
	return false
}

// GetStringValue : Get a String of the parameter from context.
// If the parameter is not used in the context. Try to get it from App layer.
func (c *ForjCli) GetStringValue(param_name string) string {
	if v := c.getValue(param_name); v != nil {
		return to_string(v)
	}
	return ""
}

// IsObjectList returns
// - true if the context is a list and is that object.
// - true if the action has a ObjectList
func (c *ForjCli) IsObjectList(obj_name string) bool {
	if c.context.list != nil {
		return true
	}
	// Search in flags if the object list has been added.

	return false
}
