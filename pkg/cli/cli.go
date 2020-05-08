package cli

import (
	"context"
)

type ArgSyntax int

const (
	GNU ArgSyntax = iota
	GoFlag
	POSIX
)

type ErrorBehavior int

const (
	ExitOnError ErrorBehavior = iota
	ContinueOnError
)

type ArgDefinition struct {
	Name       string
	ShortName  rune
	UsageText  string
	Repeatable bool
	Required   bool
}

type commandArg struct {
	name       string
	shortName  rune
	usageText  string
	repeatable bool
	required   bool
}

type boolArg struct {
	*commandArg
	defaultValue bool
	value        *bool
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

type ParsedCommand struct {
	args        []string
	argConfigs  []*ArgConfig
	Context     context.Context
	Name        string
	Operands    []string
	parentCmd   string
	parsedArgs  []*parsedArg
	Run         CommandRunFunc
	Subcommands []*ParsedCommand
}

type argParserContext struct {
	argConfigs      []*ArgConfig
	lastParsedArg   map[string][]string
	operands        []string
	parsedArgs      []*parsedArg
	terminated      bool
	terminatorIndex int
}

type argParserRule func(a string, i int, c *argParserContext) (bool, error)

type argParserInit func(a []string) *argParserContext

type ArgAdder interface {
	AddBoolArg(p *bool, v bool, a *ArgDefinition)
	AddIntArg(p *int, v int, a *ArgDefinition)
	AddIntListArg(p *[]int, v []int, a *ArgDefinition)
	AddInt64Arg(p *int64, v int64, a *ArgDefinition)
	AddInt64ListArg(p *[]int64, v []int64, a *ArgDefinition)
	AddStringArg(p *string, v string, a *ArgDefinition)
	AddStringListArg(p *[]string, v []string, a *ArgDefinition)
	AddUintArg(p *uint, v uint, a *ArgDefinition)
	AddUintListArg(p *[]uint, v []uint, a *ArgDefinition)
	AddUint64Arg(p *uint64, v uint64, a *ArgDefinition)
	AddUint64ListArg(p *[]uint64, v []uint64, a *ArgDefinition)
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
	Parse(c *CommandConfig) (*ParsedCommand, error)
}

type parser struct {
	argSyntax      ArgSyntax
	parsedCommands []*ParsedCommand
}

type Runner interface {
	Run(p *ParsedCommand) error
}

type runner struct {
	errBehavior ErrorBehavior
	version     string
}
