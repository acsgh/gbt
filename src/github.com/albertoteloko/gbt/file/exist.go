package file

import (
	"os"
)

func ExistsFolder(path string) (bool, error) {
	return exists(path, IsDirectory)
}

func ExistsFile(path string) (bool, error) {
	return exists(path, IsFile)
}


func exists(path string, predicate func(info os.FileInfo) bool) (bool, error) {
	info, err := os.Stat(path)

	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	return (info != nil) && predicate(info), nil
}
