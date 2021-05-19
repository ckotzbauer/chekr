package internal

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var excludeDirs []string

func ListFiles(dir string) ([]FileGroup, error) {
	excludeDirs = []string{"testdata/"}
	var groups []FileGroup
	var files []string

	err := filepath.WalkDir(dir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if contains(path) {
			return nil
		}

		if strings.Contains(path, "types.go") || strings.Contains(path, "register.go") {
			files = append(files, strings.TrimPrefix(path, dir+"/"))
		}

		return nil
	})

	for i := 0; i < len(files); i++ {
		if i%2 == 0 {
			groups = append(groups, FileGroup{registerFile: dir + "/" + files[i], typesFile: dir + "/" + files[i+1]})
		}
	}

	return groups, err
}

func ReadFile(path string) (string, error) {
	buf, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func WriteFile(path, content string) error {
	return ioutil.WriteFile(path, []byte(content), fs.ModeAppend)
}

func contains(path string) bool {
	for _, v := range excludeDirs {
		if strings.Contains(path, v) {
			return true
		}
	}

	return false
}
