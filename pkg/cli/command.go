package cli

import "context"

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

	for _, arg := range c.args.int64Args {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue, Int64)
		int64ArgConfigs = append(int64ArgConfigs, argConfig)
	}

	return int64ArgConfigs
}

func (c *commandConfigurer) configureStringArgs() []*ArgConfig {
	var stringArgConfigs []*ArgConfig

	for _, arg := range c.args.stringArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue, String)
		stringArgConfigs = append(stringArgConfigs, argConfig)
	}

	return stringArgConfigs
}

func (c *commandConfigurer) configureUintArgs() []*ArgConfig {
	var uintArgConfigs []*ArgConfig

	for _, arg := range c.args.uintArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue, Uint)
		uintArgConfigs = append(uintArgConfigs, argConfig)
	}

	return uintArgConfigs
}

func (c *commandConfigurer) configureUint64Args() []*ArgConfig {
	var uint64ArgConfigs []*ArgConfig

	for _, arg := range c.args.uint64Args {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue, Uint64)
		uint64ArgConfigs = append(uint64ArgConfigs, argConfig)
	}

	return uint64ArgConfigs
}

func (c *commandConfigurer) configureCommandArgType(a *commandArg, v interface{}, d interface{}, t ArgType) *ArgConfig {
	if v != nil {
		switch v.(type) {
		case *bool:
			val := v.(*bool)
			*val = d.(bool)
			v = val
		case *string:
			val := v.(*string)
			*val = d.(string)
			v = val
		}
	}

	return &ArgConfig{
		Name:       a.name,
		Repeatable: a.repeatable,
		ShortName:  a.shortName,
		Type:       t,
		UsageText:  a.usageText,
		Value:      v,
	}
}

func (c *commandConfigurer) configureSubcommands() []*CommandConfig {
	var subcommandConfigs []*CommandConfig

	for _, subCmd := range c.subcommands {
		subcommandConfigs = append(subcommandConfigs, subCmd.Configure())
	}

	return subcommandConfigs
}
