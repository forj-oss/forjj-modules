package gotrace

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"runtime"
)

var debug int

const debug_mode int = 3
const warning_mode int = 2
const error_mode int = 1
const fatal_mode int = 0

type Debug struct {
	debug int
}

var internal_debug Debug

func SetDebug() {
	internal_debug.debug = debug_mode
}

func Trace(s string, a ...interface{}) {
	if internal_debug.debug < debug_mode {
		return
	}
	internal_debug.funcprintf("DEBUG", s, a...)
}

func Warning(s string, a ...interface{}) {
	if internal_debug.debug < warning_mode {
		return
	}
	yellow := color.New(color.FgHiYellow).SprintFunc()
	internal_debug.funcprintf(yellow("WARNING !"), s, a...)
}

func Error(s string, a ...interface{}) {
	if internal_debug.debug < error_mode {
		return
	}
	red := color.New(color.FgHiRed).SprintFunc()
	internal_debug.funcprintf(red("ERROR !"), s, a...)
}

func Info(s string, a ...interface{}) {
	if internal_debug.debug < error_mode {
		return
	}
	green := color.New(color.FgGreen).SprintFunc()
	internal_debug.printf(green("INFO"), s, a...)
}

func (d *Debug) funcprintf(prefix, s string, a ...interface{}) {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	d.printf(prefix + " " + f.Name(), s, a...)
}

func (d *Debug) printf(prefix, s string, a ...interface{}) {
	txt := fmt.Sprintf("%s: %s\n", prefix, s)
	fmt.Printf(txt, a...)
}

func Test(s string, a ...interface{}) {
	internal_debug.printf("TEST", s, a...)
}

func init() {
	internal_debug.debug = warning_mode
	if os.Getenv("GOTRACE") == "true" {
		internal_debug.debug = debug_mode
	}
}
