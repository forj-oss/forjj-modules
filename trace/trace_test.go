package gotrace

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

var debug_msg string

func (d *Debug) test_printf(prefix, s string, a ...interface{}) string {
	return fmt.Sprintf(fmt.Sprintf("%s: %s", prefix, s), a...)
}

func TestDefault(t *testing.T) {
	t.Log("Expect debug to be in warning mode by default.")
	if internalDebug.debug != warningMode {
		t.Errorf("Expected to have debug mode set to 'warning'. Got %d", internalDebug.debug)
	}
}

func TestSetDebugPrintfHandler(t *testing.T) {
	t.Log("Expect SetDebugPrintfHandler to work.")

	const test = "test: toto"
	if v := internalDebug.test_printf("test", "toto"); v != test {
		t.Errorf("Internal test_printf should return '%s'. Got '%s", test, v)
		return
	}

	SetDebugPrintfHandler(internalDebug.test_printf)

	if reflect.ValueOf(internalDebug.printf).Pointer() != reflect.ValueOf(internalDebug.test_printf).Pointer() {
		t.Error("Internal printf should set. Different function found.")
	}
}

func TestSetDebug(t *testing.T) {
	t.Log("Expect SetDebug to set it in debug mode.")
	SetDebug()
	if internalDebug.debug != debugMode {
		t.Errorf("Expected to have debug mode set to 'debug'. Got %d", internalDebug.debug)
	} else if !internalDebug.defaultDebug {
		t.Errorf("Expected to have debug default mode set to 'true'. Got %t", internalDebug.defaultDebug)
	}
}

func TestSetDebugLevel(t *testing.T) {
	t.Log("Expect SetDebugLevel to set it in debug mode.")
	SetDebugLevel(1)
	if internalDebug.debug != debugMode+1 {
		t.Errorf("Expected to have debug mode set to 'debug level 1'. Got %d", internalDebug.debug)
	}
}

func TestSetWarning(t *testing.T) {
	t.Log("Expect SetDebug to set it in debug mode.")
	SetWarning()
	if internalDebug.debug != warningMode {
		t.Errorf("Expected to have debug mode set to 'warning'. Got %d", internalDebug.debug)
	}
}

func TestSetError(t *testing.T) {
	t.Log("Expect SetDebug to set it in debug mode.")
	SetError()
	if internalDebug.debug != errorMode {
		t.Errorf("Expected to have debug mode set to 'error'. Got %d", internalDebug.debug)
	}
}

func TestSetFatalError(t *testing.T) {
	t.Log("Expect SetDebug to set it in debug mode.")
	SetFatalError()
	if internalDebug.debug != fatalMode {
		t.Errorf("Expected to have debug mode set to 'fatal'. Got %d", internalDebug.debug)
	}
}

func TestIsDebugMode(t *testing.T) {
	t.Log("Expect IsDebugMode to set it in debug mode.")

	SetDebug()

	if !IsDebugMode() {
		t.Error("Expected to have debug mode detected. IsDebugMode returned false.")
	}
	if !IsWarningMode() {
		t.Error("Expected to have warning mode detected when debug is activated. IsWarningMode returned false.")
	}
	if !IsErrorMode() {
		t.Error("Expected to have error mode detected when debug is activated. IsErrorMode returned false.")
	}
	if !IsFatalMode() {
		t.Error("Expected to have fatal mode detected when debug is activated. IsFatalMode returned false.")
	}

	SetDebugLevel(5)

	if !IsDebugMode() {
		t.Error("Expected to have debug mode detected. IsDebugMode returned false.")
	}
	if !IsWarningMode() {
		t.Error("Expected to have warning mode detected when debug is activated. IsWarningMode returned false.")
	}
	if !IsErrorMode() {
		t.Error("Expected to have error mode detected when debug is activated. IsErrorMode returned false.")
	}
	if !IsFatalMode() {
		t.Error("Expected to have fatal mode detected when debug is activated. IsFatalMode returned false.")
	}
}

func TestIsDebugLevelMode(t *testing.T) {
	t.Log("Expect IsDebugLevelMode to set it in debug mode.")

	SetDebug()

	if IsDebugLevelMode(1) {
		t.Error("Expected to have debug mode detected. IsDebugLevelMode 0 returned true.")
	}
	if !IsDebugMode() {
		t.Error("Expected to have debug mode detected. IsDebugMode returned false.")
	}
	if !IsWarningMode() {
		t.Error("Expected to have warning mode detected when debug is activated. IsWarningMode returned false.")
	}
	if !IsErrorMode() {
		t.Error("Expected to have error mode detected when debug is activated. IsErrorMode returned false.")
	}
	if !IsFatalMode() {
		t.Error("Expected to have fatal mode detected when debug is activated. IsFatalMode returned false.")
	}

	SetDebugLevel(5)

	if !IsDebugLevelMode(1) {
		t.Error("Expected to have debug mode detected. IsDebugLevelMode 0 returned false.")
	}
	if !IsDebugMode() {
		t.Error("Expected to have debug mode detected. IsDebugMode returned false.")
	}
	if !IsWarningMode() {
		t.Error("Expected to have warning mode detected when debug is activated. IsWarningMode returned false.")
	}
	if !IsErrorMode() {
		t.Error("Expected to have error mode detected when debug is activated. IsErrorMode returned false.")
	}
	if !IsFatalMode() {
		t.Error("Expected to have fatal mode detected when debug is activated. IsFatalMode returned false.")
	}
}

func TestIsInfoMode(t *testing.T) {
	t.Log("Expect IsInfoMode to return appropriate mode.")

	SetInfo()

	if IsDebugMode() {
		t.Error("Expected to have info mode detected. IsDebugMode returned true.")
	}
	if !IsInfoMode() {
		t.Error("Expected to NOT have info mode detected. IsInfoMode returned true.")
	}
	if !IsWarningMode() {
		t.Error("Expected to have warning mode detected when info is activated. IsWarningMode returned false.")
	}
	if !IsErrorMode() {
		t.Error("Expected to have error mode detected when info is activated. IsErrorMode returned false.")
	}
	if !IsFatalMode() {
		t.Error("Expected to have fatal mode detected when info is activated. IsFatalMode returned false.")
	}

}

func TestIsWarningMode(t *testing.T) {
	t.Log("Expect IsWarningMode to return appropriate mode.")

	SetWarning()

	if IsDebugMode() {
		t.Error("Expected to NOT have debug mode detected. IsDebugMode returned true.")
	}
	if IsInfoMode() {
		t.Error("Expected to NOT have info mode detected. IsInfoMode returned true.")
	}
	if !IsWarningMode() {
		t.Error("Expected to have warning mode detected when warning is activated. IsWarningMode returned false.")
	}
	if !IsErrorMode() {
		t.Error("Expected to have error mode detected when warning is activated. IsErrorMode returned false.")
	}
	if !IsFatalMode() {
		t.Error("Expected to have fatal mode detected when warning is activated. IsFatalMode returned false.")
	}

}

func TestIsErrorMode(t *testing.T) {
	t.Log("Expect IsErrorMode to return appropriate mode.")

	SetError()

	if IsDebugMode() {
		t.Error("Expected to NOT have debug mode detected when error is activated. IsDebugMode returned true.")
	}
	if IsInfoMode() {
		t.Error("Expected to NOT have info mode detected. IsInfoMode returned true.")
	}
	if IsWarningMode() {
		t.Error("Expected to NOT have warning mode detected when error is activated. IsWarningMode returned true.")
	}
	if !IsErrorMode() {
		t.Error("Expected to have error mode detected when error is activated. IsErrorMode returned false.")
	}
	if !IsFatalMode() {
		t.Error("Expected to have fatal mode detected when error is activated. IsFatalMode returned false.")
	}

}

func TestIsFatalMode(t *testing.T) {
	t.Log("Expect IsFatalErrorMode to return appropriate mode.")

	SetFatalError()

	if IsDebugMode() {
		t.Error("Expected to NOT have debug mode detected when fatal error is activated. IsDebugMode returned true.")
	}
	if IsInfoMode() {
		t.Error("Expected to NOT have info mode detected. IsInfoMode returned true.")
	}
	if IsWarningMode() {
		t.Error("Expected to NOT have warning mode detected when fatal error is activated. IsWarningMode returned true.")
	}
	if IsErrorMode() {
		t.Error("Expected to have error mode detected when fatal error is activated. IsErrorMode returned true.")
	}
	if !IsFatalMode() {
		t.Error("Expected to have fatal mode detected when fatal error is activated. IsFatalMode returned false.")
	}

}

func TestTraceLevel(t *testing.T) {
	t.Log("Expect TraceLevel to display in appropriate mode.")

	SetDebugLevel(5)

	test := "DEBUG5 forjj-modules/trace.TestTraceLevel: blabla toto"
	if ret := TraceLevel(5, "blabla %s", "toto"); ret != test {
		t.Errorf("Expected TraceLevel to display '%s'. Got '%s'.", test, ret)
	}

	SetDebugLevel(4)

	test = ""
	if ret := TraceLevel(5, "blabla %s", "toto"); ret != test {
		t.Errorf("Expected TraceLevel to display '%s'. Got '%s'.", test, ret)
	}

	SetDebug()

	if ret := TraceLevel(5, "blabla %s", "toto"); ret != test {
		t.Errorf("Expected TraceLevel to display '%s'. Got '%s'.", test, ret)
	}

	SetWarning()

	if ret := TraceLevel(5, "blabla %s", "toto"); ret != test {
		t.Errorf("Expected TraceLevel to display '%s'. Got '%s'.", test, ret)
	}

	SetError()

	if ret := TraceLevel(5, "blabla %s", "toto"); ret != test {
		t.Errorf("Expected TraceLevel to display '%s'. Got '%s'.", test, ret)
	}

	SetFatalError()

	if ret := TraceLevel(5, "blabla %s", "toto"); ret != test {
		t.Errorf("Expected TraceLevel to display '%s'. Got '%s'.", test, ret)
	}
}

func TestTrace(t *testing.T) {
	t.Log("Expect Trace to display in appropriate mode.")

	SetDebugLevel(5)

	test := "DEBUG forjj-modules/trace.TestTrace: blabla toto"
	if ret := Trace("blabla %s", "toto"); ret != test {
		t.Errorf("Expected Trace to display '%s'. Got '%s'.", test, ret)
	}

	SetDebug()

	if ret := Trace("blabla %s", "toto"); ret != test {
		t.Errorf("Expected Trace to display '%s'. Got '%s'.", test, ret)
	}

	SetWarning()

	test = ""
	if ret := Trace("blabla %s", "toto"); ret != test {
		t.Errorf("Expected Trace to display '%s'. Got '%s'.", test, ret)
	}

	SetError()

	if ret := Trace("blabla %s", "toto"); ret != test {
		t.Errorf("Expected Trace to display '%s'. Got '%s'.", test, ret)
	}

	SetFatalError()

	if ret := Trace("blabla %s", "toto"); ret != test {
		t.Errorf("Expected Trace to display '%s'. Got '%s'.", test, ret)
	}
}

func TestWarning(t *testing.T) {
	t.Log("Expect warning to display in appropriate mode.")

	SetDebugLevel(5)

	test := "WARNING ! forjj-modules/trace.TestWarning: blabla toto"
	if ret := Warning("blabla %s", "toto"); ret != test {
		t.Errorf("Expected Warning to display '%s'. Got '%s'.", test, ret)
	}

	SetDebug()

	if ret := Warning("blabla %s", "toto"); ret != test {
		t.Errorf("Expected Warning to display '%s'. Got '%s'.", test, ret)
	}

	SetWarning()

	if ret := Warning("blabla %s", "toto"); ret != test {
		t.Errorf("Expected Warning to display '%s'. Got '%s'.", test, ret)
	}

	SetError()

	test = ""
	if ret := Warning("blabla %s", "toto"); ret != test {
		t.Errorf("Expected Warning to display '%s'. Got '%s'.", test, ret)
	}

	SetFatalError()

	if ret := Warning("blabla %s", "toto"); ret != test {
		t.Errorf("Expected Warning to display '%s'. Got '%s'.", test, ret)
	}
}

func TestError(t *testing.T) {
	t.Log("Expect Error to display in appropriate mode.")

	SetDebugLevel(5)

	test := "ERROR ! forjj-modules/trace.TestError: blabla toto"
	if ret := Error("blabla %s", "toto"); ret != test {
		t.Errorf("Expected Error to display '%s'. Got '%s'.", test, ret)
	}

	SetDebug()

	if ret := Error("blabla %s", "toto"); ret != test {
		t.Errorf("Expected Error to display '%s'. Got '%s'.", test, ret)
	}

	SetWarning()

	if ret := Error("blabla %s", "toto"); ret != test {
		t.Errorf("Expected Error to display '%s'. Got '%s'.", test, ret)
	}

	SetError()

	if ret := Error("blabla %s", "toto"); ret != test {
		t.Errorf("Expected Error to display '%s'. Got '%s'.", test, ret)
	}

	SetFatalError()

	test = ""
	if ret := Error("blabla %s", "toto"); ret != test {
		t.Errorf("Expected Error to display '%s'. Got '%s'.", test, ret)
	}
}

func TestFatalError(t *testing.T) {
	t.Log("Expect FatalError to display in appropriate mode.")

	SetDebugLevel(5)

	test := "FATAL ERROR ! forjj-modules/trace.TestFatalError: blabla toto"
	if ret := FatalError("blabla %s", "toto"); ret != test {
		t.Errorf("Expected FatalError to display '%s'. Got '%s'.", test, ret)
	}

	SetDebug()

	if ret := FatalError("blabla %s", "toto"); ret != test {
		t.Errorf("Expected FatalError to display '%s'. Got '%s'.", test, ret)
	}

	SetWarning()

	if ret := FatalError("blabla %s", "toto"); ret != test {
		t.Errorf("Expected FatalError to display '%s'. Got '%s'.", test, ret)
	}

	SetError()

	if ret := FatalError("blabla %s", "toto"); ret != test {
		t.Errorf("Expected FatalError to display '%s'. Got '%s'.", test, ret)
	}

	SetFatalError()

	if ret := FatalError("blabla %s", "toto"); ret != test {
		t.Errorf("Expected FatalError to display '%s'. Got '%s'.", test, ret)
	}
}

func TestGOTRACE(t *testing.T) {
	t.Log("Expect Module init to detect GOTRACE and set debug mode appropriately.")

	type valTest struct {
		value string
		mode  int
	}

	values := []valTest{
		{value: "", mode: warningMode},
		{value: "true", mode: debugMode},
		{value: "debug", mode: debugMode},
		{value: "4", mode: debugMode + 4},
		{value: "info", mode: infoMode},
		{value: "warning", mode: warningMode},
		{value: "error", mode: errorMode},
		{value: "fatal", mode: fatalMode},
	}

	for _, value := range values {
		os.Setenv("GOTRACE", value.value)
		internalDebug.debug = 0            // Force it.
		internalDebug.defaultDebug = false // Force it.
		internalDebug.init()

		if v := internalDebug.debug; v != value.mode {
			t.Errorf("Expected init with GOTRACE='%s' set to %s. Got %s.",
				value.value, internalDebug.prefix(value.mode), internalDebug.prefix(v))
		} else if internalDebug.defaultDebug && value.value != "" {
			t.Errorf("Expected Default debug to '%t'. Got '%t'.", false, internalDebug.defaultDebug)
		}

	}

}
