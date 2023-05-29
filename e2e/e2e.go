package e2e

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
	"gitlab.com/rodrigoodhin/gocure/models"
	"gitlab.com/rodrigoodhin/gocure/pkg/gocure"
	"gitlab.com/rodrigoodhin/gocure/report/html"
)

type TestSuiteGroup struct {
	featFolderRootPath       string
	reportsFolderRootPath    string
	htmlReportFolderRootPath string
	suites                   map[string]*godog.TestSuite
}

func NewTestSuiteGroup(featFolderRootPath, reportsFolderRootPath, htmlReportFolderRootPath string) *TestSuiteGroup {
	return &TestSuiteGroup{
		featFolderRootPath:       featFolderRootPath,
		reportsFolderRootPath:    reportsFolderRootPath,
		htmlReportFolderRootPath: htmlReportFolderRootPath,
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
	return e.generateSingleHTMLReport()
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

func (e *TestSuiteGroup) AddTestSuite(suite *godog.TestSuite) error {
	_, err := os.Stat(suite.Name)
	if err != nil {
		return fmt.Errorf("godog.TestSuite's name is wrong: %s", suite.Name)
	}
	rel, err := filepath.Rel(e.featFolderRootPath, suite.Name)
	if err != nil {
		return err
	}
	name := strings.TrimSuffix(rel, ".feature")
	e.suites[name] = suite
	return nil
}

func (e *TestSuiteGroup) TestSuites() map[string]*godog.TestSuite {
	// TODO use values of --godog.Paths to filter
	return e.suites
}
