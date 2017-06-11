package tasks

import (
	pd "github.com/albertoteloko/gbt/project-definition"
	"os"
	"github.com/albertoteloko/gbt/file"
)

func clean(pd pd.ProjectDefinition) error {
	err := os.RemoveAll(file.BUILD_PATH)

	if err != nil {
		err = os.MkdirAll(file.BUILD_PATH, 0777)
	}

	return err
}
