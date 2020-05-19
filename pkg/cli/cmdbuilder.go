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

func (b *commandBuilder) AddRunFunc(r RunFunc) {
	b.run = r
}

func (b *commandBuilder) AddSubcommand(cmd ...CommandBuilder) {
	b.subcommands = append(b.subcommands, cmd...)
}

func (b *commandBuilder) AddBoolArg(p *bool, a *ArgDefinition) {
	b.args.boolArgs = append(b.args.boolArgs, &boolArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (b *commandBuilder) AddFloat64Arg(p *float64, a *ArgDefinition) {
	b.args.float64Args = append(b.args.float64Args, &float64Arg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (b *commandBuilder) AddFloat64ListArg(p *[]float64, a *ArgDefinition) {
	b.args.float64ListArgs = append(b.args.float64ListArgs, &float64ListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (b *commandBuilder) AddIntArg(p *int, a *ArgDefinition) {
	b.args.intArgs = append(b.args.intArgs, &intArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (b *commandBuilder) AddIntListArg(p *[]int, a *ArgDefinition) {
	b.args.intListArgs = append(b.args.intListArgs, &intListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (b *commandBuilder) AddInt64Arg(p *int64, a *ArgDefinition) {
	b.args.int64Args = append(b.args.int64Args, &int64Arg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (b *commandBuilder) AddInt64ListArg(p *[]int64, a *ArgDefinition) {
	b.args.int64ListArgs = append(b.args.int64ListArgs, &int64ListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (b *commandBuilder) AddStringArg(p *string, a *ArgDefinition) {
	b.args.stringArgs = append(b.args.stringArgs, &stringArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (b *commandBuilder) AddStringListArg(p *[]string, a *ArgDefinition) {
	b.args.stringListArgs = append(b.args.stringListArgs, &stringListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (b *commandBuilder) AddUintArg(p *uint, a *ArgDefinition) {
	b.args.uintArgs = append(b.args.uintArgs, &uintArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (b *commandBuilder) AddUintListArg(p *[]uint, a *ArgDefinition) {
	b.args.uintListArgs = append(b.args.uintListArgs, &uintListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (b *commandBuilder) AddUint64Arg(p *uint64, a *ArgDefinition) {
	b.args.uint64Args = append(b.args.uint64Args, &uint64Arg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (b *commandBuilder) AddUint64ListArg(p *[]uint64, a *ArgDefinition) {
	b.args.uint64ListArgs = append(b.args.uint64ListArgs, &uint64ListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (b *commandBuilder) Build() *command {
	argConfigs := b.configureArgs()
	subcommands := b.configureSubcommands()

	command := &command{
		Args:        argConfigs,
		Context:     b.ctx,
		HelpFunc:    b.configureHelpFunc(argConfigs),
		Name:        b.name,
		Run:         b.run,
		Subcommands: subcommands,
	}

	for _, subCmd := range command.Subcommands {
		subCmd.Parent = command
	}

	return command
}

func (b *commandBuilder) configureArgs() []*argConfig {
	var argConfigs []*argConfig
	helpArgConfigExists := false
	versionArgExists := false

	argConfigs = append(argConfigs, b.configureBoolArgs()...)
	argConfigs = append(argConfigs, b.configureFloat64Args()...)
	argConfigs = append(argConfigs, b.configureFloat64ListArgs()...)
	argConfigs = append(argConfigs, b.configureIntArgs()...)
	argConfigs = append(argConfigs, b.configureIntListArgs()...)
	argConfigs = append(argConfigs, b.configureInt64Args()...)
	argConfigs = append(argConfigs, b.configureInt64ListArgs()...)
	argConfigs = append(argConfigs, b.configureStringArgs()...)
	argConfigs = append(argConfigs, b.configureStringListArgs()...)
	argConfigs = append(argConfigs, b.configureUintArgs()...)
	argConfigs = append(argConfigs, b.configureUintListArgs()...)
	argConfigs = append(argConfigs, b.configureUint64Args()...)
	argConfigs = append(argConfigs, b.configureUint64ListArgs()...)

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

func (b *commandBuilder) configureHelpFunc(a []*argConfig) HelpFunc {
	return func(c *command, s ArgSyntax, w io.Writer) error {
		_, err := w.Write([]byte(b.getHelpTemplate(c, s, a)))

		return err
	}
}

func (b *commandBuilder) getHelpTemplate(c *command, s ArgSyntax, a []*argConfig) string {
	var helpBuilder strings.Builder
	parent := c.Parent
	var parentNames []string
	var parentUsage string
	longestArgLine := float64(0)
	var argLines [][]string

	for parent != nil {
		parentNames = append(parentNames, c.Parent.Name)
		parent = parent.Parent
	}

	if len(parentNames) > 0 {
		parentUsage = strings.Join(parentNames, " ") + " "
	}

	helpBuilder.WriteString(`Usage:
    ` + parentUsage + b.name)

	if len(c.Subcommands) > 0 {
		helpBuilder.WriteString(` [command]

Commands:
`)
	}

	for _, cmd := range c.Subcommands {
		helpBuilder.WriteString(strings.Repeat(" ", 4) + cmd.Name)
	}

	if len(a) > 0 {
		helpBuilder.WriteString(`

Options:
    `)
	}

	for _, arg := range a {
		argLine := ""

		switch s {
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

func (b *commandBuilder) configureBoolArgs() []*argConfig {
	var boolArgConfigs []*argConfig

	for _, arg := range b.args.boolArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		boolArgConfigs = append(boolArgConfigs, argConfig)
	}

	return boolArgConfigs
}

func (b *commandBuilder) configureFloat64Args() []*argConfig {
	var float64ArgConfigs []*argConfig

	for _, arg := range b.args.float64Args {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		float64ArgConfigs = append(float64ArgConfigs, argConfig)
	}

	return float64ArgConfigs
}

func (b *commandBuilder) configureFloat64ListArgs() []*argConfig {
	var float64ListArgConfigs []*argConfig

	for _, arg := range b.args.float64ListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		float64ListArgConfigs = append(float64ListArgConfigs, argConfig)
	}

	return float64ListArgConfigs
}

func (b *commandBuilder) configureIntArgs() []*argConfig {
	var intArgConfigs []*argConfig

	for _, arg := range b.args.intArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		intArgConfigs = append(intArgConfigs, argConfig)
	}

	return intArgConfigs
}

func (b *commandBuilder) configureIntListArgs() []*argConfig {
	var intListArgConfigs []*argConfig

	for _, arg := range b.args.intListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		intListArgConfigs = append(intListArgConfigs, argConfig)
	}

	return intListArgConfigs
}

func (b *commandBuilder) configureInt64Args() []*argConfig {
	var int64ArgConfigs []*argConfig

	for _, arg := range b.args.int64Args {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		int64ArgConfigs = append(int64ArgConfigs, argConfig)
	}

	return int64ArgConfigs
}

func (b *commandBuilder) configureInt64ListArgs() []*argConfig {
	var int64ListArgConfigs []*argConfig

	for _, arg := range b.args.int64ListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		int64ListArgConfigs = append(int64ListArgConfigs, argConfig)
	}

	return int64ListArgConfigs
}

func (b *commandBuilder) configureStringArgs() []*argConfig {
	var stringArgConfigs []*argConfig

	for _, arg := range b.args.stringArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		stringArgConfigs = append(stringArgConfigs, argConfig)
	}

	return stringArgConfigs
}

func (b *commandBuilder) configureStringListArgs() []*argConfig {
	var stringListArgConfigs []*argConfig

	for _, arg := range b.args.stringListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		stringListArgConfigs = append(stringListArgConfigs, argConfig)
	}

	return stringListArgConfigs
}

func (b *commandBuilder) configureUintArgs() []*argConfig {
	var uintArgConfigs []*argConfig

	for _, arg := range b.args.uintArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		uintArgConfigs = append(uintArgConfigs, argConfig)
	}

	return uintArgConfigs
}

func (b *commandBuilder) configureUintListArgs() []*argConfig {
	var uintListArgConfigs []*argConfig

	for _, arg := range b.args.uintListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		uintListArgConfigs = append(uintListArgConfigs, argConfig)
	}

	return uintListArgConfigs
}

func (b *commandBuilder) configureUint64Args() []*argConfig {
	var uint64ArgConfigs []*argConfig

	for _, arg := range b.args.uint64Args {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		uint64ArgConfigs = append(uint64ArgConfigs, argConfig)
	}

	return uint64ArgConfigs
}

func (b *commandBuilder) configureUint64ListArgs() []*argConfig {
	var uint64ListArgConfigs []*argConfig

	for _, arg := range b.args.uint64ListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		uint64ListArgConfigs = append(uint64ListArgConfigs, argConfig)
	}

	return uint64ListArgConfigs
}

func (b *commandBuilder) configureSubcommands() []*command {
	var subcommandConfigs []*command

	for _, subCmd := range b.subcommands {
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
