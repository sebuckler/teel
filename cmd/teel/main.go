package main

import (
	"fmt"
	"github.com/sebuckler/teel/internal/executor"
	"github.com/sebuckler/teel/internal/filestreamer"
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/internal/scaffolder"
	"github.com/sebuckler/teel/internal/scaffolder/directives"
	"github.com/sebuckler/teel/internal/sighandler"
	"os"
)

var version string

func main() {
	fileStreamer := filestreamer.New(filestreamer.HOME, "teel_logs", "teel.log")
	streamErr := fileStreamer.Stream(func(f *filestreamer.StreamFile) {
		signalHandler := sighandler.New(os.Kill, os.Interrupt)
		fileLogger := logger.New(f, nil)
		siteScaffolder := scaffolder.New(directives.NewConfig())
		cmdExecutor := executor.New(fileLogger, siteScaffolder, version)

		execErr := cmdExecutor.Execute()

		if execErr != nil {
			fmt.Printf("Error: %v\n", execErr)
			os.Exit(1)
		}

		signalHandler.Handle(func(os.Signal) { _ = f.Flush() })
	})

	if streamErr != nil {
		fmt.Printf("Error: %v\n", streamErr)
		os.Exit(1)
	}
}
