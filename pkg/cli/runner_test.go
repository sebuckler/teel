package cli_test

import (
	"bufio"
	"context"
	"github.com/sebuckler/teel/pkg/cli"
	"io"
	"strings"
	"testing"
)

func TestRunner_Run(t *testing.T) {
	for name, test := range getRunnerTestCases() {
		test(t, name)
	}
}

func getRunnerTestCases() map[string]func(t *testing.T, n string) {
	return map[string]func(t *testing.T, n string){
		"should error when no root cmd parsed":                    shouldErrorWhenNoRootCmdParsed,
		"should do nothing when only root cmd parsed with no run": shouldDoNothingWhenOnlyRootCmdParsedWithNoRun,
		"should not run command when help arg exists":             shouldNotRunCommandWhenHelpArgExists,
		"should print help text when help mode is true":           shouldPrintTextWhenHelpModeIsTrue,
		"should run root cmd run":                                 shouldRunRootCmdRun,
		"should run subcommand runs":                              shouldRunSubcommandRuns,
	}
}

func shouldErrorWhenNoRootCmdParsed(t *testing.T, n string) {
	var strBuilder strings.Builder
	writer := bufio.NewWriter(&strBuilder)
	runner := cli.NewRunner("v1", writer)
	runErr := runner.Run(nil)

	if runErr == nil {
		t.Fail()
		t.Log(n + ": failed to error when root command not parsed")
	}
}

func shouldDoNothingWhenOnlyRootCmdParsedWithNoRun(t *testing.T, n string) {
	var strBuilder strings.Builder
	writer := bufio.NewWriter(&strBuilder)
	runner := cli.NewRunner("v1", writer)
	runErr := runner.Run(&cli.ParsedCommand{})

	if runErr != nil {
		t.Fail()
		t.Log(n + ": incorrectly errored on empty root cmd")
	}
}

func shouldNotRunCommandWhenHelpArgExists(t *testing.T, n string) {
	runResult := 0
	var strBuilder strings.Builder
	writer := bufio.NewWriter(&strBuilder)
	runner := cli.NewRunner("v1", writer)
	runErr := runner.Run(&cli.ParsedCommand{
		HelpFunc: func(s cli.ArgSyntax, w io.Writer) error {
			return nil
		},
		HelpMode: true,
		Run: func(context.Context, []string) {
			runResult = 1
		},
	})

	if runErr != nil || runResult == 1 {
		t.Fail()
		t.Log(n + ": failed to display usage text correctly")
	}
}

func shouldPrintTextWhenHelpModeIsTrue(t *testing.T, n string) {
	var strBuilder strings.Builder
	writer := bufio.NewWriter(&strBuilder)
	runner := cli.NewRunner("v1", writer)
	runErr := runner.Run(&cli.ParsedCommand{
		HelpFunc: func(s cli.ArgSyntax, w io.Writer) error {
			_, _ = w.Write([]byte("help called"))
			return nil
		},
		HelpMode: true,
	})
	_ = writer.Flush()

	if runErr != nil || strBuilder.String() == "" {
		t.Fail()
		t.Log(n + ": failed to run or incorrectly errored on root cmd run")
	}
}

func shouldRunRootCmdRun(t *testing.T, n string) {
	runResult := 0
	var strBuilder strings.Builder
	writer := bufio.NewWriter(&strBuilder)
	runner := cli.NewRunner("v1", writer)
	runErr := runner.Run(&cli.ParsedCommand{
		Run: func(context.Context, []string) {
			runResult = 1
		},
	})

	if runErr != nil || runResult != 1 {
		t.Fail()
		t.Log(n + ": failed to run or incorrectly errored on root cmd run")
	}
}

func shouldRunSubcommandRuns(t *testing.T, n string) {
	var runResults []int
	var strBuilder strings.Builder
	writer := bufio.NewWriter(&strBuilder)
	runner := cli.NewRunner("v1", writer)
	runErr := runner.Run(&cli.ParsedCommand{
		Subcommands: []*cli.ParsedCommand{
			{
				Run: func(context.Context, []string) {
					runResults = append(runResults, 1)
				},
			},
			{
				Run: func(context.Context, []string) {
					runResults = append(runResults, 1)
				},
			},
		},
	})

	if runErr != nil || len(runResults) != 2 {
		t.Fail()
		t.Log(n + ": failed to run incorrectly errored on subcommand runs")
	}
}
