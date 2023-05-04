package client

import (
	"golang.org/x/crypto/ssh"
	"package-manager/internal/app"
)

type adapter struct {
	client *ssh.Client
}

func New() (app.Client, error) {
	config := &ssh.ClientConfig{
		User: "adachi",
		Auth: []ssh.AuthMethod{
			ssh.Password("inadequate"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", "localhost:2222", config)
	if err != nil {
		return nil, err
	}

	return &adapter{client: client}, nil
}

func (a *adapter) Close() {
	a.client.Close()
}
