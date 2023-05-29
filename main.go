package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
)

var (
	reportPath = "e2e/reports/html"
)

func main() {
	http.HandleFunc("/reports", reportsHandler)
	port := ":8080"
	log.Printf("listen on %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func reportsHandler(w http.ResponseWriter, r *http.Request) {
	name, err := findLatestHTMLReportName(reportPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<b>There is no reports, please run e2e test first and try again</b>"))
		return
	}
	http.ServeFile(w, r, name)
}

func findLatestHTMLReportName(path string) (string, error) {
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
		return "", fmt.Errorf("no report found in %s", path)
	}
	return filepath.Join(path, name), nil
}
