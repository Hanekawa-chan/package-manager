package app

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"
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

	fmt.Println(files)

	dependencies := make([]Dependency, len(packet.Packets))

	for i, pack := range packet.Packets {
		//TODO
		dependencies[i] = Dependency{
			Name: pack.Name,
			Data: nil,
		}
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
		name := "packages/" + pack.Name
		if pack.Ver == "" {
			buf, version, err := s.client.ReceivePackageLatest(pack.Name)
			if err != nil {
				return errors.Wrap(err, "receive latest package")
			}

			err = unArchiveFiles(name+"/"+version+"/", buf)
			if err != nil {
				return errors.Wrap(err, "unArchive latest package")
			}
		} else if strings.ContainsAny(pack.Ver, "><=") {
			buf, version, err := s.client.ReceivePackageByVersionPattern(pack.Name, pack.Ver)
			if err != nil {
				return errors.Wrap(err, "receive version patterned package")
			}

			err = unArchiveFiles(name+"/"+version+"/", buf)
			if err != nil {
				return errors.Wrap(err, "unArchive version patterned package")
			}
		} else {
			buf, err := s.client.ReceivePackage(pack.Name, pack.Ver)
			if err != nil {
				return errors.Wrap(err, "receive package")
			}

			err = unArchiveFiles(name+"/"+pack.Ver+"/", buf)
			if err != nil {
				return errors.Wrap(err, "unArchive package")
			}
		}
	}

	return nil
}
