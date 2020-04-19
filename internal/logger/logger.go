package logger

import (
	"io"
	"log"
)

type Logger interface {
	Error(v ...interface{})
	Errorf(f string, v ...interface{})
	Log(v ...interface{})
	Logf(f string, v ...interface{})
}

type logger struct {
	errLogger *log.Logger
	outLogger *log.Logger
}

func New(o io.Writer, e io.Writer) Logger {
	if e == nil {
		e = o
	}

	return &logger{
		errLogger: log.New(e, "Error: ", log.LstdFlags),
		outLogger: log.New(o, "Log: ", log.LstdFlags),
	}
}

func (l *logger) Error(v ...interface{}) {
	l.errLogger.Println(v...)
}

func (l *logger) Errorf(f string, v ...interface{}) {
	l.errLogger.Printf(f, v...)
}

func (l *logger) Log(v ...interface{}) {
	l.outLogger.Println(v...)
}

func (l *logger) Logf(f string, v ...interface{}) {
	l.outLogger.Printf(f, v...)
}
