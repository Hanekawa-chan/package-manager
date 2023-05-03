package app

import (
	"github.com/pkg/errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// readFileContents читает содержимое файла и возвращает его содержимое в виде массива байтов
func readFileContents(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "read bytes of file")
	}

	return bytes, nil
}

// findPaths ищет пути до файлов, указанных в targets.path и не включает файлы указанные в targets.exclude
func findPaths(path, excludePattern string) ([]string, error) {
	var paths []string
	path, pattern := separatePathAndPattern(path)

	err := filepath.WalkDir(path, func(pathToAdd string, entry fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if ok, err := filepath.Match(pattern, entry.Name()); ok && pathToAdd != path {
			if ok, err := filepath.Match(excludePattern, entry.Name()); !ok {
				paths = append(paths, pathToAdd)
			} else if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return paths, nil
}

// separatePathAndPattern возвращает папку, в которой будет вестись поиск, в виде переменной path и паттерн в виде переменной pattern
func separatePathAndPattern(str string) (path, pattern string) {
	if strings.ContainsAny(str, "*?[]") {
		lastPathSeparator := strings.LastIndex(str, "/") + 1
		path = str[:lastPathSeparator]
		pattern = str[lastPathSeparator:]
	} else {
		path = str
	}
	return
}
