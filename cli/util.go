package cli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func copyValue(src interface{}, dest interface{}) {
	switch src.(type) {
	case *int32:
		dest_b := dest.(*int32)
		src_b := src.(*int32)
		*dest_b = *src_b
	case *byte:
		dest_b := dest.(*byte)
		src_b := src.(*byte)
		*dest_b = *src_b
	case *bool:
		dest_b := dest.(*bool)
		src_b := src.(*bool)
		*dest_b = *src_b
	case *string:
		dest_s := dest.(*string)
		src_s := src.(*string)
		*dest_s = *src_s
	}
}

// Simple function to convert a dynamic type to bool
// it returns false by default except if the internal type is:
// - bool. value as is
// - string: call https://golang.org/pkg/strconv/#ParseBool
//
func to_bool(v interface{}) bool {
	switch v.(type) {
	case *bool:
		return *v.(*bool)
	case *string:
		s := *v.(*string)
		if b, err := strconv.ParseBool(s); err == nil {
			return b
		}
		return false
	}
	return false
}

// simply extract string from the dynamic type
// otherwise the returned string is empty.
func to_string(v interface{}) (result string) {
	switch v.(type) {
	case *string:
		result = *v.(*string)
	}
	return
}

func is_string(v interface{}) bool {
	switch v.(type) {
	case *string:
		return true
	}
	return false
}

func to_byte(v interface{}) (result int32) {
	switch v.(type) {
	case byte:
		result = *v.(*int32)
	}
	return
}

func is_byte(v interface{}) bool {
	switch v.(type) {
	case *byte, *int32:
		return true
	}
	return false
}

func Split(expr, s, sep string) []string {
	re, _ := regexp.Compile(strings.Replace(expr, sep, `\\?`+sep, 1))
	re_esc, _ := regexp.Compile(`\\` + sep)
	if len(expr) > 0 && len(s) == 0 {
		return []string{""}
	}

	matches := re.FindAllStringIndex(s, -1)
	res := make([]string, 0, len(matches))

	beg := 0
	end := 0

	for _, match := range matches {
		end = match[0]
		if match[1] != 0 {
			fmt.Printf("=> '%s'\n", s[match[0]:match[1]])
			fmt.Printf("=> '%#v'\n", re_esc.FindStringIndex(s[match[0]:match[1]]))
			if re_esc.FindStringIndex(s[match[0]:match[1]]) != nil {
				continue
			}
			res = append(res, strings.Replace(s[beg:end], "\\", "", -1))
		}
		beg = match[1]
	}

	if end != len(s) {
		res = append(res, s[beg:])
	}

	return res
}
