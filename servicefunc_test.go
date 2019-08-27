package serviceexecutor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceFunc(t *testing.T) {
	var run int
	var shutdown int
	s := ServiceFunc{
		OnRun: func() error {
			run++
			return nil
		},
		OnShutdown: func(ctx context.Context) error {
			shutdown++
			return nil
		},
	}
	require.NoError(t, s.Run())
	require.Equal(t, 1, run)
	require.NoError(t, s.OnShutdown(context.Background()))
	require.Equal(t, 1, shutdown)
}
