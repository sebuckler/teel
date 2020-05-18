package cmdbuilder

import (
	"context"
	"fmt"
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/internal/scaffolder"
	"github.com/sebuckler/teel/pkg/cli"
)

type CommandBuilder interface {
	Build() cli.CommandBuilder
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

func (c *commandBuilder) Build() cli.CommandBuilder {
	var a bool
	var b string
	rootCmd := cli.NewCommand("teel", context.Background())
	rootCmd.AddBoolArg(&a, &cli.ArgDefinition{
		Name:      "a",
		ShortName: 'a',
	})
	rootCmd.AddStringArg(&b, &cli.ArgDefinition{
		Name:      "b",
		ShortName: 'b',
	})
	rootCmd.AddRunFunc(func(ctx context.Context, o []string) {
		fmt.Print("a: ", a, ", b: ", b, "\n")
		fmt.Println("welcome to the thunderdome")
	})
	subCmd := cli.NewCommand("subby", context.Background())
	file := "filename"
	subCmd.AddStringArg(&file, &cli.ArgDefinition{
		Name:      "file",
		ShortName: 'f',
	})
	subCmd.AddRunFunc(func(ctx context.Context, o []string) {
		fmt.Println("and me, " + file + "!")
	})
	rootCmd.AddSubcommand(subCmd)

	return rootCmd
}
