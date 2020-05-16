package cli

import (
	"context"
)

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

func (c *commandConfigurer) AddRunFunc(r CommandRunFunc) {
	c.run = r
}

func (c *commandConfigurer) AddSubcommand(cmd CommandConfigurer) {
	c.subcommands = append(c.subcommands, cmd)
}

func (c *commandConfigurer) AddBoolArg(p *bool, a *ArgDefinition) {
	c.args.boolArgs = append(c.args.boolArgs, &boolArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandConfigurer) AddFloat64Arg(p *float64, a *ArgDefinition) {
	c.args.float64Args = append(c.args.float64Args, &float64Arg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandConfigurer) AddFloat64ListArg(p *[]float64, a *ArgDefinition) {
	c.args.float64ListArgs = append(c.args.float64ListArgs, &float64ListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandConfigurer) AddIntArg(p *int, a *ArgDefinition) {
	c.args.intArgs = append(c.args.intArgs, &intArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandConfigurer) AddIntListArg(p *[]int, a *ArgDefinition) {
	c.args.intListArgs = append(c.args.intListArgs, &intListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandConfigurer) AddInt64Arg(p *int64, a *ArgDefinition) {
	c.args.int64Args = append(c.args.int64Args, &int64Arg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandConfigurer) AddInt64ListArg(p *[]int64, a *ArgDefinition) {
	c.args.int64ListArgs = append(c.args.int64ListArgs, &int64ListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandConfigurer) AddStringArg(p *string, a *ArgDefinition) {
	c.args.stringArgs = append(c.args.stringArgs, &stringArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandConfigurer) AddStringListArg(p *[]string, a *ArgDefinition) {
	c.args.stringListArgs = append(c.args.stringListArgs, &stringListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandConfigurer) AddUintArg(p *uint, a *ArgDefinition) {
	c.args.uintArgs = append(c.args.uintArgs, &uintArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandConfigurer) AddUintListArg(p *[]uint, a *ArgDefinition) {
	c.args.uintListArgs = append(c.args.uintListArgs, &uintListArg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandConfigurer) AddUint64Arg(p *uint64, a *ArgDefinition) {
	c.args.uint64Args = append(c.args.uint64Args, &uint64Arg{
		commandArg: newCommandArg(a),
		value:      p,
	})
}

func (c *commandConfigurer) AddUint64ListArg(p *[]uint64, a *ArgDefinition) {
	c.args.uint64ListArgs = append(c.args.uint64ListArgs, &uint64ListArg{
		commandArg: newCommandArg(a),
		value:      p,
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

func (c *commandConfigurer) configureArgs() []*ArgConfig {
	var argConfigs []*ArgConfig

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

	return argConfigs
}

func (c *commandConfigurer) configureBoolArgs() []*ArgConfig {
	var boolArgConfigs []*ArgConfig

	for _, arg := range c.args.boolArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		boolArgConfigs = append(boolArgConfigs, argConfig)
	}

	return boolArgConfigs
}

func (c *commandConfigurer) configureFloat64Args() []*ArgConfig {
	var float64ArgConfigs []*ArgConfig

	for _, arg := range c.args.float64Args {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		float64ArgConfigs = append(float64ArgConfigs, argConfig)
	}

	return float64ArgConfigs
}

func (c *commandConfigurer) configureFloat64ListArgs() []*ArgConfig {
	var float64ListArgConfigs []*ArgConfig

	for _, arg := range c.args.float64ListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		float64ListArgConfigs = append(float64ListArgConfigs, argConfig)
	}

	return float64ListArgConfigs
}

func (c *commandConfigurer) configureIntArgs() []*ArgConfig {
	var intArgConfigs []*ArgConfig

	for _, arg := range c.args.intArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		intArgConfigs = append(intArgConfigs, argConfig)
	}

	return intArgConfigs
}

func (c *commandConfigurer) configureIntListArgs() []*ArgConfig {
	var intListArgConfigs []*ArgConfig

	for _, arg := range c.args.intListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		intListArgConfigs = append(intListArgConfigs, argConfig)
	}

	return intListArgConfigs
}

func (c *commandConfigurer) configureInt64Args() []*ArgConfig {
	var int64ArgConfigs []*ArgConfig

	for _, arg := range c.args.int64Args {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		int64ArgConfigs = append(int64ArgConfigs, argConfig)
	}

	return int64ArgConfigs
}

func (c *commandConfigurer) configureInt64ListArgs() []*ArgConfig {
	var int64ListArgConfigs []*ArgConfig

	for _, arg := range c.args.int64ListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		int64ListArgConfigs = append(int64ListArgConfigs, argConfig)
	}

	return int64ListArgConfigs
}

func (c *commandConfigurer) configureStringArgs() []*ArgConfig {
	var stringArgConfigs []*ArgConfig

	for _, arg := range c.args.stringArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		stringArgConfigs = append(stringArgConfigs, argConfig)
	}

	return stringArgConfigs
}

func (c *commandConfigurer) configureStringListArgs() []*ArgConfig {
	var stringListArgConfigs []*ArgConfig

	for _, arg := range c.args.stringListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		stringListArgConfigs = append(stringListArgConfigs, argConfig)
	}

	return stringListArgConfigs
}

func (c *commandConfigurer) configureUintArgs() []*ArgConfig {
	var uintArgConfigs []*ArgConfig

	for _, arg := range c.args.uintArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		uintArgConfigs = append(uintArgConfigs, argConfig)
	}

	return uintArgConfigs
}

func (c *commandConfigurer) configureUintListArgs() []*ArgConfig {
	var uintListArgConfigs []*ArgConfig

	for _, arg := range c.args.uintListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		uintListArgConfigs = append(uintListArgConfigs, argConfig)
	}

	return uintListArgConfigs
}

func (c *commandConfigurer) configureUint64Args() []*ArgConfig {
	var uint64ArgConfigs []*ArgConfig

	for _, arg := range c.args.uint64Args {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		uint64ArgConfigs = append(uint64ArgConfigs, argConfig)
	}

	return uint64ArgConfigs
}

func (c *commandConfigurer) configureUint64ListArgs() []*ArgConfig {
	var uint64ListArgConfigs []*ArgConfig

	for _, arg := range c.args.uint64ListArgs {
		argConfig := newArgConfig(arg.commandArg, arg.value)
		uint64ListArgConfigs = append(uint64ListArgConfigs, argConfig)
	}

	return uint64ListArgConfigs
}

func (c *commandConfigurer) configureSubcommands() []*CommandConfig {
	var subcommandConfigs []*CommandConfig

	for _, subCmd := range c.subcommands {
		subcommandConfigs = append(subcommandConfigs, subCmd.Configure())
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

func newArgConfig(a *commandArg, v interface{}) *ArgConfig {
	return &ArgConfig{
		Name:       a.name,
		Repeatable: a.repeatable,
		Required:   a.required,
		ShortName:  a.shortName,
		UsageText:  a.usageText,
		Value:      v,
	}
}
