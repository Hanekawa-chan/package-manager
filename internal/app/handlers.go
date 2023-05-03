package app

import (
	"encoding/json"
	"github.com/pkg/errors"
)

func (s *service) Create(path string) error {
	fileContents, err := readFileContents(path)
	if err != nil {
		return errors.Wrap(err, "read file contents")
	}

	var packet FullPacket
	err = json.Unmarshal(fileContents, &packet)
	if err != nil {
		return errors.Wrap(err, "json unmarshal")
	}

	files := make([]string, 0)

	for _, target := range packet.Targets {
		newFiles, err := findPaths(target.Path, target.Exclude)
		if err != nil {
			return errors.Wrap(err, "find paths to targets")
		}
		files = append(files, newFiles...)
	}

	buf, err := archiveFiles(files)
	if err != nil {
		return errors.Wrap(err, "archive files")
	}

	err = s.client.SendPackage(packet.Name, packet.Ver, buf)
	if err != nil {
		return errors.Wrap(err, "send package")
	}

	return nil
}

func (s *service) Update(path string) error {
	fileContents, err := readFileContents(path)
	if err != nil {
		return errors.Wrap(err, "read file contents")
	}

	var packages Package
	err = json.Unmarshal(fileContents, &packages)
	if err != nil {
		return errors.Wrap(err, "json unmarshal")
	}

	for _, pack := range packages.Packages {
		buf, err := s.client.ReceivePackage(pack.Name, pack.Ver)
		if err != nil {
			return errors.Wrap(err, "receive package")
		}

		err = unArchiveFiles(buf)
		if err != nil {
			return errors.Wrap(err, "unArchive package")
		}
	}

	return nil
}
