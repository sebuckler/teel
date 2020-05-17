package cli

import (
	"errors"
	"io"
)

func NewRunner(v string, w io.Writer) Runner {
	return &runner{
		version: v,
		writer:  w,
	}
}

func (r *runner) Run(p *ParsedCommand) error {
	if p == nil {
		return errors.New("no root command parsed")
	}

	if p.HelpMode {
		return p.HelpFunc(p.Syntax, r.writer)
	}

	if p.VersionMode {
		_, writeErr := r.writer.Write([]byte(p.Name + " " + r.version + "\n"))

		return writeErr
	}

	if p.Run != nil {
		p.Run(p.Context, p.Operands)
	}

	for _, subCmd := range p.Subcommands {
		subCmd.Run(subCmd.Context, subCmd.Operands)
	}

	return nil
}
