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
		"should error when POSIX bool opt has opt-arg":               shouldErrorWhenPosixBoolOptHasOptArg,
		"should error when POSIX int opt has no opt-arg":             shouldErrorWhenPosixIntOptHasNoOptArg,
		"should error when POSIX int list opt has no opt-arg":        shouldErrorWhenPosixIntListOptHasNoOptArg,
		"should error when POSIX int64 opt has no opt-arg":           shouldErrorWhenPosixInt64OptHasNoOptArg,
		"should error when POSIX int64 list opt has no opt-arg":      shouldErrorWhenPosixInt64ListOptHasNoOptArg,
		"should error when POSIX string opt has no opt-arg":          shouldErrorWhenPosixStringOptHasNoOptArg,
		"should error when POSIX string list opt has no opt-arg":     shouldErrorWhenPosixStringListOptHasNoOptArg,
		"should error when POSIX uint opt has no opt-arg":            shouldErrorWhenPosixUintOptHasNoOptArg,
		"should error when POSIX uint list opt has no opt-arg":       shouldErrorWhenPosixUintListOptHasNoOptArg,
		"should error when POSIX uint64 opt has no opt-arg":          shouldErrorWhenPosixUint64OptHasNoOptArg,
		"should error when POSIX uint64 list opt has no opt-arg":     shouldErrorWhenPosixUint64ListOptHasNoOptArg,
		"should error when unsupported arg type used":                shouldErrorWhenUnsupportedArgTypeUsed,
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

func shouldErrorWhenPosixBoolOptHasOptArg(t *testing.T, n string) {
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

func shouldErrorWhenPosixIntOptHasNoOptArg(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-a"}
	parser := cli.NewParser(cli.POSIX)
	a := 1
	_, parseErr := parser.Parse(&cli.CommandConfig{
		Args: []*cli.ArgConfig{{
			Name:      "a",
			ShortName: 'a',
			Value:     &a,
		}},
	})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error on missing POSIX int option-argument")
	}
}

func shouldErrorWhenPosixIntListOptHasNoOptArg(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-a"}
	parser := cli.NewParser(cli.POSIX)
	a := []int{1}
	_, parseErr := parser.Parse(&cli.CommandConfig{
		Args: []*cli.ArgConfig{{
			Name:      "a",
			ShortName: 'a',
			Value:     &a,
		}},
	})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error on missing POSIX int list option-argument")
	}
}

func shouldErrorWhenPosixInt64OptHasNoOptArg(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-a"}
	parser := cli.NewParser(cli.POSIX)
	a := int64(1)
	_, parseErr := parser.Parse(&cli.CommandConfig{
		Args: []*cli.ArgConfig{{
			Name:      "a",
			ShortName: 'a',
			Value:     &a,
		}},
	})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error on missing POSIX int64 option-argument")
	}
}

func shouldErrorWhenPosixInt64ListOptHasNoOptArg(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-a"}
	parser := cli.NewParser(cli.POSIX)
	a := []int64{1}
	_, parseErr := parser.Parse(&cli.CommandConfig{
		Args: []*cli.ArgConfig{{
			Name:      "a",
			ShortName: 'a',
			Value:     &a,
		}},
	})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error on missing POSIX int64 list option-argument")
	}
}

func shouldErrorWhenPosixStringOptHasNoOptArg(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-a"}
	parser := cli.NewParser(cli.POSIX)
	a := "value"
	_, parseErr := parser.Parse(&cli.CommandConfig{
		Args: []*cli.ArgConfig{{
			Name:      "a",
			ShortName: 'a',
			Value:     &a,
		}},
	})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error on missing POSIX string option-argument")
	}
}

func shouldErrorWhenPosixStringListOptHasNoOptArg(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-a"}
	parser := cli.NewParser(cli.POSIX)
	a := []string{"value"}
	_, parseErr := parser.Parse(&cli.CommandConfig{
		Args: []*cli.ArgConfig{{
			Name:      "a",
			ShortName: 'a',
			Value:     &a,
		}},
	})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error on missing POSIX string list option-argument")
	}
}

func shouldErrorWhenPosixUintOptHasNoOptArg(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-a"}
	parser := cli.NewParser(cli.POSIX)
	a := uint(1)
	_, parseErr := parser.Parse(&cli.CommandConfig{
		Args: []*cli.ArgConfig{{
			Name:      "a",
			ShortName: 'a',
			Value:     &a,
		}},
	})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error on missing POSIX uint option-argument")
	}
}

func shouldErrorWhenPosixUintListOptHasNoOptArg(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-a"}
	parser := cli.NewParser(cli.POSIX)
	a := []uint{1}
	_, parseErr := parser.Parse(&cli.CommandConfig{
		Args: []*cli.ArgConfig{{
			Name:      "a",
			ShortName: 'a',
			Value:     &a,
		}},
	})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error on missing POSIX uint list option-argument")
	}
}

func shouldErrorWhenPosixUint64OptHasNoOptArg(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-a"}
	parser := cli.NewParser(cli.POSIX)
	a := uint64(1)
	_, parseErr := parser.Parse(&cli.CommandConfig{
		Args: []*cli.ArgConfig{{
			Name:      "a",
			ShortName: 'a',
			Value:     &a,
		}},
	})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error on missing POSIX uint64 option-argument")
	}
}

func shouldErrorWhenPosixUint64ListOptHasNoOptArg(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-a"}
	parser := cli.NewParser(cli.POSIX)
	a := []uint64{1}
	_, parseErr := parser.Parse(&cli.CommandConfig{
		Args: []*cli.ArgConfig{{
			Name:      "a",
			ShortName: 'a',
			Value:     &a,
		}},
	})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error on missing POSIX uint64 list option-argument")
	}
}

func shouldErrorWhenUnsupportedArgTypeUsed(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-a", "1"}
	parser := cli.NewParser(cli.POSIX)
	a := byte(1)
	_, parseErr := parser.Parse(&cli.CommandConfig{
		Args: []*cli.ArgConfig{{
			Name:      "a",
			ShortName: 'a',
			Value:     &a,
		}},
	})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error on unsupported arg type")
	}
}
