package app

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"strings"
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

func unArchiveFiles(prefix string, buf *bytes.Buffer) error {
	reader := bytes.NewReader(buf.Bytes())
	r, err := zip.NewReader(reader, reader.Size())
	if err != nil {
		return err
	}

	for _, path := range r.File {
		path.Name = strings.ReplaceAll(path.Name, "\\", "/")
		lastSlashIndex := strings.LastIndex(path.Name, "/")

		err = os.RemoveAll(prefix)
		if err != nil {
			return err
		}

		err = os.MkdirAll(prefix+path.Name[:lastSlashIndex], os.ModeAppend)
		if err != nil {
			return err
		}

		file, err := os.OpenFile(prefix+path.Name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}

		fileFromArchive, err := path.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(file, fileFromArchive)
		if err != nil {
			return err
		}
	}

	return nil
}
