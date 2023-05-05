package client

import (
	"github.com/pkg/errors"
	"sort"
	"strings"
)

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
