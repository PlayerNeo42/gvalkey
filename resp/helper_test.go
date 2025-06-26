package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeekNextInteger(t *testing.T) {
	t.Run("Valid integer", func(t *testing.T) {
		args := Array{BulkString("123")}
		val, err := peekNextInteger(args, -1)
		require.NoError(t, err)
		assert.Equal(t, int64(123), val)
	})

	t.Run("Not an integer", func(t *testing.T) {
		args := Array{BulkString("not-a-number")}
		_, err := peekNextInteger(args, -1)
		require.Error(t, err)
	})

	t.Run("Index out of bounds", func(t *testing.T) {
		args := Array{}
		_, err := peekNextInteger(args, 0)
		require.Error(t, err)
	})

	t.Run("Wrong type", func(t *testing.T) {
		args := Array{SimpleString("123")}
		_, err := peekNextInteger(args, -1)
		require.Error(t, err)
	})
}
