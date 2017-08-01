package cli

import (
	"fmt"
	"github.com/forj-oss/forjj-modules/trace"
	"strconv"
)

func (c *ForjCli) SetValue(object, instance, atype, attr string, value interface{}) (err error) {
	r := c.values[object]
	if r, err = r.set(instance, atype, attr, value); err != nil {
		return err
	}
	c.values[object] = r
	gotrace.Trace("Added instance attribute '%s/%s' to object '%s'", instance, attr, object)
	return nil
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
		d = newData(action)
		r.records[key] = d
	} else {
		d = v
		if d.attrs["action"] == "setup" && action != "" && action != "setup" {
			gotrace.Trace("object '%s' action moved from initial action 'setup' to '%s'", action)
			d.attrs["action"] = action
		}
		if d.attrs["action"] != action && action != "setup" {
			c.err = fmt.Errorf("Unable to %s AND %s attribute at the same time. "+
				"Please remove %s to one of the 2 different action and retry",
				d.attrs["action"], action, object)
			return nil
		}
	}
	return
}

type ForjRecords struct {
	records map[string]*ForjData // Collection of records identified by object key.
}

func newRecords() (r *ForjRecords) {
	r = new(ForjRecords)
	r.records = make(map[string]*ForjData)
	return
}

func (r *ForjRecords) String() (ret string) {
	ret = fmt.Sprintf("records : %d\n", len(r.records))
	for key, record := range r.records {
		ret += fmt.Sprintf("    key: %s : %d\n", key, len(record.attrs))
		for attr_name, attr_value := range record.attrs {
			if attr_value == nil {
				ret += fmt.Sprintf("        %s : Not defined\n", attr_name)
				continue
			}
			if v, ok := attr_value.(string); ok {
				ret += fmt.Sprintf("        %s : %s\n", attr_name, v)
				continue
			}
			if v, ok := attr_value.(*string); ok {
				ret += fmt.Sprintf("        %s : %s (%p) - Default\n", attr_name, *v, v)
			}
		}
	}
	return
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

func (r *ForjRecords) set(instance, atype, attr string, value interface{}) (_ *ForjRecords, err error) {
	if r == nil {
		r = newRecords()
	}
	i := r.records[instance]
	if i, err = i.set(atype, attr, value); err != nil {
		return nil, err
	}
	r.records[instance] = i
	gotrace.Trace("Added attribute '%s' to instance '%s'", attr, instance)
	return r, nil
}

type ForjData struct {
	attrs map[string]interface{} // Collection of Values per Attribute Name.
	//instance_attrs map[string]ForjInstanceData
}

func newData(defaut_action string) (r *ForjData) {
	r = new(ForjData)
	r.attrs = make(map[string]interface{})
	//r.instance_attrs = make(map[string]ForjInstanceData)
	r.set(String, "action", defaut_action)
	return
}

type ForjInstanceData map[string]interface{}

func (r *ForjData) Attrs() map[string]interface{} {
	return r.attrs
}

func (r *ForjData) Keys() (keys []string) {
	keys = make([]string, len(r.attrs))
	iCount := 0
	for key := range r.attrs {
		keys[iCount] = key
		iCount++
	}
	return
}

func (d *ForjData) set(atype, key string, value interface{}) (*ForjData, error) {
	if d == nil {
		d = newData("setup") // default action
	}
	switch atype {
	case String:
		d.attrs[key] = value
		gotrace.Trace("Added attribute '%s' value '%s'", key, value)
	case Bool:
		str := ""
		switch value.(type) {
		case *string:
			str = *value.(*string)
		case string:
			str = value.(string)
		case bool, *bool:
			d.attrs[key] = value
			return d, nil
		}

		if b, err := strconv.ParseBool(str); err != nil {
			return nil, fmt.Errorf("Unable to interpret string as boolean. %s", err)
		} else {
			d.attrs[key] = b
			gotrace.Trace("Added attribute '%s' value '%t'", key, b)
		}
	}
	return d, nil
}

func (d *ForjData) GetString(param string) (ret string) {
	if v, _, err := d.Get(param); err != nil {
		return
	} else {
		switch v.(type) {
		case string:
			ret = v.(string)
		case *string:
			ret = *v.(*string)
		}
	}
	return
}

func (d *ForjData) Get(param string) (ret interface{}, found bool, err error) {
	if v, isfound := d.attrs[param]; isfound {
		ret = v
		found = (ret != nil)
	} else {
		err = fmt.Errorf("Unable to find attribute '%s'.", param)
	}
	return
}
