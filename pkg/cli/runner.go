package cli

import "errors"

func NewRunner(v string) Runner {
	return &runner{
		version: v,
	}
}

func (r *runner) Run(p *ParsedCommand) error {
	if p == nil {
		return errors.New("no root command parsed")
	}

	if p.Run != nil {
		p.Run(p.Context, p.Operands)
	}

	for _, subCmd := range p.Subcommands {
		subCmd.Run(subCmd.Context, subCmd.Operands)
	}

	return nil
}
