package cli

import (
	"context"
	"io"
	"math"
	"strings"
)

func NewCommand(n string, c context.Context) CommandBuilder {
	return &commandBuilder{
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
		subcommands: []CommandBuilder{},
	}
}

func (c *commandBuilder) AddRunFunc(r RunFunc) {
	c.run = r
}

func (c *commandBuilder) AddSubcommand(cmd ...CommandBuilder) {
	c.subcommands = append(c.subcommands, cmd...)
}

func (c *commandBuilder) AddBoolArg(p *bool, a *ArgDefinition) {
	c.args.boolArgs = append(c.args.boolArgs, &boolArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandBuilder) AddFloat64Arg(p *float64, a *ArgDefinition) {
	c.args.float64Args = append(c.args.float64Args, &float64Arg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandBuilder) AddFloat64ListArg(p *[]float64, a *ArgDefinition) {
	c.args.float64ListArgs = append(c.args.float64ListArgs, &float64ListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandBuilder) AddIntArg(p *int, a *ArgDefinition) {
	c.args.intArgs = append(c.args.intArgs, &intArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandBuilder) AddIntListArg(p *[]int, a *ArgDefinition) {
	c.args.intListArgs = append(c.args.intListArgs, &intListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandBuilder) AddInt64Arg(p *int64, a *ArgDefinition) {
	c.args.int64Args = append(c.args.int64Args, &int64Arg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandBuilder) AddInt64ListArg(p *[]int64, a *ArgDefinition) {
	c.args.int64ListArgs = append(c.args.int64ListArgs, &int64ListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandBuilder) AddStringArg(p *string, a *ArgDefinition) {
	c.args.stringArgs = append(c.args.stringArgs, &stringArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandBuilder) AddStringListArg(p *[]string, a *ArgDefinition) {
	c.args.stringListArgs = append(c.args.stringListArgs, &stringListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandBuilder) AddUintArg(p *uint, a *ArgDefinition) {
	c.args.uintArgs = append(c.args.uintArgs, &uintArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandBuilder) AddUintListArg(p *[]uint, a *ArgDefinition) {
	c.args.uintListArgs = append(c.args.uintListArgs, &uintListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandBuilder) AddUint64Arg(p *uint64, a *ArgDefinition) {
	c.args.uint64Args = append(c.args.uint64Args, &uint64Arg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandBuilder) AddUint64ListArg(p *[]uint64, a *ArgDefinition) {
	c.args.uint64ListArgs = append(c.args.uint64ListArgs, &uint64ListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandBuilder) Build() *command {
	argConfigs := c.configureArgs()
	subcommands := c.configureSubcommands()

	command := &command{
		Args:        argConfigs,
		Context:     c.ctx,
		HelpFunc:    c.configureHelpFunc(argConfigs, subcommands),
		Name:        c.name,
		Run:         c.run,
		Subcommands: subcommands,
	}

	for _, subCmd := range command.Subcommands {
		subCmd.Parent = command
	}

	return command
}

func (c *commandBuilder) configureArgs() []*argConfig {
	var argConfigs []*argConfig
	helpArgConfigExists := false
	versionArgExists := false

	argConfigs = append(argConfigs, c.configureBoolArgs()...)
	argConfigs = append(argConfigs, c.configureFloat64Args()...)
	argConfigs = append(argConfigs, c.configureFloat64ListArgs()...)
	argConfigs = append(argConfigs, c.configureIntArgs()...)
	argConfigs = append(argConfigs, c.configureIntListArgs()...)
	argConfigs = append(argConfigs, c.configureInt64Args()...)
	argConfigs = append(argConfigs, c.configureInt64ListArgs()...)
	argConfigs = append(argConfigs, c.configureStringArgs()...)
	argConfigs = append(argConfigs, c.configureStringListArgs()...)
	argConfigs = append(argConfigs, c.configureUintArgs()...)
	argConfigs = append(argConfigs, c.configureUintListArgs()...)
	argConfigs = append(argConfigs, c.configureUint64Args()...)
	argConfigs = append(argConfigs, c.configureUint64ListArgs()...)

	for _, argConfig := range argConfigs {
		if argConfig.Name == "help" || argConfig.Name == "h" || argConfig.ShortName == 'h' {
			helpArgConfigExists = true
		}

		if argConfig.Name == "version" || argConfig.Name == "v" || argConfig.ShortName == 'v' {
			versionArgExists = true
		}
	}

	if !helpArgConfigExists {
		val := true
		argConfigs = append(argConfigs, &argConfig{
			Name:       "help",
			Repeatable: true,
			ShortName:  'h',
			UsageText:  "display usage information for this command",
			Value:      &val,
		})
	}

	if !versionArgExists {
		val := true
		argConfigs = append(argConfigs, &argConfig{
			Name:       "version",
			Repeatable: true,
			ShortName:  'v',
			UsageText:  "display the version for the utility",
			Value:      &val,
		})
	}

	return argConfigs
}

func (c *commandBuilder) configureHelpFunc(a []*argConfig, s []*command) func(s ArgSyntax, w io.Writer) error {
	return func(syntax ArgSyntax, w io.Writer) error {
		_, err := w.Write([]byte(c.getHelpTemplate(a, s, syntax)))

		return err
	}
}

func (c *commandBuilder) getHelpTemplate(a []*argConfig, s []*command, syntax ArgSyntax) string {
	var helpBuilder strings.Builder
	longestArgLine := float64(0)
	var argLines [][]string

	helpBuilder.WriteString(`Usage:
    ` + c.name)

	if len(s) > 0 {
		helpBuilder.WriteString(` [command]

Commands:
    `)
	}

	for _, cmd := range s {
		helpBuilder.WriteString(cmd.Name + `
    `)
	}

	if len(a) > 0 {
		helpBuilder.WriteString(`
Options:
    `)
	}

	for _, arg := range a {
		argLine := ""

		switch syntax {
		case GNU:
			if arg.ShortName > 0 {
				argLine = "-" + string(arg.ShortName) + ", "
			}

			if arg.Name != "" {
				if argLine == "" {
					argLine = "-" + string(arg.Name[0]) + ", "
				}

				argLine += "--" + arg.Name
			}

			strings.TrimSuffix(argLine, ", ")
		case POSIX:
			if arg.ShortName > 0 {
				argLine = "-" + string(arg.ShortName) + ", "
			}

			if argLine == "" && arg.Name != "" {
				argLine = "-" + string(arg.Name[0])
			}

			strings.TrimSuffix(argLine, ", ")
		}

		longestArgLine = math.Max(float64(len(argLine)), longestArgLine)
		argLines = append(argLines, []string{argLine, arg.UsageText})
	}

	for i, argLine := range argLines {
		helpBuilder.WriteString(argLine[0])
		helpBuilder.WriteString(strings.Repeat(" ", int(longestArgLine)-len(argLine[0])+4))
		helpBuilder.WriteString(argLine[1] + `
`)

		if i < len(argLines)-1 {
			helpBuilder.WriteString(strings.Repeat(" ", 4))
		}
	}

	return helpBuilder.String()
}

func (c *commandBuilder) configureBoolArgs() []*argConfig {
	var boolArgConfigs []*argConfig

	for _, arg := range c.args.boolArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		boolArgConfigs = append(boolArgConfigs, argConfig)
	}

	return boolArgConfigs
}

func (c *commandBuilder) configureFloat64Args() []*argConfig {
	var float64ArgConfigs []*argConfig

	for _, arg := range c.args.float64Args {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		float64ArgConfigs = append(float64ArgConfigs, argConfig)
	}

	return float64ArgConfigs
}

func (c *commandBuilder) configureFloat64ListArgs() []*argConfig {
	var float64ListArgConfigs []*argConfig

	for _, arg := range c.args.float64ListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		float64ListArgConfigs = append(float64ListArgConfigs, argConfig)
	}

	return float64ListArgConfigs
}

func (c *commandBuilder) configureIntArgs() []*argConfig {
	var intArgConfigs []*argConfig

	for _, arg := range c.args.intArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		intArgConfigs = append(intArgConfigs, argConfig)
	}

	return intArgConfigs
}

func (c *commandBuilder) configureIntListArgs() []*argConfig {
	var intListArgConfigs []*argConfig

	for _, arg := range c.args.intListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		intListArgConfigs = append(intListArgConfigs, argConfig)
	}

	return intListArgConfigs
}

func (c *commandBuilder) configureInt64Args() []*argConfig {
	var int64ArgConfigs []*argConfig

	for _, arg := range c.args.int64Args {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		int64ArgConfigs = append(int64ArgConfigs, argConfig)
	}

	return int64ArgConfigs
}

func (c *commandBuilder) configureInt64ListArgs() []*argConfig {
	var int64ListArgConfigs []*argConfig

	for _, arg := range c.args.int64ListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		int64ListArgConfigs = append(int64ListArgConfigs, argConfig)
	}

	return int64ListArgConfigs
}

func (c *commandBuilder) configureStringArgs() []*argConfig {
	var stringArgConfigs []*argConfig

	for _, arg := range c.args.stringArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		stringArgConfigs = append(stringArgConfigs, argConfig)
	}

	return stringArgConfigs
}

func (c *commandBuilder) configureStringListArgs() []*argConfig {
	var stringListArgConfigs []*argConfig

	for _, arg := range c.args.stringListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		stringListArgConfigs = append(stringListArgConfigs, argConfig)
	}

	return stringListArgConfigs
}

func (c *commandBuilder) configureUintArgs() []*argConfig {
	var uintArgConfigs []*argConfig

	for _, arg := range c.args.uintArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		uintArgConfigs = append(uintArgConfigs, argConfig)
	}

	return uintArgConfigs
}

func (c *commandBuilder) configureUintListArgs() []*argConfig {
	var uintListArgConfigs []*argConfig

	for _, arg := range c.args.uintListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		uintListArgConfigs = append(uintListArgConfigs, argConfig)
	}

	return uintListArgConfigs
}

func (c *commandBuilder) configureUint64Args() []*argConfig {
	var uint64ArgConfigs []*argConfig

	for _, arg := range c.args.uint64Args {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		uint64ArgConfigs = append(uint64ArgConfigs, argConfig)
	}

	return uint64ArgConfigs
}

func (c *commandBuilder) configureUint64ListArgs() []*argConfig {
	var uint64ListArgConfigs []*argConfig

	for _, arg := range c.args.uint64ListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		uint64ListArgConfigs = append(uint64ListArgConfigs, argConfig)
	}

	return uint64ListArgConfigs
}

func (c *commandBuilder) configureSubcommands() []*command {
	var subcommandConfigs []*command

	for _, subCmd := range c.subcommands {
		subcommandConfigs = append(subcommandConfigs, subCmd.Build())
	}

	return subcommandConfigs
}

func newCommandArg(a *ArgDefinition) *commandArg {
	if a == nil {
		return &commandArg{}
	}

	return &commandArg{
		name:       a.Name,
		shortName:  a.ShortName,
		usageText:  a.UsageText,
		repeatable: a.Repeatable,
		required:   a.Required,
	}
}

func newArgConfig(a *commandArg, v interface{}) *argConfig {
	return &argConfig{
		Name:       a.name,
		Repeatable: a.repeatable,
		Required:   a.required,
		ShortName:  a.shortName,
		UsageText:  a.usageText,
		Value:      v,
	}
}
