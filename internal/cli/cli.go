package cli

import (
	"fmt"
	"github.com/sebuckler/teel/internal/cli/commands"
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/internal/scaffolder"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

type CommandLineInterface interface {
	Execute()
}

type commandLineInterface struct {
	blogScaffolder scaffolder.Scaffolder
	logWriter      logger.Logger
	rootCmd        *cobra.Command
}

func New(v string, s scaffolder.Scaffolder, l logger.Logger, d string) CommandLineInterface {
	rootCmd := commands.NewRoot(v)

	rootCmd.AddCommand(commands.NewCreate(s, l, d).Command)
	handleSigterm(l)

	return &commandLineInterface{
		blogScaffolder: s,
		logWriter:      l,
		rootCmd:        rootCmd,
	}
}

func (c *commandLineInterface) Execute() {
	c.logWriter.Log("--commandLineInterface.Execute()")

	if err := c.rootCmd.Execute(); err != nil {
		fmt.Printf("Error: failed to execute command: %v\n", err)
		os.Exit(1)
	}

	c.logWriter.Log("--finished commandLineInterface.Execute()")
}

func handleSigterm(l logger.Logger) {
	signalChannel := make(chan os.Signal, 3)

	signal.Notify(signalChannel, os.Kill, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Kill:
		case os.Interrupt:
		case syscall.SIGTERM:
			l.Write()
		}
	}()
}
