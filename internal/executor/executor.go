package executor

import (
	"context"
	"fmt"
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/internal/scaffolder"
	"github.com/sebuckler/teel/pkg/cli"
)

type Executor interface {
	Execute() error
}

type executor struct {
	logger     logger.Logger
	runner     cli.Runner
	scaffolder scaffolder.Scaffolder
	version    string
}

func New(l logger.Logger, s scaffolder.Scaffolder, v string) Executor {
	rootCmd := cli.NewCommand("", context.Background())
	rootCmd.AddRunFunc(func(ctx context.Context) {
		fmt.Println("welcome to the thunderdome")
	})
	parser := cli.NewParser(cli.Posix, cli.Error)
	runner := cli.NewRunner(rootCmd, parser, v, cli.ExitOnError)

	return &executor{
		logger:     l,
		runner:     runner,
		scaffolder: s,
		version:    v,
	}
}

func (e *executor) Execute() error {
	return e.runner.Run()
}
