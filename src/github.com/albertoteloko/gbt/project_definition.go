package main

import (
	"encoding/json"
	"io/ioutil"
	"github.com/albertoteloko/gbt/validation"
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

type ProjectDefinitionLoader interface {
	load() (ProjectDefinition, error)
}

type ProjectDefinitionLoaderFileSystem struct {
}

func (loader ProjectDefinitionLoaderFileSystem) load() (ProjectDefinition, error) {
	var definition ProjectDefinition
	var err error

	file, err := ioutil.ReadFile(PROJECT_FILE)
	if err == nil {
		err = json.Unmarshal(file, &definition)

	}
	return definition, err
}

func (definition ProjectDefinition) validate() []error {
	var validationErrors []error = []error{}

	validationErrors = append(validationErrors, validation.ValidateString("Name", definition.Name)...)
	validationErrors = append(validationErrors, validation.ValidateString("Version", definition.Version)...)
	validationErrors = append(validationErrors, validation.ValidateString("Go", definition.GoVersion)...)
	validationErrors = append(validationErrors, validation.ValidateNil("Dependencies", definition.Dependencies)...)
	validationErrors = append(validationErrors, validation.ValidateNil("Dependencies", definition.Dependencies.Dev)...)
	validationErrors = append(validationErrors, validation.ValidateNil("Dependencies", definition.Dependencies.Prod)...)

	return validationErrors
}
