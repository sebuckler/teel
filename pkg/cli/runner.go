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
	parsedCommands, parseErr := r.parser.Parse()

	if parseErr != nil {
		return parseErr
	}

	if len(parsedCommands) == 0 {
		return errors.New("no commands parsed")
	}

	rootCmd := parsedCommands[0]

	if rootCmd == nil {
		return errors.New("no root command parsed")
	}

	if rootCmd.HelpMode {
		return rootCmd.HelpFunc(rootCmd.Syntax, r.writer)
	}

	if rootCmd.VersionMode {
		_, writeErr := r.writer.Write([]byte(rootCmd.Name + " " + r.version + "\n"))

		return writeErr
	}

	for _, cmd := range parsedCommands {
		if cmd.Run == nil {
			continue
		}

		cmd.Run(cmd.Context, cmd.Operands)
	}

	return nil
}
