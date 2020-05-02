package executor

import (
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/pkg/cli"
)

type Executor interface {
	Execute() error
}

type executor struct {
	logger     logger.Logger
	runner     cli.Runner
}

func New(l logger.Logger, r cli.Runner) Executor {
	return &executor{
		logger:     l,
		runner:     r,
	}
}

func (e *executor) Execute() error {
	return e.runner.Run()
}
