package cli

import (
	"flag"
	"github.com/pkg/errors"
	"package-manager/internal/app"
)

type adapter struct {
	service app.Service
}

func New(service app.Service) app.CLI {
	return &adapter{
		service: service,
	}
}

func (a *adapter) Listen() error {
	flag.Parse()
	command := flag.Arg(0)
	path := flag.Arg(1)
	if path == "" {
		return errors.New("empty second argument")
	}

	switch command {
	case "create":
		err := a.service.Create(path)
		if err != nil {
			return errors.Wrap(err, "create failed")
		}
	case "update":
		err := a.service.Update(path)
		if err != nil {
			return errors.Wrap(err, "update failed")
		}
	default:
		return errors.New("empty/wrong command")
	}
	return nil
}
