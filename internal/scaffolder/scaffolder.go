package scaffolder

import (
	"errors"
	"github.com/sebuckler/teel/internal/logger"
	"os"
	"path"
)

type Scaffolder interface {
	Scaffold(d string, n string) error
}

type scaffolder struct {
	logger logger.Logger
}

func New(l logger.Logger) Scaffolder {
	return &scaffolder{
		logger: l,
	}
}

func (s *scaffolder) Scaffold(d string, n string) error {
	if _, statErr := os.Stat(d); !os.IsNotExist(statErr) {
		return errors.New("directory already exists at " + d)
	}

	if n == "" {
		return errors.New("name must be provided")
	}

	mkdirErr := os.MkdirAll(d, 0755)

	if mkdirErr != nil {
		return mkdirErr
	}

	s.makeServer(d, n)

	return nil
}

func (s *scaffolder) makeServer(d string, n string) error {
	mkdirErr := os.Mkdir(path.Join(d, "server/"), 0755)

	if mkdirErr != nil {
		return mkdirErr
	}

	_, fileErr := os.OpenFile("config.json", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)

	return fileErr
}
