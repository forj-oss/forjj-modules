package clier

// Define interface against kingpin or kingpinMock (for go test)

type Applicationer interface {
	Flag(string, string) FlagClauser
	Arg(string, string) ArgClauser
	Command(string, string) CmdClauser
	IsNil() bool
	ParseContext([]string) (ParseContexter, error)
	Parse([]string) (string, error)
	Name() string
}

type FlagClauser interface {
	Stringer() string
	String() *string
	Bool() *bool
	Required() FlagClauser
	Short(rune) FlagClauser
	Hidden() FlagClauser
	Default(string) FlagClauser
	Envar(string) FlagClauser
	SetValue(Valuer) FlagClauser
}

type ArgClauser interface {
	Stringer() string
	String() *string
	Bool() *bool
	Required() ArgClauser
	Default(string) ArgClauser
	SetValue(Valuer) ArgClauser
	Envar(string) ArgClauser
}

type CmdClauser interface {
	Command(string, string) CmdClauser
	Flag(string, string) FlagClauser
	Arg(string, string) ArgClauser
	FullCommand() string
	IsEqualTo(CmdClauser) bool
}

type ParseContexter interface {
	GetArgValue(ArgClauser) (interface{}, bool)
	GetFlagValue(FlagClauser) (interface{}, bool)
	GetParam(string) (interface{}, string)
	SelectedCommands() []CmdClauser
	IsInvalidContext() bool
}

type Valuer interface {
	Set(string) error
	String() string
}
