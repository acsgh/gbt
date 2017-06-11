package tasks

import (
	pd "github.com/albertoteloko/gbt/project-definition"
	"github.com/albertoteloko/gbt/file"
	"os/exec"
	"github.com/albertoteloko/gbt/log"
)

func format(pd pd.ProjectDefinition) error {
	folders, err := file.GetFolders(file.SRC_PATH, true, file.IsGoFolder)

	if err == nil {
		for _, folder := range folders {
			formatFolder(folder)
		}
	}

	return err
}

func formatFolder(dir string) {
	files, err := file.GetFiles(dir, file.And(file.IsGoFile, file.IsGitHubPath))

	if err != nil {
		log.Errorf("Error during folder read: %s", err)
	}

	for _, f := range files {
		log.Debugf("Formatting file: %s", f)
		if err := formatFile(f); err != nil {
			log.Errorf("Error formatting file %s, %s", f, err)
		} else {
			log.Debugf("Formatted file: %s", f)
		}
	}
}

func formatFile(file string) error {
	return exec.Command("gofmt", "-w", file).Run()
}
