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

func (o *ForjOpts) Default(v string) *ForjOpts {
	o.opts["default"] = v
	return o
}

func (o *ForjOpts) Short(b byte) *ForjOpts {
	o.opts["short"] = b
	return o
}

func (o *ForjOpts) Envar(v string) *ForjOpts {
	o.opts["envar"] = v
	return o
}

func (o *ForjOpts) MergeWith(fromOpts *ForjOpts) {
	for k, opt := range fromOpts.opts {
		o.opts[k] = opt
	}
}
