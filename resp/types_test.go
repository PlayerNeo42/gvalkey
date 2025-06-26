package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleString_MarshalRESP(t *testing.T) {
	s := SimpleString("OK")
	assert.Equal(t, []byte("+OK\r\n"), s.MarshalRESP())
}

func TestSimpleError_MarshalRESP(t *testing.T) {
	e := NewSimpleError("Error message")
	assert.Equal(t, []byte("-ERR Error message\r\n"), e.MarshalRESP())
}

func TestBulkString_MarshalRESP(t *testing.T) {
	t.Run("Normal string", func(t *testing.T) {
		b := BulkString("hello")
		assert.Equal(t, []byte("$5\r\nhello\r\n"), b.MarshalRESP())
	})

	t.Run("Empty string", func(t *testing.T) {
		b := BulkString("")
		assert.Equal(t, []byte("$0\r\n\r\n"), b.MarshalRESP())
	})
}

func TestInteger_MarshalRESP(t *testing.T) {
	i := Integer(1000)
	assert.Equal(t, []byte(":1000\r\n"), i.MarshalRESP())
}

func TestNull_MarshalRESP(t *testing.T) {
	n := Null{}
	assert.Equal(t, []byte("$-1\r\n"), n.MarshalRESP())
}

func TestArray_MarshalRESP(t *testing.T) {
	a := Array{
		BulkString("hello"),
		Integer(123),
	}
	expected := []byte("*2\r\n$5\r\nhello\r\n:123\r\n")
	assert.Equal(t, expected, a.MarshalRESP())
}
