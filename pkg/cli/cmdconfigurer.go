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

func (c *commandConfigurer) AddBoolListArg(n string, s rune, p *[]bool, v []bool, u string, r bool) {
	c.args.boolListArgs = append(c.args.boolListArgs, &boolListArg{
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

func (c *commandConfigurer) AddIntListArg(n string, s rune, p *[]int, v []int, u string, r bool) {
	c.args.intListArgs = append(c.args.intListArgs, &intListArg{
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

func (c *commandConfigurer) AddInt64ListArg(n string, s rune, p *[]int64, v []int64, u string, r bool) {
	c.args.int64ListArgs = append(c.args.int64ListArgs, &int64ListArg{
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

func (c *commandConfigurer) AddStringListArg(n string, s rune, p *[]string, v []string, u string, r bool) {
	c.args.stringListArgs = append(c.args.stringListArgs, &stringListArg{
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

func (c *commandConfigurer) AddUintListArg(n string, s rune, p *[]uint, v []uint, u string, r bool) {
	c.args.uintListArgs = append(c.args.uintListArgs, &uintListArg{
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

func (c *commandConfigurer) AddUint64ListArg(n string, s rune, p *[]uint64, v []uint64, u string, r bool) {
	c.args.uint64ListArgs = append(c.args.uint64ListArgs, &uint64ListArg{
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
	argConfigs = append(argConfigs, c.configureBoolListArgs()...)
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
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue)
		boolArgConfigs = append(boolArgConfigs, argConfig)
	}

	return boolArgConfigs
}

func (c *commandConfigurer) configureBoolListArgs() []*ArgConfig {
	var boolListArgConfigs []*ArgConfig

	for _, arg := range c.args.boolListArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue)
		boolListArgConfigs = append(boolListArgConfigs, argConfig)
	}

	return boolListArgConfigs
}

func (c *commandConfigurer) configureIntArgs() []*ArgConfig {
	var intArgConfigs []*ArgConfig

	for _, arg := range c.args.intArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue)
		intArgConfigs = append(intArgConfigs, argConfig)
	}

	return intArgConfigs
}

func (c *commandConfigurer) configureIntListArgs() []*ArgConfig {
	var intListArgConfigs []*ArgConfig

	for _, arg := range c.args.intListArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue)
		intListArgConfigs = append(intListArgConfigs, argConfig)
	}

	return intListArgConfigs
}

func (c *commandConfigurer) configureInt64Args() []*ArgConfig {
	var int64ArgConfigs []*ArgConfig

	for _, arg := range c.args.int64Args {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue)
		int64ArgConfigs = append(int64ArgConfigs, argConfig)
	}

	return int64ArgConfigs
}

func (c *commandConfigurer) configureInt64ListArgs() []*ArgConfig {
	var int64ListArgConfigs []*ArgConfig

	for _, arg := range c.args.int64ListArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue)
		int64ListArgConfigs = append(int64ListArgConfigs, argConfig)
	}

	return int64ListArgConfigs
}

func (c *commandConfigurer) configureStringArgs() []*ArgConfig {
	var stringArgConfigs []*ArgConfig

	for _, arg := range c.args.stringArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue)
		stringArgConfigs = append(stringArgConfigs, argConfig)
	}

	return stringArgConfigs
}

func (c *commandConfigurer) configureStringListArgs() []*ArgConfig {
	var stringListArgConfigs []*ArgConfig

	for _, arg := range c.args.stringListArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue)
		stringListArgConfigs = append(stringListArgConfigs, argConfig)
	}

	return stringListArgConfigs
}

func (c *commandConfigurer) configureUintArgs() []*ArgConfig {
	var uintArgConfigs []*ArgConfig

	for _, arg := range c.args.uintArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue)
		uintArgConfigs = append(uintArgConfigs, argConfig)
	}

	return uintArgConfigs
}

func (c *commandConfigurer) configureUintListArgs() []*ArgConfig {
	var uintListArgConfigs []*ArgConfig

	for _, arg := range c.args.uintListArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue)
		uintListArgConfigs = append(uintListArgConfigs, argConfig)
	}

	return uintListArgConfigs
}

func (c *commandConfigurer) configureUint64Args() []*ArgConfig {
	var uint64ArgConfigs []*ArgConfig

	for _, arg := range c.args.uint64Args {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue)
		uint64ArgConfigs = append(uint64ArgConfigs, argConfig)
	}

	return uint64ArgConfigs
}

func (c *commandConfigurer) configureUint64ListArgs() []*ArgConfig {
	var uint64ListArgConfigs []*ArgConfig

	for _, arg := range c.args.uint64ListArgs {
		argConfig := c.configureCommandArgType(arg.commandArg, arg.value, arg.defaultValue)
		uint64ListArgConfigs = append(uint64ListArgConfigs, argConfig)
	}

	return uint64ListArgConfigs
}

func (c *commandConfigurer) configureCommandArgType(a *commandArg, v interface{}, d interface{}) *ArgConfig {
	switch v.(type) {
	case *bool:
		if val, ok := v.(*bool); ok && val != nil {
			*val = d.(bool)
			v = val
		}
	case *[]bool:
		if val, ok := v.(*[]bool); ok && val != nil {
			*val = d.([]bool)
			v = val
		}
	case *int:
		if val, ok := v.(*int); ok && val != nil {
			*val = d.(int)
			v = val
		}
	case *[]int:
		if val, ok := v.(*[]int); ok && val != nil {
			*val = d.([]int)
			v = val
		}
	case *int64:
		if val, ok := v.(*int64); ok && val != nil {
			*val = d.(int64)
			v = val
		}
	case *[]int64:
		if val, ok := v.(*[]int64); ok && val != nil {
			*val = d.([]int64)
			v = val
		}
	case *string:
		if val, ok := v.(*string); ok && val != nil {
			*val = d.(string)
			v = val
		}
	case *[]string:
		if val, ok := v.(*[]string); ok && val != nil {
			*val = d.([]string)
			v = val
		}
	case *uint:
		if val, ok := v.(*uint); ok && val != nil {
			*val = d.(uint)
			v = val
		}
	case *[]uint:
		if val, ok := v.(*[]uint); ok && val != nil {
			*val = d.([]uint)
			v = val
		}
	case *uint64:
		if val, ok := v.(*uint64); ok && val != nil {
			*val = d.(uint64)
			v = val
		}
	case *[]uint64:
		if val, ok := v.(*[]uint64); ok && val != nil {
			*val = d.([]uint64)
			v = val
		}
	}

	return &ArgConfig{
		Name:       a.name,
		Repeatable: a.repeatable,
		ShortName:  a.shortName,
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
