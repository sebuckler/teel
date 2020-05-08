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
	rootCmd.AddBoolArg(&a, false, &cli.ArgDefinition{
		Name:       "a",
		ShortName:  'a',
	})
	rootCmd.AddStringArg(&b, "second", &cli.ArgDefinition{
		Name:       "b",
		ShortName:  'b',
	})
	rootCmd.AddRunFunc(func(ctx context.Context, o []string) {
		fmt.Print("a: ")
		fmt.Print(a)
		fmt.Print(", b: ")
		fmt.Print(b + "\n")
		fmt.Println("welcome to the thunderdome")
	})
	subCmd := cli.NewCommand("subby", context.Background())
	file := "filename"
	subCmd.AddStringArg(&file, "", &cli.ArgDefinition{
		Name:       "file",
		ShortName:  'f',
	})
	subCmd.AddRunFunc(func(ctx context.Context, o []string) {
		fmt.Println("and me!")
	})
	rootCmd.AddSubcommand(subCmd)

	return rootCmd
}
