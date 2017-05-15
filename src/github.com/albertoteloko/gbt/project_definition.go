package main

import (
	"encoding/json"
	"io/ioutil"
)

const PROJECT_FILE = "gbt.json"

type ProjectDefinition struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	GoVersion    string `json:"go"`
	Dependencies Dependencies `json:"dependencies"`
}

type Dependencies struct {
	Dev  []Dependency `json:"dev"`
	Prod []Dependency `json:"prod"`
}

type Dependency struct {
	Url        string `json:"url"`
	Transitive bool `json:"transitive"`
}

func loadProjectDefinition() (ProjectDefinition, error) {
	var definition ProjectDefinition
	var err error

	file, err := ioutil.ReadFile(PROJECT_FILE)
	if err == nil {
		err = json.Unmarshal(file, &definition)

	}
	return definition, err
}

func validateDefinition(definition ProjectDefinition) []error {
	var validationErrors []error = []error{}

	validationErrors = append(validationErrors, validateString("Name", definition.Name)...)
	validationErrors = append(validationErrors, validateString("Version", definition.Version)...)
	validationErrors = append(validationErrors, validateString("Go", definition.GoVersion)...)
	validationErrors = append(validationErrors, validateNil("Dependencies", definition.Dependencies)...)
	validationErrors = append(validationErrors, validateNil("Dependencies", definition.Dependencies.Dev)...)
	validationErrors = append(validationErrors, validateNil("Dependencies", definition.Dependencies.Prod)...)

	return validationErrors
}