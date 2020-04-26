package sighandler

import (
	"os"
	"os/signal"
)

type Handler interface {
	Handle(f handleFunc)
}

type handler struct {
	signals []os.Signal
}

type handleFunc func(os.Signal)

func New(s ...os.Signal) Handler {
	return &handler{
		signals: s,
	}
}

func (h *handler) Handle(f handleFunc) {
	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, h.signals...)

	go func() { f(<-signalChannel) }()
}
