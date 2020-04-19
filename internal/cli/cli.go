package cli

import (
	"fmt"
	"github.com/sebuckler/teel/internal/cli/commands"
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/internal/scaffolder"
	"github.com/spf13/cobra"
	"os"
)

type CommandLineInterface interface {
	Execute()
}

type commandLineInterface struct {
	logger     logger.Logger
	rootCmd    *cobra.Command
	scaffolder scaffolder.Scaffolder
}

func New(d string, l logger.Logger, s scaffolder.Scaffolder, v string) CommandLineInterface {
	rootCmd := commands.NewRoot(v)

	rootCmd.AddCommand(commands.NewCreate(s, l, d).Command)

	return &commandLineInterface{
		logger:     l,
		rootCmd:    rootCmd,
		scaffolder: s,
	}
}

func (c *commandLineInterface) Execute() {
	c.logger.Log("--commandLineInterface.Execute()")

	if err := c.rootCmd.Execute(); err != nil {
		fmt.Printf("Error: failed to execute command: %v\n", err)
		os.Exit(1)
	}

	c.logger.Log("--finished commandLineInterface.Execute()")
}
