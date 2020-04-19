package directives

import (
	"os"
	"path"
)

type ConfigDirective struct {
}

func NewConfig() Directive {
	return &ConfigDirective{}
}

func (c *ConfigDirective) Execute(d string, n string) error {
	mkdirErr := os.MkdirAll(path.Join(d, "server"), 0755)

	if mkdirErr != nil {
		return mkdirErr
	}

	_, fileErr := os.OpenFile(path.Join(d, "config.json"), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)

	return fileErr
}
