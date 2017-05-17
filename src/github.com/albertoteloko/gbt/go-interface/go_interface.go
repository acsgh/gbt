package go_interface

import (
	"github.com/albertoteloko/gbt/log"
	"github.com/albertoteloko/gbt/file"
	pd "github.com/albertoteloko/gbt/project-definition"
	"errors"
)

const GO_URL_BASE = "https://storage.googleapis.com/golang"

type GoInterface interface {
	CheckAndDownloadGo(projectDefinition pd.ProjectDefinition) error
}

type GoInterfaceFileSystem struct {
}

func (goInterface GoInterfaceFileSystem) CheckAndDownloadGo(projectDefinition pd.ProjectDefinition) error {
	var goFolder = file.HOME_PATH + "/.gbt/go/" + projectDefinition.GoVersion
	log.Debug("Go Folder: %v", goFolder)

	var err = testGoInstallation(goFolder)

	if err == nil {
		log.Warn("Invalid GO installation, downloading")

		err = downloadGoInstallation(projectDefinition.GoVersion, goFolder)
		if err == nil {
			err = testGoInstallation(goFolder)
		}
	}

	return err
}

func downloadGoInstallation(goVersion string, goFolder string) error {
	goUrl := getGoUrl(goVersion)
	log.Info(goUrl)
	zipFile := file.TMP_PATH + "/" + getGoFile(goVersion)
	err := file.DownloadUrl(goUrl, zipFile)

	if err == nil{
		err = file.UnTarGz(zipFile, goFolder)
	}

	return err
}

func getGoUrl(goVersion string) string {
	return GO_URL_BASE + "/" + getGoFile(goVersion)
}
func getGoFile(goVersion string) string {
	return "go" + goVersion + ".linux-amd64.tar.gz"
}

func testGoInstallation(goFolder string) error {
	var err error

	var goExec = goFolder + "/bin/go"

	var exist, _ = file.ExistsFile(goExec)

	if (err == nil) && (exist) {
		err = errors.New("Not implemented yet")
	}

	return err
}
