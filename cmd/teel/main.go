package main

import (
	"fmt"
	"github.com/sebuckler/teel/internal/cmdbuilder"
	"github.com/sebuckler/teel/internal/executor"
	"github.com/sebuckler/teel/internal/filestreamer"
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/internal/scaffolder"
	"github.com/sebuckler/teel/internal/scaffolder/directives"
	"github.com/sebuckler/teel/internal/sighandler"
	"github.com/sebuckler/teel/pkg/cli"
	"os"
)

var version string

func main() {
	fileStreamer := filestreamer.New(filestreamer.HomeDir, "teel_logs", "teel.log")
	streamErr := fileStreamer.Stream(func(f *filestreamer.StreamFile) {
		signalHandler := sighandler.New(os.Kill, os.Interrupt)
		fileLogger := logger.New(f, nil)
		siteScaffolder := scaffolder.New(directives.NewConfig())
		cmdBuilder := cmdbuilder.New(fileLogger, siteScaffolder)
		parser := cli.NewParser(cli.GNU)
		runner := cli.NewRunner(version)
		cmdExecutor := executor.New(cmdBuilder, fileLogger, parser, runner)

		signalHandler.Handle(func(os.Signal) { _ = f.Flush() })

		if execErr := cmdExecutor.Execute(); execErr != nil {
			fmt.Printf("Error: %v\n", execErr)
			os.Exit(1)
		}
	})

	if streamErr != nil {
		fmt.Printf("Error: %v\n", streamErr)
		os.Exit(1)
	}
}
