package serviceexecutor

import (
	"context"
	"github.com/stretchr/testify/require"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type testService struct {
	runCount      int32
	shutdownCount int32
	setupCount    int32
	runErr        error
	shutdownErr   error
	setupErr      error
	n             Noop
}

func (t *testService) Run() error {
	atomic.AddInt32(&t.runCount, 1)
	if t.runErr != nil {
		return t.runErr
	}
	return t.n.Run()
}

func (t *testService) Shutdown(ctx context.Context) error {
	atomic.AddInt32(&t.shutdownCount, 1)
	if t.shutdownErr != nil {
		return t.shutdownErr
	}
	return t.n.Shutdown(ctx)
}

func (t *testService) Setup() error {
	atomic.AddInt32(&t.setupCount, 1)
	return t.setupErr
}

var _ Setupable = &testService{}

func TestMulti_Run(t *testing.T) {
	ts := &testService{}
	m := Multi{
		Services: []Service{
			ts,
		},
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		require.NoError(t, m.Run())
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond)
		require.NoError(t, m.Shutdown(context.Background()))
	}()
	wg.Wait()
	require.EqualValues(t, 1, ts.setupCount)
	require.EqualValues(t, 1, ts.runCount)
	require.EqualValues(t, 1, ts.shutdownCount)
}
