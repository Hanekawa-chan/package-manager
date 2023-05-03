package client

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
)

func (a *adapter) SendPackage(name, ver string, file *bytes.Buffer) error {
	session, err := a.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	pipe, err := session.StdinPipe()
	if err != nil {
		return err
	}

	if err = session.Start("cat > " + name + "-" + ver + ".zip"); err != nil {
		return err
	}

	_, err = io.Copy(pipe, file)
	if err != nil {
		return err
	}

	err = session.Signal(ssh.SIGQUIT)
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) ReceivePackage(name, ver string) error {
	session, err := a.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	pipe, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	filename := fileName(name, ver)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := session.Start("cat " + filename); err != nil {
		return err
	}

	_, err = io.Copy(file, pipe)
	if err != nil {
		return err
	}

	if err := session.Wait(); err != nil {
		return err
	}

	return nil
}

func fileName(name, ver string) string {
	return name + "-" + ver + ".zip"
}
