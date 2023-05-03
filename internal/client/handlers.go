package client

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"io"
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

func (a *adapter) ReceivePackage(name, ver string) (*bytes.Buffer, error) {
	session, err := a.client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	pipe, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := session.Start("cat " + fileName(name, ver)); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, pipe)
	if err != nil {
		return nil, err
	}

	if err := session.Wait(); err != nil {
		return nil, err
	}

	return buf, nil
}

func fileName(name, ver string) string {
	return name + "-" + ver + ".zip"
}
