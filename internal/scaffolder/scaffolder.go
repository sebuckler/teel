package scaffolder

import (
	"errors"
	"github.com/sebuckler/teel/internal/logger"
	"github.com/sebuckler/teel/internal/scaffolder/directives"
	"os"
)

type Scaffolder interface {
	Scaffold(d string, n string) error
}

type scaffolder struct {
	directives []directives.Directive
	logger     logger.Logger
}

func New(l logger.Logger, d ...directives.Directive) Scaffolder {
	return &scaffolder{
		directives: d,
		logger:     l,
	}
}

func (s *scaffolder) Scaffold(d string, n string) error {
	if _, statErr := os.Stat(d); !os.IsNotExist(statErr) {
		return errors.New("directory already exists at " + d)
	}

	if n == "" {
		return errors.New("name must be provided")
	}

	for _, directive := range s.directives {
		execErr := directive.Execute(d, n)

		if execErr != nil {
			return execErr
		}
	}

	return nil
}
