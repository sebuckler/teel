package main

import (
	"fmt"
	"github.com/sebuckler/teel/internal/cmdbuilder"
	"github.com/sebuckler/teel/internal/executor"
	"github.com/sebuckler/teel/internal/filestreamer"
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/internal/scaffolder"
	"github.com/sebuckler/teel/internal/scaffolder/directives"
	"github.com/sebuckler/teel/pkg/cli"
	"os"
)

var version string

func main() {
	fileStreamer := filestreamer.New(filestreamer.HomeDir, "teel_logs", "teel.log")
	streamErr := fileStreamer.Stream(func(f *filestreamer.StreamFile) {
		fileLogger := logger.New(f, nil)
		siteScaffolder := scaffolder.New(directives.NewConfig())
		cmdBuilder := cmdbuilder.New(fileLogger, siteScaffolder)
		parser := cli.NewParser(cli.GNU, cmdBuilder.Build())
		runner := cli.NewRunner(parser, version, os.Stdout)
		cmdExecutor := executor.New(fileLogger, runner)
		exit := cmdExecutor.Execute()
		_ = f.Flush()

		exit()
	})

	if streamErr != nil {
		fmt.Printf("Error: %v\n", streamErr)
		os.Exit(1)
	}
}
