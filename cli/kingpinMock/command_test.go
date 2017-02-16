package kingpinMock

import (
	"testing"
)

func TestNewCmd(t *testing.T) {
	t.Log("NewArg creates a new Cmd with name and help")
	a := NewCmd("test", "help")
	if a.command != "test" {
		t.Errorf("name expected to be 'test'. Got %s", a.command)
	}
	if a.help != "help" {
		t.Errorf("help expected to be 'help'. Got %s", a.help)
	}
}

func TestCmdClause_Command(t *testing.T) {
	t.Log("NewArg().Command() creates a new Arg with name and help")
	c := NewCmd("test", "help")

	c2 := c.Command("test2", "help")

	if c2 == nil {
		t.Error("Expected having a new Arg object. Got nil.")
	}

	c2_found, found := c.cmds["test2"]
	if !found {
		t.Error("Expected having a new test2 in Command layer. Not found.")
	}

	if c2_found != c2 {
		t.Error("Expected object created is the one returned.")
	}

}

func TestCmdClause_Arg(t *testing.T) {
	t.Log("NewArg().Arg() creates a new Arg with name and help")
	a := NewCmd("test", "help")

	arg := a.Arg("arg1", "help")

	if arg == nil {
		t.Error("Expect having a new Arg object. Got nil.")
	}

	if _, found := a.args["arg1"]; !found {
		t.Error("Expect having a new Arg in App layer. Not found.")
	}

}

func TestCmdClause_Flag(t *testing.T) {
	t.Log("NewArg().Flag() creates a new Flag with name and help")
	a := NewCmd("test", "help")

	flag := a.Flag("flag1", "help")

	if flag == nil {
		t.Error("Expect having a new Arg object. Got nil.")
	}
	if _, found := a.flags["flag1"]; !found {
		t.Error("Expect having a new Arg in App layer. Not found.")
	}

}
