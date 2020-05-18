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

type argConfig struct {
	Name       string
	Repeatable bool
	Required   bool
	ShortName  rune
	UsageText  string
	Value      interface{}
}

type HelpFunc func(s ArgSyntax, w io.Writer) error

type RunFunc func(ctx context.Context, o []string)

type commandConfig struct {
	Args        []*argConfig
	Context     context.Context
	HelpFunc    HelpFunc
	Name        string
	Parent      *commandConfig
	Operands    []string
	Run         RunFunc
	Subcommands []*commandConfig
}

type commandWalker struct {
	root *commandConfig
	path []*commandConfig
}

type parsedArg struct {
	bindVal  interface{}
	name     string
	rawArg   string
	required bool
	value    []string
}

type parsedCommand struct {
	args         []string
	argConfigs   []*argConfig
	config       *commandConfig
	Context      context.Context
	HelpFunc     HelpFunc
	HelpMode     bool
	Name         string
	Operands     []string
	parsedArgs   []*parsedArg
	Run          RunFunc
	Subcommands  []*parsedCommand
	Syntax       ArgSyntax
	VersionMode  bool
}

type argParserContext struct {
	argConfigs      []*argConfig
	lastParsedArg   *parsedArg
	operands        []string
	parsedArgs      []*parsedArg
	terminated      bool
	terminatorIndex int
}

type argParserRule func(a *string, i int, c *argParserContext) (bool, error)

type argParserInit func(a []string) *argParserContext

type ArgDefinition struct {
	Name       string
	ShortName  rune
	UsageText  string
	Repeatable bool
	Required   bool
}

type argAdder interface {
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
	AddSubcommand(c ...CommandConfigurer)
	AddRunFunc(r RunFunc)
	argAdder
	Configure() *commandConfig
}

type commandConfigurer struct {
	args        *commandArgs
	ctx         context.Context
	name        string
	run         RunFunc
	subcommands []CommandConfigurer
}

type Parser interface {
	Parse() (*parsedCommand, error)
}

type parser struct {
	argSyntax      ArgSyntax
	configurer     CommandConfigurer
	helpFunc       HelpFunc
	helpMode       bool
	parsedCommands []*parsedCommand
}

type Runner interface {
	Run() error
}

type runner struct {
	parser  Parser
	version string
	writer  io.Writer
}
