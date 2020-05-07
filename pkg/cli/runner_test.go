package cli_test

import (
	"context"
	"errors"
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
		"should error when no cmd configured":                     shouldErrorWhenNoCmdConfigured,
		"should error when parser errors":                         shouldErrorWhenParserErrors,
		"should error when no root cmd parsed":                    shouldErrorWhenNoRootCmdParsed,
		"should do nothing when only root cmd parsed with no run": shouldDoNothingWhenOnlyRootCmdParsedWithNoRun,
		"should run root cmd run":                                 shouldRunRootCmdRun,
		"should run subcommand runs":                              shouldRunSubcommandRuns,
	}
}

func shouldErrorWhenNoCmdConfigured(t *testing.T, n string) {
	configurer := &fakeConfigurer{
		configure: func() *cli.CommandConfig {
			return nil
		},
	}
	runner := cli.NewRunner(configurer, nil, "v1", cli.ExitOnError)
	runErr := runner.Run()

	if runErr == nil {
		t.Fail()
		t.Log(n + ": failed to error when config is nil")
	}
}

func shouldErrorWhenParserErrors(t *testing.T, n string) {
	configurer := &fakeConfigurer{
		configure: func() *cli.CommandConfig {
			return nil
		},
	}
	parser := &fakeParser{
		parse: func(*cli.CommandConfig) (*cli.ParsedCommand, error) {
			return nil, errors.New("wrecked")
		},
	}
	runner := cli.NewRunner(configurer, parser, "v1", cli.ExitOnError)
	runErr := runner.Run()

	if runErr == nil {
		t.Fail()
		t.Log(n + ": failed to error when parser errored")
	}
}

func shouldErrorWhenNoRootCmdParsed(t *testing.T, n string) {
	configurer := &fakeConfigurer{
		configure: func() *cli.CommandConfig {
			return nil
		},
	}
	parser := &fakeParser{
		parse: func(*cli.CommandConfig) (*cli.ParsedCommand, error) {
			return nil, nil
		},
	}
	runner := cli.NewRunner(configurer, parser, "v1", cli.ExitOnError)
	runErr := runner.Run()

	if runErr == nil {
		t.Fail()
		t.Log(n + ": failed to error when root command not parsed")
	}
}

func shouldDoNothingWhenOnlyRootCmdParsedWithNoRun(t *testing.T, n string) {
	configurer := &fakeConfigurer{
		configure: func() *cli.CommandConfig {
			return &cli.CommandConfig{}
		},
	}
	parser := &fakeParser{
		parse: func(*cli.CommandConfig) (*cli.ParsedCommand, error) {
			return &cli.ParsedCommand{}, nil
		},
	}
	runner := cli.NewRunner(configurer, parser, "v1", cli.ExitOnError)
	runErr := runner.Run()

	if runErr != nil {
		t.Fail()
		t.Log(n + ": incorrectly errored on empty root cmd")
	}
}

func shouldRunRootCmdRun(t *testing.T, n string) {
	runResult := 0
	configurer := &fakeConfigurer{
		configure: func() *cli.CommandConfig {
			return &cli.CommandConfig{
				Run: func(context.Context, []string) {
					runResult = 1
				},
			}
		},
	}
	parser := &fakeParser{
		parse: func(c *cli.CommandConfig) (*cli.ParsedCommand, error) {
			return &cli.ParsedCommand{
				Run: c.Run,
			}, nil
		},
	}
	runner := cli.NewRunner(configurer, parser, "v1", cli.ExitOnError)
	runErr := runner.Run()

	if runErr != nil || runResult != 1 {
		t.Fail()
		t.Log(n + ": failed to run or incorrectly errored on root cmd run")
	}
}

func shouldRunSubcommandRuns(t *testing.T, n string) {
	var runResults []int
	configurer := &fakeConfigurer{
		configure: func() *cli.CommandConfig {
			return &cli.CommandConfig{
				Subcommands: []*cli.CommandConfig{
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
			}
		},
	}
	parser := &fakeParser{
		parse: func(c *cli.CommandConfig) (*cli.ParsedCommand, error) {
			return &cli.ParsedCommand{
				Subcommands: []*cli.ParsedCommand{
					{
						Run: c.Subcommands[0].Run,
					},
					{
						Run: c.Subcommands[1].Run,
					},
				},
			}, nil
		},
	}
	runner := cli.NewRunner(configurer, parser, "v1", cli.ExitOnError)
	runErr := runner.Run()

	if runErr != nil || len(runResults) != 2 {
		t.Fail()
		t.Log(n + ": failed to run incorrectly errored on subcommand runs")
	}
}

type fakeConfigurer struct {
	configure func() *cli.CommandConfig
}

func (f *fakeConfigurer) AddSubcommand(cli.CommandConfigurer)                              {}
func (f *fakeConfigurer) AddRunFunc(cli.CommandRunFunc)                                    {}
func (f *fakeConfigurer) AddBoolArg(string, rune, *bool, bool, string, bool)               {}
func (f *fakeConfigurer) AddIntArg(string, rune, *int, int, string, bool)                  {}
func (f *fakeConfigurer) AddIntListArg(string, rune, *[]int, []int, string, bool)          {}
func (f *fakeConfigurer) AddInt64Arg(string, rune, *int64, int64, string, bool)            {}
func (f *fakeConfigurer) AddInt64ListArg(string, rune, *[]int64, []int64, string, bool)    {}
func (f *fakeConfigurer) AddStringArg(string, rune, *string, string, string, bool)         {}
func (f *fakeConfigurer) AddStringListArg(string, rune, *[]string, []string, string, bool) {}
func (f *fakeConfigurer) AddUintArg(string, rune, *uint, uint, string, bool)               {}
func (f *fakeConfigurer) AddUintListArg(string, rune, *[]uint, []uint, string, bool)       {}
func (f *fakeConfigurer) AddUint64Arg(string, rune, *uint64, uint64, string, bool)         {}
func (f *fakeConfigurer) AddUint64ListArg(string, rune, *[]uint64, []uint64, string, bool) {}

func (f *fakeConfigurer) Configure() *cli.CommandConfig {
	return f.configure()
}

type fakeParser struct {
	parse func(c *cli.CommandConfig) (*cli.ParsedCommand, error)
}

func (f *fakeParser) Parse(c *cli.CommandConfig) (*cli.ParsedCommand, error) {
	return f.parse(c)
}
