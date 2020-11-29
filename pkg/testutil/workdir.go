package testutil

import (
	"os"
	"path/filepath"
)

func SetRootWorkDir() {
	os.Chdir(findWorkDir())
}

func findWorkDir() string {
	fullpath, _ := os.Getwd()
	for !fileExists(filepath.Join(fullpath, "go.mod")) && fullpath != "/" {
		fullpath = filepath.Dir(fullpath)
	}
	return fullpath
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
