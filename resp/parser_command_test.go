package resp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSetArgs(t *testing.T) {
	t.Run("Simple SET", func(t *testing.T) {
		args := Array{
			BulkString("SET"),
			BulkString("key"),
			BulkString("value"),
		}
		parsed, err := ParseSetArgs(args)
		require.NoError(t, err)
		assert.Equal(t, BulkString("key"), parsed.Key)
		assert.Equal(t, BulkString("value"), parsed.Value)
		assert.False(t, parsed.NX)
		assert.False(t, parsed.XX)
		assert.False(t, parsed.Get)
		assert.Zero(t, parsed.Expire)
	})

	t.Run("SET with EX", func(t *testing.T) {
		args := Array{
			BulkString("SET"),
			BulkString("key"),
			BulkString("value"),
			BulkString("EX"),
			BulkString("10"),
		}
		parsed, err := ParseSetArgs(args)
		require.NoError(t, err)
		assert.Equal(t, BulkString("key"), parsed.Key)
		assert.Equal(t, BulkString("value"), parsed.Value)
		assert.WithinDuration(t, time.Now().Add(10*time.Second), time.UnixMilli(parsed.Expire), 50*time.Millisecond)
	})

	t.Run("SET with PX", func(t *testing.T) {
		args := Array{
			BulkString("SET"),
			BulkString("key"),
			BulkString("value"),
			BulkString("PX"),
			BulkString("1234"),
		}
		parsed, err := ParseSetArgs(args)
		require.NoError(t, err)
		assert.Equal(t, BulkString("key"), parsed.Key)
		assert.Equal(t, BulkString("value"), parsed.Value)
		assert.WithinDuration(t, time.Now().Add(1234*time.Millisecond), time.UnixMilli(parsed.Expire), 50*time.Millisecond)
	})

	t.Run("SET with NX", func(t *testing.T) {
		args := Array{
			BulkString("SET"),
			BulkString("key"),
			BulkString("value"),
			BulkString("NX"),
		}
		parsed, err := ParseSetArgs(args)
		require.NoError(t, err)
		assert.True(t, parsed.NX)
	})

	t.Run("SET with XX", func(t *testing.T) {
		args := Array{
			BulkString("SET"),
			BulkString("key"),
			BulkString("value"),
			BulkString("XX"),
		}
		parsed, err := ParseSetArgs(args)
		require.NoError(t, err)
		assert.True(t, parsed.XX)
	})

	t.Run("SET with GET", func(t *testing.T) {
		args := Array{
			BulkString("SET"),
			BulkString("key"),
			BulkString("value"),
			BulkString("GET"),
		}
		parsed, err := ParseSetArgs(args)
		require.NoError(t, err)
		assert.True(t, parsed.Get)
	})

	t.Run("SET with multiple options", func(t *testing.T) {
		args := Array{
			BulkString("SET"),
			BulkString("key"),
			BulkString("value"),
			BulkString("NX"),
			BulkString("GET"),
			BulkString("PX"),
			BulkString("500"),
		}
		parsed, err := ParseSetArgs(args)
		require.NoError(t, err)
		assert.True(t, parsed.NX)
		assert.True(t, parsed.Get)
		assert.WithinDuration(t, time.Now().Add(500*time.Millisecond), time.UnixMilli(parsed.Expire), 50*time.Millisecond)
	})

	t.Run("Error: NX and XX", func(t *testing.T) {
		args := Array{
			BulkString("SET"),
			BulkString("key"),
			BulkString("value"),
			BulkString("NX"),
			BulkString("XX"),
		}
		_, err := ParseSetArgs(args)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "syntax error")
	})

	t.Run("Error: EX and PX", func(t *testing.T) {
		args := Array{
			BulkString("SET"),
			BulkString("key"),
			BulkString("value"),
			BulkString("EX"),
			BulkString("10"),
			BulkString("PX"),
			BulkString("10000"),
		}
		_, err := ParseSetArgs(args)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "syntax error")
	})

	t.Run("Error: Invalid EX value", func(t *testing.T) {
		args := Array{
			BulkString("SET"),
			BulkString("key"),
			BulkString("value"),
			BulkString("EX"),
			BulkString("not-a-number"),
		}
		_, err := ParseSetArgs(args)
		require.Error(t, err)
	})
}
