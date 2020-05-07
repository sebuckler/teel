package cli

import "errors"

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

	if cmdConfig == nil {
		return errors.New("no root command configured")
	}

	parsedCmd, parseErr := r.parser.Parse(cmdConfig)

	if parseErr != nil {
		return parseErr
	}

	if parsedCmd == nil {
		return errors.New("no root command parsed")
	}

	if parsedCmd.Run != nil {
		parsedCmd.Run(cmdConfig.Context, cmdConfig.Operands)
	}

	for _, subCmd := range parsedCmd.Subcommands {
		subCmd.Run(subCmd.Context, subCmd.Operands)
	}

	return nil
}
