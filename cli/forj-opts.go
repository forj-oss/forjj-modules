package cli

import (
	"os"
	"forjj-modules/cli/tools"
)

// Flag/Arg options management

type ForjOpts struct {
	opts map[string]interface{}
}

func Opts() *ForjOpts {
	o := new(ForjOpts)
	o.opts = make(map[string]interface{})
	return o
}

func (o *ForjOpts) Required() *ForjOpts {
	o.opts["required"] = true
	return o
}

func (o *ForjOpts) NotRequired() *ForjOpts {
	delete(o.opts, "required")
	return o
}

func (o *ForjOpts) IsRequired() bool {
	if v, found := o.opts["required"]; found {
		return v.(bool)
	}
	return false
}

func (o *ForjOpts) Default(v string) *ForjOpts {
	o.opts["default"] = v
	return o
}

func (o *ForjOpts) NoDefault() *ForjOpts {
	delete(o.opts, "default")
	return o
}

func (o *ForjOpts) Short(b byte) *ForjOpts {
	o.opts["short"] = b
	return o
}

func (o *ForjOpts) NoShort() *ForjOpts {
	delete(o.opts, "short")
	return o
}

func (o *ForjOpts) Envar(v string) *ForjOpts {
	o.opts["envar"] = v
	return o
}

func (o *ForjOpts) NoEnvar() *ForjOpts {
	delete(o.opts, "envar")
	return o
}

func (o *ForjOpts) HasEnvar() (bool, string) {
	if v, found := o.opts["envar"]; found {
		return true, v.(string)
	}
	return false, ""
}

// GetDefault return the default value from defined options.
// Used to set single object attribute default value
// It must return a pointer to a pType value type (*string, *bool, ...)
func (o *ForjOpts) GetDefault(pType string) (interface{}) {
	if o == nil {
		return nil
	}
	switch pType {
	case String:
		s := ""
		if found, v := o.HasEnvar() ; found {
			s = os.Getenv(v)
		}
		if s != "" {
			return &s
		}
		if v, found2 := o.opts["default"] ; found2 && v != "" {
			return &v
		}
	case Bool:
		if found, v := o.HasEnvar() ; found {
			s := os.Getenv(v)
			if s != "" {
				if v, err := tools.ToBoolWithAddr(&s) ; err == nil {
					return v
				}
			}
		}
		if v, found2 := o.opts["default"] ; found2 && v != "" {
			if b, err := tools.ToBoolWithAddr(&v) ; err == nil {
				return b
			}
			return nil
		}
	}
	return nil
}

func (o *ForjOpts) MergeWith(fromOpts *ForjOpts) {
	for k, opt := range fromOpts.opts {
		o.opts[k] = opt
	}
}
