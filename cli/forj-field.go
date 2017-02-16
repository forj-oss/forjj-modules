package cli

import "fmt"

type ForjField struct {
	name       string      // name
	help       string      // help
	value_type string      // Expected value type
	key        bool        // true if this field is a key for list.
	obj        *ForjObject // Object where this field is attached.
	options    *ForjOpts   // Default field options

	found     bool                 // True if the flag was used.
	plugins   []string             // List of plugins that use this flag.
	inActions map[string]ForjParam // Collection of flags linked to Main actions. From
	// AddActionFlagsFromObjectAction
	regexp string // Regexp to validate input.
}

func (f *ForjField) String() string {
	ret := fmt.Sprintf("Field (%p):\n", f)
	ret += fmt.Sprintf("  name: '%s'\n", f.name)
	ret += fmt.Sprintf("  help: '%s'\n", f.help)
	ret += fmt.Sprintf("  value_type: '%s'\n", f.value_type)
	ret += fmt.Sprintf("  found: '%s'\n", f.found)
	return ret
}

func NewField(o *ForjObject, pIntType, name, help, re string, opts *ForjOpts) *ForjField {
	return &ForjField{
		name:       name,
		help:       help,
		value_type: pIntType,
		inActions:  make(map[string]ForjParam),
		regexp:     re,
		obj:        o,
		options:    opts,
	}
}
