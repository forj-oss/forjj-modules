package gotrace

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"runtime"
	"regexp"
	"strconv"
)

const (
	fatal_mode int = 0
	error_mode = 1 + fatal_mode
	warning_mode = 1 + error_mode
	info_mode = 1 + warning_mode
	debug_mode = 1 + info_mode
	debug_level_mode = 1 + debug_mode
)

type Debug struct {
	debug int
	printf func(prefix, s string, a ...interface{})(string)
}

var internal_debug Debug

func SetDebugPrintfHandler(printf func(prefix, s string, a ...interface{})(string)) {
	internal_debug.printf = printf
}

func SetDebug() {
	internal_debug.debug = debug_mode
}

func SetError() {
	internal_debug.debug = error_mode
}

func SetFatalError() {
	internal_debug.debug = fatal_mode
}

func SetWarning() {
	internal_debug.debug = warning_mode
}

func SetInfo() {
	internal_debug.debug = info_mode
}

func SetDebugLevel(level int) {
	internal_debug.debug = debug_mode + level
}

func IsDebugMode() bool {
	return (internal_debug.debug >= debug_mode)
}

func IsInfoMode() bool {
	return (internal_debug.debug >= info_mode)
}

func IsWarningMode() bool {
	return (internal_debug.debug >= warning_mode)
}

func IsErrorMode() bool {
	return (internal_debug.debug >= error_mode )
}

func IsFatalMode() bool {
	return (internal_debug.debug >= fatal_mode)
}

func (Debug)prefix(mode int) string {
	values := []string{ "FATAL ERROR !", "ERROR !", "WARNING !", "INFO", "DEBUG", "DEBUG" }

	if mode > debug_level_mode {
		return values[debug_level_mode] + strconv.Itoa(mode - debug_mode)
	} else {
		return values[mode]
	}
}

func Trace(s string, a ...interface{}) (_ string) {
	mymode := debug_mode
	if internal_debug.debug < mymode {
		return
	}
	return internal_debug.funcprintf(internal_debug.prefix(mymode ), s, a...)
}

func TraceLevel(level int, s string, a ...interface{}) (_ string) {
	if level < 0 {
		level = 0
	}
	mymode := debug_mode + level
	if internal_debug.debug < mymode {
		return
	}
	return internal_debug.funcprintf(internal_debug.prefix(mymode), s, a...)
}

func Warning(s string, a ...interface{}) (_ string) {
	mymode := warning_mode
	if internal_debug.debug < mymode {
		return
	}
	yellow := color.New(color.FgHiYellow).SprintFunc()
	return internal_debug.funcprintf(yellow(internal_debug.prefix(mymode)), s, a...)
}

func Error(s string, a ...interface{}) (_ string) {
	mymode := error_mode
	if internal_debug.debug < mymode {
		return
	}
	red := color.New(color.FgHiRed).SprintFunc()
	return internal_debug.funcprintf(red(internal_debug.prefix(mymode)), s, a...)
}

func FatalError(s string, a ...interface{}) (_ string) {
	mymode := fatal_mode
	if internal_debug.debug < mymode {
		return
	}
	red := color.New(color.FgHiRed).SprintFunc()
	return internal_debug.funcprintf(red(internal_debug.prefix(mymode)), s, a...)
}

func Info(s string, a ...interface{}) (_ string) {
	mymode := info_mode
	if internal_debug.debug < mymode {
		return
	}
	green := color.New(color.FgGreen).SprintFunc()
	return internal_debug.printf(green(internal_debug.prefix(mymode)), s, a...)
}

func (d *Debug) funcprintf(prefix, s string, a ...interface{}) string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	if d.printf != nil {
		return d.printf(prefix + " " + f.Name(), s, a...)
	} else {
		return d.internal_printf(prefix + " " + f.Name(), s, a...)
	}
}

func (d *Debug) internal_printf(prefix, s string, a ...interface{}) (ret string) {
	txt := fmt.Sprintf("%s: %s\n", prefix, s)
	ret = fmt.Sprintf(txt, a...)
	fmt.Print(ret)
	return
}

func Test(s string, a ...interface{}) (_ string){
	return internal_debug.printf("TEST", s, a...)
}

func (d *Debug)init() {
	d.debug = warning_mode
	SetDebugPrintfHandler(d.internal_printf)
	debug := os.Getenv("GOTRACE")
	if debug == "true" || debug == "debug" {
		d.debug = debug_mode
	} else if found, _  := regexp.MatchString("[0-9]+", debug) ; found {
		if v, err := strconv.Atoi(debug) ; err != nil {
			d.debug = debug_mode
			d.printf("DEBUG CONF", "Invalid GOTRACE number %s", debug)
		} else {
			d.debug = debug_mode + v
		}
	} else if debug == "info" {
		d.debug = info_mode
	} else if debug == "warning" {
		d.debug = warning_mode
	} else if debug == "error" {
		d.debug = error_mode
	} else if debug == "fatal" {
		d.debug = fatal_mode
	}
}

func init() {
	internal_debug.init()
}
