# serviceexecutor
[![CircleCI](https://circleci.com/gh/cep21/serviceexecutor.svg)](https://circleci.com/gh/cep21/serviceexecutor)
[![GoDoc](https://godoc.org/github.com/cep21/serviceexecutor?status.svg)](https://godoc.org/github.com/cep21/serviceexecutor)
[![codecov](https://codecov.io/gh/cep21/serviceexecutor/branch/master/graph/badge.svg)](https://codecov.io/gh/cep21/serviceexecutor)

Serviceexecutor can manage long running service goroutines in go.

A best practice for Go libraries is to push concurrency control up the call stack.  There is more
information about this at both this [Gophercon](https://www.youtube.com/watch?v=5v2fqm_8jYI&feature=youtu.be&t=1852)
talk and the [Synchronous Functions](https://github.com/golang/go/wiki/CodeReviewComments#synchronous-functions) part of
of the Go wiki.

This means if your microservice needs daemon goroutines, you should write these mini services as **blocking** daemons
and push the spawning of these daemon goroutines up your call stack, ideally up to main.go.  These daemon processes
will need a way to run themselves, a way to shutdown cleanly, and sometimes, a way to setup state before running.

All of this is encapsulated inside this application's Service interface.  This package contains helpers to both
spawn and track many of these services at once, as well as helpers to trigger gracefull shutdown on os signals, such as
SIGTERM.

# Usage

```go
func ExampleMulti() {
	// Start with some services.  These services obey the serviceexecutor.Service contract
	service1 := &serviceexecutor.Noop{}
	service2 := &serviceexecutor.Noop{}
	// Create a multi for all of these services
	m := &serviceexecutor.Multi{
		Services: []serviceexecutor.Service{service1, service2},
	}
	// Register a Shutdown signal on your serviceexecutor.  You can also use SignalWatcher directly.
	serviceexecutor.ShutdownOnSignals(m, time.Second, syscall.SIGTERM)
	// After placing all your service inside Multi, run them.  This should end when Shutdown is called on Multi.
	err := m.Run()
	// Err is returned after all Run functions of Services complete.  An error is the
	fmt.Println(err)
	// Output: <nil>
}
```

# Design Rational

## Primary abstraction
The primary design is inside the choice of API for Service and Setupable, copied below.

```go
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
```

### Run
There should be a function to start the service.  The name `Start` is generally reserved for services that spawn
themselves in a goroutine.  The use of `Run` is already used by [os.exec](https://golang.org/pkg/os/exec/#Cmd.Run)
so it seems like a natural name for something that runs and blocks till done.

### Shutdown
There should be a way to stop the service.  The initial choice was `Close(error)` which matches
[io.Closer](https://golang.org/pkg/io/#Closer).  However, there's no way to gracefully shutdown services with `Close`
in a way that won't block forever.  Since context.Context is a natural way to escape functions early, the choice of
Shutdown(ctx) was made to match [http.Server.Shutdown](https://golang.org/pkg/net/http/#Server.Shutdown).

### Setup
Many services do not need a setup phase, so it is optional for Multi.  However, for services that do the general pattern
is that the main application should fail to start if the setup phase does not complete.  This puts `Setup()` as a blocking,
optional pre step to Multi startup.

# Contributing

Contributions welcome!  Submit a pull request on github and make sure your code passes `make lint test`.  For
large changes, I strongly recommend [creating an issue](https://github.com/cep21/serviceexecutor/issues) on GitHub first to
confirm your change will be accepted before writing a lot of code.  GitHub issues are also recommended, at your discretion,
for smaller changes or questions.

# License

This library is licensed under the Apache 2.0 License.