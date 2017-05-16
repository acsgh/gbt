package file

import (
	"os"
	"github.com/albertoteloko/gbt/log"
	"strings"
	"io/ioutil"
)

func And(f1, f2 func(string) bool) func(string) bool {
	return func(file string) bool {
		return f1(file) && f2(file)
	}
}

func GetFolderName(baseName string, folder string) string {
	return strings.Replace(folder, baseName, "", -1)
}

func IsDirectory(info os.FileInfo) bool {
	return info.IsDir()
}

func IsFile(info os.FileInfo) bool {
	return !info.IsDir()
}

func IsGoFile(file string) bool {
	return strings.HasSuffix(file, ".go")
}

func IsGoTestFile(file string) bool {
	return strings.HasSuffix(file, "_test.go")
}

func IsGitHubPath(file string) bool {
	return strings.Contains(file, "github.com/albertoteloko")
}

func IsGoFolder(name string) bool {
	files, err := ioutil.ReadDir(name)

	if err != nil {
		return false
	}

	for _, f := range files {
		fileName := name + "/" + f.Name()
		if !f.IsDir() && IsGoFile(fileName) {
			return true
		}
	}
	return false
}

func IsGoMainFolder(name string) bool {
	files, err := ioutil.ReadDir(name)

	if err != nil {
		return false
	}

	for _, f := range files {
		fileName := name + "/" + f.Name()
		if !f.IsDir() && IsGoMainFile(fileName) {
			return true
		}
	}
	return false
}

func IsGoMainFile(name string) bool {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		log.Error("Unable to read file %s: %s", name, err)
		return false
	}

	return strings.HasPrefix(strings.Trim(string(b), " \n\r"), "package main")
}