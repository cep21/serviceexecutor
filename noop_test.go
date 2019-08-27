package serviceexecutor

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoop(t *testing.T) {
	n := Noop{}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		require.NoError(t, n.Run())
	}()
	go func() {
		defer wg.Done()
		require.NoError(t, n.Shutdown(context.Background()))
	}()
	wg.Wait()
}
