package serviceexecutor

import "context"

// ServiceFunc allows turning functional callbacks into services.  It is the idea equivalent of http.HandlerFunc
type ServiceFunc struct {
	// OnSetup is allowed to be nil: -> Setup()
	OnSetup func() error
	// OnRun -> Run()
	OnRun func() error
	// OnShutdown -> Shutdown(ctx)
	OnShutdown func(ctx context.Context) error
}

// Setup simply executes OnSetup if it is non nil.  If it is nil, it returns nil.
func (f *ServiceFunc) Setup() error {
	if f.OnSetup == nil {
		return nil
	}
	return f.OnSetup()
}

// Run simply executes OnRun
func (f *ServiceFunc) Run() error {
	return f.OnRun()
}

// Shutdown simply executes OnShutdown
func (f *ServiceFunc) Shutdown(ctx context.Context) error {
	return f.OnShutdown(ctx)
}
