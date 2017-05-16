package file

import (
	"os"
	"net/http"
	"io"
	"errors"
)

func DownloadUrl(url string, output string) error {
	response, err := http.Get(url)
	if err == nil {
		defer response.Body.Close()
		if response.StatusCode == 200 {
			file, err := os.Create(output)
			defer file.Close()

			if err == nil {
				_, err = io.Copy(file, response.Body)
			}
		} else if response.StatusCode == 404 {
			err = errors.New(response.Status)
		}
	}

	return err
}
