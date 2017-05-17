package file

import (
	"os"
	"net/http"
	"io"
	"errors"
	"github.com/albertoteloko/gbt/log"
)

func DownloadUrl(url string, output string) error {
	return log.LogTimeWithError("Downloaded: "+url, func() error {
		response, err := http.Get(url)

		if err == nil {
			defer response.Body.Close()
			if response.StatusCode == 200 {
				file, err := os.Create(output)
				defer file.Close()

				reader := &progressReader{Reader: response.Body, length: response.ContentLength}

				if err == nil {
					_, err = io.Copy(file, reader)
					log.Debug("Downloaded %v/%v (%.2f %%)", response.ContentLength, response.ContentLength, float64(100))
				}
			} else if response.StatusCode == 404 {
				err = errors.New(response.Status)
			}
		}
		return err
	})
}

type progressReader struct {
	io.Reader
	total              int64
	length             int64
	lastPercentageShow float64
}

func (pt *progressReader) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	if n > 0 {
		pt.total += int64(n)
		percentage := float64(pt.total) / float64(pt.length) * float64(100)

		if percentage-pt.lastPercentageShow >= 5 {
			log.Debug("Downloaded %v/%v (%.2f %%)", pt.total, pt.length, percentage)
			pt.lastPercentageShow = percentage
		}
	}

	return n, err
}
