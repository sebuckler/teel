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
		"should error when repeated arg is not repeatable":           shouldErrorWhenRepeatedArgIsNotRepeatable,
		"should error when GoFlag first arg is invalid format":       shouldErrorWhenGoFlagFirstArgIsInvalidFormat,
		"should parse when GoFlag args provided correctly":           shouldParseWhenGoFlagArgsProvidedCorrectly,
		"should error when GNU first arg is invalid format":          shouldErrorWhenGnuFirstArgIsInvalidFormat,
		"should error when configured GNU arg name is invalid":       shouldErrorWhenConfiguredGnuArgNameIsInvalid,
		"should error when GNU optional opt-arg is invalid format":   shouldErrorWhenGnuOptionalOptArgIsInvalidFormat,
		"should parse when GNU args provided correctly":              shouldParseWhenGnuArgsProvidedCorrectly,
		"should error when POSIX first arg is invalid format":        shouldErrorWhenPosixFirstArgIsInvalidFormat,
		"should error when configured POSIX arg name is invalid":     shouldErrorWhenConfiguredPosixArgNameIsInvalid,
		"should parse when POSIX args provided correctly":            shouldParseWhenPosixArgsProvidedCorrectly,
		"should error when bool opt has opt-arg":                     shouldErrorWhenBoolOptHasOptArg,
		"should error when required float64 opt has no opt-arg":      shouldErrorWhenRequiredFloat64OptHasNoOptArg,
		"should error when required float64 list opt has no opt-arg": shouldErrorWhenRequiredFloat64ListOptHasNoOptArg,
		"should error when required int opt has no opt-arg":          shouldErrorWhenRequiredIntOptHasNoOptArg,
		"should error when required int list opt has no opt-arg":     shouldErrorWhenRequiredIntListOptHasNoOptArg,
		"should error when required int64 opt has no opt-arg":        shouldErrorWhenRequiredInt64OptHasNoOptArg,
		"should error when required int64 list opt has no opt-arg":   shouldErrorWhenRequiredInt64ListOptHasNoOptArg,
		"should error when required string opt has no opt-arg":       shouldErrorWhenRequiredStringOptHasNoOptArg,
		"should error when required string list opt has no opt-arg":  shouldErrorWhenRequiredStringListOptHasNoOptArg,
		"should error when required uint opt has no opt-arg":         shouldErrorWhenRequiredUintOptHasNoOptArg,
		"should error when required uint list opt has no opt-arg":    shouldErrorWhenRequiredUintListOptHasNoOptArg,
		"should error when required uint64 opt has no opt-arg":       shouldErrorWhenRequiredUint64OptHasNoOptArg,
		"should error when required uint64 list opt has no opt-arg":  shouldErrorWhenRequiredUint64ListOptHasNoOptArg,
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

func shouldErrorWhenRepeatedArgIsNotRepeatable(t *testing.T, n string) {
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

func shouldErrorWhenGoFlagFirstArgIsInvalidFormat(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-a=value"}
	parser := cli.NewParser(cli.GoFlag)
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
		t.Log(n + ": did not return error with invalid GNU argument")
	}
}

func shouldParseWhenGoFlagArgsProvidedCorrectly(t *testing.T, n string) {
	cmd := "testcmd"
	testCases := map[string]struct {
		args  []string
		value func() []*cli.ArgConfig
	}{
		"single bool short option": {[]string{cmd, "-a"}, func() []*cli.ArgConfig {
			val := false
			return []*cli.ArgConfig{{Name: "a", Value: &val}}
		}},
		"single bool short option with value": {[]string{cmd, "-a=true"}, func() []*cli.ArgConfig {
			val := false
			return []*cli.ArgConfig{{Name: "a", Value: &val}}
		}},
		"single bool long option": {[]string{cmd, "--aaa"}, func() []*cli.ArgConfig {
			val := false
			return []*cli.ArgConfig{{Name: "aaa", Value: &val}}
		}},
		"single bool long option with value": {[]string{cmd, "--aaa=true"}, func() []*cli.ArgConfig {
			val := false
			return []*cli.ArgConfig{{Name: "aaa", Value: &val}}
		}},
		"single float64 option": {[]string{cmd, "--aaa", "1.0"}, func() []*cli.ArgConfig {
			val := float64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}}
		}},
		"multiple float64 options": {[]string{cmd, "--aaa", "1.0", "--bbb", "2.0"}, func() []*cli.ArgConfig {
			val := float64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}, {Name: "bbb", Required: true, Value: &val}}
		}},
		"single int option": {[]string{cmd, "--aaa", "1"}, func() []*cli.ArgConfig {
			val := 0
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}}
		}},
		"single int option with bool multiple args": {[]string{cmd, "--aaa", "--bbb", "1"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := 0
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &v1}, {Name: "bbb", Required: true, Value: &v2}}
		}},
		"multiple int options": {[]string{cmd, "--aaa", "1", "--bbb", "2"}, func() []*cli.ArgConfig {
			val := 0
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}, {Name: "bbb", Required: true, Value: &val}}
		}},
		"single int64 option": {[]string{cmd, "--aaa", "1"}, func() []*cli.ArgConfig {
			val := int64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}}
		}},
		"single int64 option with bool multiple args": {[]string{cmd, "--aaa", "--bbb", "1"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := int64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &v1}, {Name: "bbb", Required: true, Value: &v2}}
		}},
		"multiple int64 options": {[]string{cmd, "--aaa", "1", "--bbb", "2"}, func() []*cli.ArgConfig {
			val := int64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}, {Name: "bbb", Required: true, Value: &val}}
		}},
		"single string option": {[]string{cmd, "--aaa", "foo"}, func() []*cli.ArgConfig {
			val := "foobar"
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}}
		}},
		"single string option with bool multiple args": {[]string{cmd, "--aaa", "--bbb", "foo"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := "foobar"
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &v1}, {Name: "bbb", Required: true, Value: &v2}}
		}},
		"multiple string options": {[]string{cmd, "--aaa", "foo", "--bbb", "bar"}, func() []*cli.ArgConfig {
			val := "foobar"
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}, {Name: "bbb", Required: true, Value: &val}}
		}},
		"single uint option": {[]string{cmd, "--aaa", "1"}, func() []*cli.ArgConfig {
			val := uint(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}}
		}},
		"single uint option with bool multiple args": {[]string{cmd, "--aaa", "--bbb", "1"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := uint(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &v1}, {Name: "bbb", Required: true, Value: &v2}}
		}},
		"multiple uint options": {[]string{cmd, "--aaa", "1", "--bbb", "2"}, func() []*cli.ArgConfig {
			val := uint(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}, {Name: "bbb", Required: true, Value: &val}}
		}},
		"single uint64 option": {[]string{cmd, "--aaa", "1"}, func() []*cli.ArgConfig {
			val := uint64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}}
		}},
		"single uint64 option with bool multiple args": {[]string{cmd, "--aaa", "--bbb", "1"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := uint64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &v1}, {Name: "bbb", Required: true, Value: &v2}}
		}},
		"multiple uint64 options": {[]string{cmd, "--aaa", "1", "--bbb", "2"}, func() []*cli.ArgConfig {
			val := uint64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}, {Name: "bbb", Required: true, Value: &val}}
		}},
	}

	for name, test := range testCases {
		os.Args = test.args
		_, parseErr := cli.NewParser(cli.GoFlag).Parse(&cli.CommandConfig{Args: test.value()})

		if parseErr != nil {
			t.Fail()
			t.Log(n + ": " + name)
		}
	}
}

func shouldErrorWhenGnuFirstArgIsInvalidFormat(t *testing.T, n string) {
	os.Args = []string{"testcmd", "a"}
	parser := cli.NewParser(cli.GNU)
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
		t.Log(n + ": did not return error with invalid GNU argument")
	}
}

func shouldErrorWhenConfiguredGnuArgNameIsInvalid(t *testing.T, n string) {
	os.Args = []string{"testcmd", "--="}
	parser := cli.NewParser(cli.GNU)
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
		t.Log(n + ": did not error on incorrectly configured GNU arg name")
	}
}

func shouldErrorWhenGnuOptionalOptArgIsInvalidFormat(t *testing.T, n string) {
	os.Args = []string{"testcmd", "--alphabet", "abc"}
	parser := cli.NewParser(cli.GNU)
	a := "xyz"
	_, parseErr := parser.Parse(&cli.CommandConfig{
		Args: []*cli.ArgConfig{{
			Name:      "alphabet",
			ShortName: 'a',
			Value:     &a,
		}},
	})

	if parseErr == nil {
		t.Fail()
		t.Log(n + ": did not error on GNU optional option-argument not separated by '='")
	}
}

func shouldParseWhenGnuArgsProvidedCorrectly(t *testing.T, n string) {
	cmd := "testcmd"
	testCases := map[string]struct {
		args  []string
		value func() []*cli.ArgConfig
	}{
		"single bool option": {[]string{cmd, "--aaa"}, func() []*cli.ArgConfig {
			val := false
			return []*cli.ArgConfig{{Name: "aaa", Value: &val}}
		}},
		"multiple bool options separate arg": {[]string{cmd, "--aaa", "--bbb"}, func() []*cli.ArgConfig {
			val := false
			return []*cli.ArgConfig{{Name: "aaa", Value: &val}, {Name: "bbb", Required: true, Value: &val}}
		}},
		"single float64 option": {[]string{cmd, "--aaa", "1.0"}, func() []*cli.ArgConfig {
			val := float64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}}
		}},
		"single float64 option with bool multiple args": {[]string{cmd, "--aaa", "--bbb", "1.0"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := float64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &v1}, {Name: "bbb", Required: true, Value: &v2}}
		}},
		"multiple float64 options": {[]string{cmd, "--aaa", "1.0", "--bbb", "2.0"}, func() []*cli.ArgConfig {
			val := float64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}, {Name: "bbb", Required: true, Value: &val}}
		}},
		"single int option": {[]string{cmd, "--aaa", "1"}, func() []*cli.ArgConfig {
			val := 0
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}}
		}},
		"single int option with bool multiple args": {[]string{cmd, "--aaa", "--bbb", "1"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := 0
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &v1}, {Name: "bbb", Required: true, Value: &v2}}
		}},
		"multiple int options": {[]string{cmd, "--aaa", "1", "--bbb", "2"}, func() []*cli.ArgConfig {
			val := 0
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}, {Name: "bbb", Required: true, Value: &val}}
		}},
		"single int64 option": {[]string{cmd, "--aaa", "1"}, func() []*cli.ArgConfig {
			val := int64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}}
		}},
		"single int64 option with bool multiple args": {[]string{cmd, "--aaa", "--bbb", "1"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := int64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &v1}, {Name: "bbb", Required: true, Value: &v2}}
		}},
		"multiple int64 options": {[]string{cmd, "--aaa", "1", "--bbb", "2"}, func() []*cli.ArgConfig {
			val := int64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}, {Name: "bbb", Required: true, Value: &val}}
		}},
		"single string option": {[]string{cmd, "--aaa", "foo"}, func() []*cli.ArgConfig {
			val := "foobar"
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}}
		}},
		"single string option with bool multiple args": {[]string{cmd, "--aaa", "--bbb", "foo"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := "foobar"
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &v1}, {Name: "bbb", Required: true, Value: &v2}}
		}},
		"multiple string options": {[]string{cmd, "--aaa", "foo", "--bbb", "bar"}, func() []*cli.ArgConfig {
			val := "foobar"
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}, {Name: "bbb", Required: true, Value: &val}}
		}},
		"single uint option": {[]string{cmd, "--aaa", "1"}, func() []*cli.ArgConfig {
			val := uint(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}}
		}},
		"single uint option with bool multiple args": {[]string{cmd, "--aaa", "--bbb", "1"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := uint(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &v1}, {Name: "bbb", Required: true, Value: &v2}}
		}},
		"multiple uint options": {[]string{cmd, "--aaa", "1", "--bbb", "2"}, func() []*cli.ArgConfig {
			val := uint(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}, {Name: "bbb", Required: true, Value: &val}}
		}},
		"single uint64 option": {[]string{cmd, "--aaa", "1"}, func() []*cli.ArgConfig {
			val := uint64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}}
		}},
		"single uint64 option with bool multiple args": {[]string{cmd, "--aaa", "--bbb", "1"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := uint64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &v1}, {Name: "bbb", Required: true, Value: &v2}}
		}},
		"multiple uint64 options": {[]string{cmd, "--aaa", "1", "--bbb", "2"}, func() []*cli.ArgConfig {
			val := uint64(0)
			return []*cli.ArgConfig{{Name: "aaa", Required: true, Value: &val}, {Name: "bbb", Required: true, Value: &val}}
		}},
	}

	for name, test := range testCases {
		os.Args = test.args
		_, parseErr := cli.NewParser(cli.GNU).Parse(&cli.CommandConfig{Args: test.value()})

		if parseErr != nil {
			t.Fail()
			t.Log(n + ": " + name)
		}
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

func shouldParseWhenPosixArgsProvidedCorrectly(t *testing.T, n string) {
	cmd := "testcmd"
	testCases := map[string]struct {
		args  []string
		value func() []*cli.ArgConfig
	}{
		"single bool option": {[]string{cmd, "-a"}, func() []*cli.ArgConfig {
			val := false
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &val}}
		}},
		"multiple bool options single arg": {[]string{cmd, "-ab"}, func() []*cli.ArgConfig {
			val := false
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &val}, {Name: "b", ShortName: 'b', Value: &val}}
		}},
		"multiple bool options separate arg": {[]string{cmd, "-a", "-b"}, func() []*cli.ArgConfig {
			val := false
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &val}, {Name: "b", ShortName: 'b', Value: &val}}
		}},
		"single float64 option": {[]string{cmd, "-a", "1.0"}, func() []*cli.ArgConfig {
			val := float64(0)
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &val}}
		}},
		"single float64 option with bool single arg": {[]string{cmd, "-ab", "1.0"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := float64(0)
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &v1}, {Name: "b", ShortName: 'b', Value: &v2}}
		}},
		"single float64 option with bool multiple args": {[]string{cmd, "-a", "-b", "1.0"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := float64(0)
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &v1}, {Name: "b", ShortName: 'b', Value: &v2}}
		}},
		"multiple float64 options": {[]string{cmd, "-a", "1.0", "-b", "2.0"}, func() []*cli.ArgConfig {
			val := float64(0)
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &val}, {Name: "b", ShortName: 'b', Value: &val}}
		}},
		"single int option": {[]string{cmd, "-a", "1"}, func() []*cli.ArgConfig {
			val := 0
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &val}}
		}},
		"single int option with bool single arg": {[]string{cmd, "-ab", "1"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := 0
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &v1}, {Name: "b", ShortName: 'b', Value: &v2}}
		}},
		"single int option with bool multiple args": {[]string{cmd, "-a", "-b", "1"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := 0
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &v1}, {Name: "b", ShortName: 'b', Value: &v2}}
		}},
		"multiple int options": {[]string{cmd, "-a", "1", "-b", "2"}, func() []*cli.ArgConfig {
			val := 0
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &val}, {Name: "b", ShortName: 'b', Value: &val}}
		}},
		"single int64 option": {[]string{cmd, "-a", "1"}, func() []*cli.ArgConfig {
			val := int64(0)
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &val}}
		}},
		"single int64 option with bool single arg": {[]string{cmd, "-ab", "1"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := int64(0)
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &v1}, {Name: "b", ShortName: 'b', Value: &v2}}
		}},
		"single int64 option with bool multiple args": {[]string{cmd, "-a", "-b", "1"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := int64(0)
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &v1}, {Name: "b", ShortName: 'b', Value: &v2}}
		}},
		"multiple int64 options": {[]string{cmd, "-a", "1", "-b", "2"}, func() []*cli.ArgConfig {
			val := int64(0)
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &val}, {Name: "b", ShortName: 'b', Value: &val}}
		}},
		"single string option": {[]string{cmd, "-a", "foo"}, func() []*cli.ArgConfig {
			val := "foobar"
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &val}}
		}},
		"single string option with bool single arg": {[]string{cmd, "-ab", "foo"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := "foobar"
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &v1}, {Name: "b", ShortName: 'b', Value: &v2}}
		}},
		"single string option with bool multiple args": {[]string{cmd, "-a", "-b", "foo"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := "foobar"
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &v1}, {Name: "b", ShortName: 'b', Value: &v2}}
		}},
		"multiple string options": {[]string{cmd, "-a", "foo", "-b", "bar"}, func() []*cli.ArgConfig {
			val := "foobar"
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &val}, {Name: "b", ShortName: 'b', Value: &val}}
		}},
		"single uint option": {[]string{cmd, "-a", "1"}, func() []*cli.ArgConfig {
			val := uint(0)
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &val}}
		}},
		"single uint option with bool single arg": {[]string{cmd, "-ab", "1"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := uint(0)
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &v1}, {Name: "b", ShortName: 'b', Value: &v2}}
		}},
		"single uint option with bool multiple args": {[]string{cmd, "-a", "-b", "1"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := uint(0)
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &v1}, {Name: "b", ShortName: 'b', Value: &v2}}
		}},
		"multiple uint options": {[]string{cmd, "-a", "1", "-b", "2"}, func() []*cli.ArgConfig {
			val := uint(0)
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &val}, {Name: "b", ShortName: 'b', Value: &val}}
		}},
		"single uint64 option": {[]string{cmd, "-a", "1"}, func() []*cli.ArgConfig {
			val := uint64(0)
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &val}}
		}},
		"single uint64 option with bool single arg": {[]string{cmd, "-ab", "1"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := uint64(0)
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &v1}, {Name: "b", ShortName: 'b', Value: &v2}}
		}},
		"single uint64 option with bool multiple args": {[]string{cmd, "-a", "-b", "1"}, func() []*cli.ArgConfig {
			v1 := false
			v2 := uint64(0)
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &v1}, {Name: "b", ShortName: 'b', Value: &v2}}
		}},
		"multiple uint64 options": {[]string{cmd, "-a", "1", "-b", "2"}, func() []*cli.ArgConfig {
			val := uint64(0)
			return []*cli.ArgConfig{{Name: "a", ShortName: 'a', Value: &val}, {Name: "b", ShortName: 'b', Value: &val}}
		}},
	}

	for name, test := range testCases {
		os.Args = test.args
		_, parseErr := cli.NewParser(cli.POSIX).Parse(&cli.CommandConfig{Args: test.value()})

		if parseErr != nil {
			t.Fail()
			t.Log(n + ": " + name)
		}
	}
}

func shouldErrorWhenBoolOptHasOptArg(t *testing.T, n string) {
	testCases := map[string]struct {
		syntax cli.ArgSyntax
		arg    string
	}{
		"GoFlag": {cli.GoFlag, "-a=value"},
		"GNU":    {cli.GNU, "-a"},
		"POSIX":  {cli.POSIX, "-a"},
	}

	for syntaxName, test := range testCases {
		os.Args = []string{"testcmd", test.arg, "value"}
		parser := cli.NewParser(test.syntax)
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
			t.Log(n + ": did not error on " + syntaxName + " bool option-argument")
		}
	}
}

func shouldErrorWhenRequiredFloat64OptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GoFlag": cli.GoFlag,
		"GNU":    cli.GNU,
		"POSIX":  cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		parser := cli.NewParser(syntax)
		a := float64(1)
		_, parseErr := parser.Parse(&cli.CommandConfig{
			Args: []*cli.ArgConfig{{
				Name:      "a",
				ShortName: 'a',
				Value:     &a,
			}},
		})

		if parseErr == nil {
			t.Fail()
			t.Log(n + ": did not error on missing " + syntaxName + " float64 option-argument")
		}
	}
}

func shouldErrorWhenRequiredFloat64ListOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GoFlag": cli.GoFlag,
		"GNU":    cli.GNU,
		"POSIX":  cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		parser := cli.NewParser(syntax)
		a := []float64{1}
		_, parseErr := parser.Parse(&cli.CommandConfig{
			Args: []*cli.ArgConfig{{
				Name:      "a",
				ShortName: 'a',
				Value:     &a,
			}},
		})

		if parseErr == nil {
			t.Fail()
			t.Log(n + ": did not error on missing " + syntaxName + " float64 option-argument")
		}
	}
}

func shouldErrorWhenRequiredIntOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GoFlag": cli.GoFlag,
		"GNU":    cli.GNU,
		"POSIX":  cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		parser := cli.NewParser(syntax)
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
			t.Log(n + ": did not error on missing " + syntaxName + " int option-argument")
		}
	}
}

func shouldErrorWhenRequiredIntListOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GoFlag": cli.GoFlag,
		"GNU":    cli.GNU,
		"POSIX":  cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		parser := cli.NewParser(syntax)
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
			t.Log(n + ": did not error on missing " + syntaxName + " int list option-argument")
		}
	}
}

func shouldErrorWhenRequiredInt64OptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GoFlag": cli.GoFlag,
		"GNU":    cli.GNU,
		"POSIX":  cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		parser := cli.NewParser(syntax)
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
			t.Log(n + ": did not error on missing " + syntaxName + " int64 option-argument")
		}
	}
}

func shouldErrorWhenRequiredInt64ListOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GoFlag": cli.GoFlag,
		"GNU":    cli.GNU,
		"POSIX":  cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		parser := cli.NewParser(syntax)
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
			t.Log(n + ": did not error on missing " + syntaxName + " int64 list option-argument")
		}
	}
}

func shouldErrorWhenRequiredStringOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GoFlag": cli.GoFlag,
		"GNU":    cli.GNU,
		"POSIX":  cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		parser := cli.NewParser(syntax)
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
			t.Log(n + ": did not error on missing " + syntaxName + " string option-argument")
		}
	}
}

func shouldErrorWhenRequiredStringListOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GoFlag": cli.GoFlag,
		"GNU":    cli.GNU,
		"POSIX":  cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		parser := cli.NewParser(syntax)
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
			t.Log(n + ": did not error on missing " + syntaxName + " string list option-argument")
		}
	}
}

func shouldErrorWhenRequiredUintOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GoFlag": cli.GoFlag,
		"GNU":    cli.GNU,
		"POSIX":  cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		parser := cli.NewParser(syntax)
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
			t.Log(n + ": did not error on missing " + syntaxName + " uint option-argument")
		}
	}
}

func shouldErrorWhenRequiredUintListOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GoFlag": cli.GoFlag,
		"GNU":    cli.GNU,
		"POSIX":  cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		parser := cli.NewParser(syntax)
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
			t.Log(n + ": did not error on missing " + syntaxName + " uint list option-argument")
		}
	}
}

func shouldErrorWhenRequiredUint64OptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GoFlag": cli.GoFlag,
		"GNU":    cli.GNU,
		"POSIX":  cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		parser := cli.NewParser(syntax)
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
			t.Log(n + ": did not error on missing " + syntaxName + " uint64 option-argument")
		}
	}
}

func shouldErrorWhenRequiredUint64ListOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GoFlag": cli.GoFlag,
		"GNU":    cli.GNU,
		"POSIX":  cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		parser := cli.NewParser(syntax)
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
			t.Log(n + ": did not error on missing " + syntaxName + " uint64 list option-argument")
		}
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
