package serviceexecutor

import (
	"context"
	"os"
	"os/signal"
	"time"
)

// ShutdownOnSignals is a helper that wraps a SiganlWatcher for your Multi service.  It will setup signal watcher
// that call Shutdown on toShutdown when any of signals is seen by the signal watcher.  The Shutdown context
// is timed out by shutdownTimeout.
func ShutdownOnSignals(toShutdown *Multi, shutdownTimeout time.Duration, signals ...os.Signal) {
	w := SignalWatcher{
		Service:         toShutdown,
		Signals:         signals,
		ShutdownTimeout: shutdownTimeout,
	}
	toShutdown.Services = append(toShutdown.Services, &w)
}

// SignalWatcher allows easily hooking into os signals to shutdown long running services.
type SignalWatcher struct {
	// Service is the Service that should receive a Shutdown signal.
	Service Service
	// Signals are the OS signals we should wait and signal on.  If this is empty, all signals will trigger Shutdown
	Signals []os.Signal
	// ShutdownContext creates the context that is passed to Shutdown.  Optional: if nil, will use context.Background.
	ShutdownContext func() context.Context
	// ShutdownTimeout is an optional timeout to attach to the Shutdown context.  Optional: if empty, no timeout will be used.
	ShutdownTimeout time.Duration
	signalNotify    func(c chan<- os.Signal, sig ...os.Signal)
	signalStop      func(c chan<- os.Signal)
	ch              chan os.Signal
}

// Setup ensures the signal channel is created and registered with the signal notifier.
func (w *SignalWatcher) Setup() error {
	if w.ch == nil {
		w.ch = make(chan os.Signal, 1)
	}
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
		ctx := w.shutdownContext()
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
	if w.signalStop == nil {
		signal.Stop(w.ch)
	} else {
		w.signalStop(w.ch)
	}
	close(w.ch)
	return nil
}

func (w *SignalWatcher) shutdownContext() context.Context {
	if w.ShutdownContext == nil {
		return context.Background()
	}
	return w.ShutdownContext()
}

var _ Setupable = &SignalWatcher{}
