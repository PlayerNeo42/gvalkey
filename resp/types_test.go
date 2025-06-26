package resp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimpleString_MarshalRESP(t *testing.T) {
	s := SimpleString("OK")
	require.Equal(t, []byte("+OK\r\n"), s.MarshalRESP())
}

func TestSimpleError_MarshalRESP(t *testing.T) {
	e := NewSimpleError("Error message")
	require.Equal(t, []byte("-ERR Error message\r\n"), e.MarshalRESP())
}

func TestBulkString_MarshalRESP(t *testing.T) {
	t.Run("Normal string", func(t *testing.T) {
		b := BulkString("hello")
		require.Equal(t, []byte("$5\r\nhello\r\n"), b.MarshalRESP())
	})

	t.Run("Empty string", func(t *testing.T) {
		b := BulkString("")
		require.Equal(t, []byte("$0\r\n\r\n"), b.MarshalRESP())
	})
}

func TestInteger_MarshalRESP(t *testing.T) {
	i := Integer(1000)
	require.Equal(t, []byte(":1000\r\n"), i.MarshalRESP())
}

func TestNull_MarshalRESP(t *testing.T) {
	n := Null{}
	require.Equal(t, []byte("$-1\r\n"), n.MarshalRESP())
}

func TestArray_MarshalRESP(t *testing.T) {
	a := Array{
		BulkString("hello"),
		Integer(123),
	}
	expected := []byte("*2\r\n$5\r\nhello\r\n:123\r\n")
	require.Equal(t, expected, a.MarshalRESP())
}
