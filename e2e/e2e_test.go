package e2e_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

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
		Format: "progress",
	}
	testSuiteGroup = e2e.NewTestSuiteGroup(featFolderPath, jsonFolderPath, htmlFolderPath)
)

func init() {
	godog.BindCommandLineFlags("godog.", &defaultOptions) // godog v0.11.0 and later
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
	t.Parallel()
	// TODO 看看使用parallel后乱序生成的html报告是否可读
	// TODO 如果可读,则考虑将format: progress,cucumber:%s.json, pretty在并行运行测试的表现不好,会乱序输出提示信息
	// TODO 还要看看, 对于新定义的feature,在parallel模式下,是否会自动生成steps,生成的steps定义是否可读
	// 一旦添加t.Parallel(),go框架会自行决定是否并发,命令行参数只能控制并发度--parallel=16

	addGodogTestSuitesToE2ETestSuite(t, testSuiteGroup)

	for name, suite := range testSuiteGroup.TestSuites() {
		suite := suite
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			options := defaultOptions
			options.TestingT = t
			options.Randomize = time.Now().UnixNano()
			options.Paths = []string{suite.Name}
			options.Format = fmt.Sprintf("pretty,cucumber:%s.json", filepath.Join(jsonFolderPath, name))
			suite.Options = &options
			require.Equal(t, 0, suite.Run(), "0 - success\n1 - failed\n2 - command line usage error\n128 - or higher, os signal related error exit codes")
		})
	}
}

func addGodogTestSuitesToE2ETestSuite(t *testing.T, e *e2e.TestSuiteGroup) {
	t.Helper()

	require.NoError(t, e.AddTestSuite(users.GetAPI()))
	require.NoError(t, e.AddTestSuite(godogs.Eat()))
}
