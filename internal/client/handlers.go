package client

import (
	"bytes"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"io"
	"sort"
	"strings"
)

func (a *adapter) SendPackage(name, ver string, file *bytes.Buffer) error {
	session, err := a.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	err = session.Run("mkdir -p " + name)
	if err != nil {
		return err
	}

	session, err = a.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	err = session.Run("mkdir -p " + name + "/" + ver)
	if err != nil {
		return err
	}

	session, err = a.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	pipe, err := session.StdinPipe()
	if err != nil {
		return err
	}

	if err = session.Start("cat > " + name + "/" + ver + "/data.zip"); err != nil {
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
	cmd := getPackage(name, ver)

	buf, err := a.startCommand(cmd)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (a *adapter) ReceivePackageByVersionPattern(name, ver string) (*bytes.Buffer, string, error) {
	cmd := findFiles(name)

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

	cmd = getPackage(name, latest)

	buf, err = a.startCommand(cmd)
	if err != nil {
		return nil, "", err
	}

	return buf, latest, nil
}

func (a *adapter) ReceivePackageLatest(name string) (*bytes.Buffer, string, error) {
	cmd := findFiles(name)

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

	cmd = getPackage(name, latest)

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

func getPackage(name, ver string) string {
	return "cat " + name + "/" + ver + "/data.zip"
}

func findFiles(name string) string {
	return "find ./" + name + " -maxdepth 1 -name '*.*'"
}

func findLatest(versions []string) string {
	sort.Strings(versions)
	latest := versions[len(versions)-1]
	index := strings.LastIndex(latest, "/")

	return latest[index+1:]
}

func findLatestByPattern(versions []string, pattern string) (string, error) {
	sort.Strings(versions)
	index := strings.LastIndex(versions[0], "/")
	latest := ""

	if strings.HasPrefix(pattern, ">=") {
		for i := len(versions) - 1; i >= 0; i-- {
			if versions[i][index+1:] >= pattern[2:] {
				latest = versions[i]
				break
			}
		}

		if latest == "" {
			return "", errors.New("couldn't find version needed")
		}
	} else if strings.HasPrefix(pattern, ">") {
		for i := len(versions) - 1; i >= 0; i-- {
			if versions[i][index+1:] > pattern[1:] {
				latest = versions[i]
				break
			}
		}

		if latest == "" {
			return "", errors.New("couldn't find version needed")
		}
	} else if strings.HasPrefix(pattern, "<=") {
		for i := len(versions) - 1; i >= 0; i-- {
			if versions[i][index+1:] <= pattern[2:] {
				latest = versions[i]
				break
			}
		}

		if latest == "" {
			return "", errors.New("couldn't find version needed")
		}
	} else if strings.HasPrefix(pattern, "<") {
		for i := len(versions) - 1; i >= 0; i-- {
			if versions[i][index+1:] < pattern[1:] {
				latest = versions[i]
				break
			}
		}

		if latest == "" {
			return "", errors.New("couldn't find version needed")
		}
	}

	return latest[index+1:], nil
}
