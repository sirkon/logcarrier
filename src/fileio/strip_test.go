package fileio

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStripMutualPrefix(t *testing.T) {
	path := "123456789abcdef"
	sample := "12345abcde"
	require.Equal(t, "6789abcdef", stripMutualPrefix(path, sample))

	path = "123456789abcdef"
	sample = "123456789abcde"
	require.Equal(t, "f", stripMutualPrefix(path, sample))

	path = "123456789abcdef"
	sample = "123456789abc"
	require.Equal(t, "def", stripMutualPrefix(path, sample))
}
