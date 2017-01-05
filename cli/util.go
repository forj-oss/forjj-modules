package cli

import (
    "strconv"
    "regexp"
    "strings"
    "fmt"
)

// Simple function to convert a dynamic type to bool
// it returns false by default except if the internal type is:
// - bool. value as is
// - string: call https://golang.org/pkg/strconv/#ParseBool
//
func to_bool(v interface{}) bool {
    switch v.(type) {
    case bool:
        return v.(bool)
    case string:
        s := v.(string)
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
    case string:
        return v.(string)
    }
    return
}

func to_byte(v interface{}) (result int32) {
    switch v.(type) {
    case byte:
        return v.(int32)
    }
    return
}

func is_byte(v interface{}) bool {
    if get_itype(v) == "string" {
        return true
    }
    return false
}

func get_itype(v interface{}) string {
    switch v.(type) {
    case byte:
        return "byte"
    case string:
        return "string"
    case bool:
        return "bool"
    }
    return ""
}

func Split(expr, s , sep string) []string {
    re, _ := regexp.Compile(strings.Replace(expr, sep, `\\?` + sep, 1))
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
