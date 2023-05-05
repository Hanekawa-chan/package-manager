package client

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
	"strings"
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

	if err = session.Start(createDirCmd(name, ver) + " && " + sendFileCmd(name, ver)); err != nil {
		return err
	}

	_, err = io.Copy(pipe, file)
	if err != nil {
		return err
	}

	err = pipe.Close()
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) ReceivePackage(name, ver string) (*bytes.Buffer, error) {
	cmd := getPackageCmd(name, ver)

	buf, err := a.startCommand(cmd)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (a *adapter) ReceivePackageByVersionPattern(name, ver string) (*bytes.Buffer, string, error) {
	cmd := findFilesCmd(name)

	buf, err := a.startCommand(cmd)
	if err != nil {
		return nil, "", err
	}

	versions := buf.String()
	versionsSlice := strings.Split(versions, "\n")

	if versionsSlice[len(versionsSlice)-1] == "" {
		versionsSlice = versionsSlice[:len(versionsSlice)-1]
	}

	latest, err := findLatestByPattern(versionsSlice, ver)
	if err != nil {
		return nil, "", err
	}

	cmd = getPackageCmd(name, latest)

	buf, err = a.startCommand(cmd)
	if err != nil {
		return nil, "", err
	}

	return buf, latest, nil
}

func (a *adapter) ReceivePackageLatest(name string) (*bytes.Buffer, string, error) {
	cmd := findFilesCmd(name)

	buf, err := a.startCommand(cmd)
	if err != nil {
		return nil, "", err
	}

	versions := buf.String()
	versionsSlice := strings.Split(versions, "\n")

	if versionsSlice[len(versionsSlice)-1] == "" {
		versionsSlice = versionsSlice[:len(versionsSlice)-1]
	}

	latest := findLatest(versionsSlice)

	cmd = getPackageCmd(name, latest)

	buf, err = a.startCommand(cmd)
	if err != nil {
		return nil, "", err
	}

	return buf, latest, nil
}

func (a *adapter) startCommand(cmd string) (*bytes.Buffer, error) {
	session, err := a.client.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "new session")
	}
	defer session.Close()

	pipe, err := session.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "get stdout pipe")
	}

	buf := new(bytes.Buffer)

	if err := session.Start(cmd); err != nil {
		return nil, errors.Wrap(err, "start command: "+cmd)
	}

	_, err = io.Copy(buf, pipe)
	if err != nil {
		return nil, errors.Wrap(err, "copy output from pipe")
	}

	if err = session.Wait(); err != nil {
		return nil, errors.Wrap(err, "wait session")
	}

	return buf, nil
}
