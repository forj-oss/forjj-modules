package cli

import (
	"fmt"
	"github.com/forj-oss/forjj-modules/trace"
	"strconv"
)

type ForjRecords struct {
	records map[string]*ForjData // Collection of records identified by object key.
}

type ForjData struct {
	attrs map[string]interface{} // Collection of Values per Attribute Name.
}

func (r *ForjRecords) String() (ret string) {
	ret = fmt.Sprint("records : \n")
	for key, record := range r.records {
		ret += fmt.Sprintf("    key: %s : \n", key)
		for attr_name, attr_value := range record.attrs {
			ret += fmt.Sprintf("        %s : %s\n", attr_name, attr_value)
		}
	}
	return
}

func (r *ForjData) Attrs() map[string]interface{} {
	return r.attrs
}

// GetFrom, get the param value from the defined context.
// If no context exists, it tries to get from App Flag layer
// It search in action_object and then action.
// If the value context is a list, it moves to get it from the App layer directly.
func (r *ForjRecords) Get(key, param string) (ret interface{}, found bool, err error) {
	if key == "" {
		if ls := len(r.records); ls == 1 {
			for k := range r.records {
				key = k
				gotrace.Trace("defined record key to '%s'.", k)
				break
			}
		} else {
			if ls > 1 {
				err = fmt.Errorf("Unable to identify a uniq record key. Found %d keys.", ls)
			} else {
				err = fmt.Errorf("Unable to find one record.")
			}
			return
		}
	}

	if v, isfound := r.records[key]; isfound {
		ret, found, err = v.Get(param)
	} else {
		err = fmt.Errorf("Unable to find record identified by key '%s'", key)
	}
	return
}

func (d *ForjData) set(atype, key string, value interface{}) error {
	switch atype {
	case String:
		d.attrs[key] = value
	case Bool:
		str := ""
		switch value.(type) {
		case *string:
			str = *value.(*string)
		case string:
			str = value.(string)
		}
		if b, err := strconv.ParseBool(str); err != nil {
			return fmt.Errorf("Unable to interpret string as boolean. %s", err)
		} else {
			d.attrs[key] = b
		}
	}
	return nil
}

func (d *ForjData) GetString(param string) (ret string) {
	if v, _, err := d.Get(param); err != nil {
		return
	} else {
		if v2, ok := v.(string); ok {
			ret = v2
		}
	}
	return
}

func (d *ForjData) Get(param string) (ret interface{}, found bool, err error) {
	if v, isfound := d.attrs[param]; isfound {
		ret = v
		found = true
	} else {
		err = fmt.Errorf("Unable to find attribute '%s'.", param)
	}
	return
}

func (c *ForjCli) setObjectAttributes(action, object, key string) (d *ForjData) {
	var r *ForjRecords
	if v, found := c.values[object]; !found {
		r = new(ForjRecords)
		r.records = make(map[string]*ForjData)
		c.values[object] = r
	} else {
		r = v
	}

	if v, found := r.records[key]; !found {
		d = new(ForjData)
		d.attrs = make(map[string]interface{})
		d.attrs["action"] = action
		r.records[key] = d
	} else {
		d = v
		if d.attrs["action"] != action {
			c.err = fmt.Errorf("Unable to %s AND %s attribute at the same time. "+
				"Please remove %s to one of the 2 different action and retry",
				d.attrs["action"], action, object)
			return nil
		}
	}
	return
}
