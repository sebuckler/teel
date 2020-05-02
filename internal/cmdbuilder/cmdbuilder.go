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
	logger logger.Logger
	scaffolder scaffolder.Scaffolder
}

func New(l logger.Logger, s scaffolder.Scaffolder) CommandBuilder {
	return &commandBuilder{
		logger: l,
		scaffolder: s,
	}
}

func (c *commandBuilder) Build() cli.CommandConfigurer {
	rootCmd := cli.NewCommand("", context.Background())
	rootCmd.AddRunFunc(func(ctx context.Context) {
		fmt.Println("welcome to the thunderdome")
	})
	subCmd := cli.NewCommand("subby", context.Background())
	subCmd.AddRunFunc(func(ctx context.Context) {
		fmt.Println("and me!")
	})
	rootCmd.AddSubcommand(subCmd)

	return rootCmd
}
