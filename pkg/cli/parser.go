package cli

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

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

	lastParsedArg := map[string][]string{}
	var operands []string
	var parsedArgs []*parsedArg
	terminated := false
	terminatorIndex := getPosixTerminatorIndex(c.args)

	for argIndex, arg := range c.args {
		if argIndex == 0 && !strings.HasPrefix(arg, "-") {
			return errors.New("invalid POSIX option: " + arg)
		}

		if arg == "--" && argIndex == terminatorIndex {
			terminated = true

			continue
		}

		if strings.HasPrefix(arg, "-") && len(arg) > 1 && arg != "--" {
			option := strings.TrimPrefix(arg, "-")

			for _, char := range option {
				optName := string(char)

				for _, argConfig := range c.argConfigs {
					if optName != argConfig.Name && char != argConfig.ShortName {
						continue
					}

					if !isValidPosixOptionName(argConfig.Name, argConfig.ShortName) {
						return errors.New("invalid POSIX option name: -" + option)
					}

					for _, pArg := range parsedArgs {
						if option == pArg.name && !argConfig.Repeatable {
							return errors.New("non-repeatable POSIX option: -" + option)
						}
					}

					parsedArgs = append(parsedArgs, &parsedArg{
						argType: argConfig.Type,
						bindVal: argConfig.Value,
						name:    optName,
						value:   []string{},
					})

					if argConfig.Type == Bool {
						lastParsedArg = map[string][]string{optName: {""}}

						break
					}

					lastParsedArg = map[string][]string{optName: {}}

					break
				}
			}

			continue
		}

		if terminated {
			operands = append(operands, arg)

			continue
		}

		if len(lastParsedArg) > 0 {
			for _, pArg := range parsedArgs {
				if _, ok := lastParsedArg[pArg.name]; ok {
					pArg.value = append(pArg.value, arg)
					lastParsedArg[pArg.name] = append(lastParsedArg[pArg.name], arg)

					break
				}
			}

			continue
		}

		return errors.New("invalid POSIX option: " + arg)
	}

	c.parsedArgs = parsedArgs
	c.operands = operands

	return p.bindArgs(c)
}

func (p *parser) bindArgs(c *parsedCommand) error {
	if len(c.parsedArgs) == 0 {
		return nil
	}

	for _, arg := range c.parsedArgs {
		switch arg.argType {
		case Bool:
			if len(arg.value) > 0 && arg.value[0] != "" {
				return errors.New(
					"invalid POSIX option-argument: '" + strings.Join(arg.value, ",") +
						"' for option: -" + arg.name,
				)
			}

			bindVal := arg.bindVal.(*bool)
			*bindVal = true
		case Int:
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
		case IntList:
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
		case Int64:
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
		case Int64List:
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
		case String:
			if err := isValidPosixNonlistArg(arg); err != nil {
				return err
			}

			bindVal := arg.bindVal.(*string)
			*bindVal = arg.value[0]
		case StringList:
			if len(arg.value) == 0 {
				return errors.New("invalid POSIX option-argument: '" + strings.Join(arg.value, ",") +
					"' for option: -" + arg.name,
				)
			}

			bindVal := arg.bindVal.(*[]string)
			*bindVal = arg.value
		case Uint:
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
		case UintList:
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
		case Uint64:
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
		case Uint64List:
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
	return (s == "" || len(s) > 1) || ((s[0] < 'a' || s[0] > 'z') && (s[0] < 'A' || s[0] > 'Z')) ||
		((r < 'a' || r > 'z') && (r < 'A' || r > 'Z')) || (s[0] == 'W' || r == 'W')
}
