//go:build docs

package main

import (
	"bytes"
	_ "embed"
	"flag"
	"github.com/rad12000/list-agent/cmd"
	"github.com/spf13/cobra/doc"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var (
	//go:embed md.html
	mdHTMLTemplate string
)

func main() {
	docsDir := flag.String("dir", "./docs", "directory to write docs to")
	if err := doc.GenMarkdownTree(cmd.ListAgentCmd, *docsDir); err != nil {
		panic(err)
	}

	tmpl, err := template.New("md-template").Parse(mdHTMLTemplate)
	if err != nil {
		panic(err)
	}

	dirEntries, err := os.ReadDir(*docsDir)
	if err != nil {
		panic(err)
	}

	for _, dirEntry := range dirEntries {
		fileInfo, err := dirEntry.Info()
		if err != nil {
			panic(err)
		}

		if fileInfo.IsDir() {
			continue
		}

		if !strings.HasSuffix(dirEntry.Name(), ".md") {
			continue
		}

		fullPath := filepath.Join(*docsDir, dirEntry.Name())
		file, err := os.OpenFile(fullPath, os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}

		readme, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}

		mdRegex := regexp.MustCompile(`(?m)(\[.*])\((.*)\.md\)`)
		readme = mdRegex.ReplaceAll(readme, []byte("${1}(${2}.html)"))
		if err := file.Truncate(0); err != nil {
			panic(err)
		}

		if _, err = file.Seek(0, io.SeekStart); err != nil {
			panic(err)
		}

		if err := tmpl.Execute(file, string(bytes.TrimSpace(readme))); err != nil {
			panic(err)
		}

		if err := file.Close(); err != nil {
			panic(err)
		}

		newName := strings.Replace(fullPath, ".md", ".html", 1)
		if err := os.Rename(fullPath, newName); err != nil {
			panic(err)
		}
	}
}
