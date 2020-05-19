package executor

import (
	"fmt"
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/pkg/cli"
	"os"
	"os/signal"
)

type ExitFunc func()

type Executor interface {
	Execute() ExitFunc
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

func (e *executor) Execute() ExitFunc {
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
		fmt.Printf("Error: %v\n", err)

		return func() { os.Exit(1) }
	}

	return func() { os.Exit(0) }
}
