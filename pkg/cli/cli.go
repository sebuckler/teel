package cli

import (
	"context"
	"io"
)

type ArgSyntax int

const (
	GNU ArgSyntax = iota
	POSIX
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
	value *bool
}

type float64Arg struct {
	*commandArg
	value *float64
}

type float64ListArg struct {
	*commandArg
	value *[]float64
}

type intArg struct {
	*commandArg
	value *int
}

type intListArg struct {
	*commandArg
	value *[]int
}

type int64Arg struct {
	*commandArg
	value *int64
}

type int64ListArg struct {
	*commandArg
	value *[]int64
}

type stringArg struct {
	*commandArg
	value *string
}

type stringListArg struct {
	*commandArg
	value *[]string
}

type uintArg struct {
	*commandArg
	value *uint
}

type uintListArg struct {
	*commandArg
	value *[]uint
}

type uint64Arg struct {
	*commandArg
	value *uint64
}

type uint64ListArg struct {
	*commandArg
	value *[]uint64
}

type commandArgs struct {
	boolArgs        []*boolArg
	float64Args     []*float64Arg
	float64ListArgs []*float64ListArg
	intArgs         []*intArg
	intListArgs     []*intListArg
	int64Args       []*int64Arg
	int64ListArgs   []*int64ListArg
	stringArgs      []*stringArg
	stringListArgs  []*stringListArg
	uintArgs        []*uintArg
	uintListArgs    []*uintListArg
	uint64Args      []*uint64Arg
	uint64ListArgs  []*uint64ListArg
}

type CommandRunFunc func(ctx context.Context, o []string)

type ArgConfig struct {
	Name       string
	Repeatable bool
	Required   bool
	ShortName  rune
	UsageText  string
	Value      interface{}
}

type HelpFunc func(s ArgSyntax, w io.Writer) error

type CommandConfig struct {
	Args        []*ArgConfig
	Context     context.Context
	HelpFunc    HelpFunc
	Name        string
	Operands    []string
	Run         CommandRunFunc
	Subcommands []*CommandConfig
}

type parsedArg struct {
	bindVal  interface{}
	name     string
	rawArg   string
	required bool
	value    []string
}

type ParsedCommand struct {
	args        []string
	argConfigs  []*ArgConfig
	Context     context.Context
	HelpFunc    HelpFunc
	HelpMode    bool
	Name        string
	Operands    []string
	parentCmd   string
	parsedArgs  []*parsedArg
	Run         CommandRunFunc
	Subcommands []*ParsedCommand
	Syntax      ArgSyntax
	VersionMode bool
}

type argParserContext struct {
	argConfigs      []*ArgConfig
	lastParsedArg   *parsedArg
	operands        []string
	parsedArgs      []*parsedArg
	terminated      bool
	terminatorIndex int
}

type argParserRule func(a *string, i int, c *argParserContext) (bool, error)

type argParserInit func(a []string) *argParserContext

type ArgAdder interface {
	AddBoolArg(p *bool, a *ArgDefinition)
	AddFloat64Arg(p *float64, a *ArgDefinition)
	AddFloat64ListArg(p *[]float64, a *ArgDefinition)
	AddIntArg(p *int, a *ArgDefinition)
	AddIntListArg(p *[]int, a *ArgDefinition)
	AddInt64Arg(p *int64, a *ArgDefinition)
	AddInt64ListArg(p *[]int64, a *ArgDefinition)
	AddStringArg(p *string, a *ArgDefinition)
	AddStringListArg(p *[]string, a *ArgDefinition)
	AddUintArg(p *uint, a *ArgDefinition)
	AddUintListArg(p *[]uint, a *ArgDefinition)
	AddUint64Arg(p *uint64, a *ArgDefinition)
	AddUint64ListArg(p *[]uint64, a *ArgDefinition)
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
	helpFunc       HelpFunc
	helpMode       bool
	parsedCommands []*ParsedCommand
}

type Runner interface {
	Run(p *ParsedCommand) error
}

type runner struct {
	version string
	writer  io.Writer
}
