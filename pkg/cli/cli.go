package cli

import (
	"context"
	"errors"
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

type ArgConfig struct {
	Name       string
	Repeatable bool
	ShortName  rune
	Type       ArgType
	UsageText  string
	Value      interface{}
}

type CommandConfig struct {
	Args        []*ArgConfig
	Context     context.Context
	Name        string
	Run         CommandRunFunc
	Subcommands []*CommandConfig
}

type parsedArg struct {
	name  string
	value string
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
	return &CommandConfig{
		Args:        c.configureArgs(),
		Context:     c.ctx,
		Name:        c.name,
		Run:         c.run,
		Subcommands: c.configureSubcommands(),
	}
}

func (c *commandConfigurer) newCommandArg(n string, s rune, u string, r bool) *commandArg {
	return &commandArg{
		name:       n,
		shortName:  s,
		usageText:  u,
		repeatable: r,
	}
}

func (c *commandConfigurer) configureArgs() []*ArgConfig {
	var argConfigs []*ArgConfig

	argConfigs = append(argConfigs, c.configureBoolArgs()...)
	argConfigs = append(argConfigs, c.configureIntArgs()...)
	argConfigs = append(argConfigs, c.configureInt64Args()...)
	argConfigs = append(argConfigs, c.configureStringArgs()...)
	argConfigs = append(argConfigs, c.configureUintArgs()...)
	argConfigs = append(argConfigs, c.configureUint64Args()...)

	return argConfigs
}

func (c *commandConfigurer) configureBoolArgs() []*ArgConfig {
	var boolArgConfigs []*ArgConfig

	for _, arg := range c.args.boolArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue, Bool)
		boolArgConfigs = append(boolArgConfigs, argConfig)
	}

	return boolArgConfigs
}

func (c *commandConfigurer) configureIntArgs() []*ArgConfig {
	var intArgConfigs []*ArgConfig

	for _, arg := range c.args.intArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue, Int)
		intArgConfigs = append(intArgConfigs, argConfig)
	}

	return intArgConfigs
}

func (c *commandConfigurer) configureInt64Args() []*ArgConfig {
	var int64ArgConfigs []*ArgConfig

	for _, arg := range c.args.intArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue, Int64)
		int64ArgConfigs = append(int64ArgConfigs, argConfig)
	}

	return int64ArgConfigs
}

func (c *commandConfigurer) configureStringArgs() []*ArgConfig {
	var stringArgConfigs []*ArgConfig

	for _, arg := range c.args.intArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue, String)
		stringArgConfigs = append(stringArgConfigs, argConfig)
	}

	return stringArgConfigs
}

func (c *commandConfigurer) configureUintArgs() []*ArgConfig {
	var uintArgConfigs []*ArgConfig

	for _, arg := range c.args.intArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue, Uint)
		uintArgConfigs = append(uintArgConfigs, argConfig)
	}

	return uintArgConfigs
}

func (c *commandConfigurer) configureUint64Args() []*ArgConfig {
	var uint64ArgConfigs []*ArgConfig

	for _, arg := range c.args.intArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue, Uint64)
		uint64ArgConfigs = append(uint64ArgConfigs, argConfig)
	}

	return uint64ArgConfigs
}

func (c *commandConfigurer) configureCommandArgType(a *commandArg, v interface{}, d interface{}, t ArgType) *ArgConfig {
	argValue := v

	if argValue == nil {
		argValue = &d
	}

	return &ArgConfig{
		Name:       a.name,
		Repeatable: a.repeatable,
		ShortName:  a.shortName,
		Type:       t,
		UsageText:  a.usageText,
		Value:      argValue,
	}
}

func (c *commandConfigurer) configureSubcommands() []*CommandConfig {
	var subcommandConfigs []*CommandConfig

	for _, subCmd := range c.subcommands {
		subcommandConfigs = append(subcommandConfigs, subCmd.Configure())
	}

	return subcommandConfigs
}

func NewParser(a ArgSyntax, d DuplicateSubcommands) Parser {
	return &parser{
		argSyntax:      a,
		dupSubcmd:      d,
		parsedCommands: []*parsedCommand{},
	}
}

func (p *parser) Parse(c *CommandConfig) error {
	args := os.Args[1:]
	rootCmd := p.parseRootCmd(c)
	p.parsedCommands = append(p.parsedCommands, rootCmd)
	p.mapSubcommands(args, c, rootCmd)

	for _, cmd := range p.parsedCommands {
		if argErr := p.parseArgs(cmd); argErr != nil {
			return argErr
		}
	}

	return nil
}

func (p *parser) mapSubcommands(a []string, c *CommandConfig, l *parsedCommand) {
	if len(a) == 0 || l == nil {
		return
	}

	arg := a[0]
	argMapped := false
	var lastParsedCmd *parsedCommand

	for _, cmd := range c.Subcommands {
		if arg == cmd.Name {
			parsedCmd := &parsedCommand{
				argConfigs: cmd.Args,
				context:    cmd.Context,
				name:       cmd.Name,
				parentCmd:  c.Name,
				run:        cmd.Run,
			}
			p.parsedCommands = append(p.parsedCommands, parsedCmd)
			lastParsedCmd = parsedCmd
			argMapped = true

			break
		}

		if l.name == cmd.Name && len(cmd.Subcommands) > 0 {
			p.mapSubcommands(a, cmd, l)
		}
	}

	if lastParsedCmd == nil {
		lastParsedCmd = l
	}

	if !argMapped {
		lastParsedCmd.args = append(lastParsedCmd.args, arg)
	}

	if len(a) == 1 {
		return
	}

	p.mapSubcommands(a[1:], c, lastParsedCmd)
}

func (p *parser) parseRootCmd(c *CommandConfig) *parsedCommand {
	return &parsedCommand{
		args:       []string{},
		argConfigs: c.Args,
		context:    c.Context,
		run:        c.Run,
	}
}

func (p *parser) parseArgs(c *parsedCommand) error {
	switch p.argSyntax {
	case Gnu:
		return p.parseGnuArgs(c.args)
	case GoFlag:
		return p.parseGoFlagArgs(c.args)
	case Posix:
		return p.parsePosixArgs(c)
	default:
		return errors.New("unsupported ArgSyntax")
	}
}

func (p *parser) parseGnuArgs(a []string) error {
	panic("implement me")
}

func (p *parser) parseGoFlagArgs(a []string) error {
	panic("implement me")
}

func (p *parser) parsePosixArgs(c *parsedCommand) error {
	if len(c.args) == 0 {
		return nil
	}

	var lastParsedArg string
	var operands []string
	var parsedArgs []*parsedArg
	terminated := false

	for argIndex, arg := range c.args {
		if argIndex == 0 && !strings.HasPrefix(arg, "-") {
			return errors.New("invalid POSIX option: " + arg)
		}

		if strings.HasPrefix(arg, "-") && len(arg) > 0 {
			option := strings.TrimPrefix(arg, "-")

			for charIndex, char := range option {
				optName := string(char)

				if charIndex+1 < len(option) && char == '-' {
					return errors.New("invalid POSIX option: " + optName)
				}

				if char == '-' {
					terminated = true

					break
				}

				parsedArgs = append(parsedArgs, &parsedArg{
					name:  optName,
					value: "",
				})
				lastParsedArg = optName
			}

			continue
		}

		if terminated {
			operands = append(operands, arg)

			continue
		}

		if lastParsedArg != "" {
			for _, pArg := range parsedArgs {
				if pArg.name == lastParsedArg {
					pArg.value = arg

					break
				}
			}

			continue
		}

		return errors.New("invalid POSIX option: " + arg)
	}

	c.parsedArgs = parsedArgs
	c.operands = operands

	return nil
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
	cmdConfig := r.cmd.Configure()
	parseErr := r.parser.Parse(cmdConfig)

	if parseErr != nil {
		return parseErr
	}

	if cmdConfig.Run != nil {
		cmdConfig.Run(cmdConfig.Context)
	}

	for _, subCmd := range cmdConfig.Subcommands {
		subCmd.Run(subCmd.Context)
	}

	return nil
}
