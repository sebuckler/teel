package executor

import (
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/internal/scaffolder"
)

type Executor interface {
	Execute() error
}

type executor struct {
	logger     logger.Logger
	scaffolder scaffolder.Scaffolder
	version    string
}

func New(l logger.Logger, s scaffolder.Scaffolder, v string) Executor {
	return &executor{
		logger:     l,
		scaffolder: s,
		version:    v,
	}
}

func (e *executor) Execute() error {
	return nil
}
