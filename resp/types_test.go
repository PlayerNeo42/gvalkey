package resp

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimpleString_Reader(t *testing.T) {
	s := SimpleString("OK")
	reader := s.RESPReader()
	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, "+OK\r\n", string(data))
}

func TestSimpleError_Reader(t *testing.T) {
	e := NewSimpleError("Error message")
	reader := e.RESPReader()
	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, "-ERR Error message\r\n", string(data))
}

func TestBulkString_Reader(t *testing.T) {
	t.Run("Normal string", func(t *testing.T) {
		b := BulkString("hello")
		reader := b.RESPReader()
		data, err := io.ReadAll(reader)
		require.NoError(t, err)
		require.Equal(t, "$5\r\nhello\r\n", string(data))
	})

	t.Run("Empty string", func(t *testing.T) {
		b := BulkString("")
		reader := b.RESPReader()
		data, err := io.ReadAll(reader)
		require.NoError(t, err)
		require.Equal(t, "$0\r\n\r\n", string(data))
	})
}

func TestInteger_Reader(t *testing.T) {
	i := Integer(1000)
	reader := i.RESPReader()
	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, ":1000\r\n", string(data))
}

func TestNull_Reader(t *testing.T) {
	n := Null{}
	reader := n.RESPReader()
	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, "$-1\r\n", string(data))
}

func TestArray_Reader(t *testing.T) {
	a := Array{
		BulkString("hello"),
		Integer(123),
	}
	reader := a.RESPReader()
	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	expected := "*2\r\n$5\r\nhello\r\n:123\r\n"
	require.Equal(t, expected, string(data))
}
