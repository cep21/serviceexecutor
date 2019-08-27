package serviceexecutor

import (
	"context"
	"sync"
)

// Hooks allow optional logging and metric gathering for service calls.
type MultiHooks struct {
	OnServiceRunStarted       func(s Service)
	OnServiceRunFinished      func(s Service, err error)
	OnServiceShutdownStarted  func(s Service)
	OnServiceShutdownFinished func(s Service, err error)
}

// States: (init) -> (setup ran) -> (run) -> (shutdown) ->

// Multi tracks multiple running services.
type Multi struct {
	Services    []Service
	Hooks       MultiHooks
	setupCalled bool
}

var _ Setupable = &Multi{}

// Setup calls Setup if all Services. You do not have to call this manually.  It will be called for you by Run.
func (m *Multi) Setup() error {
	if m.setupCalled {
		// This could be an error, but that feels wrong.  There is no valid use of Multi that would allow Setup
		// to be called twice.
		panic("logic error: do not call setup twice")
	}
	m.setupCalled = true
	for _, s := range m.Services {
		if asSetupable, ok := s.(Setupable); ok {
			if err := asSetupable.Setup(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Run blocks until all services finish.
func (m *Multi) Run() error {
	if !m.setupCalled {
		if err := m.Setup(); err != nil {
			return err
		}
	}
	wg := sync.WaitGroup{}
	errs := make([]error, 0, len(m.Services))
	for i, s := range m.Services {
		wg.Add(1)
		i := i
		s := s
		go func() {
			defer wg.Done()
			m.Hooks.OnServiceRunStarted(s)
			err := s.Run()
			m.Hooks.OnServiceRunFinished(s, err)
			errs[i] = err
		}()
	}
	wg.Wait()
	return errFromManyErrors(errs)
}

// Shutdown ends all running services and prevents future calls to Run.
func (m *Multi) Shutdown(ctx context.Context) error {
	errs := make([]error, 0, len(m.Services))
	for i := len(m.Services); i >= 0; i-- {
		s := m.Services[i]
		m.Hooks.OnServiceShutdownStarted(s)
		err := s.Shutdown(ctx)
		m.Hooks.OnServiceShutdownFinished(s, err)
		errs[i] = err
	}
	return errFromManyErrors(errs)
}
