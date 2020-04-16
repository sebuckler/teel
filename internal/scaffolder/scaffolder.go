package scaffolder

import (
	"errors"
	"github.com/sebuckler/teel/internal/logger"
	"os"
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

func (s *scaffolder) Scaffold(d string, _ string) error {
	if _, statErr := os.Stat(d); !os.IsNotExist(statErr) {
		return errors.New("directory already exists at " + d)
	}

	return os.MkdirAll(d, 0777)
}
