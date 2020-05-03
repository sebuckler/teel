package cli

import (
	"context"
)

type ArgSyntax int

const (
	Gnu ArgSyntax = iota
	GoFlag
	Posix
)

type DuplicateSubcommands int

const (
	Error DuplicateSubcommands = iota
	BreadthFirst
	DepthFirst
)

type ErrorBehavior int

const (
	ExitOnError ErrorBehavior = iota
	ContinueOnError
)

type commandArg struct {
	name       string
	shortName  rune
	usageText  string
	repeatable bool
}

type boolArg struct {
	*commandArg
	defaultValue bool
	value        *bool
}

type boolListArg struct {
	*commandArg
	defaultValue []bool
	value        *[]bool
}

type intArg struct {
	*commandArg
	defaultValue int
	value        *int
}

type intListArg struct {
	*commandArg
	defaultValue []int
	value        *[]int
}

type int64Arg struct {
	*commandArg
	defaultValue int64
	value        *int64
}

type int64ListArg struct {
	*commandArg
	defaultValue []int64
	value        *[]int64
}

type stringArg struct {
	*commandArg
	defaultValue string
	value        *string
}

type stringListArg struct {
	*commandArg
	defaultValue []string
	value        *[]string
}

type uintArg struct {
	*commandArg
	defaultValue uint
	value        *uint
}

type uintListArg struct {
	*commandArg
	defaultValue []uint
	value        *[]uint
}

type uint64Arg struct {
	*commandArg
	defaultValue uint64
	value        *uint64
}

type uint64ListArg struct {
	*commandArg
	defaultValue []uint64
	value        *[]uint64
}

type commandArgs struct {
	boolArgs       []*boolArg
	boolListArgs   []*boolListArg
	intArgs        []*intArg
	intListArgs    []*intListArg
	int64Args      []*int64Arg
	int64ListArgs  []*int64ListArg
	stringArgs     []*stringArg
	stringListArgs []*stringListArg
	uintArgs       []*uintArg
	uintListArgs   []*uintListArg
	uint64Args     []*uint64Arg
	uint64ListArgs []*uint64ListArg
}

type CommandRunFunc func(ctx context.Context, o []string)

type ArgConfig struct {
	Name       string
	Repeatable bool
	ShortName  rune
	UsageText  string
	Value      interface{}
}

type CommandConfig struct {
	Args        []*ArgConfig
	Context     context.Context
	Name        string
	Operands    []string
	Run         CommandRunFunc
	Subcommands []*CommandConfig
}

type parsedArg struct {
	bindVal interface{}
	name    string
	value   []string
}

type parsedCommand struct {
	args        []string
	argConfigs  []*ArgConfig
	context     context.Context
	name        string
	operands    []string
	parentCmd   string
	parsedArgs  []*parsedArg
	run         CommandRunFunc
	subcommands []*parsedCommand
}

type ArgAdder interface {
	AddBoolArg(n string, s rune, p *bool, v bool, u string, r bool)
	AddBoolListArg(n string, s rune, p *[]bool, v []bool, u string, r bool)
	AddIntArg(n string, s rune, p *int, v int, u string, r bool)
	AddIntListArg(n string, s rune, p *[]int, v []int, u string, r bool)
	AddInt64Arg(n string, s rune, p *int64, v int64, u string, r bool)
	AddInt64ListArg(n string, s rune, p *[]int64, v []int64, u string, r bool)
	AddStringArg(n string, s rune, p *string, v string, u string, r bool)
	AddStringListArg(n string, s rune, p *[]string, v []string, u string, r bool)
	AddUintArg(n string, s rune, p *uint, v uint, u string, r bool)
	AddUintListArg(n string, s rune, p *[]uint, v []uint, u string, r bool)
	AddUint64Arg(n string, s rune, p *uint64, v uint64, u string, r bool)
	AddUint64ListArg(n string, s rune, p *[]uint64, v []uint64, u string, r bool)
}

type CommandConfigurer interface {
	AddSubcommand(c CommandConfigurer)
	AddRunFunc(r CommandRunFunc)
	ArgAdder
	Configure() *CommandConfig
}

type commandConfigurer struct {
	args        *commandArgs
	ctx         context.Context
	name        string
	run         CommandRunFunc
	subcommands []CommandConfigurer
}

type Parser interface {
	Parse(c *CommandConfig) error
}

type parser struct {
	argSyntax      ArgSyntax
	dupSubcmd      DuplicateSubcommands
	parsedCommands []*parsedCommand
}

type Runner interface {
	Run() error
}

type runner struct {
	cmd         CommandConfigurer
	errBehavior ErrorBehavior
	parser      Parser
	version     string
}
