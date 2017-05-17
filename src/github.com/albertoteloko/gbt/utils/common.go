package utils

import (
	"path/filepath"
	"strings"
)

func ReplaceChars(input string, replaceBy string, values ...string) string {
	result := input

	for _, value := range values {
		result = strings.Replace(result, value, replaceBy, -1)
	}

	return result
}

func Max(v1, v2 int) int {
	if v1 >= v2 {
		return v1
	} else {
		return v2
	}
}

func FixWidth(value string, width int) string {
	result := value

	for len(result) < width {
		result += " "
	}

	return result
}


func GetBaseDir() string {
	path, _ := filepath.Abs(filepath.Dir("."))
	return path
}
