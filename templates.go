package main

import (
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var templates *template.Template

func init() {
	files := loadTemplates("templates")
	templates = template.Must(template.ParseFiles(files...))
}

func loadTemplates(templateDir string) []string {
	result := make([]string, 0)
	fileSystem := os.DirFS(templateDir)
	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".html") {
			result = append(result, filepath.Join(templateDir, path))
		}

		return nil
	})

	return result
}
