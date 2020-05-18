package executor

import (
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/pkg/cli"
	"os"
	"os/signal"
)

type Executor interface {
	Execute()
}

type executor struct {
	logger logger.Logger
	runner cli.Runner
}

func New(l logger.Logger, r cli.Runner) Executor {
	return &executor{
		logger: l,
		runner: r,
	}
}

func (e *executor) Execute() {
	sigChan := make(chan os.Signal, 1)
	done := make(chan error, 1)

	signal.Notify(sigChan, os.Interrupt, os.Kill)

	go func() {
		<-sigChan
		done <- nil
	}()

	go func() {
		done <- e.runner.Run()
	}()

	if err := <-done; err != nil {
		e.logger.Errorf("Error: %v\n", err)
		os.Exit(1)
	}
}
