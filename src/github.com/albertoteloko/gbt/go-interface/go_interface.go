package go_interface

import (
	"github.com/albertoteloko/gbt/log"
	"github.com/albertoteloko/gbt/file"
	"github.com/albertoteloko/gbt/utils"
	pd "github.com/albertoteloko/gbt/project-definition"
	"errors"
	"os"
	"fmt"
	"strings"
)

const GO_URL_BASE = "https://storage.googleapis.com/golang"

type GoInterface interface {
	CheckAndDownloadGo(projectDefinition pd.ProjectDefinition) error
}

type GoInterfaceFileSystem struct {
}

func (goInterface GoInterfaceFileSystem) CheckAndDownloadGo(projectDefinition pd.ProjectDefinition) error {
	var goFolder = file.HOME_PATH + "/.gbt/go/" + projectDefinition.GoVersion
	log.Debugf("Go Folder: %v", goFolder)

	var err = testGoInstallation(goFolder, projectDefinition.GoVersion)

	if err != nil {
		log.Warnf("Invalid GO installation, downloading")

		err = downloadGoInstallation(projectDefinition.GoVersion, goFolder)
		if err != nil {
			err = testGoInstallation(goFolder, projectDefinition.GoVersion)
		}
	}

	return err
}

func downloadGoInstallation(goVersion string, goFolder string) error {
	goUrl := getGoUrl(goVersion)
	zipFile := file.TMP_PATH + "/" + getGoFile(goVersion)
	err := file.DownloadUrl(goUrl, zipFile)

	if err == nil {
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

func testGoInstallation(goFolder string, goVersion string) error {
	var err error

	var goExec = goFolder + "/go/bin/go"

	exist, err := file.ExistsFile(goExec)

	if err != nil {
		return err
	}

	if !exist {
		return errors.New("Go executable not found")
	} else {
		env := os.Environ()
		env = append(env, fmt.Sprintf("GOROOT=%s", goFolder+"/go"))
		args := []string{"version"}
		output, _, err := file.ExecWithEnv(goExec, args, env)

		if err != nil {
			return err
		}

		output = utils.ReplaceChars(output, "", "\r", "\n", "go", "version", "   ", "  ")

		log.Debugf("Version: %v", output)
		log.Debugf("Go Version: %v", goVersion)

		if !strings.Contains(output, goVersion) {
			return errors.New("GO folder version does not match with the expected: " + goVersion)
		}

		return nil
	}
}
