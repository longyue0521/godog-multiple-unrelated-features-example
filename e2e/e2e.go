package e2e

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cucumber/godog"
	"gitlab.com/rodrigoodhin/gocure/models"
	"gitlab.com/rodrigoodhin/gocure/pkg/gocure"
	"gitlab.com/rodrigoodhin/gocure/report/html"
)

var (
	errReportNotFound = errors.New("e2e: report not found")
)

type TestSuiteGroup struct {
	featFolderRootPath       string
	reportsFolderRootPath    string
	htmlReportFolderRootPath string
	customizedOptions        *godog.Options
	suites                   map[string]*godog.TestSuite
}

func NewTestSuiteGroup(featFolderRootPath, reportsFolderRootPath, htmlReportFolderRootPath string, options *godog.Options) *TestSuiteGroup {
	return &TestSuiteGroup{
		featFolderRootPath:       featFolderRootPath,
		reportsFolderRootPath:    reportsFolderRootPath,
		htmlReportFolderRootPath: htmlReportFolderRootPath,
		customizedOptions:        options,
		suites:                   make(map[string]*godog.TestSuite)}
}

func (e *TestSuiteGroup) Before() error {
	if err := e.deleteReports(e.reportsFolderRootPath); err != nil {
		return err
	}
	if err := e.deleteReports(e.htmlReportFolderRootPath); err != nil {
		return err
	}
	if err := e.syncSubdirectoriesFromFeaturesFolderToReportsFolder(); err != nil {
		return err
	}
	return nil
}

func (e *TestSuiteGroup) deleteReports(root string) error {
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

func (e *TestSuiteGroup) syncSubdirectoriesFromFeaturesFolderToReportsFolder() error {
	dirs, err := e.getSubdirectoriesNameFromFeaturesFolder()
	if err != nil {
		return err
	}
	err = e.deleteUnknownSubdirectoriesInReportsFolder(dirs)
	if err != nil {
		return err
	}
	return e.createMissingSubdirectoriesInReportsFolder(dirs)
}

func (e *TestSuiteGroup) getSubdirectoriesNameFromFeaturesFolder() (map[string]struct{}, error) {
	dirs := make(map[string]struct{})
	return dirs, filepath.WalkDir(e.featFolderRootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if path == e.featFolderRootPath {
			return nil
		}
		relPath, err := filepath.Rel(e.featFolderRootPath, path)
		if err != nil {
			return err
		}
		dirs[relPath] = struct{}{}
		return nil
	})
}

func (e *TestSuiteGroup) deleteUnknownSubdirectoriesInReportsFolder(dirs map[string]struct{}) error {
	return filepath.WalkDir(e.reportsFolderRootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return os.Remove(path)
		}
		if _, ok := dirs[d.Name()]; !ok && path != e.reportsFolderRootPath {
			return os.RemoveAll(path)
		}
		return nil
	})
}

func (e *TestSuiteGroup) createMissingSubdirectoriesInReportsFolder(dirs map[string]struct{}) error {
	for relPath := range dirs {
		if err := os.MkdirAll(filepath.Join(e.reportsFolderRootPath, relPath), 0755); err != nil {
			return err
		}
	}
	return nil
}

func (e *TestSuiteGroup) After() error {
	if err := e.generateSingleHTMLReport(); err != nil {
		return err
	}
	if err := e.printReportFileLink(); err != nil {
		return err
	}
	return nil
}

func (e *TestSuiteGroup) generateSingleHTMLReport() error {
	HTML := gocure.HTML{
		Config: html.Data{
			Title:            "report",
			MergeFiles:       true,
			InputFolderPath:  e.reportsFolderRootPath,
			OutputHtmlFolder: e.htmlReportFolderRootPath,
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

func (e *TestSuiteGroup) printReportFileLink() error {
	filePath, err := e.findLatestHTMLReportName(e.htmlReportFolderRootPath)
	if err != nil {
		if errors.Is(err, errReportNotFound) {
			fmt.Println("No Reports, Please run e2e tests first!")
		}
		return err
	}
	content := fmt.Sprintf("Full Report - file://%s", filePath)
	_, err = fmt.Println(e.formatted(content))
	return err
}

func (e *TestSuiteGroup) findLatestHTMLReportName(path string) (string, error) {
	var name string
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && name < d.Name() {
			name = d.Name()
			return nil
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if name == "" {
		return "", fmt.Errorf("%w : in %s", errReportNotFound, path)
	}
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, path, name), nil
}

func (e *TestSuiteGroup) formatted(content string) string {
	width := len(content) + 2
	symbol := "+"
	var b bytes.Buffer
	_, _ = b.WriteString("\n")
	for i := 0; i < width; i++ {
		_, _ = b.WriteString(symbol)
	}
	_, _ = b.WriteString("\n")
	_, _ = b.WriteString(" " + content)
	_, _ = b.WriteString("\n")
	for i := 0; i < width; i++ {
		_, _ = b.WriteString(symbol)
	}
	_, _ = b.WriteString("\n")
	return b.String()
}

// TODO using Option pattern to control it with Flag - openWithBrowser
func (e *TestSuiteGroup) openReportWithDefaultBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func (e *TestSuiteGroup) AddTestSuite(suite *godog.TestSuite) error {
	_, err := os.Stat(suite.Name)
	if err != nil {
		return fmt.Errorf("godog.TestSuite's name is wrong: %s", suite.Name)
	}
	e.suites[e.readableName(suite.Name)] = suite
	return e.setOptions(suite)
}

func (e *TestSuiteGroup) readableName(path string) string {
	return strings.Replace(strings.TrimPrefix(path, e.featFolderRootPath), ".", "_", 1)
}

func (e *TestSuiteGroup) setOptions(suite *godog.TestSuite) error {
	options := *e.customizedOptions
	if options.Randomize == 0 {
		options.Randomize = -1
	}
	options.Paths = []string{suite.Name}
	format, err := e.rewriteCucumberReportPath(&options)
	if err != nil {
		return err
	}
	options.Format = format
	suite.Options = &options
	return nil
}

func (e *TestSuiteGroup) rewriteCucumberReportPath(options *godog.Options) (string, error) {
	var newFormats []string
	for _, f := range strings.Split(options.Format, ",") {
		if strings.HasPrefix(f, "cucumber:") {
			relativePath, err := filepath.Rel(e.featFolderRootPath, strings.TrimSuffix(options.Paths[0], ".feature"))
			if err != nil {
				return "", err
			}
			f = fmt.Sprintf("cucumber:%s.json", filepath.Join(e.reportsFolderRootPath, relativePath))
		}
		newFormats = append(newFormats, f)
	}
	return strings.Join(newFormats, ","), nil
}

func (e *TestSuiteGroup) TestSuites() map[string]*godog.TestSuite {
	suites := make(map[string]*godog.TestSuite)
	for _, path := range e.customizedOptions.Paths {
		name := e.readableName(path)
		if suite, ok := e.suites[name]; ok {
			suites[name] = suite
		}
	}
	if len(e.customizedOptions.Paths) == 0 && len(suites) == 0 {
		suites = e.suites
	}
	return suites
}
