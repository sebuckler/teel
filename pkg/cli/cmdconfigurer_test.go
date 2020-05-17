package cli_test

import (
	"context"
	"github.com/sebuckler/teel/pkg/cli"
	"testing"
)

func TestCommandConfigurer_Configure(t *testing.T) {
	for name, test := range getConfigurerTestCases() {
		test(t, name)
	}
}

func getConfigurerTestCases() map[string]func(t *testing.T, n string) {
	return map[string]func(t *testing.T, n string){
		"should have only command defined":                        shouldHaveOnlyCommandDefined,
		"should have command with only bool arg":                  shouldHaveCommandWithOnlyBoolArg,
		"should have command with only float64 arg":               shouldHaveCommandWithOnlyFloat64Arg,
		"should have command with only float64 list arg":          shouldHaveCommandWithOnlyFloat64ListArg,
		"should have command with only int arg":                   shouldHaveCommandWithOnlyIntArg,
		"should have command with only int list arg":              shouldHaveCommandWithOnlyIntListArg,
		"should have command with only int64 arg":                 shouldHaveCommandWithOnlyInt64Arg,
		"should have command with only int64 list arg":            shouldHaveCommandWithOnlyInt64ListArg,
		"should have command with only string arg":                shouldHaveCommandWithOnlyStringArg,
		"should have command with only string list arg":           shouldHaveCommandWithOnlyStringListArg,
		"should have command with only uint arg":                  shouldHaveCommandWithOnlyUintArg,
		"should have command with only uint list arg":             shouldHaveCommandWithOnlyUintListArg,
		"should have command with only uint64 arg":                shouldHaveCommandWithOnlyUint64Arg,
		"should have command with only uint64 list arg":           shouldHaveCommandWithOnlyUint64ListArg,
		"should have command with run function":                   shouldHaveCommandWithRunFunction,
		"should have subcommands":                                 shouldHaveSubcommands,
		"should have only help command when no other arg defined": shouldHaveOnlyHelpCommandWhenNoOtherArgDefined,
	}
}

func shouldHaveOnlyCommandDefined(t *testing.T, n string) {
	cmd := cli.NewCommand("foo", context.Background())
	config := cmd.Configure()

	if config.Name != "foo" || config.Run != nil || len(config.Subcommands) > 0 {
		t.Fail()
		t.Log(n + ": command incorrectly configured")
	}
}

func shouldHaveCommandWithOnlyBoolArg(t *testing.T, n string) {
	testCases := map[string]func() (*bool, *cli.ArgDefinition){
		"nil pointer value": func() (*bool, *cli.ArgDefinition) {
			return nil, &cli.ArgDefinition{Name: "bar", ShortName: 'b'}
		},
		"set pointer value": func() (*bool, *cli.ArgDefinition) {
			val := true
			return &val, &cli.ArgDefinition{Name: "bar", ShortName: 'b'}
		},
	}

	for name, test := range testCases {
		cmd := cli.NewCommand("foo", context.Background())
		val, def := test()
		cmd.AddBoolArg(val, def)
		config := cmd.Configure()
		_, ok := config.Args[0].Value.(*bool)

		if config.Args[0].Name != "bar" || !ok {
			t.Fail()
			t.Log(n + ": args incorrectly configured for " + name)
		}
	}
}

func shouldHaveCommandWithOnlyFloat64Arg(t *testing.T, n string) {
	argDef := &cli.ArgDefinition{Name: "bar", ShortName: 'b'}
	testCases := map[string]func(c cli.CommandConfigurer) func(v *float64) bool{
		"nil pointer value": func(c cli.CommandConfigurer) func(v *float64) bool {
			c.AddFloat64Arg(nil, argDef)
			return func(v *float64) bool { return v == nil }
		},
		"set pointer value": func(c cli.CommandConfigurer) func(v *float64) bool {
			val := float64(1)
			c.AddFloat64Arg(&val, argDef)
			return func(v *float64) bool { return *v == val }
		},
	}

	for name, test := range testCases {
		cmd := cli.NewCommand("foo", context.Background())
		success := test(cmd)
		config := cmd.Configure()
		val, ok := config.Args[0].Value.(*float64)

		if config.Args[0].Name != "bar" || !ok || !success(val) {
			t.Fail()
			t.Log(n + ": args incorrectly configured for " + name)
		}
	}
}

func shouldHaveCommandWithOnlyFloat64ListArg(t *testing.T, n string) {
	argDef := &cli.ArgDefinition{Name: "bar", ShortName: 'b'}
	testCases := map[string]func(c cli.CommandConfigurer) func(v *[]float64) bool{
		"nil pointer value": func(c cli.CommandConfigurer) func(v *[]float64) bool {
			c.AddFloat64ListArg(nil, argDef)
			return func(v *[]float64) bool { return v == nil }
		},
		"set pointer value": func(c cli.CommandConfigurer) func(v *[]float64) bool {
			val := []float64{1}
			c.AddFloat64ListArg(&val, argDef)
			return func(v *[]float64) bool { return len(*v) == len(val) && (*v)[0] == val[0] }
		},
	}

	for name, test := range testCases {
		cmd := cli.NewCommand("foo", context.Background())
		success := test(cmd)
		config := cmd.Configure()
		val, ok := config.Args[0].Value.(*[]float64)

		if config.Args[0].Name != "bar" || !ok || !success(val) {
			t.Fail()
			t.Log(n + ": args incorrectly configured for " + name)
		}
	}
}

func shouldHaveCommandWithOnlyIntArg(t *testing.T, n string) {
	argDef := &cli.ArgDefinition{Name: "bar", ShortName: 'b'}
	testCases := map[string]func(c cli.CommandConfigurer) func(v *int) bool{
		"nil pointer value": func(c cli.CommandConfigurer) func(v *int) bool {
			c.AddIntArg(nil, argDef)
			return func(v *int) bool { return v == nil }
		},
		"set pointer value": func(c cli.CommandConfigurer) func(v *int) bool {
			val := 1
			c.AddIntArg(&val, argDef)
			return func(v *int) bool { return *v == val }
		},
	}

	for name, test := range testCases {
		cmd := cli.NewCommand("foo", context.Background())
		success := test(cmd)
		config := cmd.Configure()
		val, ok := config.Args[0].Value.(*int)

		if config.Args[0].Name != "bar" || !ok || !success(val) {
			t.Fail()
			t.Log(n + ": args incorrectly configured for " + name)
		}
	}
}

func shouldHaveCommandWithOnlyIntListArg(t *testing.T, n string) {
	argDef := &cli.ArgDefinition{Name: "bar", ShortName: 'b'}
	testCases := map[string]func(c cli.CommandConfigurer) func(v *[]int) bool{
		"nil pointer value": func(c cli.CommandConfigurer) func(v *[]int) bool {
			c.AddIntListArg(nil, argDef)
			return func(v *[]int) bool { return v == nil }
		},
		"set pointer value": func(c cli.CommandConfigurer) func(v *[]int) bool {
			val := []int{1}
			c.AddIntListArg(&val, argDef)
			return func(v *[]int) bool { return len(*v) == len(val) && (*v)[0] == val[0] }
		},
	}

	for name, test := range testCases {
		cmd := cli.NewCommand("foo", context.Background())
		success := test(cmd)
		config := cmd.Configure()
		val, ok := config.Args[0].Value.(*[]int)

		if config.Args[0].Name != "bar" || !ok || !success(val) {
			t.Fail()
			t.Log(n + ": args incorrectly configured for " + name)
		}
	}
}

func shouldHaveCommandWithOnlyInt64Arg(t *testing.T, n string) {
	argDef := &cli.ArgDefinition{Name: "bar", ShortName: 'b'}
	testCases := map[string]func(c cli.CommandConfigurer) func(v *int64) bool{
		"nil pointer value": func(c cli.CommandConfigurer) func(v *int64) bool {
			c.AddInt64Arg(nil, argDef)
			return func(v *int64) bool { return v == nil }
		},
		"set pointer value": func(c cli.CommandConfigurer) func(v *int64) bool {
			val := int64(1)
			c.AddInt64Arg(&val, argDef)
			return func(v *int64) bool { return *v == val }
		},
	}

	for name, test := range testCases {
		cmd := cli.NewCommand("foo", context.Background())
		success := test(cmd)
		config := cmd.Configure()
		val, ok := config.Args[0].Value.(*int64)

		if config.Args[0].Name != "bar" || !ok || !success(val) {
			t.Fail()
			t.Log(n + ": args incorrectly configured for " + name)
		}
	}
}

func shouldHaveCommandWithOnlyInt64ListArg(t *testing.T, n string) {
	argDef := &cli.ArgDefinition{Name: "bar", ShortName: 'b'}
	testCases := map[string]func(c cli.CommandConfigurer) func(v *[]int64) bool{
		"nil pointer value": func(c cli.CommandConfigurer) func(v *[]int64) bool {
			c.AddInt64ListArg(nil, argDef)
			return func(v *[]int64) bool { return v == nil }
		},
		"set pointer value": func(c cli.CommandConfigurer) func(v *[]int64) bool {
			val := []int64{1}
			c.AddInt64ListArg(&val, argDef)
			return func(v *[]int64) bool { return len(*v) == len(val) && (*v)[0] == val[0] }
		},
	}

	for name, test := range testCases {
		cmd := cli.NewCommand("foo", context.Background())
		success := test(cmd)
		config := cmd.Configure()
		val, ok := config.Args[0].Value.(*[]int64)

		if config.Args[0].Name != "bar" || !ok || !success(val) {
			t.Fail()
			t.Log(n + ": args incorrectly configured for " + name)
		}
	}
}

func shouldHaveCommandWithOnlyStringArg(t *testing.T, n string) {
	argDef := &cli.ArgDefinition{Name: "bar", ShortName: 'b'}
	testCases := map[string]func(c cli.CommandConfigurer) func(v *string) bool{
		"nil pointer value": func(c cli.CommandConfigurer) func(v *string) bool {
			c.AddStringArg(nil, argDef)
			return func(v *string) bool { return v == nil }
		},
		"set pointer value": func(c cli.CommandConfigurer) func(v *string) bool {
			val := "value"
			c.AddStringArg(&val, argDef)
			return func(v *string) bool { return *v == val }
		},
	}

	for name, test := range testCases {
		cmd := cli.NewCommand("foo", context.Background())
		success := test(cmd)
		config := cmd.Configure()
		val, ok := config.Args[0].Value.(*string)

		if config.Args[0].Name != "bar" || !ok || !success(val) {
			t.Fail()
			t.Log(n + ": args incorrectly configured for " + name)
		}
	}
}

func shouldHaveCommandWithOnlyStringListArg(t *testing.T, n string) {
	argDef := &cli.ArgDefinition{Name: "bar", ShortName: 'b'}
	testCases := map[string]func(c cli.CommandConfigurer) func(v *[]string) bool{
		"nil pointer value": func(c cli.CommandConfigurer) func(v *[]string) bool {
			c.AddStringListArg(nil, argDef)
			return func(v *[]string) bool { return v == nil }
		},
		"set pointer value": func(c cli.CommandConfigurer) func(v *[]string) bool {
			val := []string{"value"}
			c.AddStringListArg(&val, argDef)
			return func(v *[]string) bool { return len(*v) == len(val) && (*v)[0] == val[0] }
		},
	}

	for name, test := range testCases {
		cmd := cli.NewCommand("foo", context.Background())
		success := test(cmd)
		config := cmd.Configure()
		val, ok := config.Args[0].Value.(*[]string)

		if config.Args[0].Name != "bar" || !ok || !success(val) {
			t.Fail()
			t.Log(n + ": args incorrectly configured for " + name)
		}
	}
}

func shouldHaveCommandWithOnlyUintArg(t *testing.T, n string) {
	argDef := &cli.ArgDefinition{Name: "bar", ShortName: 'b'}
	testCases := map[string]func(c cli.CommandConfigurer) func(v *uint) bool{
		"nil pointer value": func(c cli.CommandConfigurer) func(v *uint) bool {
			c.AddUintArg(nil, argDef)
			return func(v *uint) bool { return v == nil }
		},
		"set pointer value": func(c cli.CommandConfigurer) func(v *uint) bool {
			val := uint(1)
			c.AddUintArg(&val, argDef)
			return func(v *uint) bool { return *v == val }
		},
	}

	for name, test := range testCases {
		cmd := cli.NewCommand("foo", context.Background())
		success := test(cmd)
		config := cmd.Configure()
		val, ok := config.Args[0].Value.(*uint)

		if config.Args[0].Name != "bar" || !ok || !success(val) {
			t.Fail()
			t.Log(n + ": args incorrectly configured for " + name)
		}
	}
}

func shouldHaveCommandWithOnlyUintListArg(t *testing.T, n string) {
	argDef := &cli.ArgDefinition{Name: "bar", ShortName: 'b'}
	testCases := map[string]func(c cli.CommandConfigurer) func(v *[]uint) bool{
		"nil pointer value": func(c cli.CommandConfigurer) func(v *[]uint) bool {
			c.AddUintListArg(nil, argDef)
			return func(v *[]uint) bool { return v == nil }
		},
		"set pointer value": func(c cli.CommandConfigurer) func(v *[]uint) bool {
			val := []uint{1}
			c.AddUintListArg(&val, argDef)
			return func(v *[]uint) bool { return len(*v) == len(val) && (*v)[0] == val[0] }
		},
	}

	for name, test := range testCases {
		cmd := cli.NewCommand("foo", context.Background())
		success := test(cmd)
		config := cmd.Configure()
		val, ok := config.Args[0].Value.(*[]uint)

		if config.Args[0].Name != "bar" || !ok || !success(val) {
			t.Fail()
			t.Log(n + ": args incorrectly configured for " + name)
		}
	}
}

func shouldHaveCommandWithOnlyUint64Arg(t *testing.T, n string) {
	argDef := &cli.ArgDefinition{Name: "bar", ShortName: 'b'}
	testCases := map[string]func(c cli.CommandConfigurer) func(v *uint64) bool{
		"nil pointer value": func(c cli.CommandConfigurer) func(v *uint64) bool {
			c.AddUint64Arg(nil, argDef)
			return func(v *uint64) bool { return v == nil }
		},
		"set pointer value": func(c cli.CommandConfigurer) func(v *uint64) bool {
			val := uint64(1)
			c.AddUint64Arg(&val, argDef)
			return func(v *uint64) bool { return *v == val }
		},
	}

	for name, test := range testCases {
		cmd := cli.NewCommand("foo", context.Background())
		success := test(cmd)
		config := cmd.Configure()
		val, ok := config.Args[0].Value.(*uint64)

		if config.Args[0].Name != "bar" || !ok || !success(val) {
			t.Fail()
			t.Log(n + ": args incorrectly configured for " + name)
		}
	}
}

func shouldHaveCommandWithOnlyUint64ListArg(t *testing.T, n string) {
	argDef := &cli.ArgDefinition{Name: "bar", ShortName: 'b'}
	testCases := map[string]func(c cli.CommandConfigurer) func(v *[]uint64) bool{
		"nil pointer value": func(c cli.CommandConfigurer) func(v *[]uint64) bool {
			c.AddUint64ListArg(nil, argDef)
			return func(v *[]uint64) bool { return v == nil }
		},
		"set pointer value": func(c cli.CommandConfigurer) func(v *[]uint64) bool {
			val := []uint64{1}
			c.AddUint64ListArg(&val, argDef)
			return func(v *[]uint64) bool { return len(*v) == len(val) && (*v)[0] == val[0] }
		},
	}

	for name, test := range testCases {
		cmd := cli.NewCommand("foo", context.Background())
		success := test(cmd)
		config := cmd.Configure()
		val, ok := config.Args[0].Value.(*[]uint64)

		if config.Args[0].Name != "bar" || !ok || !success(val) {
			t.Fail()
			t.Log(n + ": args incorrectly configured for " + name)
		}
	}
}

func shouldHaveCommandWithRunFunction(t *testing.T, n string) {
	cmd := cli.NewCommand("foo", context.Background())
	cmd.AddRunFunc(func(ctx context.Context, o []string) {
		t.Fail()
		t.Log(n + ": should not have run command")
	})
	config := cmd.Configure()

	if config.Run == nil {
		t.Fail()
		t.Log(n + ": run function incorrectly configured")
	}
}

func shouldHaveSubcommands(t *testing.T, n string) {
	ctx := context.Background()
	cmd := cli.NewCommand("foo", ctx)
	cmd.AddSubcommand(cli.NewCommand("bar", ctx))
	cmd.AddSubcommand(cli.NewCommand("bar2", ctx))
	cmd.AddSubcommand(cli.NewCommand("bar3", ctx))
	config := cmd.Configure()

	if len(config.Subcommands) != 3 {
		t.Fail()
		t.Log(n + ": subcommands incorrectly configured")
	}

	for _, subCmd := range config.Subcommands {
		switch subCmd.Name {
		case "bar":
		case "bar2":
		case "bar3":
		default:
			t.Fail()
			t.Log(n + ": subcommand names incorrectly configured")
		}
	}
}

func shouldHaveOnlyHelpCommandWhenNoOtherArgDefined(t *testing.T, n string) {
	ctx := context.Background()
	cmd := cli.NewCommand("foo", ctx)
	cmd.AddBoolArg(nil, nil)
	config := cmd.Configure()

	if len(config.Args) > 2 || config.Args[1].Name != "help" {
		t.Fail()
		t.Log(n + ": help arg not configured")
	}
}
