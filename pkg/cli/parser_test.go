package cli_test

import (
	"github.com/sebuckler/teel/pkg/cli"
	"os"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	for name, test := range getParserTestCases() {
		test(t, name)
	}
}

func getParserTestCases() map[string]func(t *testing.T, n string) {
	return map[string]func(t *testing.T, n string){
		"should error when unsupported parse syntax used":            shouldErrorWhenUnsupportedParseSyntaxUsed,
		"should error when arg passed with no args configured":       shouldErrorWhenArgPassedWithNoArgsConfigured,
		"should error when repeated arg is not marked as repeatable": shouldErrorWhenRepeatedArgIsNotMarkedAsRepeatable,
		"should error when POSIX first arg is invalid format":        shouldErrorWhenPosixFirstArgIsInvalidFormat,
		"should error when configured POSIX arg name is invalid":     shouldErrorWhenConfiguredPosixArgNameIsInvalid,
		"should error when POSIX bool option has option-argument":    shouldErrorWhenPosixBoolOptionHasOptionArgument,
	}
}

func shouldErrorWhenUnsupportedParseSyntaxUsed(t *testing.T, n string) {
	os.Args = []string{"testcmd"}
	parser := cli.NewParser(99)
	_, parseErr := parser.Parse(&cli.CommandConfig{})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error on unsupported parse syntax")
	}
}

func shouldErrorWhenArgPassedWithNoArgsConfigured(t *testing.T, n string) {
	os.Args = []string{"testcmd", "a"}
	parser := cli.NewParser(cli.POSIX)
	_, parseErr := parser.Parse(&cli.CommandConfig{})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error when args passed in with no args configured")
	}
}

func shouldErrorWhenRepeatedArgIsNotMarkedAsRepeatable(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-a", "-a"}
	parser := cli.NewParser(cli.POSIX)
	a := true
	_, parseErr := parser.Parse(&cli.CommandConfig{
		Args: []*cli.ArgConfig{{
			Name:       "a",
			Repeatable: false,
			ShortName:  'a',
			Value:      &a,
		}},
	})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error on non-repeatable arg being repeated")
	}
}

func shouldErrorWhenPosixFirstArgIsInvalidFormat(t *testing.T, n string) {
	os.Args = []string{"testcmd", "a"}
	parser := cli.NewParser(cli.POSIX)
	a := false
	_, parseErr := parser.Parse(&cli.CommandConfig{
		Args: []*cli.ArgConfig{{
			Name:      "a",
			ShortName: 'a',
			Value:     &a,
		}},
	})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not return error with invalid POSIX argument")
	}
}

func shouldErrorWhenConfiguredPosixArgNameIsInvalid(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-="}
	parser := cli.NewParser(cli.POSIX)
	a := false
	_, parseErr := parser.Parse(&cli.CommandConfig{
		Args: []*cli.ArgConfig{{
			Name:      "=",
			ShortName: '=',
			Value:     &a,
		}},
	})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error on incorrectly configured POSIX arg name")
	}
}

func shouldErrorWhenPosixBoolOptionHasOptionArgument(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-a", "value"}
	parser := cli.NewParser(cli.POSIX)
	a := false
	_, parseErr := parser.Parse(&cli.CommandConfig{
		Args: []*cli.ArgConfig{{
			Name:      "a",
			ShortName: 'a',
			Value:     &a,
		}},
	})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error on POSIX bool option-argument")
	}
}
