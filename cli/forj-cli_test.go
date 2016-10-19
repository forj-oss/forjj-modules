package cli

import (
	"github.com/forj-oss/forjj-modules/cli/kingpinMock"
	"testing"
)

var app = kingpinMock.NewMock([]string{}, "")

const (
	w_f  = `([a-z]+[a-z0-9_-]*)`
	ft_f = `([A-Za-z0-9_ !:/.-]+)`
)

func mustPanic(t *testing.T, f func()) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("Panic expected: No panic returned.")
		}
	}()

	f()
}

func mustNotPanic(t *testing.T, f func()) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("Panic NOT expected: Panic returned. Message : '%s'", err)
		}
	}()

	f()
}

func TestNewForjCli(t *testing.T) {
	var app_nil *kingpinMock.Application

	t.Log("Expect an exception if the App is nil.")
	mustPanic(t, func() {
		NewForjCli(app_nil)
	})

	t.Log("Expect application to be registered.")
	mustNotPanic(t, func() {
		c := NewForjCli(app)
		if c.App != app {
			t.Fail()
		}
	})
}

func TestAddFieldListCapture(t *testing.T) {
	t.Log("Expect AddFieldListCapture to add capture list.")
	mustNotPanic(t, func() {
		c := NewForjCli(app)
		c.AddFieldListCapture("w", w_f)
		c.AddFieldListCapture("ft", ft_f)

		if v, found := c.filters["w"]; !found || v != w_f {
			t.Fail()
		}
		if v, found := c.filters["ft"]; !found || v != ft_f {
			t.Fail()
		}
	})

}

func TestAddAppFlag(t *testing.T) {

	t.Log("Expect AddAppFlag to create a Flag at App level.")

	mustNotPanic(t, func() {
		c := NewForjCli(app)
		c.AddFieldListCapture("w", w_f)
		c.AddFieldListCapture("ft", ft_f)

		if v, found := c.filters["w"]; !found || v != w_f {
			t.Fail()
		}
		if v, found := c.filters["ft"]; !found || v != ft_f {
			t.Fail()
		}
	})

}
