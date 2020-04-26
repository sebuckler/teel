package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
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

type ArgType int

const (
	Bool ArgType = iota
	Int
	Int64
	String
	Uint
	Uint64
)

type ErrorBehavior int

const (
	ExitOnError ErrorBehavior = iota
	ContinueOnError
)

type ArgConfig struct {
	Name       string
	Repeatable bool
	Type       ArgType
}

type CommandConfig struct {
	Args        []*ArgConfig
	Name        string
	Subcommands []*CommandConfig
}

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

type intArg struct {
	*commandArg
	defaultValue int
	value        *int
}

type int64Arg struct {
	*commandArg
	defaultValue int64
	value        *int64
}

type stringArg struct {
	*commandArg
	defaultValue string
	value        *string
}

type uintArg struct {
	*commandArg
	defaultValue uint
	value        *uint
}

type uint64Arg struct {
	*commandArg
	defaultValue uint64
	value        *uint64
}

type commandArgs struct {
	boolArgs   []*boolArg
	intArgs    []*intArg
	int64Args  []*int64Arg
	stringArgs []*stringArg
	uintArgs   []*uintArg
	uint64Args []*uint64Arg
}

type CommandRunFunc func(ctx context.Context)

type ArgAdder interface {
	AddBoolArg(n string, s rune, p *bool, v bool, u string, r bool)
	AddIntArg(n string, s rune, p *int, v int, u string, r bool)
	AddInt64Arg(n string, s rune, p *int64, v int64, u string, r bool)
	AddStringArg(n string, s rune, p *string, v string, u string, r bool)
	AddUintArg(n string, s rune, p *uint, v uint, u string, r bool)
	AddUint64Arg(n string, s rune, p *uint64, v uint64, u string, r bool)
}

type CommandConfigurer interface {
	AddOperand()
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
	argSyntax ArgSyntax
	dupSubcmd DuplicateSubcommands
	registry  context.Context
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

func NewCommand(n string, c context.Context) CommandConfigurer {
	return &commandConfigurer{
		args: &commandArgs{
			boolArgs:   []*boolArg{},
			intArgs:    []*intArg{},
			int64Args:  []*int64Arg{},
			stringArgs: []*stringArg{},
			uintArgs:   []*uintArg{},
			uint64Args: []*uint64Arg{},
		},
		ctx:         c,
		name:        n,
		subcommands: []CommandConfigurer{},
	}
}

func (c *commandConfigurer) AddOperand() {
	panic("implement me")
}

func (c *commandConfigurer) AddRunFunc(r CommandRunFunc) {
	c.run = r
}

func (c *commandConfigurer) AddSubcommand(cmd CommandConfigurer) {
	c.subcommands = append(c.subcommands, cmd)
}

func (c *commandConfigurer) AddBoolArg(n string, s rune, p *bool, v bool, u string, r bool) {
	c.args.boolArgs = append(c.args.boolArgs, &boolArg{
		commandArg:   c.newCommandArg(n, s, u, r),
		defaultValue: v,
		value:        p,
	})
}

func (c *commandConfigurer) AddIntArg(n string, s rune, p *int, v int, u string, r bool) {
	c.args.intArgs = append(c.args.intArgs, &intArg{
		commandArg:   c.newCommandArg(n, s, u, r),
		defaultValue: v,
		value:        p,
	})
}

func (c *commandConfigurer) AddInt64Arg(n string, s rune, p *int64, v int64, u string, r bool) {
	c.args.int64Args = append(c.args.int64Args, &int64Arg{
		commandArg:   c.newCommandArg(n, s, u, r),
		defaultValue: v,
		value:        p,
	})
}

func (c *commandConfigurer) AddStringArg(n string, s rune, p *string, v string, u string, r bool) {
	c.args.stringArgs = append(c.args.stringArgs, &stringArg{
		commandArg:   c.newCommandArg(n, s, u, r),
		defaultValue: v,
		value:        p,
	})
}

func (c *commandConfigurer) AddUintArg(n string, s rune, p *uint, v uint, u string, r bool) {
	c.args.uintArgs = append(c.args.uintArgs, &uintArg{
		commandArg:   c.newCommandArg(n, s, u, r),
		defaultValue: v,
		value:        p,
	})
}

func (c *commandConfigurer) AddUint64Arg(n string, s rune, p *uint64, v uint64, u string, r bool) {
	c.args.uint64Args = append(c.args.uint64Args, &uint64Arg{
		commandArg:   c.newCommandArg(n, s, u, r),
		defaultValue: v,
		value:        p,
	})
}

func (c *commandConfigurer) Configure() *CommandConfig {
	panic("implement me")
}

func (c *commandConfigurer) newCommandArg(n string, s rune, u string, r bool) *commandArg {
	return &commandArg{
		name:       n,
		shortName:  s,
		usageText:  u,
		repeatable: r,
	}
}

func NewParser(a ArgSyntax, d DuplicateSubcommands) Parser {
	return &parser{
		argSyntax: a,
		dupSubcmd: d,
	}
}

func (p *parser) Parse(c *CommandConfig) error {
	args := os.Args[1:]

	switch p.argSyntax {
	case Gnu:
		return p.parseGnuArgs(args)
	case GoFlag:
		return p.parseGoFlagArgs(args)
	case Posix:
		return p.parsePosixArgs(args)
	default:
		return errors.New("unsupported ArgSyntax")
	}
}

func (p *parser) mapSubcommands() {
	panic("implement me")
}

func (p *parser) parseGnuArgs(a []string) error {
	panic("implement me")
}

func (p *parser) parseGoFlagArgs(a []string) error {
	panic("implement me")
}

func (p *parser) parsePosixArgs(a []string) error {
	if len(a) == 0 {
		return nil
	}
	fmt.Println(a)
	for _, arg := range a {
		if strings.HasPrefix(arg, "-") {
			args, argErr := p.parsePosixOption(arg)

			if argErr != nil {
				return argErr
			}

			if len(args) == 0 {
				return nil
			}

			fmt.Println(args)
		}

		p.parseSubcommand(arg)
	}

	return nil
}

func (p *parser) parsePosixOption(a string) ([]string, error) {
	if strings.HasPrefix(a, "--") && len(a) > 2 {
		return nil, errors.New("invalid POSIX argument syntax")
	}

	if strings.HasPrefix(a, "--") && len(a) == 2 {
		return nil, nil
	}

	posixArg := strings.TrimPrefix(a, "-")

	if posixArg == "W" {
		return nil, errors.New("-W is a reserved vendor argument")
	}

	if len(posixArg) == 1 {
		return []string{posixArg}, nil
	}

	var argChain []string
	argValueSplit := strings.Split(posixArg, "=")

	for _, arg := range argValueSplit[0] {
		argChar := string(arg)

		if argChar == "-" {
			return nil, errors.New("invalid POSIX argument syntax")
		}

		argChain = append(argChain, argChar)
	}

	return argChain, nil
}

func (p *parser) parseSubcommand(a string) error {
	panic("implement me")
}

func NewRunner(c CommandConfigurer, p Parser, v string, e ErrorBehavior) Runner {
	return &runner{
		cmd:         c,
		errBehavior: e,
		parser:      p,
		version:     v,
	}
}

func (r *runner) Run() error {
	parseErr := r.parser.Parse(r.cmd.Configure())

	if parseErr != nil {
		return parseErr
	}

	return nil
}
