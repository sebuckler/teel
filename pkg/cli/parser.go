package cli

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func NewParser(a ArgSyntax) Parser {
	return &parser{
		argSyntax:      a,
		parsedCommands: []*ParsedCommand{},
	}
}

func (p *parser) Parse(c *CommandConfig) (*ParsedCommand, error) {
	args := os.Args[1:]
	rootCmd := p.parseRootCmd(c)
	p.parsedCommands = append(p.parsedCommands, rootCmd)
	p.mapSubcommands(args, c, rootCmd)

	for _, cmd := range p.parsedCommands {
		if argErr := p.parseArgs(cmd); argErr != nil {
			return nil, argErr
		}

		if bindErr := p.bindArgs(cmd); bindErr != nil {
			return nil, bindErr
		}
	}

	return rootCmd, nil
}

func (p *parser) mapSubcommands(a []string, c *CommandConfig, l *ParsedCommand) {
	if len(a) == 0 || l == nil {
		return
	}

	arg := a[0]
	argMapped := false
	var lastParsedCmd *ParsedCommand

	for _, cmd := range c.Subcommands {
		if arg == cmd.Name {
			parsedCmd := &ParsedCommand{
				argConfigs: cmd.Args,
				Context:    cmd.Context,
				Name:       cmd.Name,
				parentCmd:  c.Name,
				Run:        cmd.Run,
			}
			p.parsedCommands = append(p.parsedCommands, parsedCmd)
			l.Subcommands = append(l.Subcommands, parsedCmd)
			lastParsedCmd = parsedCmd
			argMapped = true

			break
		}

		if l.Name == cmd.Name && len(cmd.Subcommands) > 0 {
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

func (p *parser) parseRootCmd(c *CommandConfig) *ParsedCommand {
	return &ParsedCommand{
		args:       []string{},
		argConfigs: c.Args,
		Context:    c.Context,
		Run:        c.Run,
	}
}

func (p *parser) parseArgs(c *ParsedCommand) error {
	switch p.argSyntax {
	case GNU:
		return p.parseGnuArgs(c.args)
	case GoFlag:
		return p.parseGoFlagArgs(c.args)
	case POSIX:
		return p.parsePosixArgs(c)
	default:
		return errors.New("unsupported argument parsing syntax")
	}
}

func (p *parser) parseGnuArgs(a []string) error {
	panic("implement me")
}

func (p *parser) parseGoFlagArgs(a []string) error {
	panic("implement me")
}

func getPosixArgParserContext(a []string) *argParserContext {
	return &argParserContext{
		lastParsedArg:   map[string][]string{},
		terminated:      false,
		terminatorIndex: getPosixTerminatorIndex(a),
	}
}

func getPosixRules() []func(a string, i int, c *argParserContext) (bool, error) {
	return []func(a string, i int, c *argParserContext) (bool, error){
		checkPosixOptionValidity,
		checkPosixArgsTerminated,
		checkPosixArgIsOperand,
		checkPosixArgIsOption,
		checkPosixArgIsOptionArgument,
	}
}

func checkPosixOptionValidity(a string, i int, c *argParserContext) (bool, error) {
	if i == 0 && !strings.HasPrefix(a, "-") {
		return false, errors.New("invalid POSIX option: " + a)
	}

	return false, nil
}

func checkPosixArgsTerminated(a string, i int, c *argParserContext) (bool, error) {
	if a == "--" && i == c.terminatorIndex {
		c.terminated = true

		return true, nil
	}

	return false, nil
}

func checkPosixArgIsOperand(a string, i int, c *argParserContext) (bool, error) {
	if c.terminated {
		c.operands = append(c.operands, a)

		return true, nil
	}

	return false, nil
}

func checkPosixArgIsOption(a string, i int, c *argParserContext) (bool, error) {
	if strings.HasPrefix(a, "-") && len(a) > 1 && a != "--" {
		option := strings.TrimPrefix(a, "-")

		for _, char := range option {
			optName := string(char)

			for _, argConfig := range c.argConfigs {
				if optName != argConfig.Name && char != argConfig.ShortName {
					continue
				}

				if !isValidPosixOptionName(argConfig.Name, argConfig.ShortName) {
					return false, errors.New("invalid POSIX option name: -" + option)
				}
				for _, pArg := range c.parsedArgs {
					if option == pArg.name && !argConfig.Repeatable {
						return false, errors.New("non-repeatable POSIX option: -" + option)
					}
				}

				c.parsedArgs = append(c.parsedArgs, &parsedArg{
					bindVal: argConfig.Value,
					name:    optName,
					value:   []string{},
				})

				if _, ok := argConfig.Value.(*bool); ok {
					c.lastParsedArg = map[string][]string{optName: {""}}

					break
				}

				c.lastParsedArg = map[string][]string{optName: {}}

				break
			}
		}

		return true, nil
	}

	return false, nil
}

func checkPosixArgIsOptionArgument(a string, i int, c *argParserContext) (bool, error) {
	if len(c.lastParsedArg) > 0 {
		for _, pArg := range c.parsedArgs {
			if _, ok := c.lastParsedArg[pArg.name]; ok {
				pArg.value = append(pArg.value, a)
				c.lastParsedArg[pArg.name] = append(c.lastParsedArg[pArg.name], a)

				break
			}
		}

		return true, nil
	}

	return false, nil
}

func (p *parser) parsePosixArgs(c *ParsedCommand) error {
	if len(c.args) == 0 {
		return nil
	}

	context := getPosixArgParserContext(c.args)
	context.argConfigs = c.argConfigs

	for argIndex, arg := range c.args {
		var skip bool
		var err error

		for _, rule := range getPosixRules() {
			skip, err = rule(arg, argIndex, context)

			if err != nil {
				return err
			}

			if skip {
				break
			}
		}

		if skip {
			continue
		}

		return errors.New("failed to parse argument: " + arg)
	}

	c.parsedArgs = context.parsedArgs
	c.Operands = context.operands

	return nil
}

func (p *parser) bindArgs(c *ParsedCommand) error {
	if len(c.parsedArgs) == 0 {
		return nil
	}

	for _, arg := range c.parsedArgs {
		switch arg.bindVal.(type) {
		case *bool:
			if len(arg.value) > 0 && arg.value[0] != "" {
				return errors.New(
					"invalid POSIX option-argument: '" + strings.Join(arg.value, ",") +
						"' for option: -" + arg.name,
				)
			}

			bindVal := arg.bindVal.(*bool)
			*bindVal = true
		case *int:
			if err := isValidPosixNonlistArg(arg); err != nil {
				return err
			}

			argVal := arg.value[0]
			intVal, intErr := strconv.Atoi(argVal)

			if intErr != nil || argVal == "" {
				return errors.New("invalid POSIX option-argument: '" + argVal + "' for option: -" + arg.name)
			}

			bindVal := arg.bindVal.(*int)
			*bindVal = intVal
		case *[]int:
			if listArgErr := isValidPosixListArg(arg); listArgErr != nil {
				return listArgErr
			}

			var intVals []int

			for _, argVal := range arg.value {
				intVal, intErr := strconv.Atoi(argVal)

				if intErr != nil || argVal == "" {
					return errors.New("invalid POSIX option-argument: '" + argVal + "' for option: -" + arg.name)
				}

				intVals = append(intVals, intVal)
			}

			bindVal := arg.bindVal.(*[]int)
			*bindVal = intVals
		case *int64:
			if err := isValidPosixNonlistArg(arg); err != nil {
				return err
			}

			argVal := arg.value[0]
			int64Val, int64Err := strconv.ParseInt(argVal, 10, 64)

			if int64Err != nil || argVal == "" {
				return errors.New("invalid POSIX option-argument: '" + argVal + "' for option: -" + arg.name)
			}

			bindVal := arg.bindVal.(*int64)
			*bindVal = int64Val
		case *[]int64:
			if listArgErr := isValidPosixListArg(arg); listArgErr != nil {
				return listArgErr
			}

			var int64Vals []int64

			for _, argVal := range arg.value {
				int64Val, int64Err := strconv.ParseInt(argVal, 10, 64)

				if int64Err != nil || argVal == "" {
					return errors.New("invalid POSIX option-argument: '" + argVal + "' for option: -" + arg.name)
				}

				int64Vals = append(int64Vals, int64Val)
			}

			bindVal := arg.bindVal.(*[]int64)
			*bindVal = int64Vals
		case *string:
			if err := isValidPosixNonlistArg(arg); err != nil {
				return err
			}

			bindVal := arg.bindVal.(*string)
			*bindVal = arg.value[0]
		case *[]string:
			if listArgErr := isValidPosixListArg(arg); listArgErr != nil {
				return listArgErr
			}

			if len(arg.value) == 0 {
				return errors.New("invalid POSIX option-argument: '" + strings.Join(arg.value, ",") +
					"' for option: -" + arg.name,
				)
			}

			bindVal := arg.bindVal.(*[]string)
			*bindVal = arg.value
		case *uint:
			if err := isValidPosixNonlistArg(arg); err != nil {
				return err
			}

			argVal := arg.value[0]
			uintVal, uintErr := strconv.ParseUint(argVal, 10, 0)

			if uintErr != nil || argVal == "" {
				return errors.New("invalid POSIX option-argument: '" + argVal + "' for option: -" + arg.name)
			}

			bindVal := arg.bindVal.(*uint)
			*bindVal = uint(uintVal)
		case *[]uint:
			if listArgErr := isValidPosixListArg(arg); listArgErr != nil {
				return listArgErr
			}

			var uintVals []uint

			for _, argVal := range arg.value {
				uintVal, uintErr := strconv.ParseUint(argVal, 10, 32)

				if uintErr != nil || argVal == "" {
					return errors.New("invalid POSIX option-argument: '" + argVal + "' for option: -" + arg.name)
				}

				uintVals = append(uintVals, uint(uintVal))
			}

			bindVal := arg.bindVal.(*[]uint)
			*bindVal = uintVals
		case *uint64:
			if err := isValidPosixNonlistArg(arg); err != nil {
				return err
			}

			argVal := arg.value[0]
			uint64Val, uint64Err := strconv.ParseUint(argVal, 10, 64)

			if uint64Err != nil || argVal == "" {
				return errors.New("invalid POSIX option-argument: '" + argVal + "' for option: -" + arg.name)
			}

			bindVal := arg.bindVal.(*uint64)
			*bindVal = uint64Val
		case *[]uint64:
			if listArgErr := isValidPosixListArg(arg); listArgErr != nil {
				return listArgErr
			}

			var uint64Vals []uint64

			for _, argVal := range arg.value {
				uint64Val, uint64Err := strconv.ParseUint(argVal, 10, 64)

				if uint64Err != nil || argVal == "" {
					return errors.New("invalid POSIX option-argument: '" + argVal + "' for option: -" + arg.name)
				}

				uint64Vals = append(uint64Vals, uint64Val)
			}

			bindVal := arg.bindVal.(*[]uint64)
			*bindVal = uint64Vals
		default:
			return errors.New("invalid POSIX option: -" + arg.name)
		}
	}

	return nil
}

func isValidPosixListArg(arg *parsedArg) error {
	if len(arg.value) == 0 {
		return errors.New("no POSIX option-arguments provided for option: -" + arg.name)
	}
	return nil
}

func isValidPosixNonlistArg(arg *parsedArg) error {
	if len(arg.value) != 1 || arg.value[0] == "" {
		return errors.New("invalid POSIX option-argument: '" + strings.Join(arg.value, ",") +
			"' for option: -" + arg.name,
		)
	}

	return nil
}

func getPosixTerminatorIndex(a []string) int {
	lastIndex := -1

	for i, arg := range a {
		if arg == "--" {
			lastIndex = i
		}
	}

	return lastIndex
}

func isValidPosixOptionName(s string, r rune) bool {
	return !((s == "" || len(s) > 1) || ((s[0] < 'a' || s[0] > 'z') && (s[0] < 'A' || s[0] > 'Z')) ||
		((r < 'a' || r > 'z') && (r < 'A' || r > 'Z')) || (s[0] == 'W' || r == 'W'))
}
