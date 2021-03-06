package cli

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func NewParser(a ArgSyntax, c CommandBuilder) Parser {
	return &parser{
		argSyntax:      a,
		builder:        c,
		parsedCommands: []*parsedCommand{},
	}
}

func (p *parser) Parse() ([]*parsedCommand, error) {
	args := os.Args[1:]
	rootCmd := p.parseCommands(args, p.builder.Build())

	for _, cmd := range p.parsedCommands {
		if argErr := p.parseArgs(cmd); argErr != nil {
			return nil, argErr
		}
	}

	if p.helpMode {
		rootCmd.HelpMode = p.helpMode
		rootCmd.HelpCommand = p.HelpCommand
	}

	return p.parsedCommands, nil
}

func (p *parser) parseCommands(a []string, c *command) *parsedCommand {
	rootCmd := p.newParsedCommand(c)
	p.parsedCommands = append(p.parsedCommands, rootCmd)
	lastParsed := rootCmd
	walker := newWalker(c)

	if len(a) == 0 {
		return rootCmd
	}

	for _, arg := range a {
		if found := walker.Walk(arg); found != nil {
			parsed := p.newParsedCommand(found)
			p.addParsedCommand(parsed)
			lastParsed = parsed

			continue
		}

		lastParsed.args = append(lastParsed.args, arg)
	}

	return rootCmd
}

func (p *parser) newParsedCommand(c *command) *parsedCommand {
	return &parsedCommand{
		args:        []string{},
		argConfigs:  c.Args,
		command:     c,
		Context:     c.Context,
		HelpCommand: c,
		Name:        c.Name,
		Run:         c.Run,
		Syntax:      p.argSyntax,
	}
}

func (p *parser) addParsedCommand(c *parsedCommand) {
	for _, parsed := range p.parsedCommands {
		if c.command.Parent == parsed.command {
			parsed.Subcommands = append(parsed.Subcommands, c)
		}
	}

	p.parsedCommands = append(p.parsedCommands, c)
}

func (p *parser) parseArgs(c *parsedCommand) error {
	switch p.argSyntax {
	case GNU:
		return p.parseArgRules(c, getGnuRules(), getPosixArgParserContext)
	case POSIX:
		return p.parseArgRules(c, getPosixRules(), getPosixArgParserContext)
	default:
		return errors.New("unsupported argument parsing syntax")
	}
}

func (p *parser) parseArgRules(c *parsedCommand, r []argParserRule, i argParserInit) error {
	if len(c.args) == 0 {
		return nil
	}

	context := i(c.args)
	context.argConfigs = c.argConfigs

	for argIndex, arg := range c.args {
		var skip bool
		var err error

		for _, rule := range r {
			skip, err = rule(&arg, argIndex, context)

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

	return p.bindArgs(c)
}

func (p *parser) bindArgs(c *parsedCommand) error {
	if len(c.parsedArgs) == 0 {
		return nil
	}

	for _, arg := range c.parsedArgs {
		if arg.name == "help" || arg.name == "h" {
			if c.HelpCommand != nil {
				p.HelpCommand = c.command
			}

			p.helpMode = true
		}

		if arg.name == "version" || arg.name == "v" {
			c.VersionMode = true
		}

		if argErr := setArgValue(arg); argErr != nil {
			return argErr
		}
	}

	return nil
}

func newWalker(c *command) *commandWalker {
	return &commandWalker{
		root: c,
		path: c.Subcommands,
	}
}

func (w *commandWalker) Walk(a string) *command {
	for _, cmd := range w.path {
		if a == cmd.Name {
			w.updatePath(cmd)

			return cmd
		}
	}

	return nil
}

func (w *commandWalker) updatePath(c *command) {
	walkablePath := append([]*command{}, c.Subcommands...)
	parent := c.Parent

	for parent != nil {
		walkablePath = append(walkablePath, parent.Subcommands...)
		parent = parent.Parent
	}

	w.path = walkablePath
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

func checkGnuOptionValidity(a *string, i int, _ *argParserContext) (bool, error) {
	if i == 0 && !strings.HasPrefix(*a, "-") && !strings.HasPrefix(*a, "--") {
		return false, errors.New("invalid GNU option: " + *a)
	}

	return false, nil
}

func checkGnuArgIsLongOption(a *string, _ int, c *argParserContext) (bool, error) {
	argParsed := false

	if !strings.HasPrefix(*a, "--") || len(*a) < 3 {
		return false, nil
	}

	option := strings.TrimPrefix(*a, "--")
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

		updateArgParserContext(argConfig, option, *a, c)
		c.lastParsedArg.value = optArgValues
		argParsed = true

		break
	}

	return argParsed, nil
}

func checkGnuArgIsLongOptionArgument(a *string, _ int, c *argParserContext) (bool, error) {
	if c.lastParsedArg == nil {
		return false, nil
	}

	for _, pArg := range c.parsedArgs {
		if c.lastParsedArg != pArg || !strings.HasPrefix(pArg.rawArg, "--") {
			continue
		}

		if !pArg.required && len(pArg.value) == 0 {
			return false, errors.New(
				"optional GNU option-argument '" + *a + "' must be provided with option '--" + pArg.name + "' separated by '='",
			)
		}

		pArg.value = append(pArg.value, *a)

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

func checkPosixOptionValidity(a *string, i int, _ *argParserContext) (bool, error) {
	if i == 0 && !strings.HasPrefix(*a, "-") {
		return false, errors.New("invalid POSIX option: " + *a)
	}

	return false, nil
}

func checkPosixArgsTerminated(a *string, i int, c *argParserContext) (bool, error) {
	if *a == "--" && i == c.terminatorIndex {
		c.terminated = true

		return true, nil
	}

	return false, nil
}

func checkPosixArgIsOperand(a *string, _ int, c *argParserContext) (bool, error) {
	if c.terminated {
		c.operands = append(c.operands, *a)

		return true, nil
	}

	return false, nil
}

func checkPosixArgIsOption(a *string, _ int, c *argParserContext) (bool, error) {
	argParsed := false

	if !strings.HasPrefix(*a, "-") || len(*a) < 2 || *a == "--" {
		return false, nil
	}

	option := strings.TrimPrefix(*a, "-")
	restOfArgs := option

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

			updateArgParserContext(argConfig, optName, *a, c)
			argParsed = true
			restOfArgs = strings.TrimPrefix(restOfArgs, optName)
			*a = strings.TrimPrefix(restOfArgs, optName)

			break
		}
	}

	return argParsed, nil
}

func checkPosixArgIsOptionArgument(a *string, _ int, c *argParserContext) (bool, error) {
	if c.lastParsedArg != nil {
		for _, pArg := range c.parsedArgs {
			if c.lastParsedArg == pArg {
				pArg.value = append(pArg.value, *a)

				return true, nil
			}
		}
	}

	return false, nil
}

func updateArgParserContext(a *argConfig, o string, r string, c *argParserContext) {
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
