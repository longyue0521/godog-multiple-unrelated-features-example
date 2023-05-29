package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	http.HandleFunc("/reports", func(w http.ResponseWriter, r *http.Request) {
		name, err := findLatestHTMLReportName("e2e/reports/html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if name == "" {
			w.Write([]byte("there is no reports, please run e2e test first and try again"))
			return
		}
		fmt.Println(name, err)
		http.ServeFile(w, r, name)
	})
	log.Println("listen on :8080")
	log.Println(http.ListenAndServe(":8080", nil))
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
	return filepath.Join(path, name), nil
}
