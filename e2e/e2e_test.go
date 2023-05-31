package e2e_test

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/longyue0521/godog-multiple-unrelated-features-example/e2e"
	"github.com/longyue0521/godog-multiple-unrelated-features-example/e2e/features/godogs"
	"github.com/longyue0521/godog-multiple-unrelated-features-example/e2e/features/users"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

var (
	featFolderPath = "features/"
	jsonFolderPath = "reports/json/"
	htmlFolderPath = "reports/html/"

	defaultOptions = godog.Options{
		Output: colors.Colored(os.Stdout),
	}
	testSuiteGroup = e2e.NewTestSuiteGroup(featFolderPath, jsonFolderPath, htmlFolderPath, &defaultOptions)
)

func init() {
	godog.BindCommandLineFlags("godog.", &defaultOptions)                                                                             // godog v0.11.0 and later
	pflag.StringSliceVar(&defaultOptions.Paths, "feature", []string{}, "list of relative path from features folder to feature files") // godog not support --godog.path features/xx.feature
}

func TestMain(m *testing.M) {
	pflag.Parse()

	if err := testSuiteGroup.Before(); err != nil {
		os.Exit(101)
	}

	code := m.Run()

	if err := testSuiteGroup.After(); err != nil {
		os.Exit(102)
	}

	os.Exit(code)
}

func TestE2E(t *testing.T) {
	// t.Parallel()
	// When using the "pretty" format, step definitions are generated for undefined steps,
	// but if t.Parallel() is used, the output content will be printed in an unordered manner to the terminal.
	addGodogTestSuitesToTestSuiteGroup(t, testSuiteGroup)

	for name, suite := range testSuiteGroup.TestSuites() {
		name, suite := name, suite
		t.Run(name, func(t *testing.T) {
			// t.Parallel()
			suite.Options.TestingT = t
			require.Equal(t, 0, suite.Run(), "0 - success\n1 - failed\n2 - command line usage error\n128 - or higher, os signal related error exit codes")
		})
	}
}

func addGodogTestSuitesToTestSuiteGroup(t *testing.T, e *e2e.TestSuiteGroup) {
	t.Helper()

	require.NoError(t, e.AddTestSuite(users.GetAPI()))
	require.NoError(t, e.AddTestSuite(godogs.Eat()))
}
