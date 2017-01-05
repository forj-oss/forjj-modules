package cli

import (
    "github.com/alecthomas/kingpin"
    "testing"
)

func mustPanic(t *testing.T, f func()) {
    defer func() {
        if err := recover(); err == nil {
            t.Errorf("Panic expected: No panic returned", err)
        }
    }()

    f()
}

func mustNotPanic(t *testing.T, f func()) {
    defer func() {
        if err := recover(); err != nil {
            t.Errorf("Panic expected: No panic returned", err)
        }
    }()

    f()
}

func TestNewForjCli(t *testing.T) {
    var app *kingpin.Application

    t.Log("Expect an exception if the App is nil.")
    mustPanic(t, func() {
        NewForjCli(app)
    })

    t.Log("Expect application to be registered.")
    mustNotPanic(t, func() {
        app := kingpin.New("test", "test")
        c := NewForjCli(app)
        if c.App != app {
            t.Fail()
        }
    })

}
