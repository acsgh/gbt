package main

import (
	"path/filepath"
	"strings"
)

var recursive = true

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


func getBaseDir() string {
	path, _ := filepath.Abs(filepath.Dir("."))
	return path
}
