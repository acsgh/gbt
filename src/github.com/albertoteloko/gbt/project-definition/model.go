package project_definition

import "github.com/albertoteloko/gbt/validation"

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
	Load() (ProjectDefinition, error)
}

type ProjectDefinitionLoaderFileSystem struct {
}


func (definition ProjectDefinition) Validate() []error {
	var validationErrors []error = []error{}

	validationErrors = append(validationErrors, validation.ValidateString("Name", definition.Name)...)
	validationErrors = append(validationErrors, validation.ValidateString("Version", definition.Version)...)
	validationErrors = append(validationErrors, validation.ValidateString("Go", definition.GoVersion)...)
	validationErrors = append(validationErrors, validation.ValidateNil("Dependencies", definition.Dependencies)...)
	validationErrors = append(validationErrors, validation.ValidateNil("Dependencies", definition.Dependencies.Dev)...)
	validationErrors = append(validationErrors, validation.ValidateNil("Dependencies", definition.Dependencies.Prod)...)

	return validationErrors
}
