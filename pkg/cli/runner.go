package cli

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

	if parseErr := r.parser.Parse(cmdConfig); parseErr != nil {
		return parseErr
	}

	if cmdConfig.Run != nil {
		cmdConfig.Run(cmdConfig.Context)
	}

	for _, subCmd := range cmdConfig.Subcommands {
		subCmd.Run(subCmd.Context)
	}

	return nil
}
