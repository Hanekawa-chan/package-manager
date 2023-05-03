package app

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
)

func archiveFiles(filePaths []string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)

	w := zip.NewWriter(buf)
	defer w.Close()

	for _, path := range filePaths {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		fileContent, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		f, err := w.Create(file.Name())
		if err != nil {
			return nil, err
		}
		_, err = f.Write(fileContent)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}
