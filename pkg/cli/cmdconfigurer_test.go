package cli_test

import (
	"context"
	"github.com/sebuckler/teel/pkg/cli"
	"testing"
)

func TestCommandConfigurer_Configure(t *testing.T) {
	testCases := getConfigurerTestCases()

	for name, test := range testCases {
		test(t, name)
	}
}

func getConfigurerTestCases() map[string]func(t *testing.T, n string) {
	return map[string]func(t *testing.T, n string){
		"should have only command defined":              shouldHaveOnlyCommandDefined,
		"should have command with only bool arg":        shouldHaveCommandWithOnlyBoolArg,
		"should have command with only bool list arg":   shouldHaveCommandWithOnlyBoolListArg,
		"should have command with only int arg":         shouldHaveCommandWithOnlyIntArg,
		"should have command with only int list arg":    shouldHaveCommandWithOnlyIntListArg,
		"should have command with only int64 arg":       shouldHaveCommandWithOnlyInt64Arg,
		"should have command with only int64 list arg":  shouldHaveCommandWithOnlyInt64ListArg,
		"should have command with only string arg":      shouldHaveCommandWithOnlyStringArg,
		"should have command with only string list arg": shouldHaveCommandWithOnlyStringListArg,
		"should have command with only uint arg":        shouldHaveCommandWithOnlyUintArg,
		"should have command with only uint list arg":   shouldHaveCommandWithOnlyUintListArg,
		"should have command with only uint64 arg":      shouldHaveCommandWithOnlyUint64Arg,
		"should have command with only uint64 list arg": shouldHaveCommandWithOnlyUint64ListArg,
		"should have command with run function":         shouldHaveCommandWithRunFunction,
		"should have subcommands":                       shouldHaveSubcommands,
	}
}

func shouldHaveOnlyCommandDefined(t *testing.T, n string) {
	cmd := cli.NewCommand("foo", context.Background())
	config := cmd.Configure()

	if config.Name != "foo" || len(config.Args) > 0 || config.Run != nil || len(config.Subcommands) > 0 {
		t.Fail()
		t.Log(n + ": command incorrectly configured")
	}
}

func shouldHaveCommandWithOnlyBoolArg(t *testing.T, n string) {
	cmd := cli.NewCommand("foo", context.Background())
	cmd.AddBoolArg("bar", 'b', nil, false, "", false)
	config := cmd.Configure()
	_, ok := config.Args[0].Value.(*bool)

	if len(config.Args) != 1 || config.Args[0].Name != "bar" || !ok {
		t.Fail()
		t.Log(n + ": args incorrectly configured")
	}
}

func shouldHaveCommandWithOnlyBoolListArg(t *testing.T, n string) {
	cmd := cli.NewCommand("foo", context.Background())
	cmd.AddBoolListArg("bar", 'b', nil, []bool{false}, "", false)
	config := cmd.Configure()
	_, ok := config.Args[0].Value.(*[]bool)

	if len(config.Args) != 1 || config.Args[0].Name != "bar" || !ok {
		t.Fail()
		t.Log(n + ": args incorrectly configured")
	}
}

func shouldHaveCommandWithOnlyIntArg(t *testing.T, n string) {
	cmd := cli.NewCommand("foo", context.Background())
	cmd.AddIntArg("bar", 'b', nil, 0, "", false)
	config := cmd.Configure()
	_, ok := config.Args[0].Value.(*int)

	if len(config.Args) != 1 || config.Args[0].Name != "bar" || !ok {
		t.Fail()
		t.Log(n + ": args incorrectly configured")
	}
}

func shouldHaveCommandWithOnlyIntListArg(t *testing.T, n string) {
	cmd := cli.NewCommand("foo", context.Background())
	cmd.AddIntListArg("bar", 'b', nil, []int{0}, "", false)
	config := cmd.Configure()
	_, ok := config.Args[0].Value.(*[]int)

	if len(config.Args) != 1 || config.Args[0].Name != "bar" || !ok {
		t.Fail()
		t.Log(n + ": args incorrectly configured")
	}
}

func shouldHaveCommandWithOnlyInt64Arg(t *testing.T, n string) {
	cmd := cli.NewCommand("foo", context.Background())
	cmd.AddInt64Arg("bar", 'b', nil, 0, "", false)
	config := cmd.Configure()
	_, ok := config.Args[0].Value.(*int64)

	if len(config.Args) != 1 || config.Args[0].Name != "bar" || !ok {
		t.Fail()
		t.Log(n + ": args incorrectly configured")
	}
}

func shouldHaveCommandWithOnlyInt64ListArg(t *testing.T, n string) {
	cmd := cli.NewCommand("foo", context.Background())
	cmd.AddInt64ListArg("bar", 'b', nil, []int64{0}, "", false)
	config := cmd.Configure()
	_, ok := config.Args[0].Value.(*[]int64)

	if len(config.Args) != 1 || config.Args[0].Name != "bar" || !ok {
		t.Fail()
		t.Log(n + ": args incorrectly configured")
	}
}

func shouldHaveCommandWithOnlyStringArg(t *testing.T, n string) {
	cmd := cli.NewCommand("foo", context.Background())
	cmd.AddStringArg("bar", 'b', nil, "", "", false)
	config := cmd.Configure()
	_, ok := config.Args[0].Value.(*string)

	if len(config.Args) != 1 || config.Args[0].Name != "bar" || !ok {
		t.Fail()
		t.Log(n + ": args incorrectly configured")
	}
}

func shouldHaveCommandWithOnlyStringListArg(t *testing.T, n string) {
	cmd := cli.NewCommand("foo", context.Background())
	cmd.AddStringListArg("bar", 'b', nil, []string{""}, "", false)
	config := cmd.Configure()
	_, ok := config.Args[0].Value.(*[]string)

	if len(config.Args) != 1 || config.Args[0].Name != "bar" || !ok {
		t.Fail()
		t.Log(n + ": args incorrectly configured")
	}
}

func shouldHaveCommandWithOnlyUintArg(t *testing.T, n string) {
	cmd := cli.NewCommand("foo", context.Background())
	cmd.AddUintArg("bar", 'b', nil, 0, "", false)
	config := cmd.Configure()
	_, ok := config.Args[0].Value.(*uint)

	if len(config.Args) != 1 || config.Args[0].Name != "bar" || !ok {
		t.Fail()
		t.Log(n + ": args incorrectly configured")
	}
}

func shouldHaveCommandWithOnlyUintListArg(t *testing.T, n string) {
	cmd := cli.NewCommand("foo", context.Background())
	cmd.AddUintListArg("bar", 'b', nil, []uint{0}, "", false)
	config := cmd.Configure()
	_, ok := config.Args[0].Value.(*[]uint)

	if len(config.Args) != 1 || config.Args[0].Name != "bar" || !ok {
		t.Fail()
		t.Log(n + ": args incorrectly configured")
	}
}

func shouldHaveCommandWithOnlyUint64Arg(t *testing.T, n string) {
	cmd := cli.NewCommand("foo", context.Background())
	cmd.AddUint64Arg("bar", 'b', nil, 0, "", false)
	config := cmd.Configure()
	_, ok := config.Args[0].Value.(*uint64)

	if len(config.Args) != 1 || config.Args[0].Name != "bar" || !ok {
		t.Fail()
		t.Log(n + ": args incorrectly configured")
	}
}

func shouldHaveCommandWithOnlyUint64ListArg(t *testing.T, n string) {
	cmd := cli.NewCommand("foo", context.Background())
	cmd.AddUint64ListArg("bar", 'b', nil, []uint64{0}, "", false)
	config := cmd.Configure()
	_, ok := config.Args[0].Value.(*[]uint64)

	if len(config.Args) != 1 || config.Args[0].Name != "bar" || !ok {
		t.Fail()
		t.Log(n + ": args incorrectly configured")
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
