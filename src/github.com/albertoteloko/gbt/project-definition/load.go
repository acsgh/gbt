package project_definition

import (
	"encoding/json"
	"io/ioutil"
)

const PROJECT_FILE = "gbt.json"

func (loader ProjectDefinitionLoaderFileSystem) Load() (ProjectDefinition, error) {
	var definition ProjectDefinition
	var err error

	file, err := ioutil.ReadFile(PROJECT_FILE)
	if err == nil {
		err = json.Unmarshal(file, &definition)

	}
	return definition, err
}
