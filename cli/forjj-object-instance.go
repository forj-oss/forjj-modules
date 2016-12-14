package cli

import (
	"fmt"
	"github.com/kr/text"
)

type ForjObjectInstance struct {
	name              string // Instance name
	additional_fields map[string]*ForjField
}

func NewObjectInstance(name string) *ForjObjectInstance {
	return &ForjObjectInstance{
		name:              name,
		additional_fields: make(map[string]*ForjField),
	}
}

func (i *ForjObjectInstance) String() string {
	ret := fmt.Sprintf("Object Instance (%p):\n", i)
	ret += fmt.Sprintf("  name: '%s'\n", i.name)
	ret += fmt.Sprint("  fields (map):\n")
	for key, field := range i.additional_fields {
		ret += fmt.Sprintf("    %s: \n", key)
		ret += text.Indent(field.String(), "      ")
	}
	return ret
}

func (i *ForjObjectInstance) hasField(name string) (found bool) {
	_, found = i.additional_fields[name]
	return
}

func (i *ForjObjectInstance) addField(o *ForjObject, pIntType, name, help, re string, opts *ForjOpts) *ForjObjectInstance {
	if i == nil {
		return nil
	}
	if i.hasField(name) {
		o.err = fmt.Errorf("Additionnal field '%s' already exist.", name)
		return nil
	}
	f := NewField(o, pIntType, name, help, re, opts)
	i.additional_fields[name] = f

	return i
}
