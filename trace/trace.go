package gotrace

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const (
	fatalMode      int = 0
	errorMode          = 1 + fatalMode
	warningMode        = 1 + errorMode
	infoMode           = 1 + warningMode
	debugMode          = 1 + infoMode
	debugLevelMode     = 1 + debugMode
)

// Debug implement a debug control structure
type Debug struct {
	debug         int
	defaultDebug  bool
	formatFunc    func(prefix, s string, a ...interface{}) string
	printFunc     func(a ...interface{}) (n int, err error)
	hideSecrets   bool
	secretsToHide []string
}

var internalDebug Debug

// SetDebugPrintfHandler define a different logger function to format differently
func SetDebugPrintfHandler(formatFunc func(prefix, s string, a ...interface{}) string, printFunc func(a ...interface{}) (int, error)) {
	internalDebug.formatFunc = formatFunc
	internalDebug.printFunc = printFunc
}

// SetDebug move the default debug mode to Debug
func SetDebug() {
	if internalDebug.defaultDebug {
		internalDebug.debug = debugMode
	}
}

// SetDebugMode define the overall app debug level to print.
func SetDebugMode(debug string) {
	if internalDebug.defaultDebug {
		internalDebug.setDebugMode(debug)
	}
}

// SetError move the default debug mode to Error
func SetError() {
	if internalDebug.defaultDebug {
		internalDebug.debug = errorMode
	}
}

// SetFatalError move the default debug mode to FatalError
func SetFatalError() {
	if internalDebug.defaultDebug {
		internalDebug.debug = fatalMode
	}
}

// SetWarning move the default debug mode to Warning
func SetWarning() {
	if internalDebug.defaultDebug {
		internalDebug.debug = warningMode
	}
}

// SetInfo move the default debug mode to Info
func SetInfo() {
	if internalDebug.defaultDebug {
		internalDebug.debug = infoMode
	}
}

// SetDebugLevel move the default debug mode to Debug at the requested level
func SetDebugLevel(level int) {
	if internalDebug.defaultDebug {
		internalDebug.debug = debugMode + level
	}
}

// IsDebugMode return true if we are at debug mode
func IsDebugMode() bool {
	return (internalDebug.debug >= debugMode)
}

// IsDebugLevelMode return true if we are at debug mode level
func IsDebugLevelMode(level int) bool {
	return (internalDebug.debug >= debugMode+level)
}

// IsInfoMode return true if we are at info mode
func IsInfoMode() bool {
	return (internalDebug.debug >= infoMode)
}

// IsWarningMode return true if we are at warning mode
func IsWarningMode() bool {
	return (internalDebug.debug >= warningMode)
}

// IsErrorMode return true if we are at error mode
func IsErrorMode() bool {
	return (internalDebug.debug >= errorMode)
}

// IsFatalMode return true if we are at fatal error mode
func IsFatalMode() bool {
	return (internalDebug.debug >= fatalMode)
}

func AddSecrets(secrets ...string) {
	internalDebug.addSecrets(secrets...)
}

func (Debug) prefix(mode int) string {
	values := []string{"FATAL ERROR !", "ERROR !", "WARNING !", "INFO", "DEBUG", "DEBUG"}

	if mode > debugLevelMode {
		return values[debugLevelMode] + strconv.Itoa(mode-debugMode)
	}
	return values[mode]
}

// Trace log a debug message
func Trace(s string, a ...interface{}) (_ string) {
	mymode := debugMode
	if internalDebug.debug < mymode {
		return
	}
	return internalDebug.print(internalDebug.prefix(mymode), s, a...)
}

// TraceLevel log a debug message at given level
func TraceLevel(level int, s string, a ...interface{}) (_ string) {
	if level < 0 {
		level = 0
	}
	mymode := debugMode + level
	if internalDebug.debug < mymode {
		return
	}
	return internalDebug.print(internalDebug.prefix(mymode), s, a...)
}

// Warning log a warning message
func Warning(s string, a ...interface{}) (_ string) {
	mymode := warningMode
	if internalDebug.debug < mymode {
		return
	}
	yellow := color.New(color.FgHiYellow).SprintFunc()
	return internalDebug.print(yellow(internalDebug.prefix(mymode)), s, a...)
}

// Error log an error message
func Error(s string, a ...interface{}) (_ string) {
	mymode := errorMode
	if internalDebug.debug < mymode {
		return
	}
	red := color.New(color.FgHiRed).SprintFunc()
	return internalDebug.print(red(internalDebug.prefix(mymode)), s, a...)
}

// FatalError log a fatal error message
func FatalError(s string, a ...interface{}) (_ string) {
	mymode := fatalMode
	if internalDebug.debug < mymode {
		return
	}
	red := color.New(color.FgHiRed).SprintFunc()
	return internalDebug.print(red(internalDebug.prefix(mymode)), s, a...)
}

// Info log an info message
func Info(s string, a ...interface{}) (_ string) {
	mymode := infoMode
	if internalDebug.debug < mymode {
		return
	}
	green := color.New(color.FgGreen).SprintFunc()
	return internalDebug.print(green(internalDebug.prefix(mymode)), s, a...)
}

// -------------------------------------- Internal Debug functions

func (d *Debug) addSecrets(secrets ...string) {
	if d == nil {
		return
	}

	if len(secrets) == 0 {
		return
	}

	d.secretsToHide = append(d.secretsToHide, secrets...)
}

func (d *Debug) doHideSecretsOn(value string) (ret string) {
	ret = value
	for _, secret := range d.secretsToHide {
		ret = strings.Replace(ret, secret, "***", -1)
	}
	return
}

// setDebugMode define the overall app debug level to print.
func (d *Debug) setDebugMode(debug string) {
	if debug == "true" || debug == "debug" {
		d.debug = debugMode
	} else if found, _ := regexp.MatchString("[0-9]+", debug); found {
		if v, err := strconv.Atoi(debug); err != nil {
			d.debug = debugMode
			d.print("DEBUG CONF", "Invalid GOTRACE number %s", debug)
		} else {
			d.debug = debugMode + v
		}
	} else if debug == "info" {
		d.debug = infoMode
	} else if debug == "warning" {
		d.debug = warningMode
	} else if debug == "error" {
		d.debug = errorMode
	} else if debug == "fatal" {
		d.debug = fatalMode
	} else {
		d.defaultDebug = true
	}
}

func (d *Debug) print(prefix, s string, a ...interface{}) (ret string) {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	if d.formatFunc != nil {
		ret = d.formatFunc(prefix+" "+f.Name(), s, a...)
	} else {
		ret = d.internalSprintf(prefix+" "+f.Name(), s, a...)
	}

	ret = d.doHideSecretsOn(ret)
	d.printFunc(ret)
	return
}

func (d *Debug) internalSprintf(prefix, s string, a ...interface{}) string {
	txt := fmt.Sprintf("%s: %s\n", prefix, s)
	return fmt.Sprintf(txt, a...)
}

// Test log a permanent test message (not filtered by debug mode)
func Test(s string, a ...interface{}) (_ string) {
	return internalDebug.print("TEST", s, a...)
}

func (d *Debug) init() {
	d.debug = warningMode
	SetDebugPrintfHandler(d.internalSprintf, fmt.Print)
	d.hideSecrets = true
	d.setDebugMode(os.Getenv("GOTRACE"))
	d.secretsToHide = make([]string, 0, 5)
}

func init() {
	internalDebug.init()
}
