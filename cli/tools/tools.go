package tools

import (
	"strconv"
	"fmt"
)

// ToBoolWithAddr interpret a collection of type in Bool
// If value type is a pointer, it will return a *bool
// else it will return a bool
// or at the env a nil.
func ToBoolWithAddr(value interface{}) (interface{}, error) {
	str := ""
	addr := false
	switch value.(type) {
	case *string:
		str = *value.(*string)
		addr = true
	case string:
		str = value.(string)
	case bool:
		return value.(bool), nil
	case *bool:

		return value.(*bool), nil
	}

	if b, err := strconv.ParseBool(str); err != nil {
		return nil, fmt.Errorf("Unable to interpret string as boolean. %s", err)
	} else {
		if addr {
			return &b, nil
		}
		return b, nil
	}

}

// ToBoolWithAddr interpret a collection of type in Bool
// If value type is a pointer, it will return a *bool
// else it will return a bool
// or at the env a nil.
func ToBool(value interface{}) (bool, error) {
	if b, err := ToBoolWithAddr(value) ; err != nil {
		return false, err
	} else {

		switch b.(type) {
		case bool:
			return b.(bool), nil
		case *bool:
			return *b.(*bool), nil
		}
	}
	return false, nil
}
