package serviceexecutor_test

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/cep21/serviceexecutor"
)

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
	// Fake a shutdown so our example finishes
	time.AfterFunc(time.Millisecond*10, func() {
		if err := m.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	})
	// After placing all your service inside Multi, run them.  This should end when Shutdown is called on Multi.
	err := m.Run()
	// Err is returned after all Run functions of Services complete.  An error is the
	fmt.Println(err)
	// Output: <nil>
}
