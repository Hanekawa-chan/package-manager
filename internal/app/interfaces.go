package app

import "bytes"

type Service interface {
	Create(path string) error
	Update(path string) error
}

type CLI interface {
	Listen() error
}

type Client interface {
	SendPackage(name, ver string, file *bytes.Buffer) error
	ReceivePackage(name, ver string) (*bytes.Buffer, error)
	ReceivePackageByVersionPattern(name, ver string) (buf *bytes.Buffer, version string, err error)
	ReceivePackageLatest(name string) (buf *bytes.Buffer, version string, err error)
	Close()
}
