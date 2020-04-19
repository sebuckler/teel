package sighandler

import (
	"os"
	"os/signal"
)

type Handler interface {
	Handle(h func(s os.Signal), s ...os.Signal)
}

type handler struct{}

func New() Handler {
	return &handler{}
}

func (*handler) Handle(h func(os.Signal), s ...os.Signal) {
	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, s...)

	go func() { h(<-signalChannel) }()
}
