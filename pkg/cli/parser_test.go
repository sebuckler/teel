package cli_test

import (
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
		"should error when arg passed with no args configured":       shouldErrorWhenArgPassedWithNoArgsConfigured,
		"should error when repeated arg is not repeatable":           shouldErrorWhenRepeatedArgIsNotRepeatable,
		"should error when GoFlag first arg is invalid format":       shouldErrorWhenGoFlagFirstArgIsInvalidFormat,
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
		"should error when unsupported arg type used":                shouldErrorWhenUnsupportedArgTypeUsed,
		"should parse when args provided correctly":                  shouldParseWhenPosixArgsProvidedCorrectly,
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

func shouldParseSubcommands(t *testing.T, n string) {
	os.Args = []string{"testcmd", "foo", "bar"}
	parser := cli.NewParser(cli.POSIX)
	parsedCmd, parseErr := parser.Parse(&cli.CommandConfig{
		Subcommands: []*cli.CommandConfig{{
			Name: "foo",
			Subcommands: []*cli.CommandConfig{{
				Name: "bar",
			}}},
		},
	})

	if parseErr != nil || len(parsedCmd.Subcommands) == 0 || len(parsedCmd.Subcommands[0].Subcommands) == 0 {
		t.Fail()
		t.Log(n + ": did not parse subcommands properly")
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
	testCases := map[cli.ArgSyntax][]string{
		cli.GoFlag: {"testcmd", "-aaa", "-aaa"},
		cli.GNU:    {"testcmd", "--aaa", "--aaa"},
		cli.POSIX:  {"testcmd", "-a", "-a"},
	}

	for syntax, args := range testCases {
		os.Args = args
		parser := cli.NewParser(syntax)
		a := true
		_, parseErr := parser.Parse(&cli.CommandConfig{
			Args: []*cli.ArgConfig{{
				Name:       strings.ReplaceAll(args[1], "-", ""),
				Repeatable: false,
				ShortName:  rune(strings.ReplaceAll(args[1], "-", "")[0]),
				Value:      &a,
			}},
		})

		if parseErr == nil {
			t.Fail()
			t.Log(n + ": did not error on non-repeatable arg being repeated")
		}
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
	testCases := map[string][]string{
		"invalid characters":          {"testcmd", "--="},
		"invalid multi-part name '='": {"testcmd", "--foo-="},
		"invalid multi-part name '`'": {"testcmd", "--foo-`"},
	}

	for name, args := range testCases {
		os.Args = args
		parser := cli.NewParser(cli.GNU)
		a := false
		_, parseErr := parser.Parse(&cli.CommandConfig{
			Args: []*cli.ArgConfig{{
				Name:      strings.TrimPrefix(args[1], "--"),
				ShortName: rune(strings.TrimPrefix(args[1], "--")[0]),
				Value:     &a,
			}},
		})

		if parseErr == nil {
			t.Fail()
			t.Log(n + ": did not error on incorrectly configured GNU arg name: " + name)
		}
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

func shouldParseWhenPosixArgsProvidedCorrectly(t *testing.T, n string) {
	cmd := "testcmd"
	testCases := map[string]func() func() bool{
		"bool options": func() func() bool {
			tests := map[cli.ArgSyntax]map[*[]string][]string{
				cli.GoFlag: {
					&[]string{"aaa"}:                        {cmd, "-aaa"},
					&[]string{"aaa", "bbb"}:                 {cmd, "-aaa", "-bbb"},
					&[]string{"a", "b", "c", "d", "e", "f"}: {cmd, "-a=1", "-b=0", "-c=t", "-d=f", "-e=T", "-f=F"},
					&[]string{"a", "b", "c"}:                {cmd, "-a=true", "-b=false", "-c=TRUE"},
					&[]string{"a", "b", "c"}:                {cmd, "-a=FALSE", "-b=True", "-c=False"},
				},
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
					var configs []*cli.ArgConfig

					for _, name := range *argNames {
						configs = append(configs, &cli.ArgConfig{Name: name, ShortName: rune(name[0]), Value: &val})
					}

					if _, parseErr := cli.NewParser(syntax).Parse(&cli.CommandConfig{Args: configs}); parseErr != nil {
						errs = append(errs, parseErr)
					}
				}
			}

			return func() bool { return len(errs) == 0 && val }
		},
		"float64 options": func() func() bool {
			tests := map[cli.ArgSyntax]map[*[]string][]string{
				cli.GoFlag: {
					&[]string{"aaa"}:        {cmd, "-aaa=1.0"},
					&[]string{"aaa", "bbb"}: {cmd, "-aaa=1.0", "-bbb=2.0"},
				},
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
					var configs []*cli.ArgConfig

					for i, name := range *argNames {
						val := float64(i)
						bindVal := &val
						configs = append(configs, &cli.ArgConfig{Name: name, ShortName: rune(name[0]), Value: bindVal})
						vals[val+1] = bindVal
					}

					if _, parseErr := cli.NewParser(syntax).Parse(&cli.CommandConfig{Args: configs}); parseErr != nil {
						errs = append(errs, parseErr)
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
				cli.GoFlag: {
					&[]string{"aaa"}:        {cmd, "-aaa=1.0,2.0"},
					&[]string{"aaa", "bbb"}: {cmd, "-aaa=1.0,2.0", "-bbb=3.0,4.0"},
				},
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
					var configs []*cli.ArgConfig

					for i, name := range *argNames {
						val := []float64{float64(i)}
						bindVal := &val
						configs = append(configs, &cli.ArgConfig{Name: name, ShortName: rune(name[0]), Value: bindVal})
						results = append(results, result{[]float64{float64((i * 2) + 1), float64((i * 2) + 2)}, bindVal})
					}

					if _, parseErr := cli.NewParser(syntax).Parse(&cli.CommandConfig{Args: configs}); parseErr != nil {
						errs = append(errs, parseErr)
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
				cli.GoFlag: {
					&[]string{"aaa"}:        {cmd, "-aaa=1"},
					&[]string{"aaa", "bbb"}: {cmd, "-aaa=1", "-bbb=2"},
				},
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
					var configs []*cli.ArgConfig

					for i, name := range *argNames {
						val := i
						bindVal := &val
						configs = append(configs, &cli.ArgConfig{Name: name, ShortName: rune(name[0]), Value: bindVal})
						vals[val+1] = bindVal
					}

					if _, parseErr := cli.NewParser(syntax).Parse(&cli.CommandConfig{Args: configs}); parseErr != nil {
						errs = append(errs, parseErr)
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
				cli.GoFlag: {
					&[]string{"aaa"}:        {cmd, "-aaa=1,2"},
					&[]string{"aaa", "bbb"}: {cmd, "-aaa=1,2", "-bbb=3,4"},
				},
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
					var configs []*cli.ArgConfig

					for i, name := range *argNames {
						val := []int{i}
						bindVal := &val
						configs = append(configs, &cli.ArgConfig{Name: name, ShortName: rune(name[0]), Value: bindVal})
						results = append(results, result{[]int{(i * 2) + 1, (i * 2) + 2}, bindVal})
					}

					if _, parseErr := cli.NewParser(syntax).Parse(&cli.CommandConfig{Args: configs}); parseErr != nil {
						errs = append(errs, parseErr)
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
				cli.GoFlag: {
					&[]string{"aaa"}:        {cmd, "-aaa=1"},
					&[]string{"aaa", "bbb"}: {cmd, "-aaa=1", "-bbb=2"},
				},
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
					var configs []*cli.ArgConfig

					for i, name := range *argNames {
						val := int64(i)
						bindVal := &val
						configs = append(configs, &cli.ArgConfig{Name: name, ShortName: rune(name[0]), Value: bindVal})
						vals[val+1] = bindVal
					}

					if _, parseErr := cli.NewParser(syntax).Parse(&cli.CommandConfig{Args: configs}); parseErr != nil {
						errs = append(errs, parseErr)
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
				cli.GoFlag: {
					&[]string{"aaa"}:        {cmd, "-aaa=1,2"},
					&[]string{"aaa", "bbb"}: {cmd, "-aaa=1,2", "-bbb=3,4"},
				},
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
					var configs []*cli.ArgConfig

					for i, name := range *argNames {
						val := []int64{int64(i)}
						bindVal := &val
						configs = append(configs, &cli.ArgConfig{Name: name, ShortName: rune(name[0]), Value: bindVal})
						results = append(results, result{[]int64{int64((i * 2) + 1), int64((i * 2) + 2)}, bindVal})
					}

					if _, parseErr := cli.NewParser(syntax).Parse(&cli.CommandConfig{Args: configs}); parseErr != nil {
						errs = append(errs, parseErr)
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
				cli.GoFlag: {
					&[]string{"aaa"}:        {cmd, "-aaa=1"},
					&[]string{"aaa", "bbb"}: {cmd, "-aaa=1", "-bbb=2"},
				},
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
					var configs []*cli.ArgConfig

					for i, name := range *argNames {
						val := strconv.Itoa(i + 1)
						bindVal := &val
						configs = append(configs, &cli.ArgConfig{Name: name, ShortName: rune(name[0]), Value: bindVal})
						vals[*bindVal] = bindVal
					}

					if _, parseErr := cli.NewParser(syntax).Parse(&cli.CommandConfig{Args: configs}); parseErr != nil {
						errs = append(errs, parseErr)
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
				cli.GoFlag: {
					&[]string{"aaa"}:        {cmd, "-aaa=1,2"},
					&[]string{"aaa", "bbb"}: {cmd, "-aaa=1,2", "-bbb=3,4"},
				},
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
					var configs []*cli.ArgConfig

					for i, name := range *argNames {
						val := []string{strconv.Itoa(i + 1)}
						bindVal := &val
						configs = append(configs, &cli.ArgConfig{Name: name, ShortName: rune(name[0]), Value: bindVal})
						results = append(
							results,
							result{[]string{strconv.Itoa((i * 2) + 1), strconv.Itoa((i * 2) + 2)}, bindVal},
						)
					}

					if _, parseErr := cli.NewParser(syntax).Parse(&cli.CommandConfig{Args: configs}); parseErr != nil {
						errs = append(errs, parseErr)
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
				cli.GoFlag: {
					&[]string{"aaa"}:        {cmd, "-aaa=1"},
					&[]string{"aaa", "bbb"}: {cmd, "-aaa=1", "-bbb=2"},
				},
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
					var configs []*cli.ArgConfig

					for i, name := range *argNames {
						val := uint(i)
						bindVal := &val
						configs = append(configs, &cli.ArgConfig{Name: name, ShortName: rune(name[0]), Value: bindVal})
						vals[val+1] = bindVal
					}

					if _, parseErr := cli.NewParser(syntax).Parse(&cli.CommandConfig{Args: configs}); parseErr != nil {
						errs = append(errs, parseErr)
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
				cli.GoFlag: {
					&[]string{"aaa"}:        {cmd, "-aaa=1,2"},
					&[]string{"aaa", "bbb"}: {cmd, "-aaa=1,2", "-bbb=3,4"},
				},
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
					var configs []*cli.ArgConfig

					for i, name := range *argNames {
						val := []uint{uint(i)}
						bindVal := &val
						configs = append(configs, &cli.ArgConfig{Name: name, ShortName: rune(name[0]), Value: bindVal})
						results = append(results, result{[]uint{uint((i * 2) + 1), uint((i * 2) + 2)}, bindVal})
					}

					if _, parseErr := cli.NewParser(syntax).Parse(&cli.CommandConfig{Args: configs}); parseErr != nil {
						errs = append(errs, parseErr)
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
				cli.GoFlag: {
					&[]string{"aaa"}:        {cmd, "-aaa=1"},
					&[]string{"aaa", "bbb"}: {cmd, "-aaa=1", "-bbb=2"},
				},
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
					var configs []*cli.ArgConfig

					for i, name := range *argNames {
						val := uint64(i)
						bindVal := &val
						configs = append(configs, &cli.ArgConfig{Name: name, ShortName: rune(name[0]), Value: bindVal})
						vals[val+1] = bindVal
					}

					if _, parseErr := cli.NewParser(syntax).Parse(&cli.CommandConfig{Args: configs}); parseErr != nil {
						errs = append(errs, parseErr)
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
				cli.GoFlag: {
					&[]string{"aaa"}:        {cmd, "-aaa=1,2"},
					&[]string{"aaa", "bbb"}: {cmd, "-aaa=1,2", "--bbb=3,4"},
				},
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
					var configs []*cli.ArgConfig

					for i, name := range *argNames {
						val := []uint64{uint64(i)}
						bindVal := &val
						configs = append(configs, &cli.ArgConfig{Name: name, ShortName: rune(name[0]), Value: bindVal})
						results = append(results, result{[]uint64{uint64((i * 2) + 1), uint64((i * 2) + 2)}, bindVal})
					}

					if _, parseErr := cli.NewParser(syntax).Parse(&cli.CommandConfig{Args: configs}); parseErr != nil {
						errs = append(errs, parseErr)
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
