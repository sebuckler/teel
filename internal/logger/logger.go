package logger

import (
	"bufio"
	"log"
	"os"
)

type writeMode string

const (
	BUFFERED  writeMode = "BUFFERED"
	IMMEDIATE writeMode = "IMMEDIATE"
)

type Logger interface {
	Error(v ...interface{})
	Errorf(f string, v ...interface{})
	Log(v ...interface{})
	Logf(f string, v ...interface{})
	Write()
}

type logger struct {
	flushErr func() error
	flushOut func() error
	stderr   *log.Logger
	stdout   *log.Logger
	mode     writeMode
}

func NewFile(f *os.File, m writeMode) Logger {
	fileWriter := bufio.NewWriter(f)

	return newLogger(fileWriter, nil, m)
}

func NewStandard(m writeMode) Logger {
	stdoutWriter := bufio.NewWriter(os.Stdout)
	stderrWriter := bufio.NewWriter(os.Stderr)

	return newLogger(stdoutWriter, stderrWriter, m)
}

func (l *logger) Error(v ...interface{}) {
	l.stderr.Println(v...)
	l.tryFlush()
}

func (l *logger) Errorf(f string, v ...interface{}) {
	l.stderr.Printf(f, v...)
	l.tryFlush()
}

func (l *logger) Log(v ...interface{}) {
	l.stdout.Println(v...)
	l.tryFlush()
}

func (l *logger) Logf(f string, v ...interface{}) {
	l.stdout.Printf(f, v...)
	l.tryFlush()
}

func (l *logger) Write() {
	l.flush()
}

func (l *logger) flush() {
	_ = l.flushErr()
	_ = l.flushOut()
}

func (l *logger) tryFlush() {
	if l.mode == IMMEDIATE {
		l.flush()
	}
}

func newLogger(o *bufio.Writer, e *bufio.Writer, m writeMode) Logger {
	if e == nil {
		e = o
	}

	stderrLogger := log.New(e, "", log.LstdFlags)
	stdoutLogger := log.New(o, "", log.LstdFlags)

	return &logger{
		flushErr: e.Flush,
		flushOut: o.Flush,
		stderr:   stderrLogger,
		stdout:   stdoutLogger,
		mode:     m,
	}
}
