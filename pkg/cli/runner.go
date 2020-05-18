package cli

import (
	"errors"
	"io"
)

func NewRunner(p Parser, v string, w io.Writer) Runner {
	return &runner{
		parser:  p,
		version: v,
		writer:  w,
	}
}

func (r *runner) Run() error {
	parsedCmd, parseErr := r.parser.Parse()

	if parseErr != nil {
		return parseErr
	}

	if parsedCmd == nil {
		return errors.New("no root command parsed")
	}

	if parsedCmd.HelpMode {
		return parsedCmd.HelpFunc(parsedCmd.Syntax, r.writer)
	}

	if parsedCmd.VersionMode {
		_, writeErr := r.writer.Write([]byte(parsedCmd.Name + " " + r.version + "\n"))

		return writeErr
	}

	if parsedCmd.Run != nil {
		parsedCmd.Run(parsedCmd.Context, parsedCmd.Operands)
	}

	for _, subCmd := range parsedCmd.Subcommands {
		subCmd.Run(subCmd.Context, subCmd.Operands)
	}

	return nil
}
