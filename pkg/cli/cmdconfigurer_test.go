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
		"should have only command defined": shouldHaveOnlyCommandDefined,
		"should have command with only bool arg": shouldHaveCommandWithOnlyBoolArg,
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
		t.Log(n + ": command incorrectly configured")
	}
}
