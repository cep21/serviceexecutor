package serviceexecutor

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// Hooks allow optional logging and metric gathering for service calls.
type MultiHooks struct {
	OnServiceRunStarted       func(s Service)
	OnServiceRunFinished      func(s Service, err error)
	OnServiceShutdownStarted  func(s Service)
	OnServiceShutdownFinished func(s Service, err error)
}

// States: (init) -> (setup ran) -> (run) -> (shutdown) ->

// Multi tracks multiple running services.  It is an error to modify Multi after you have called any method on it.
type Multi struct {
	// Services are managed by Multi with Run and Shutdown
	Services []Service
	// Hooks allow logging at lifecycle stages of Multi
	Hooks MultiHooks
	// setupCalled is set true after setup is called.
	setupCalled bool
	// shutdownCalled enforces that Shutdown cannot be called twice
	shutdownCalled int32
	// runCalled enforces that run cannot be called twice
	runCalled int32
	// runOnce allows us to enforce that the only services we call Shutdown on are those that we've already
	// called Run on
	runOnce sync.Once
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

// Run blocks until all services finish.  It is an error to call run twice.  If Shutdown is called before Run can spawn
// any goroutines, none will spawn.
func (m *Multi) Run() error {
	if atomic.SwapInt32(&m.runCalled, 1) == 1 {
		return errors.New("run called twice")
	}
	if !m.setupCalled {
		if err := m.Setup(); err != nil {
			return err
		}
	}
	wg := sync.WaitGroup{}
	errs := make([]error, 0, len(m.Services))
	m.runOnce.Do(func() {
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
	})
	wg.Wait()
	return errFromManyErrors(errs)
}

// Shutdown ends all running services.  It is an error to call Shutdown twice.  If
// Shutdown is called before we can call "Run" on services, it does nothing and returns nil.
func (m *Multi) Shutdown(ctx context.Context) error {
	if atomic.SwapInt32(&m.shutdownCalled, 1) == 1 {
		return errors.New("shutdown called twice")
	}
	services := m.Services
	m.runOnce.Do(func() {
		services = nil
	})
	errs := make([]error, 0, len(services))
	for i := len(services); i >= 0; i-- {
		s := services[i]
		m.Hooks.OnServiceShutdownStarted(s)
		err := s.Shutdown(ctx)
		m.Hooks.OnServiceShutdownFinished(s, err)
		errs[i] = err
	}
	return errFromManyErrors(errs)
}
