package main

import (
	"github.com/albertoteloko/gbt/log"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

func replaceChars(input string, replaceBy string, values ...string) string {
	result := input

	for _, value := range values {
		result = strings.Replace(result, value, replaceBy, -1)
	}

	return result
}

func max(v1, v2 int) int {
	if v1 >= v2 {
		return v1
	} else {
		return v2
	}
}

func fixWidth(value string, width int) string {
	result := value

	for len(result) < width {
		result += " "
	}

	return result
}

func and(f1, f2 func(string) bool) func(string) bool {
	return func(file string) bool {
		return f1(file) && f2(file)
	}
}

func isGoFile(file string) bool {
	return strings.HasSuffix(file, ".go")
}

func isGoTestFile(file string) bool {
	return strings.HasSuffix(file, "_test.go")
}

func isGitHubPath(file string) bool {
	return strings.Contains(file, "github.com/albertoteloko")
}

func getFiles(name string, filter func(string) bool) ([]string, error) {
	result := []string{}

	files, err := ioutil.ReadDir(name)

	if err != nil {
		return result, err
	}

	for _, f := range files {
		fileName := name + "/" + f.Name()
		if filter(fileName) {
			result = append(result, fileName)
		}
	}

	return result, nil
}

func getFolders(name string, filter func(string) bool) ([]string, error) {
	result := []string{}

	if filter(name) {
		result = append(result, name)
	}

	files, err := ioutil.ReadDir(name)

	if err != nil {
		return result, err
	}

	for _, f := range files {
		fileName := name + "/" + f.Name()
		if f.IsDir() {
			if recursive {
				childFolders, err := getFolders(fileName, filter)

				if err != nil {
					return result, err
				}
				result = append(result, childFolders...)
			}
			if filter(fileName) {
				result = append(result, fileName)
			}
		}
	}

	return result, nil
}

func getFolderName(folder string) string {
	return strings.Replace(folder, GO_PATH+"/src/", "", -1)
}

func isGoFolder(name string) bool {
	files, err := ioutil.ReadDir(name)

	if err != nil {
		return false
	}

	for _, f := range files {
		fileName := name + "/" + f.Name()
		if !f.IsDir() && isGoFile(fileName) {
			return true
		}
	}
	return false
}

func isGoMainFolder(name string) bool {
	files, err := ioutil.ReadDir(name)

	if err != nil {
		return false
	}

	for _, f := range files {
		fileName := name + "/" + f.Name()
		if !f.IsDir() && isGoMainFile(fileName) {
			return true
		}
	}
	return false
}

func isGoMainFile(name string) bool {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		log.Error("Unable to read file %s: %s", name, err)
		return false
	}

	return strings.HasPrefix(strings.Trim(string(b), " \n\r"), "package main")
}

func logTime(taskName string, task func()) {
	startTime := time.Now().UnixNano()
	log.Debug("%s start", taskName)
	task()
	log.Info("%s in %d ms", taskName, (time.Now().UnixNano()-startTime)/1000000)
}

func getBaseDir() string {
	if directory != "" {
		return directory
	} else {
		path, _ := filepath.Abs(filepath.Dir("."))
		return path
	}
}
