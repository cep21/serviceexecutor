package serviceexecutor

import (
	"context"
	"errors"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

func TestMultiBasic(t *testing.T) {
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

func TestMultiFull(t *testing.T) {
	ts := &testService{}
	m := Multi{
		Services: []Service{
			ts,
		},
	}
	w := SignalWatcher{
		Service: &m,
		Signals: []os.Signal{syscall.SIGQUIT},
	}
	require.NoError(t, w.Setup())
	m.Services = append(m.Services, &w)
	go func() {
		time.Sleep(time.Millisecond)
		w.ch <- syscall.SIGQUIT
	}()
	require.NoError(t, m.Run())
	require.EqualValues(t, 1, ts.setupCount)
	require.EqualValues(t, 1, ts.runCount)
	require.EqualValues(t, 1, ts.shutdownCount)
}

func TestMultiSetupErr(t *testing.T) {
	ts := &testService{
		setupErr: errors.New("bad"),
	}
	m := Multi{
		Services: []Service{
			ts,
		},
	}
	require.Equal(t, ts.setupErr, m.Run())
	require.EqualValues(t, 1, ts.setupCount)
	// Never run because a single service failed to setup
	require.EqualValues(t, 0, ts.runCount)
	require.EqualValues(t, 0, ts.shutdownCount)
}

func TestMultiDoubleRun(t *testing.T) {
	m := Multi{
		Services: []Service{
			&testService{},
		},
	}
	var runErr error
	var correctCount int
	wg := sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := m.Run(); err != nil {
				runErr = err
			} else {
				correctCount++
			}
		}()
	}
	time.Sleep(time.Millisecond)
	require.Nil(t, m.Shutdown(context.Background()))
	wg.Wait()
	require.Equal(t, 1, correctCount)
	require.IsType(t, &repeatedCalls{}, runErr)
	require.IsType(t, &repeatedCalls{}, m.Shutdown(context.Background()))
}

func TestMultiShutdownFirst(t *testing.T) {
	m := Multi{
		Services: []Service{
			&Noop{},
		},
	}
	require.Nil(t, m.Shutdown(context.Background()))
	require.Nil(t, m.Run())
}
