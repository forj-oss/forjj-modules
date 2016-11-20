package cli

import (
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
func (c *ForjCli) GetAppBoolValue(paramValue string) bool {
	var f *ForjFlag

	if v, found := c.flags[paramValue]; found {
		f = v
	}

	if c.parse {
		return to_bool(f.flagv)
	}

	// Get from Parse time
	if c.cli_context.context == nil {
		return false
	}

	var (
		value string
		found bool
	)
	if value, found = c.cli_context.context.GetFlagValue(f.flag); found {
		return to_bool(value)
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
	var f *ForjFlag

	if v, found := c.flags[paramValue]; found {
		f = v
	}

	if c.parse {
		return to_string(f.flagv)
	}
	// Get from Parse time
	if c.cli_context.context == nil {
		return ""
	}
	if v, found := c.cli_context.context.GetFlagValue(f.flag); found {
		return to_string(v)
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
func (c *ForjCli) GetBoolValue(object, key, param_name string) (bool, bool, error) {

	if v, found, err := c.getValue(object, key, param_name); found {
		if c.parse {
			return to_bool(v), true, nil
		}
		// Get from Parse time
		if c.cli_context.context == nil {
			return false, false, nil
		}

		p := c.getContextParam(object, key, param_name)
		switch p.(type) {
		case *ForjFlag:
			f := p.(*ForjFlag)
			if v, found := c.cli_context.context.GetFlagValue(f.flag); found {
				return to_bool(v), true, nil
			}
		case *ForjArg:
			a := p.(*ForjArg)
			if v, found := c.cli_context.context.GetArgValue(a.arg); found {
				return to_bool(v), true, nil
			}
		}
	} else {
		return false, false, err
	}
	return false, false, nil
}

// GetStringValue : Get a String of the parameter from cli.
//
// Get data from object defined.
// if object == "application", it will get data from the Application layer
func (c *ForjCli) GetStringValue(object, key, param_name string) (string, bool, error) {
	if v, found, err := c.getValue(object, key, param_name); found {
		return to_string(v), true, nil
	} else {
		return "", false, err
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
