package main

import (
	"github.com/sebuckler/teel/internal/cmdbuilder"
	"github.com/sebuckler/teel/internal/executor"
	"github.com/sebuckler/teel/internal/fs"
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/internal/scaffolder"
	"github.com/sebuckler/teel/internal/scaffolder/directives"
	"github.com/sebuckler/teel/pkg/cli"
	"os"
)

var version string

func main() {
	file, fileErr := fs.OpenBufferedFileWriter("teel.log")

	if fileErr != nil {
		executor.Exit(fileErr)
	}

	fileLogger := logger.New(file, nil)
	siteScaffolder := scaffolder.New(directives.NewConfig())
	cmdBuilder := cmdbuilder.New(fileLogger, siteScaffolder)
	parser := cli.NewParser(cli.GNU, cmdBuilder.Build())
	runner := cli.NewRunner(parser, version, os.Stdout)
	cmdExecutor := executor.New(fileLogger, runner)
	execErr := cmdExecutor.Execute()
	_ = file.Close()

	executor.Exit(execErr)
}
