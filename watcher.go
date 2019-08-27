package serviceexecutor

import (
	"context"
	"os"
	"os/signal"
	"time"
)

// SignalWatcher allows easily hooking into os signals to shutdown long running services.
type SignalWatcher struct {
	// ShutdownContext creates the context that is passed to Shutdown.  If nil, will use context.Background
	ShutdownContext func() context.Context
	// ShutdownTimeout is an optional timeout to attach to the Shutdown context.  If empty, no timeout will be used.
	ShutdownTimeout time.Duration
	// Service is the Service that should receive a Shutdown signal.
	Service Service
	// Signals are the OS signals we should wait and signal on
	Signals      []os.Signal
	signalNotify func(c chan<- os.Signal, sig ...os.Signal)
	ch           chan os.Signal
}

// Setup ensures the signal channel is created and registered with the signal notifier.
func (w *SignalWatcher) Setup() error {
	w.ch = make(chan os.Signal, 1)
	if w.signalNotify != nil {
		w.signalNotify(w.ch, w.Signals...)
	} else {
		signal.Notify(w.ch, w.Signals...)
	}
	return nil
}

// Run blocks until a signal is received on the os signal channel.
func (w *SignalWatcher) Run() error {
	_, ok := <-w.ch
	if ok {
		ctx := w.ShutdownContext()
		if w.ShutdownTimeout != 0 {
			var onDone func()
			ctx, onDone = context.WithTimeout(ctx, w.ShutdownTimeout)
			defer onDone()
		}
		return w.Service.Shutdown(ctx)
	}
	return nil
}

// Shutdown closes the run signal's channel
func (w *SignalWatcher) Shutdown(ctx context.Context) error {
	signal.Stop(w.ch)
	close(w.ch)
	return nil
}

var _ Setupable = &SignalWatcher{}