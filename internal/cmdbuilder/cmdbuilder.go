package cmdbuilder

import (
	"context"
	"fmt"
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/internal/scaffolder"
	"github.com/sebuckler/teel/pkg/cli"
)

type CommandBuilder interface {
	Build() cli.CommandConfigurer
}

type commandBuilder struct {
	logger     logger.Logger
	scaffolder scaffolder.Scaffolder
}

func New(l logger.Logger, s scaffolder.Scaffolder) CommandBuilder {
	return &commandBuilder{
		logger:     l,
		scaffolder: s,
	}
}

func (c *commandBuilder) Build() cli.CommandConfigurer {
	var a bool
	var b string
	rootCmd := cli.NewCommand("", context.Background())
	rootCmd.AddBoolArg("a", 'a', &a, false, "Is first letter of alphabet", false)
	rootCmd.AddStringArg("b", 'b', &b, "second", "Second letter of alphabet", false)
	rootCmd.AddRunFunc(func(ctx context.Context) {
		fmt.Print("a: ")
		fmt.Print(a)
		fmt.Print(", b: ")
		fmt.Print(b + "\n")
		fmt.Println("welcome to the thunderdome")
	})
	subCmd := cli.NewCommand("subby", context.Background())
	subCmd.AddRunFunc(func(ctx context.Context) {
		fmt.Println("and me!")
	})
	rootCmd.AddSubcommand(subCmd)

	return rootCmd
}
