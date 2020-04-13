package main

import (
	"fmt"
	"github.com/sebuckler/teel/internal/cli"
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/internal/scaffolder"
	"os"
	"path"
)

var version string

func main() {
	cwd, cwdErr := os.Getwd()

	if cwdErr != nil {
		fmt.Println("Error: unreachable working directory")
		os.Exit(1)
	}

	homeDir, homeDirErr := os.UserHomeDir()

	if homeDirErr != nil {
		homeDir = cwd
	}

	logDirPath := path.Join(homeDir, "teel_logs")
	_ = os.MkdirAll(logDirPath, 0777)
	logFilePath := path.Join(logDirPath, "teel.log")
	logFile, fileErr := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if fileErr != nil {
		fmt.Printf("Error: opening log file: %v\n", fileErr)
		os.Exit(1)
	}

	fileLogger := logger.NewFile(logFile, logger.BUFFERED)
	blogScaffolder := scaffolder.New(fileLogger)

	cli.New(version, blogScaffolder, fileLogger, cwd).Execute()
	fileLogger.Write()
}
