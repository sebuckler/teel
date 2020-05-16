package cli

import (
	"errors"
	"flag"
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
	rootCmd := parseRootCmd(c)
	p.parsedCommands = append(p.parsedCommands, rootCmd)
	p.parseSubcommands(args, c, rootCmd)

	for _, cmd := range p.parsedCommands {
		if argErr := p.parseArgs(cmd); argErr != nil {
			return nil, argErr
		}
	}

	return rootCmd, nil
}

func (p *parser) parseSubcommands(a []string, c *CommandConfig, l *ParsedCommand) bool {
	if len(a) == 0 || l == nil {
		return true
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
			argMapped = p.parseSubcommands(a, cmd, l)
		}
	}

	if lastParsedCmd == nil {
		lastParsedCmd = l
	}

	if !argMapped {
		lastParsedCmd.args = append(lastParsedCmd.args, arg)
	}

	if len(a) == 1 {
		return true
	}

	return p.parseSubcommands(a[1:], c, lastParsedCmd)
}

func (p *parser) parseArgs(c *ParsedCommand) error {
	switch p.argSyntax {
	case GoFlag:
		return p.parseGoFlags()
	case GNU:
		return parseArgRules(c, getGnuRules(), getPosixArgParserContext)
	case POSIX:
		return parseArgRules(c, getPosixRules(), getPosixArgParserContext)
	default:
		return errors.New("unsupported argument parsing syntax")
	}
}

func (p *parser) parseGoFlags() error {
	for _, cmd := range p.parsedCommands {
		flagSet := flag.NewFlagSet(cmd.Name, flag.ContinueOnError)

		for _, argConfig := range cmd.argConfigs {
			flagSet.Var(&goFlagArgValue{
				arg: &parsedArg{
					bindVal:  argConfig.Value,
					name:     argConfig.Name,
					required: argConfig.Required,
				},
			}, argConfig.Name, argConfig.UsageText)
		}

		if parseErr := flagSet.Parse(cmd.args); parseErr != nil {
			return parseErr
		}
	}

	return nil
}

func (g *goFlagArgValue) IsBoolFlag() bool {
	_, ok := g.arg.bindVal.(*bool)

	return ok
}

func (g *goFlagArgValue) Set(v string) error {
	if g.IsBoolFlag() && g.isValidBoolVal(v) {
		g.arg.value = []string{""}
	} else {
		g.arg.value = []string{v}
	}

	return setArgValue(g.arg)
}

func (g *goFlagArgValue) String() string {
	return ""
}

func (g *goFlagArgValue) isValidBoolVal(v string) bool {
	boolVals := []string{"1", "0", "t", "f", "T", "F", "true", "false", "TRUE", "FALSE", "True", "False"}

	for _, val := range boolVals {
		if val == v {
			return true
		}
	}

	return false
}

func parseRootCmd(c *CommandConfig) *ParsedCommand {
	return &ParsedCommand{
		args:       []string{},
		argConfigs: c.Args,
		Context:    c.Context,
		Run:        c.Run,
	}
}

func parseArgRules(c *ParsedCommand, r []argParserRule, i argParserInit) error {
	if len(c.args) == 0 {
		return nil
	}

	context := i(c.args)
	context.argConfigs = c.argConfigs

	for argIndex, arg := range c.args {
		var skip bool
		var err error

		for _, rule := range r {
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

	return bindArgs(c)
}

func bindArgs(c *ParsedCommand) error {
	if len(c.parsedArgs) == 0 {
		return nil
	}

	for _, arg := range c.parsedArgs {
		if argErr := setArgValue(arg); argErr != nil {
			return argErr
		}
	}

	return nil
}

func getGnuRules() []argParserRule {
	return []argParserRule{
		checkGnuOptionValidity,
		checkPosixArgsTerminated,
		checkPosixArgIsOperand,
		checkGnuArgIsLongOption,
		checkGnuArgIsLongOptionArgument,
		checkPosixArgIsOption,
		checkPosixArgIsOptionArgument,
	}
}

func checkGnuOptionValidity(a string, i int, _ *argParserContext) (bool, error) {
	if i == 0 && !strings.HasPrefix(a, "-") && !strings.HasPrefix(a, "--") {
		return false, errors.New("invalid GNU option: " + a)
	}

	return false, nil
}

func checkGnuArgIsLongOption(a string, _ int, c *argParserContext) (bool, error) {
	argParsed := false

	if !strings.HasPrefix(a, "--") || len(a) < 3 {
		return false, nil
	}

	option := strings.TrimPrefix(a, "--")
	optArgValues := strings.Split(option, "=")

	if len(optArgValues) > 0 && optArgValues[0] != "" {
		option = optArgValues[0]
		optArgValues = optArgValues[1:]
	}

	if len(optArgValues) > 1 {
		return false, errors.New(
			"invalid GNU option argument: '" + strings.Join(optArgValues[1:], ",") + "' for option: --" + option,
		)
	}

	for _, argConfig := range c.argConfigs {
		if option != argConfig.Name {
			continue
		}

		for _, argNamePart := range strings.Split(argConfig.Name, "-") {
			for _, char := range argNamePart {
				if !isValidPosixOptionName(string(char), char) {
					return false, errors.New("invalid GNU option name: --" + option)
				}
			}
		}

		for _, pArg := range c.parsedArgs {
			if option == pArg.name && !argConfig.Repeatable {
				return false, errors.New("non-repeatable GNU option: --" + option)
			}
		}

		updateArgParserContext(argConfig, option, a, c)
		c.lastParsedArg.value = optArgValues
		argParsed = true

		break
	}

	return argParsed, nil
}

func checkGnuArgIsLongOptionArgument(a string, _ int, c *argParserContext) (bool, error) {
	if c.lastParsedArg == nil {
		return false, nil
	}

	for _, pArg := range c.parsedArgs {
		if c.lastParsedArg != pArg || !strings.HasPrefix(pArg.rawArg, "--") {
			continue
		}

		if !pArg.required && len(pArg.value) == 0 {
			return false, errors.New(
				"optional GNU option-argument '" + a + "' must be provided with option '--" + pArg.name + "' separated by '='",
			)
		}

		pArg.value = append(pArg.value, a)

		return true, nil
	}

	return false, nil
}

func getPosixArgParserContext(a []string) *argParserContext {
	terminatorIndex := -1

	for i, arg := range a {
		if arg == "--" {
			terminatorIndex = i
		}
	}

	return &argParserContext{
		terminatorIndex: terminatorIndex,
	}
}

func getPosixRules() []argParserRule {
	return []argParserRule{
		checkPosixOptionValidity,
		checkPosixArgsTerminated,
		checkPosixArgIsOperand,
		checkPosixArgIsOption,
		checkPosixArgIsOptionArgument,
	}
}

func checkPosixOptionValidity(a string, i int, _ *argParserContext) (bool, error) {
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

func checkPosixArgIsOperand(a string, _ int, c *argParserContext) (bool, error) {
	if c.terminated {
		c.operands = append(c.operands, a)

		return true, nil
	}

	return false, nil
}

func checkPosixArgIsOption(a string, _ int, c *argParserContext) (bool, error) {
	argParsed := false

	if !strings.HasPrefix(a, "-") || len(a) < 2 || a == "--" {
		return false, nil
	}

	option := strings.TrimPrefix(a, "-")

	for _, char := range option {
		optName := string(char)

		for _, argConfig := range c.argConfigs {
			if optName != argConfig.Name && char != argConfig.ShortName {
				argParsed = false

				continue
			}

			argConfig.Required = true

			if !isValidPosixOptionName(argConfig.Name, argConfig.ShortName) {
				return false, errors.New("invalid POSIX option name: -" + option)
			}

			for _, pArg := range c.parsedArgs {
				if option == pArg.name && !argConfig.Repeatable {
					return false, errors.New("non-repeatable POSIX option: -" + option)
				}
			}

			updateArgParserContext(argConfig, optName, a, c)
			argParsed = true

			break
		}
	}

	return argParsed, nil
}

func checkPosixArgIsOptionArgument(a string, _ int, c *argParserContext) (bool, error) {
	if c.lastParsedArg != nil {
		for _, pArg := range c.parsedArgs {
			if c.lastParsedArg == pArg {
				pArg.value = append(pArg.value, a)

				return true, nil
			}
		}
	}

	return false, nil
}

func updateArgParserContext(a *ArgConfig, o string, r string, c *argParserContext) {
	pArg := &parsedArg{
		bindVal:  a.Value,
		name:     o,
		rawArg:   r,
		required: a.Required,
		value:    []string{},
	}
	c.lastParsedArg = pArg
	c.parsedArgs = append(c.parsedArgs, pArg)
}

func isValidPosixListArg(a *parsedArg) error {
	if len(a.value) == 0 {
		return errors.New("no POSIX option-arguments provided for option: -" + a.name)
	}

	return nil
}

func isValidPosixNonlistArg(a *parsedArg) error {
	if a.required && len(a.value) == 0 {
		return errors.New("missing option-argument for required option: " + a.name)
	}

	if a.required && len(a.value) > 1 {
		return errors.New("invalid POSIX option-argument: '" + strings.Join(a.value, ",") +
			"' for option: -" + a.name,
		)
	}

	return nil
}

func isValidPosixOptionName(s string, r rune) bool {
	return (s != "" && len(s) == 1 && ((s[0] >= 'a' && s[0] <= 'z') || (s[0] >= 'A' && s[0] <= 'Z')) && s[0] != 'W') ||
		(((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) && r != 'W')
}

func setArgValue(p *parsedArg) error {
	switch p.bindVal.(type) {
	case *bool:
		if len(p.value) > 0 && p.value[0] != "" {
			return errors.New(
				"invalid option-argument: '" + strings.Join(p.value, ",") +
					"' for option: " + p.name,
			)
		}

		*(p.bindVal.(*bool)) = true
	case *float64:
		if err := isValidPosixNonlistArg(p); err != nil {
			return err
		}

		if len(p.value) == 0 {
			return nil
		}

		argVal := p.value[0]
		float64Val, float64Err := strconv.ParseFloat(argVal, 64)

		if float64Err != nil || argVal == "" {
			return errors.New("invalid option-argument: '" + argVal + "' for option: " + p.name)
		}

		*(p.bindVal.(*float64)) = float64Val
	case *[]float64:
		if listArgErr := isValidPosixListArg(p); listArgErr != nil {
			return listArgErr
		}

		var float64Vals []float64

		for _, argVal := range p.value {
			csv := strings.Split(argVal, ",")

			for _, val := range csv {
				float64Val, float64Err := strconv.ParseFloat(strings.TrimSpace(val), 64)

				if float64Err != nil || val == "" {
					return errors.New("invalid option-argument: '" + val + "' for option: " + p.name)
				}

				float64Vals = append(float64Vals, float64Val)
			}
		}

		*(p.bindVal.(*[]float64)) = float64Vals
	case *int:
		if err := isValidPosixNonlistArg(p); err != nil {
			return err
		}

		if len(p.value) == 0 {
			return nil
		}

		argVal := p.value[0]
		intVal, intErr := strconv.Atoi(argVal)

		if intErr != nil || argVal == "" {
			return errors.New("invalid option-argument: '" + argVal + "' for option: " + p.name)
		}

		*(p.bindVal.(*int)) = intVal
	case *[]int:
		if listArgErr := isValidPosixListArg(p); listArgErr != nil {
			return listArgErr
		}

		var intVals []int

		for _, argVal := range p.value {
			csv := strings.Split(argVal, ",")

			for _, val := range csv {
				intVal, intErr := strconv.Atoi(strings.TrimSpace(val))

				if intErr != nil || val == "" {
					return errors.New("invalid option-argument: '" + argVal + "' for option: " + p.name)
				}

				intVals = append(intVals, intVal)
			}
		}

		*(p.bindVal.(*[]int)) = intVals
	case *int64:
		if err := isValidPosixNonlistArg(p); err != nil {
			return err
		}

		if len(p.value) == 0 {
			return nil
		}

		argVal := p.value[0]
		int64Val, int64Err := strconv.ParseInt(argVal, 10, 64)

		if int64Err != nil || argVal == "" {
			return errors.New("invalid option-argument: '" + argVal + "' for option: " + p.name)
		}

		*(p.bindVal.(*int64)) = int64Val
	case *[]int64:
		if listArgErr := isValidPosixListArg(p); listArgErr != nil {
			return listArgErr
		}

		var int64Vals []int64

		for _, argVal := range p.value {
			csv := strings.Split(argVal, ",")

			for _, val := range csv {
				int64Val, int64Err := strconv.ParseInt(strings.TrimSpace(val), 10, 64)

				if int64Err != nil || val == "" {
					return errors.New("invalid option-argument: '" + val + "' for option: " + p.name)
				}

				int64Vals = append(int64Vals, int64Val)
			}
		}

		*(p.bindVal.(*[]int64)) = int64Vals
	case *string:
		if err := isValidPosixNonlistArg(p); err != nil {
			return err
		}

		if len(p.value) == 0 {
			return nil
		}

		*(p.bindVal.(*string)) = p.value[0]
	case *[]string:
		if listArgErr := isValidPosixListArg(p); listArgErr != nil {
			return listArgErr
		}

		var stringVals []string

		for _, argVal := range p.value {
			csv := strings.Split(argVal, ",")

			for _, val := range csv {
				stringVals = append(stringVals, val)
			}
		}

		*(p.bindVal.(*[]string)) = stringVals
	case *uint:
		if err := isValidPosixNonlistArg(p); err != nil {
			return err
		}

		if len(p.value) == 0 {
			return nil
		}

		argVal := p.value[0]
		uintVal, uintErr := strconv.ParseUint(argVal, 10, 0)

		if uintErr != nil || argVal == "" {
			return errors.New("invalid option-argument: '" + argVal + "' for option: " + p.name)
		}

		*(p.bindVal.(*uint)) = uint(uintVal)
	case *[]uint:
		if listArgErr := isValidPosixListArg(p); listArgErr != nil {
			return listArgErr
		}

		var uintVals []uint

		for _, argVal := range p.value {
			csv := strings.Split(argVal, ",")

			for _, val := range csv {
				uintVal, uintErr := strconv.ParseUint(strings.TrimSpace(val), 10, 0)

				if uintErr != nil || val == "" {
					return errors.New("invalid option-argument: '" + val + "' for option: " + p.name)
				}

				uintVals = append(uintVals, uint(uintVal))
			}
		}

		*(p.bindVal.(*[]uint)) = uintVals
	case *uint64:
		if err := isValidPosixNonlistArg(p); err != nil {
			return err
		}

		if len(p.value) == 0 {
			return nil
		}

		argVal := p.value[0]
		uint64Val, uint64Err := strconv.ParseUint(argVal, 10, 64)

		if uint64Err != nil || argVal == "" {
			return errors.New("invalid option-argument: '" + argVal + "' for option: " + p.name)
		}

		*(p.bindVal.(*uint64)) = uint64Val
	case *[]uint64:
		if listArgErr := isValidPosixListArg(p); listArgErr != nil {
			return listArgErr
		}

		var uint64Vals []uint64

		for _, argVal := range p.value {
			csv := strings.Split(argVal, ",")

			for _, val := range csv {
				uint64Val, uint64Err := strconv.ParseUint(strings.TrimSpace(val), 10, 64)

				if uint64Err != nil || val == "" {
					return errors.New("invalid option-argument: '" + val + "' for option: " + p.name)
				}

				uint64Vals = append(uint64Vals, uint64Val)
			}
		}

		*(p.bindVal.(*[]uint64)) = uint64Vals
	default:
		return errors.New("invalid option: " + p.name)
	}

	return nil
}
