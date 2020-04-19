package main

import (
	"fmt"
	"github.com/sebuckler/teel/internal/cli"
	"github.com/sebuckler/teel/internal/filestreamer"
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/internal/scaffolder"
	"github.com/sebuckler/teel/internal/sighandler"
	"os"
)

var version string

func main() {
	fileStreamer := filestreamer.New(filestreamer.HOME, "teel_logs", "teel.log")
	streamErr := fileStreamer.Stream(func(f *filestreamer.StreamFile) {
		signalHandler := sighandler.New()
		fileLogger := logger.New(f, nil)
		blogScaffolder := scaffolder.New(fileLogger)

		signalHandler.Handle(func(os.Signal) { _ = f.Flush() }, os.Kill, os.Interrupt)
		cli.New(f.Cwd, fileLogger, blogScaffolder, version).Execute()
	})

	if streamErr != nil {
		fmt.Printf("Error: %v\n", streamErr)
		os.Exit(1)
	}
}
