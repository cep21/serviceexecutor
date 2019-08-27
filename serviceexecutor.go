package serviceexecutor

import (
	"context"
)

// Service is intended to be a long running daemon style goroutine of your program.
type Service interface {
	// Run should block and not return until the service either has a fatal error during execution (returned via err)
	// or Shutdown is called on the service.
	Run() error
	// Shutdown should end a current Run() call.  Shutdown should try to exit only when it is
	// sure that Run() has finished gracefully.  It should not persist longer than the length of ctx.
	Shutdown(ctx context.Context) error
}

// Setupable is an optional interface of Service.  Services that implement Setupable have their Setup function called
// by Multi before Run is executed.
type Setupable interface {
	Service
	// Setup is expected to be called once and before Run
	Setup() error
}
