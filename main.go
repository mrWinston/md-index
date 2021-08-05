package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func findMarkdownFilesInDirectory(dir string) ([]string, error) {
	markdownFiles := []string{}
	dircontent, err := ioutil.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	for _, file := range dircontent {
		fullPath := filepath.Join(dir, file.Name())
		if filepath.Ext(fullPath) == ".md" {
			markdownFiles = append(markdownFiles, fullPath)
		}
	}
	return markdownFiles, nil
}

func main() {
	fmt.Println("hello World")
}
