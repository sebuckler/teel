package filestreamer

import (
	"bufio"
	"os"
	"path"
)

type Streamer interface {
	Stream(s streamFunc) error
}

type StreamFile struct {
	*os.File
	Cwd   string
	Flush func() error
}

type streamer struct {
	base BasePath
	path string
}

type streamFunc func(f *StreamFile)

type BasePath int

const (
	Cwd BasePath = iota
	HomeDir
	RootDir
)

func New(b BasePath, p ...string) Streamer {
	return &streamer{
		base: b,
		path: path.Join(p...),
	}
}

func (f *streamer) Stream(s streamFunc) error {
	var basePath, cwd, home string
	var err error
	cwd, err = os.Getwd()

	switch f.base {
	case Cwd:
		basePath = cwd
	case HomeDir:
		home, err = os.UserHomeDir()
		basePath = home
	case RootDir:
		basePath = ""
	}

	if err != nil {
		return err
	}

	filePath := path.Join(basePath, f.path)
	fileDir := path.Dir(filePath)
	_ = os.MkdirAll(fileDir, 0755)
	file, fileErr := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)

	if fileErr != nil {
		return fileErr
	}

	writer := bufio.NewWriter(file)

	defer func() {
		_ = writer.Flush()
		_ = file.Close()
	}()

	s(&StreamFile{
		File:  file,
		Cwd:   cwd,
		Flush: writer.Flush,
	})

	return nil
}
