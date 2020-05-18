package cli_test

import (
	"context"
	"github.com/sebuckler/teel/pkg/cli"
	"os"
	"strconv"
	"strings"
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
		"should parse subcommands":                                   shouldParseSubcommands,
		"should set help mode true if help arg exists":               shouldSetHelpModeTrueIfHelpArgExists,
		"should error when arg passed with no args configured":       shouldErrorWhenArgPassedWithNoArgsConfigured,
		"should error when repeated arg is not repeatable":           shouldErrorWhenRepeatedArgIsNotRepeatable,
		"should error when GNU first arg is invalid format":          shouldErrorWhenGnuFirstArgIsInvalidFormat,
		"should error when configured GNU arg name is invalid":       shouldErrorWhenConfiguredGnuArgNameIsInvalid,
		"should error when GNU optional opt-arg is invalid format":   shouldErrorWhenGnuOptionalOptArgIsInvalidFormat,
		"should error when POSIX first arg is invalid format":        shouldErrorWhenPosixFirstArgIsInvalidFormat,
		"should error when configured POSIX arg name is invalid":     shouldErrorWhenConfiguredPosixArgNameIsInvalid,
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
		"should parse when operands provided correctly":              shouldParseWhenOperandsProvidedCorrectly,
		"should parse when args provided correctly":                  shouldParseWhenPosixArgsProvidedCorrectly,
	}
}

func shouldErrorWhenUnsupportedParseSyntaxUsed(t *testing.T, n string) {
	os.Args = []string{"testcmd"}
	cmd := cli.NewCommand("testcmd", context.Background())
	parser := cli.NewParser(99, cmd)
	_, err := parser.Parse()

	if err == nil {
		t.Fail()
		t.Log(n + ": did not error on unsupported parse syntax")
	}
}

func shouldParseSubcommands(t *testing.T, n string) {
	os.Args = []string{"testcmd", "foo", "bar"}
	cmd := cli.NewCommand("testcmd", context.Background())
	sub1 := cli.NewCommand("foo", context.Background())
	sub2 := cli.NewCommand("bar", context.Background())
	sub1.AddSubcommand(sub2)
	cmd.AddSubcommand(sub1)
	parser := cli.NewParser(cli.POSIX, cmd)
	parsedCmd, err := parser.Parse()

	if err != nil || len(parsedCmd.Subcommands) == 0 || len(parsedCmd.Subcommands[0].Subcommands) == 0 {
		t.Fail()
		t.Log(n + ": did not parse subcommands properly")
	}
}

func shouldSetHelpModeTrueIfHelpArgExists(t *testing.T, n string) {
	testCases := map[string]map[cli.ArgSyntax][][]string{
		"GNU help added":   {cli.GNU: {{"testcmd", "--help"}, {"testcmd", "-h"}}},
		"POSIX help added": {cli.POSIX: {{"testcmd", "-h"}}},
	}

	for name, test := range testCases {
		for syntax, argSet := range test {
			for _, args := range argSet {
				os.Args = args
				cmd := cli.NewCommand("testcmd", context.Background())
				parser := cli.NewParser(syntax, cmd)
				parsedCommand, err := parser.Parse()

				if err != nil || (parsedCommand != nil && !parsedCommand.HelpMode) {
					t.Fail()
					t.Log(n + ": failed to parse help args: " + name)
				}
			}
		}
	}
}

func shouldErrorWhenArgPassedWithNoArgsConfigured(t *testing.T, n string) {
	os.Args = []string{"testcmd", "a"}
	cmd := cli.NewCommand("testcmd", context.Background())
	parser := cli.NewParser(cli.POSIX, cmd)
	_, err := parser.Parse()

	if err == nil {
		t.Fail()
		t.Log(n + ": did not error when args passed in with no args configured")
	}
}

func shouldErrorWhenRepeatedArgIsNotRepeatable(t *testing.T, n string) {
	testCases := map[cli.ArgSyntax][]string{
		cli.GNU:   {"testcmd", "--aaa", "--aaa"},
		cli.POSIX: {"testcmd", "-a", "-a"},
	}

	for syntax, args := range testCases {
		os.Args = args
		cmd := cli.NewCommand("testcmd", context.Background())

		for _, arg := range args {
			val := false
			cmd.AddBoolArg(&val, &cli.ArgDefinition{
				Name:      strings.TrimLeft(arg, "-"),
				ShortName: rune(strings.TrimLeft(arg, "-")[0]),
			})
		}

		parser := cli.NewParser(syntax, cmd)
		_, err := parser.Parse()

		if err == nil {
			t.Fail()
			t.Log(n + ": did not error on non-repeatable arg being repeated")
		}
	}
}

func shouldErrorWhenGnuFirstArgIsInvalidFormat(t *testing.T, n string) {
	os.Args = []string{"testcmd", "a"}
	cmd := cli.NewCommand("testcmd", context.Background())
	val := false
	cmd.AddBoolArg(&val, &cli.ArgDefinition{Name: "a", ShortName: 'a'})
	parser := cli.NewParser(cli.GNU, cmd)
	_, err := parser.Parse()

	if err == nil {
		t.Fail()
		t.Log(n + ": did not return error with invalid GNU argument")
	}
}

func shouldErrorWhenConfiguredGnuArgNameIsInvalid(t *testing.T, n string) {
	testCases := map[string][]string{
		"invalid characters":          {"testcmd", "--="},
		"invalid multi-part name '='": {"testcmd", "--foo-="},
		"invalid multi-part name '`'": {"testcmd", "--foo-`"},
	}

	for name, args := range testCases {
		os.Args = args
		cmd := cli.NewCommand("testcmd", context.Background())

		for _, arg := range args {
			val := false
			cmd.AddBoolArg(&val, &cli.ArgDefinition{
				Name:      strings.TrimLeft(arg, "-"),
				ShortName: rune(strings.TrimLeft(arg, "-")[0]),
			})
		}
		parser := cli.NewParser(cli.GNU, cmd)
		_, err := parser.Parse()

		if err == nil {
			t.Fail()
			t.Log(n + ": did not error on incorrectly configured GNU arg name: " + name)
		}
	}
}

func shouldErrorWhenGnuOptionalOptArgIsInvalidFormat(t *testing.T, n string) {
	os.Args = []string{"testcmd", "--foo", "bar"}
	cmd := cli.NewCommand("testcmd", context.Background())
	val := "baz"
	cmd.AddStringArg(&val, &cli.ArgDefinition{Name: "foo", ShortName: 'f'})
	parser := cli.NewParser(cli.GNU, cmd)
	_, err := parser.Parse()

	if err == nil {
		t.Fail()
		t.Log(n + ": did not error on GNU optional option-argument not separated by '='")
	}
}

func shouldErrorWhenPosixFirstArgIsInvalidFormat(t *testing.T, n string) {
	os.Args = []string{"testcmd", "a"}
	cmd := cli.NewCommand("testcmd", context.Background())
	val := false
	cmd.AddBoolArg(&val, &cli.ArgDefinition{Name: "a", ShortName: 'a'})
	parser := cli.NewParser(cli.POSIX, cmd)
	_, err := parser.Parse()

	if err == nil {
		t.Fail()
		t.Log(n + ": did not return error with invalid POSIX argument")
	}
}

func shouldErrorWhenConfiguredPosixArgNameIsInvalid(t *testing.T, n string) {
	os.Args = []string{"testcmd", "-="}
	cmd := cli.NewCommand("testcmd", context.Background())
	val := false
	cmd.AddBoolArg(&val, &cli.ArgDefinition{Name: "=", ShortName: '='})
	parser := cli.NewParser(cli.POSIX, cmd)
	_, err := parser.Parse()

	if err == nil {
		t.Fail()
		t.Log(n + ": did not error on incorrectly configured POSIX arg name")
	}
}

func shouldErrorWhenBoolOptHasOptArg(t *testing.T, n string) {
	testCases := map[string]struct {
		syntax cli.ArgSyntax
		arg    string
	}{
		"GNU":   {cli.GNU, "-a"},
		"POSIX": {cli.POSIX, "-a"},
	}

	for syntaxName, test := range testCases {
		os.Args = []string{"testcmd", test.arg, "value"}
		cmd := cli.NewCommand("testcmd", context.Background())
		val := false
		cmd.AddBoolArg(&val, &cli.ArgDefinition{Name: "a", ShortName: 'a'})
		parser := cli.NewParser(test.syntax, cmd)
		_, err := parser.Parse()

		if err == nil {
			t.Fail()
			t.Log(n + ": did not error on " + syntaxName + " bool option-argument")
		}
	}
}

func shouldErrorWhenRequiredFloat64OptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GNU":   cli.GNU,
		"POSIX": cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		cmd := cli.NewCommand("testcmd", context.Background())
		val := float64(1)
		cmd.AddFloat64Arg(&val, &cli.ArgDefinition{Name: "a", ShortName: 'a'})
		parser := cli.NewParser(syntax, cmd)
		_, err := parser.Parse()

		if err == nil {
			t.Fail()
			t.Log(n + ": did not error on missing " + syntaxName + " float64 option-argument")
		}
	}
}

func shouldErrorWhenRequiredFloat64ListOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GNU":   cli.GNU,
		"POSIX": cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		cmd := cli.NewCommand("testcmd", context.Background())
		val := []float64{1}
		cmd.AddFloat64ListArg(&val, &cli.ArgDefinition{Name: "a", ShortName: 'a'})
		parser := cli.NewParser(syntax, cmd)
		_, err := parser.Parse()

		if err == nil {
			t.Fail()
			t.Log(n + ": did not error on missing " + syntaxName + " float64 option-argument")
		}
	}
}

func shouldErrorWhenRequiredIntOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GNU":   cli.GNU,
		"POSIX": cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		cmd := cli.NewCommand("testcmd", context.Background())
		val := 1
		cmd.AddIntArg(&val, &cli.ArgDefinition{Name: "a", ShortName: 'a'})
		parser := cli.NewParser(syntax, cmd)
		_, err := parser.Parse()

		if err == nil {
			t.Fail()
			t.Log(n + ": did not error on missing " + syntaxName + " int option-argument")
		}
	}
}

func shouldErrorWhenRequiredIntListOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GNU":   cli.GNU,
		"POSIX": cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		cmd := cli.NewCommand("testcmd", context.Background())
		val := []int{1}
		cmd.AddIntListArg(&val, &cli.ArgDefinition{Name: "a", ShortName: 'a'})
		parser := cli.NewParser(syntax, cmd)
		_, err := parser.Parse()

		if err == nil {
			t.Fail()
			t.Log(n + ": did not error on missing " + syntaxName + " int list option-argument")
		}
	}
}

func shouldErrorWhenRequiredInt64OptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GNU":   cli.GNU,
		"POSIX": cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		cmd := cli.NewCommand("testcmd", context.Background())
		val := int64(1)
		cmd.AddInt64Arg(&val, &cli.ArgDefinition{Name: "a", ShortName: 'a'})
		parser := cli.NewParser(syntax, cmd)
		_, err := parser.Parse()

		if err == nil {
			t.Fail()
			t.Log(n + ": did not error on missing " + syntaxName + " int64 option-argument")
		}
	}
}

func shouldErrorWhenRequiredInt64ListOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GNU":   cli.GNU,
		"POSIX": cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		cmd := cli.NewCommand("testcmd", context.Background())
		val := []int64{1}
		cmd.AddInt64ListArg(&val, &cli.ArgDefinition{Name: "a", ShortName: 'a'})
		parser := cli.NewParser(syntax, cmd)
		_, err := parser.Parse()

		if err == nil {
			t.Fail()
			t.Log(n + ": did not error on missing " + syntaxName + " int64 list option-argument")
		}
	}
}

func shouldErrorWhenRequiredStringOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GNU":   cli.GNU,
		"POSIX": cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		cmd := cli.NewCommand("testcmd", context.Background())
		val := "value"
		cmd.AddStringArg(&val, &cli.ArgDefinition{Name: "a", ShortName: 'a'})
		parser := cli.NewParser(syntax, cmd)
		_, err := parser.Parse()

		if err == nil {
			t.Fail()
			t.Log(n + ": did not error on missing " + syntaxName + " string option-argument")
		}
	}
}

func shouldErrorWhenRequiredStringListOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GNU":   cli.GNU,
		"POSIX": cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		cmd := cli.NewCommand("testcmd", context.Background())
		val := []string{"value"}
		cmd.AddStringListArg(&val, &cli.ArgDefinition{Name: "a", ShortName: 'a'})
		parser := cli.NewParser(syntax, cmd)
		_, err := parser.Parse()

		if err == nil {
			t.Fail()
			t.Log(n + ": did not error on missing " + syntaxName + " string list option-argument")
		}
	}
}

func shouldErrorWhenRequiredUintOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GNU":   cli.GNU,
		"POSIX": cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		cmd := cli.NewCommand("testcmd", context.Background())
		val := uint(1)
		cmd.AddUintArg(&val, &cli.ArgDefinition{Name: "a", ShortName: 'a'})
		parser := cli.NewParser(syntax, cmd)
		_, err := parser.Parse()

		if err == nil {
			t.Fail()
			t.Log(n + ": did not error on missing " + syntaxName + " uint option-argument")
		}
	}
}

func shouldErrorWhenRequiredUintListOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GNU":   cli.GNU,
		"POSIX": cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		cmd := cli.NewCommand("testcmd", context.Background())
		val := []uint{1}
		cmd.AddUintListArg(&val, &cli.ArgDefinition{Name: "a", ShortName: 'a'})
		parser := cli.NewParser(syntax, cmd)
		_, err := parser.Parse()

		if err == nil {
			t.Fail()
			t.Log(n + ": did not error on missing " + syntaxName + " uint list option-argument")
		}
	}
}

func shouldErrorWhenRequiredUint64OptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GNU":   cli.GNU,
		"POSIX": cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		cmd := cli.NewCommand("testcmd", context.Background())
		val := uint64(1)
		cmd.AddUint64Arg(&val, &cli.ArgDefinition{Name: "a", ShortName: 'a'})
		parser := cli.NewParser(syntax, cmd)
		_, err := parser.Parse()

		if err == nil {
			t.Fail()
			t.Log(n + ": did not error on missing " + syntaxName + " uint64 option-argument")
		}
	}
}

func shouldErrorWhenRequiredUint64ListOptHasNoOptArg(t *testing.T, n string) {
	testCases := map[string]cli.ArgSyntax{
		"GNU":   cli.GNU,
		"POSIX": cli.POSIX,
	}

	for syntaxName, syntax := range testCases {
		os.Args = []string{"testcmd", "-a"}
		cmd := cli.NewCommand("testcmd", context.Background())
		val := []uint64{1}
		cmd.AddUint64ListArg(&val, &cli.ArgDefinition{Name: "a", ShortName: 'a'})
		parser := cli.NewParser(syntax, cmd)
		_, err := parser.Parse()

		if err == nil {
			t.Fail()
			t.Log(n + ": did not error on missing " + syntaxName + " uint64 list option-argument")
		}
	}
}

func shouldParseWhenOperandsProvidedCorrectly(t *testing.T, n string) {
	testCases := map[cli.ArgSyntax][]string{
		cli.GNU:   {"testcmd", "--aaa", "--", "+foo"},
		cli.POSIX: {"testcmd", "-a", "--", "+foo"},
	}

	for syntax, args := range testCases {
		os.Args = args
		cmd := cli.NewCommand("testcmd", context.Background())
		option := strings.TrimLeft(args[1], "-")
		val := false
		cmd.AddBoolArg(&val, &cli.ArgDefinition{Name: option, ShortName: rune(option[0])})
		parser := cli.NewParser(syntax, cmd)
		_, err := parser.Parse()

		if err != nil {
			t.Fail()
			t.Log(n + ": failed to parse operands")
		}
	}
}

func shouldParseWhenPosixArgsProvidedCorrectly(t *testing.T, n string) {
	cmd := "testcmd"
	testCases := map[string]func() func() bool{
		"bool options": func() func() bool {
			tests := map[cli.ArgSyntax]map[*[]string][]string{
				cli.GNU: {
					&[]string{"aaa"}:        {cmd, "--aaa"},
					&[]string{"aaa", "bbb"}: {cmd, "--aaa", "--bbb"},
					&[]string{"a"}:          {cmd, "-a"},
					&[]string{"a", "b"}:     {cmd, "-ab"},
					&[]string{"a", "b"}:     {cmd, "-a", "-b"},
				},
				cli.POSIX: {
					&[]string{"a"}:                {cmd, "-a"},
					&[]string{"a", "b"}:           {cmd, "-ab"},
					&[]string{"a", "b"}:           {cmd, "-a", "-b"},
					&[]string{"a", "b", "c", "d"}: {cmd, "-ab", "-cd"},
				},
			}
			val := false
			var errs []error

			for syntax, argSets := range tests {
				for argNames, args := range argSets {
					os.Args = args
					cmd := cli.NewCommand("testcmd", context.Background())

					for _, name := range *argNames {
						cmd.AddBoolArg(&val, &cli.ArgDefinition{Name: name, ShortName: rune(name[0])})
					}

					if _, err := cli.NewParser(syntax, cmd).Parse(); err != nil {
						errs = append(errs, err)
					}
				}
			}

			return func() bool { return len(errs) == 0 && val }
		},
		"float64 options": func() func() bool {
			tests := map[cli.ArgSyntax]map[*[]string][]string{
				cli.GNU: {
					&[]string{"aaa"}:        {cmd, "--aaa=1.0"},
					&[]string{"aaa", "bbb"}: {cmd, "--aaa=1.0", "--bbb=2.0"},
					&[]string{"a"}:          {cmd, "-a", "1.0"},
					&[]string{"a", "b"}:     {cmd, "-a", "1.0", "-b", "2.0"},
				},
				cli.POSIX: {
					&[]string{"a"}:      {cmd, "-a", "1.0"},
					&[]string{"a", "b"}: {cmd, "-a", "1.0", "-b", "2.0"},
				},
			}
			var errs []error
			vals := map[float64]*float64{}

			for syntax, argSets := range tests {
				for argNames, args := range argSets {
					os.Args = args
					cmd := cli.NewCommand("testcmd", context.Background())

					for i, name := range *argNames {
						val := float64(i)
						bindVal := &val
						cmd.AddFloat64Arg(bindVal, &cli.ArgDefinition{Name: name, ShortName: rune(name[0])})
						vals[val+1] = bindVal
					}

					if _, err := cli.NewParser(syntax, cmd).Parse(); err != nil {
						errs = append(errs, err)
					}
				}
			}

			return func() bool {
				for val, bindVal := range vals {
					if val != *bindVal {
						return false
					}
				}

				return len(errs) == 0
			}
		},
		"float64 list options": func() func() bool {
			type result struct {
				val     []float64
				bindVal *[]float64
			}
			tests := map[cli.ArgSyntax]map[*[]string][]string{
				cli.GNU: {
					&[]string{"aaa"}:        {cmd, "--aaa=1.0,2.0"},
					&[]string{"aaa", "bbb"}: {cmd, "--aaa=1.0,2.0", "--bbb=3.0,4.0"},
					&[]string{"a"}:          {cmd, "-a", "1.0,2.0"},
					&[]string{"a", "b"}:     {cmd, "-a", "1.0,2.0", "-b", "3.0,4.0"},
				},
				cli.POSIX: {
					&[]string{"a"}:      {cmd, "-a", "1.0,2.0"},
					&[]string{"a", "b"}: {cmd, "-a", "1.0,2.0", "-b", "3.0,4.0"},
				},
			}
			var results []result
			var errs []error

			for syntax, argSets := range tests {
				for argNames, args := range argSets {
					os.Args = args
					cmd := cli.NewCommand("testcmd", context.Background())

					for i, name := range *argNames {
						val := []float64{float64(i)}
						bindVal := &val
						cmd.AddFloat64ListArg(bindVal, &cli.ArgDefinition{Name: name, ShortName: rune(name[0])})
						results = append(results, result{[]float64{float64((i * 2) + 1), float64((i * 2) + 2)}, bindVal})
					}

					if _, err := cli.NewParser(syntax, cmd).Parse(); err != nil {
						errs = append(errs, err)
					}
				}
			}

			return func() bool {
				for _, result := range results {
					for i, val := range result.val {
						if val != (*(result.bindVal))[i] {
							return false
						}
					}
				}

				return len(errs) == 0
			}
		},
		"int options": func() func() bool {
			tests := map[cli.ArgSyntax]map[*[]string][]string{
				cli.GNU: {
					&[]string{"aaa"}:        {cmd, "--aaa=1"},
					&[]string{"aaa", "bbb"}: {cmd, "--aaa=1", "--bbb=2"},
					&[]string{"a"}:          {cmd, "-a", "1"},
					&[]string{"a", "b"}:     {cmd, "-a", "1", "-b", "2"},
				},
				cli.POSIX: {
					&[]string{"a"}:      {cmd, "-a", "1"},
					&[]string{"a", "b"}: {cmd, "-a", "1", "-b", "2"},
				},
			}
			var errs []error
			vals := map[int]*int{}

			for syntax, argSets := range tests {
				for argNames, args := range argSets {
					os.Args = args
					cmd := cli.NewCommand("testcmd", context.Background())

					for i, name := range *argNames {
						val := i
						bindVal := &val
						cmd.AddIntArg(bindVal, &cli.ArgDefinition{Name: name, ShortName: rune(name[0])})
						vals[val+1] = bindVal
					}

					if _, err := cli.NewParser(syntax, cmd).Parse(); err != nil {
						errs = append(errs, err)
					}
				}
			}

			return func() bool {
				for val, bindVal := range vals {
					if val != *bindVal {
						return false
					}
				}

				return len(errs) == 0
			}
		},
		"int list options": func() func() bool {
			type result struct {
				val     []int
				bindVal *[]int
			}
			tests := map[cli.ArgSyntax]map[*[]string][]string{
				cli.GNU: {
					&[]string{"aaa"}:        {cmd, "--aaa=1,2"},
					&[]string{"aaa", "bbb"}: {cmd, "--aaa=1,2", "--bbb=3,4"},
					&[]string{"a"}:          {cmd, "-a", "1,2"},
					&[]string{"a", "b"}:     {cmd, "-a", "1,2", "-b", "3,4"},
				},
				cli.POSIX: {
					&[]string{"a"}:      {cmd, "-a", "1,2"},
					&[]string{"a", "b"}: {cmd, "-a", "1,2", "-b", "3,4"},
				},
			}
			var results []result
			var errs []error

			for syntax, argSets := range tests {
				for argNames, args := range argSets {
					os.Args = args
					cmd := cli.NewCommand("testcmd", context.Background())

					for i, name := range *argNames {
						val := []int{i}
						bindVal := &val
						cmd.AddIntListArg(bindVal, &cli.ArgDefinition{Name: name, ShortName: rune(name[0])})
						results = append(results, result{[]int{(i * 2) + 1, (i * 2) + 2}, bindVal})
					}

					if _, err := cli.NewParser(syntax, cmd).Parse(); err != nil {
						errs = append(errs, err)
					}
				}
			}

			return func() bool {
				for _, result := range results {
					for i, val := range result.val {
						if val != (*(result.bindVal))[i] {
							return false
						}
					}
				}

				return len(errs) == 0
			}
		},
		"int64 options": func() func() bool {
			tests := map[cli.ArgSyntax]map[*[]string][]string{
				cli.GNU: {
					&[]string{"aaa"}:        {cmd, "--aaa=1"},
					&[]string{"aaa", "bbb"}: {cmd, "--aaa=1", "--bbb=2"},
					&[]string{"a"}:          {cmd, "-a", "1"},
					&[]string{"a", "b"}:     {cmd, "-a", "1", "-b", "2"},
				},
				cli.POSIX: {
					&[]string{"a"}:      {cmd, "-a", "1"},
					&[]string{"a", "b"}: {cmd, "-a", "1", "-b", "2"},
				},
			}
			var errs []error
			vals := map[int64]*int64{}

			for syntax, argSets := range tests {
				for argNames, args := range argSets {
					os.Args = args
					cmd := cli.NewCommand("testcmd", context.Background())

					for i, name := range *argNames {
						val := int64(i)
						bindVal := &val
						cmd.AddInt64Arg(bindVal, &cli.ArgDefinition{Name: name, ShortName: rune(name[0])})
						vals[val+1] = bindVal
					}

					if _, err := cli.NewParser(syntax, cmd).Parse(); err != nil {
						errs = append(errs, err)
					}
				}
			}

			return func() bool {
				for val, bindVal := range vals {
					if val != *bindVal {
						return false
					}
				}

				return len(errs) == 0
			}
		},
		"int64 list options": func() func() bool {
			type result struct {
				val     []int64
				bindVal *[]int64
			}
			tests := map[cli.ArgSyntax]map[*[]string][]string{
				cli.GNU: {
					&[]string{"aaa"}:        {cmd, "--aaa=1,2"},
					&[]string{"aaa", "bbb"}: {cmd, "--aaa=1,2", "--bbb=3,4"},
					&[]string{"a"}:          {cmd, "-a", "1,2"},
					&[]string{"a", "b"}:     {cmd, "-a", "1,2", "-b", "3,4"},
				},
				cli.POSIX: {
					&[]string{"a"}:      {cmd, "-a", "1,2"},
					&[]string{"a", "b"}: {cmd, "-a", "1,2", "-b", "3,4"},
				},
			}
			var results []result
			var errs []error

			for syntax, argSets := range tests {
				for argNames, args := range argSets {
					os.Args = args
					cmd := cli.NewCommand("testcmd", context.Background())

					for i, name := range *argNames {
						val := []int64{int64(i)}
						bindVal := &val
						cmd.AddInt64ListArg(bindVal, &cli.ArgDefinition{Name: name, ShortName: rune(name[0])})
						results = append(results, result{[]int64{int64((i * 2) + 1), int64((i * 2) + 2)}, bindVal})
					}

					if _, err := cli.NewParser(syntax, cmd).Parse(); err != nil {
						errs = append(errs, err)
					}
				}
			}

			return func() bool {
				for _, result := range results {
					for i, val := range result.val {
						if val != (*(result.bindVal))[i] {
							return false
						}
					}
				}

				return len(errs) == 0
			}
		},
		"string options": func() func() bool {
			tests := map[cli.ArgSyntax]map[*[]string][]string{
				cli.GNU: {
					&[]string{"aaa"}:        {cmd, "--aaa=1"},
					&[]string{"aaa", "bbb"}: {cmd, "--aaa=1", "--bbb=2"},
					&[]string{"a"}:          {cmd, "-a", "1"},
					&[]string{"a", "b"}:     {cmd, "-a", "1", "-b", "2"},
				},
				cli.POSIX: {
					&[]string{"a"}:      {cmd, "-a", "1"},
					&[]string{"a", "b"}: {cmd, "-a", "1", "-b", "2"},
				},
			}
			var errs []error
			vals := map[string]*string{}

			for syntax, argSets := range tests {
				for argNames, args := range argSets {
					os.Args = args
					cmd := cli.NewCommand("testcmd", context.Background())

					for i, name := range *argNames {
						val := strconv.Itoa(i + 1)
						bindVal := &val
						cmd.AddStringArg(bindVal, &cli.ArgDefinition{Name: name, ShortName: rune(name[0])})
						vals[*bindVal] = bindVal
					}

					if _, err := cli.NewParser(syntax, cmd).Parse(); err != nil {
						errs = append(errs, err)
					}
				}
			}

			return func() bool {
				for val, bindVal := range vals {
					if val != *bindVal {
						return false
					}
				}

				return len(errs) == 0
			}
		},
		"string list options": func() func() bool {
			type result struct {
				val     []string
				bindVal *[]string
			}
			tests := map[cli.ArgSyntax]map[*[]string][]string{
				cli.GNU: {
					&[]string{"aaa"}:        {cmd, "--aaa=1,2"},
					&[]string{"aaa", "bbb"}: {cmd, "--aaa=1,2", "--bbb=3,4"},
					&[]string{"a"}:          {cmd, "-a", "1,2"},
					&[]string{"a", "b"}:     {cmd, "-a", "1,2", "-b", "3,4"},
				},
				cli.POSIX: {
					&[]string{"a"}:      {cmd, "-a", "1,2"},
					&[]string{"a", "b"}: {cmd, "-a", "1,2", "-b", "3,4"},
				},
			}
			var results []result
			var errs []error

			for syntax, argSets := range tests {
				for argNames, args := range argSets {
					os.Args = args
					cmd := cli.NewCommand("testcmd", context.Background())

					for i, name := range *argNames {
						val := []string{strconv.Itoa(i + 1)}
						bindVal := &val
						cmd.AddStringListArg(bindVal, &cli.ArgDefinition{Name: name, ShortName: rune(name[0])})
						results = append(
							results,
							result{[]string{strconv.Itoa((i * 2) + 1), strconv.Itoa((i * 2) + 2)}, bindVal},
						)
					}

					if _, err := cli.NewParser(syntax, cmd).Parse(); err != nil {
						errs = append(errs, err)
					}
				}
			}

			return func() bool {
				for _, result := range results {
					for i, val := range result.val {
						if val != (*(result.bindVal))[i] {
							return false
						}
					}
				}

				return len(errs) == 0
			}
		},
		"uint options": func() func() bool {
			tests := map[cli.ArgSyntax]map[*[]string][]string{
				cli.GNU: {
					&[]string{"aaa"}:        {cmd, "--aaa=1"},
					&[]string{"aaa", "bbb"}: {cmd, "--aaa=1", "--bbb=2"},
					&[]string{"a"}:          {cmd, "-a", "1"},
					&[]string{"a", "b"}:     {cmd, "-a", "1", "-b", "2"},
				},
				cli.POSIX: {
					&[]string{"a"}:      {cmd, "-a", "1"},
					&[]string{"a", "b"}: {cmd, "-a", "1", "-b", "2"},
				},
			}
			var errs []error
			vals := map[uint]*uint{}

			for syntax, argSets := range tests {
				for argNames, args := range argSets {
					os.Args = args
					cmd := cli.NewCommand("testcmd", context.Background())

					for i, name := range *argNames {
						val := uint(i)
						bindVal := &val
						cmd.AddUintArg(bindVal, &cli.ArgDefinition{Name: name, ShortName: rune(name[0])})
						vals[val+1] = bindVal
					}

					if _, err := cli.NewParser(syntax, cmd).Parse(); err != nil {
						errs = append(errs, err)
					}
				}
			}

			return func() bool {
				for val, bindVal := range vals {
					if val != *bindVal {
						return false
					}
				}

				return len(errs) == 0
			}
		},
		"uint list options": func() func() bool {
			type result struct {
				val     []uint
				bindVal *[]uint
			}
			tests := map[cli.ArgSyntax]map[*[]string][]string{
				cli.GNU: {
					&[]string{"aaa"}:        {cmd, "--aaa=1,2"},
					&[]string{"aaa", "bbb"}: {cmd, "--aaa=1,2", "--bbb=3,4"},
					&[]string{"a"}:          {cmd, "-a", "1,2"},
					&[]string{"a", "b"}:     {cmd, "-a", "1,2", "-b", "3,4"},
				},
				cli.POSIX: {
					&[]string{"a"}:      {cmd, "-a", "1,2"},
					&[]string{"a", "b"}: {cmd, "-a", "1,2", "-b", "3,4"},
				},
			}
			var results []result
			var errs []error

			for syntax, argSets := range tests {
				for argNames, args := range argSets {
					os.Args = args
					cmd := cli.NewCommand("testcmd", context.Background())

					for i, name := range *argNames {
						val := []uint{uint(i)}
						bindVal := &val
						cmd.AddUintListArg(bindVal, &cli.ArgDefinition{Name: name, ShortName: rune(name[0])})
						results = append(results, result{[]uint{uint((i * 2) + 1), uint((i * 2) + 2)}, bindVal})
					}

					if _, err := cli.NewParser(syntax, cmd).Parse(); err != nil {
						errs = append(errs, err)
					}
				}
			}

			return func() bool {
				for _, result := range results {
					for i, val := range result.val {
						if val != (*(result.bindVal))[i] {
							return false
						}
					}
				}

				return len(errs) == 0
			}
		},
		"uint64 options": func() func() bool {
			tests := map[cli.ArgSyntax]map[*[]string][]string{
				cli.GNU: {
					&[]string{"aaa"}:        {cmd, "--aaa=1"},
					&[]string{"aaa", "bbb"}: {cmd, "--aaa=1", "--bbb=2"},
					&[]string{"a"}:          {cmd, "-a", "1"},
					&[]string{"a", "b"}:     {cmd, "-a", "1", "-b", "2"},
				},
				cli.POSIX: {
					&[]string{"a"}:      {cmd, "-a", "1"},
					&[]string{"a", "b"}: {cmd, "-a", "1", "-b", "2"},
				},
			}
			var errs []error
			vals := map[uint64]*uint64{}

			for syntax, argSets := range tests {
				for argNames, args := range argSets {
					os.Args = args
					cmd := cli.NewCommand("testcmd", context.Background())

					for i, name := range *argNames {
						val := uint64(i)
						bindVal := &val
						cmd.AddUint64Arg(bindVal, &cli.ArgDefinition{Name: name, ShortName: rune(name[0])})
						vals[val+1] = bindVal
					}

					if _, err := cli.NewParser(syntax, cmd).Parse(); err != nil {
						errs = append(errs, err)
					}
				}
			}

			return func() bool {
				for val, bindVal := range vals {
					if val != *bindVal {
						return false
					}
				}

				return len(errs) == 0
			}
		},
		"uint64 list options": func() func() bool {
			type result struct {
				val     []uint64
				bindVal *[]uint64
			}
			tests := map[cli.ArgSyntax]map[*[]string][]string{
				cli.GNU: {
					&[]string{"aaa"}:        {cmd, "--aaa=1,2"},
					&[]string{"aaa", "bbb"}: {cmd, "--aaa=1,2", "--bbb=3,4"},
					&[]string{"a"}:          {cmd, "-a", "1,2"},
					&[]string{"a", "b"}:     {cmd, "-a", "1,2", "-b", "3,4"},
				},
				cli.POSIX: {
					&[]string{"a"}:      {cmd, "-a", "1,2"},
					&[]string{"a", "b"}: {cmd, "-a", "1,2", "-b", "3,4"},
				},
			}
			var results []result
			var errs []error

			for syntax, argSets := range tests {
				for argNames, args := range argSets {
					os.Args = args
					cmd := cli.NewCommand("testcmd", context.Background())

					for i, name := range *argNames {
						val := []uint64{uint64(i)}
						bindVal := &val
						cmd.AddUint64ListArg(bindVal, &cli.ArgDefinition{Name: name, ShortName: rune(name[0])})
						results = append(results, result{[]uint64{uint64((i * 2) + 1), uint64((i * 2) + 2)}, bindVal})
					}

					if _, err := cli.NewParser(syntax, cmd).Parse(); err != nil {
						errs = append(errs, err)
					}
				}
			}

			return func() bool {
				for _, result := range results {
					for i, val := range result.val {
						if val != (*(result.bindVal))[i] {
							return false
						}
					}
				}

				return len(errs) == 0
			}
		},
	}

	for name, runtTest := range testCases {
		assertTest := runtTest()

		if !assertTest() {
			t.Fail()
			t.Log(n + ": " + name)
		}
	}
}
