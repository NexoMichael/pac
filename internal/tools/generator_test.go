package tools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	ids := make(chan uint)
	require.Len(t, ids, 0)

	stop := NewGenerator(ids)
	for i := 1; i < 1000; i++ {
		value := <-ids
		require.Equal(t, i, int(value))
	}
	stop <- true
	value, ok := <-ids
	require.False(t, ok)
	require.Zero(t, value)
}
