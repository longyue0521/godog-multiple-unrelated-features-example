package e2e_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/longyue0521/godog-multiple-unrelated-features-example/e2e/features/godogs"
	"github.com/longyue0521/godog-multiple-unrelated-features-example/e2e/features/users"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"

	"gitlab.com/rodrigoodhin/gocure/models"
	"gitlab.com/rodrigoodhin/gocure/pkg/gocure"
	"gitlab.com/rodrigoodhin/gocure/report/html"
)

var (
	featFolderPath = "features/"
	jsonFolderPath = "reports/json/"
	htmlFolderPath = "reports/html/"

	defaultOptions = godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "progress",
	}
)

func init() {
	godog.BindCommandLineFlags("godog.", &defaultOptions) // godog v0.11.0 and later
}

type E2ETestSuite struct {
	suites map[string]*godog.TestSuite
}

func NewE2ETestSuite() *E2ETestSuite {
	return &E2ETestSuite{suites: make(map[string]*godog.TestSuite)}
}

func (e *E2ETestSuite) AddTestSuite(suite *godog.TestSuite) error {
	_, err := os.Stat(suite.Name)
	if err != nil {
		return fmt.Errorf("godog.TestSuite's name is wrong: %s", suite.Name)
	}
	rel, err := filepath.Rel(featFolderPath, suite.Name)
	if err != nil {
		return err
	}
	name := strings.TrimSuffix(rel, ".feature")
	e.suites[name] = suite
	return nil
}

func (e *E2ETestSuite) GodogTestSuites() map[string]*godog.TestSuite {
	// TODO use values of --godog.Paths to filter
	return e.suites
}

func TestMain(m *testing.M) {
	// TODO 将m.Run()前面的重构为 E2ETestSuite.Setup/Before
	parseCommandLineFlagsIntoDefaultOptions()

	if deleteReports(jsonFolderPath) != nil || deleteReports(htmlFolderPath) != nil {
		os.Exit(101)
	}

	if createSubdirectories(jsonFolderPath, featFolderPath) != nil {
		os.Exit(102)
	}

	code := m.Run()
	// // TODO 将m.Run()后面的重构为 E2ETestSuite.Teardown/After
	if generateHTMLReport("report", jsonFolderPath, htmlFolderPath) != nil {
		os.Exit(103)
	}

	os.Exit(code)

}

func parseCommandLineFlagsIntoDefaultOptions() {
	pflag.Parse()
}

func deleteReports(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return os.Remove(path)
		}
		return nil
	})
}

func createSubdirectories(dst string, src string) error {
	// TODO 同步src与dst中的子目录,删除dst中与src不同的子目录
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		return os.MkdirAll(filepath.Join(dst, relPath), 0755)
	})
}

func generateHTMLReport(name string, jsonFolderPath string, htmlFolderPath string) error {
	HTML := gocure.HTML{
		Config: html.Data{
			Title:            name,
			MergeFiles:       true,
			InputFolderPath:  jsonFolderPath,
			OutputHtmlFolder: htmlFolderPath,
			Metadata: models.Metadata{
				AppVersion:      "0.8.7",
				TestEnvironment: "development",
				Browser:         "Google Chrome",
				Platform:        "Linux",
				Parallel:        "Scenarios",
				Executed:        "Remote",
			},
		},
	}
	return HTML.Generate()
}

func TestE2E(t *testing.T) {

	e := NewE2ETestSuite()

	addGodogTestSuitesToE2ETestSuite(t, e)

	for name, suite := range e.GodogTestSuites() {
		// suite := suite
		t.Run(name, func(t *testing.T) {
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

func addGodogTestSuitesToE2ETestSuite(t *testing.T, e *E2ETestSuite) {
	t.Helper()

	require.NoError(t, e.AddTestSuite(users.GetAPI()))
	require.NoError(t, e.AddTestSuite(godogs.Eat()))
}
