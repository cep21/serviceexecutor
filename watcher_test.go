package serviceexecutor

import (
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fakeSignal struct {
	m map[chan<- os.Signal][]os.Signal
}

func (f *fakeSignal) signalNotify(c chan<- os.Signal, sig ...os.Signal) {
	if f.m == nil {
		f.m = make(map[chan<- os.Signal][]os.Signal)
	}
	f.m[c] = sig
}

func (f *fakeSignal) signalStop(c chan<- os.Signal) {
	delete(f.m, c)
}

func (f *fakeSignal) sendSignal(sig os.Signal) {
	for c, arr := range f.m {
		if len(arr) == 0 {
			c <- sig
		}
		for _, ch := range arr {
			if ch.String() == sig.String() {
				c <- sig
			}
		}
	}
}

func TestSignalWatcher(t *testing.T) {
	n := &Noop{}
	var f fakeSignal
	w := SignalWatcher{
		Service:      n,
		Signals:      []os.Signal{syscall.SIGINT},
		signalNotify: f.signalNotify,
		signalStop:   f.signalStop,
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		require.NoError(t, w.Run())
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond)
		f.sendSignal(syscall.SIGINT)
	}()
	wg.Wait()
}
