package executor

import (
	"github.com/sebuckler/teel/internal/cmdbuilder"
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/pkg/cli"
)

type Executor interface {
	Execute() error
}

type executor struct {
	cmdBuilder cmdbuilder.CommandBuilder
	logger     logger.Logger
	parser     cli.Parser
	runner     cli.Runner
}

func New(c cmdbuilder.CommandBuilder, l logger.Logger, p cli.Parser, r cli.Runner) Executor {
	return &executor{
		cmdBuilder: c,
		logger:     l,
		parser:     p,
		runner:     r,
	}
}

func (e *executor) Execute() error {
	cmd := e.cmdBuilder.Build()
	parsedCmd, parseErr := e.parser.Parse(cmd.Configure())

	if parseErr != nil {
		return parseErr
	}

	return e.runner.Run(parsedCmd)
}
