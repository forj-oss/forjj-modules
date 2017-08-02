package cli

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


func (o *ForjOpts) MergeWith(fromOpts *ForjOpts) {
	for k, opt := range fromOpts.opts {
		o.opts[k] = opt
	}
}
