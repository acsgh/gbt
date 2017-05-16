package file

import (
	"io/ioutil"
)

func GetFiles(name string, filter func(string) bool) ([]string, error) {
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

func GetFolders(name string, recursive bool, filter func(string) bool) ([]string, error) {
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
				childFolders, err := GetFolders(fileName, recursive, filter)

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
