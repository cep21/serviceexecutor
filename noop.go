package serviceexecutor

import (
	"context"
	"sync"
)

// Noop is a service that does nothing
type Noop struct {
	once sync.Once
	ch   chan struct{}
}

func (n *Noop) init() {
	n.once.Do(func() {
		n.ch = make(chan struct{})
	})
}

// Run blocks until shutdown
func (n *Noop) Run() error {
	n.init()
	<-n.ch
	return nil
}

// Shutdown stops Run
func (n *Noop) Shutdown(ctx context.Context) error {
	n.init()
	close(n.ch)
	return nil
}

var _ Service = &Noop{}
