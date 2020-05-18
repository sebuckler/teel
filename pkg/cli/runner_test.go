package cli_test

import (
	"bufio"
	"context"
	"github.com/sebuckler/teel/pkg/cli"
	"os"
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
		"should do nothing when only root cmd parsed with no run": shouldDoNothingWhenOnlyRootCmdParsedWithNoRun,
		"should not run command when help arg exists":             shouldNotRunCommandWhenHelpArgExists,
		"should print help text when help mode is true":           shouldPrintTextWhenHelpModeIsTrue,
		"should run root cmd run":                                 shouldRunRootCmdRun,
		"should run subcommand runs":                              shouldRunSubcommandRuns,
	}
}

func shouldDoNothingWhenOnlyRootCmdParsedWithNoRun(t *testing.T, n string) {
	os.Args = []string{"testcmd"}
	var strBuilder strings.Builder
	writer := bufio.NewWriter(&strBuilder)
	parser := cli.NewParser(cli.GNU, cli.NewCommand("testcmd", context.Background()))
	runner := cli.NewRunner(parser, "v1", writer)
	runErr := runner.Run()

	if runErr != nil {
		t.Fail()
		t.Log(n + ": incorrectly errored on empty root cmd")
	}
}

func shouldNotRunCommandWhenHelpArgExists(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-h"}
	runResult := 0
	var strBuilder strings.Builder
	writer := bufio.NewWriter(&strBuilder)
	cmd := cli.NewCommand("testcmd", context.Background())
	cmd.AddRunFunc(func(context.Context, []string) { runResult = 1 })
	parser := cli.NewParser(cli.GNU, cmd)
	runner := cli.NewRunner(parser, "v1", writer)
	runErr := runner.Run()

	if runErr != nil || runResult == 1 {
		t.Fail()
		t.Log(n + ": failed to display usage text correctly")
	}
}

func shouldPrintTextWhenHelpModeIsTrue(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-h"}
	var strBuilder strings.Builder
	writer := bufio.NewWriter(&strBuilder)
	parser := cli.NewParser(cli.GNU, cli.NewCommand("testcmd", context.Background()))
	runner := cli.NewRunner(parser, "v1", writer)
	runErr := runner.Run()
	_ = writer.Flush()

	if runErr != nil || strBuilder.String() == "" {
		t.Fail()
		t.Log(n + ": failed to run or incorrectly errored on root cmd run")
	}
}

func shouldRunRootCmdRun(t *testing.T, n string) {
	os.Args = []string{"testcmd"}
	runResult := 0
	var strBuilder strings.Builder
	writer := bufio.NewWriter(&strBuilder)
	cmd := cli.NewCommand("testcmd", context.Background())
	cmd.AddRunFunc(func(context.Context, []string) { runResult = 1 })
	parser := cli.NewParser(cli.GNU, cmd)
	runner := cli.NewRunner(parser, "v1", writer)
	runErr := runner.Run()

	if runErr != nil || runResult != 1 {
		t.Fail()
		t.Log(n + ": failed to run or incorrectly errored on root cmd run")
	}
}

func shouldRunSubcommandRuns(t *testing.T, n string) {
	os.Args = []string{"testcmd", "foo", "bar"}
	var runResults []int
	var strBuilder strings.Builder
	writer := bufio.NewWriter(&strBuilder)
	cmd := cli.NewCommand("testcmd", context.Background())
	sub1 := cli.NewCommand("foo", context.Background())
	sub1.AddRunFunc(func(context.Context, []string) {
		runResults = append(runResults, 1)
	})
	sub2 := cli.NewCommand("bar", context.Background())
	sub2.AddRunFunc(func(context.Context, []string) {
		runResults = append(runResults, 1)
	})
	cmd.AddSubcommand(sub1, sub2)
	parser := cli.NewParser(cli.GNU, cmd)
	runner := cli.NewRunner(parser, "v1", writer)
	runErr := runner.Run()

	if runErr != nil || len(runResults) != 2 {
		t.Fail()
		t.Log(n + ": failed to run incorrectly errored on subcommand runs")
	}
}
