package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const (
	TEMPLATE = `# {{ .FolderName }}

## Files

{{ range .MarkdownFiles -}}
- [{{ . }}](./{{ . }})
{{ end }}
## Subfolders

{{ range .SubDirs -}}
- [{{ . }}](./{{ . }}/index.md)
{{ end -}}


- [UP](../index.md)
`
)

type MarkdownDir struct {
	FolderName    string
	SubDirs       []string
	MarkdownFiles []string
}

func errPanic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetFilesInDirectory(dir string) ([]string, error) {
	files := []string{}
	dircontent, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range dircontent {
		files = append(files, file.Name())
	}
	return files, nil
}

func RenderInFolder(dir string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	errPanic(err)
	markdownDir := &MarkdownDir{
		FolderName:    filepath.Base(dir),
		MarkdownFiles: []string{},
		SubDirs:       []string{},
	}
	log.Printf("%+v", markdownDir)

	for _, file := range files {

		if file.IsDir() {
			markdownDir.SubDirs = append(markdownDir.SubDirs, file.Name())
		} else {
			markdownDir.MarkdownFiles = append(markdownDir.MarkdownFiles, file.Name())
		}
	}
	tmpl, err := template.New("test").Parse(TEMPLATE)
	errPanic(err)
	file, err := os.Create(filepath.Join(dir, "index.md"))
	errPanic(err)
	defer file.Close()
	err = tmpl.Execute(file, markdownDir)
	errPanic(err)
	errPanic(file.Sync())
	return files, nil
}

func getDirsFromFileInfos(infos []os.FileInfo) []os.FileInfo {
	res := []os.FileInfo{}
	for _, info := range infos {
		if info.IsDir() {
			res = append(res, info)
		}
	}
	return res
}

func main() {
	cwd, err := os.Getwd()
	errPanic(err)
	log.Printf("CWD is %s", cwd)

	dirsToSearch := []string{cwd}

	for len(dirsToSearch) > 0 {
		log.Printf("Cur Length: %d", len(dirsToSearch))
		curDir := dirsToSearch[0]
		log.Printf("Popped: %s", curDir)
		dirsToSearch = dirsToSearch[1:]

		fileInfos, err := RenderInFolder(curDir)
		errPanic(err)

		dirInfos := getDirsFromFileInfos(fileInfos)
		for _, dirInfo := range dirInfos {
			dirsToSearch = append(dirsToSearch, filepath.Join(curDir, dirInfo.Name()))
		}
		log.Printf("New Folders: %d", len(dirInfos))
	}
}
