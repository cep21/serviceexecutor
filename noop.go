package serviceexecutor

import "context"

// Noop is a service that does nothing
type Noop struct {
	ch chan struct{}
}

func (n Noop) Run() error {
	n.ch = make(chan struct{})
	<-n.ch
	return nil
}

func (n Noop) Shutdown(ctx context.Context) error {
	close(n.ch)
	return nil
}

var _ Service = Noop{}
