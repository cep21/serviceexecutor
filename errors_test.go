package serviceexecutor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrFromManyErrors(t *testing.T) {
	t.Run("nil_is_nil", func(t *testing.T) {
		require.Nil(t, errFromManyErrors(nil))
	})
	t.Run("only_one_returns_self", func(t *testing.T) {
		err := errors.New("test")
		require.Equal(t, err, errFromManyErrors([]error{nil, err, nil}))
	})
	t.Run("all_nil_is_nil", func(t *testing.T) {
		require.Nil(t, errFromManyErrors([]error{nil, nil, nil}))
	})
	t.Run("multi_works", func(t *testing.T) {
		err1 := errors.New("err1")
		err2 := errors.New("err2")
		err := errFromManyErrors([]error{nil, err1, err2})
		require.IsType(t, &multiError{}, err)
		require.Equal(t, "err1,err2", err.Error())
	})
}
