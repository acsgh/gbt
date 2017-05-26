package file

import (
	"os"
	"archive/zip"
	"path/filepath"
	"io"
	"compress/gzip"
	"archive/tar"
	"github.com/albertoteloko/gbt/log"
	"errors"
	"fmt"
)

func Unzip(zipFile, baseFolder string) error {
	return log.LogTimeWithError("Extracted: "+zipFile, func() error {
		reader, err := zip.OpenReader(zipFile)
		if err != nil {
			return err
		}
		defer reader.Close()

		os.MkdirAll(baseFolder, 0755)

		for _, file := range reader.File {
			err := extractZipEntry(baseFolder, file)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func UnTarGz(zipFileUrl, baseFolder string) error {
	return log.LogTimeWithError("Extracted: "+zipFileUrl, func() error {
		zipFile, err := os.Open(zipFileUrl)

		if err != nil {
			return err
		}
		defer zipFile.Close()

		gzipReader, err := gzip.NewReader(zipFile)
		if err != nil {
			return err
		}

		os.MkdirAll(baseFolder, 0755)

		tarReader := tar.NewReader(gzipReader)

		for {
			header, err := tarReader.Next()

			if err == io.EOF {
				break
			}

			if err != nil {
				return err
			}

			err = extractTarEntry(baseFolder, header, tarReader)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func extractTarEntry(baseFolder string, header *tar.Header, tarReader *tar.Reader) error {
	path := filepath.Join(baseFolder, header.Name)

	log.Debugf("Extracting %v", path)

	switch header.Typeflag {
	case tar.TypeReg:
		return extractTarFile(path, header, tarReader)
	case tar.TypeDir:
		return extractTarDir(path, header)
	default:
		return errors.New(fmt.Sprintf("Unable to know the type %v of entry: %v", header.Typeflag, header.Name))
	}
}

func extractTarFile(path string, header *tar.Header, tarReader *tar.Reader) error {
	err := os.MkdirAll(filepath.Dir(path), os.FileMode(header.Mode))

	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))

	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, tarReader)
	return err
}

func extractTarDir(path string, header *tar.Header) error {
	return os.MkdirAll(path, os.FileMode(header.Mode))
}

func extractZipEntry(baseFolder string, zipEntry *zip.File) error {
	path := filepath.Join(baseFolder, zipEntry.Name)

	reader, err := zipEntry.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	if zipEntry.FileInfo().IsDir() {
		return extractZipDir(path, zipEntry)
	} else {
		return extractZipFile(path, zipEntry, reader)
	}
}

func extractZipDir(path string, zipEntry *zip.File) error {
	return os.MkdirAll(path, zipEntry.Mode())
}

func extractZipFile(path string, zipEntry *zip.File, reader io.ReadCloser) error {
	err := os.MkdirAll(filepath.Dir(path), zipEntry.Mode())

	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipEntry.Mode())

	if err == nil {
		defer f.Close()
		_, err = io.Copy(f, reader)
	}

	return err
}
