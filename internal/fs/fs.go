package fs

import (
	"bufio"
	"io"
	"os"
)

type BufferedFileWriter interface {
	Close() error
	io.Writer
}

type bufferedFileWriter struct {
	file *os.File
	writer *bufio.Writer
}

func OpenBufferedFileWriter(p string) (*bufferedFileWriter, error) {
	file, fileErr := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)

	if fileErr != nil {
		return nil, fileErr
	}

	return &bufferedFileWriter{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (b *bufferedFileWriter) Close() error {
	if flushErr := b.writer.Flush(); flushErr != nil {
		return flushErr
	}

	return b.file.Close()
}

func (b *bufferedFileWriter) Write(p []byte) (int, error) {
	return b.writer.Write(p)
}
