package main

import "errors"

type GoInterface interface {
	checkAndDownloadGo(projectDefinition ProjectDefinition) error
}

type GoInterfaceFileSystem struct {
}

func (goInterface GoInterfaceFileSystem) checkAndDownloadGo(projectDefinition ProjectDefinition) error {
	return errors.New("Go version not found :(")
}
