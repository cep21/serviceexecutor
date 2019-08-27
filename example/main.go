package main

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/cep21/serviceexecutor"
)

type Printer struct {
	shutdownSignal chan struct{}
	runDone        chan struct{}
}

func (p *Printer) Setup() error {
	p.shutdownSignal = make(chan struct{})
	p.runDone = make(chan struct{})
	return nil
}

func (p *Printer) Run() error {
	defer close(p.runDone)
	for {
		select {
		case <-p.shutdownSignal:
			return nil
		case <-time.After(time.Second):
			fmt.Println("hello", time.Now())
		}
	}
}

func (p *Printer) Shutdown(ctx context.Context) error {
	close(p.shutdownSignal)
	<-p.runDone
	return nil
}

var _ serviceexecutor.Setupable = &Printer{}

func main() {
	m := &serviceexecutor.Multi{
		Services: []serviceexecutor.Service{&Printer{}},
	}
	serviceexecutor.ShutdownOnSignals(m, time.Second, syscall.SIGTERM, syscall.SIGINT)
	fmt.Println("Starting to run.  PID=", os.Getpid())
	if err := m.Run(); err != nil {
		panic("Should not panic")
	}
	fmt.Println("Done running")
}
