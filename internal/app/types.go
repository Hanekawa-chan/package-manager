package app

import (
	"encoding/json"
)

type FullPacket struct {
	Packet
	PackageData
}

type PackageData struct {
	Targets []Target `json:"targets"`
	Packets []Packet `json:"packets"`
}

type Target struct {
	Path    string `json:"path"`
	Exclude string `json:"exclude,omitempty"`
}

type Packet struct {
	Name string `json:"name"`
	Ver  string `json:"ver"`
}

type Package struct {
	Packages []PackagePacket `json:"packages"`
}

type PackagePacket struct {
	Name string `json:"name"`
	Ver  string `json:"ver,omitempty"`
}

func (t *Target) UnmarshalJSON(data []byte) error {
	fakeTarget := struct {
		Path    string `json:"path"`
		Exclude string `json:"exclude,omitempty"`
	}{}

	err := json.Unmarshal(data, &fakeTarget)
	if err != nil {
		str := ""
		err := json.Unmarshal(data, &str)
		if err != nil {
			return err
		}

		fakeTarget.Path = str
	}

	t.Path = fakeTarget.Path
	t.Exclude = fakeTarget.Exclude

	return nil
}
