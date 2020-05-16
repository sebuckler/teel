package cli_test

import (
	"context"
	"github.com/sebuckler/teel/pkg/cli"
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
		"should run root cmd run":                                 shouldRunRootCmdRun,
		"should run subcommand runs":                              shouldRunSubcommandRuns,
	}
}

func shouldErrorWhenNoRootCmdParsed(t *testing.T, n string) {
	runner := cli.NewRunner("v1")
	runErr := runner.Run(nil)

	if runErr == nil {
		t.Fail()
		t.Log(n + ": failed to error when root command not parsed")
	}
}

func shouldDoNothingWhenOnlyRootCmdParsedWithNoRun(t *testing.T, n string) {
	runner := cli.NewRunner("v1")
	runErr := runner.Run(&cli.ParsedCommand{})

	if runErr != nil {
		t.Fail()
		t.Log(n + ": incorrectly errored on empty root cmd")
	}
}

func shouldRunRootCmdRun(t *testing.T, n string) {
	runResult := 0
	runner := cli.NewRunner("v1")
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
	runner := cli.NewRunner("v1")
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
