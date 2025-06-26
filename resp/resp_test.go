package resp_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/PlayerNeo42/gvalkey/resp"

	"github.com/stretchr/testify/require"
)

// Test integration of RESP marshaling functionality
func TestMarshalerIntegration(t *testing.T) {
	t.Run("Mixed types array marshaling", func(t *testing.T) {
		array := resp.Array{
			resp.SimpleString("OK"),
			resp.Integer(123),
			resp.BulkString("hello"),
			resp.NewSimpleError("test error"),
			resp.Null{},
		}

		result := array.MarshalRESP()
		expected := "*5\r\n+OK\r\n:123\r\n$5\r\nhello\r\n-ERR test error\r\n$-1\r\n"
		require.Equal(t, expected, string(result))
	})

	t.Run("Nested array marshaling", func(t *testing.T) {
		innerArray := resp.Array{
			resp.BulkString("inner"),
			resp.Integer(456),
		}
		outerArray := resp.Array{
			resp.BulkString("outer"),
			innerArray,
		}

		result := outerArray.MarshalRESP()
		expected := "*2\r\n$5\r\nouter\r\n*2\r\n$5\r\ninner\r\n:456\r\n"
		require.Equal(t, expected, string(result))
	})

	t.Run("Empty array", func(t *testing.T) {
		array := resp.Array{}
		result := array.MarshalRESP()
		expected := "*0\r\n"
		require.Equal(t, expected, string(result))
	})
}

// Test integration functionality of parser
func TestParserIntegration(t *testing.T) {
	t.Run("Parse simple string", func(t *testing.T) {
		data := "+OK\r\n"
		parser := resp.NewParser(strings.NewReader(data))
		result, err := parser.Parse()
		require.NoError(t, err)
		require.Equal(t, resp.SimpleString("OK"), result)
	})

	t.Run("Parse integer", func(t *testing.T) {
		data := ":1000\r\n"
		parser := resp.NewParser(strings.NewReader(data))
		result, err := parser.Parse()
		require.NoError(t, err)
		require.Equal(t, resp.Integer(1000), result)
	})

	t.Run("Parse bulk string", func(t *testing.T) {
		data := "$5\r\nhello\r\n"
		parser := resp.NewParser(strings.NewReader(data))
		result, err := parser.Parse()
		require.NoError(t, err)
		require.Equal(t, resp.BulkString("hello"), result)
	})

	t.Run("Parse null bulk string", func(t *testing.T) {
		data := "$-1\r\n"
		parser := resp.NewParser(strings.NewReader(data))
		result, err := parser.Parse()
		require.NoError(t, err)
		require.Equal(t, resp.BulkString(""), result)
	})

	t.Run("Parse empty bulk string", func(t *testing.T) {
		data := "$0\r\n\r\n"
		parser := resp.NewParser(strings.NewReader(data))
		result, err := parser.Parse()
		require.NoError(t, err)
		require.Equal(t, resp.BulkString(""), result)
	})

	t.Run("Parse array", func(t *testing.T) {
		data := "*2\r\n$5\r\nhello\r\n:123\r\n"
		parser := resp.NewParser(strings.NewReader(data))
		result, err := parser.Parse()
		require.NoError(t, err)

		array, ok := result.(resp.Array)
		require.True(t, ok)
		require.Len(t, array, 2)
		require.Equal(t, resp.BulkString("hello"), array[0])
		require.Equal(t, resp.Integer(123), array[1])
	})

	t.Run("Parse empty array", func(t *testing.T) {
		data := "*0\r\n"
		parser := resp.NewParser(strings.NewReader(data))
		result, err := parser.Parse()
		require.NoError(t, err)

		array, ok := result.(resp.Array)
		require.True(t, ok)
		require.Len(t, array, 0)
	})

	t.Run("Parse complex nested array", func(t *testing.T) {
		data := "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
		parser := resp.NewParser(strings.NewReader(data))
		result, err := parser.Parse()
		require.NoError(t, err)

		array, ok := result.(resp.Array)
		require.True(t, ok)
		require.Len(t, array, 3)
		require.Equal(t, resp.BulkString("SET"), array[0])
		require.Equal(t, resp.BulkString("key"), array[1])
		require.Equal(t, resp.BulkString("value"), array[2])
	})
}

// Test parser error handling
func TestParserErrorHandling(t *testing.T) {
	t.Run("Invalid array length", func(t *testing.T) {
		data := "*abc\r\n"
		parser := resp.NewParser(strings.NewReader(data))
		_, err := parser.Parse()
		require.Error(t, err)
		require.Contains(t, err.Error(), "parse array length failed")
	})

	t.Run("Invalid bulk string length", func(t *testing.T) {
		data := "$abc\r\n"
		parser := resp.NewParser(strings.NewReader(data))
		_, err := parser.Parse()
		require.Error(t, err)
		require.Contains(t, err.Error(), "parse bulk string length failed")
	})

	t.Run("Invalid integer", func(t *testing.T) {
		data := ":abc\r\n"
		parser := resp.NewParser(strings.NewReader(data))
		_, err := parser.Parse()
		require.Error(t, err)
		require.Contains(t, err.Error(), "parse integer failed")
	})

	t.Run("Unsupported RESP type", func(t *testing.T) {
		data := "?unknown\r\n"
		parser := resp.NewParser(strings.NewReader(data))
		_, err := parser.Parse()
		require.Error(t, err)
		require.Contains(t, err.Error(), "unsupported RESP type")
	})

	t.Run("RESP3 type not supported", func(t *testing.T) {
		data := "_\r\n"
		parser := resp.NewParser(strings.NewReader(data))
		_, err := parser.Parse()
		require.Error(t, err)
		require.Contains(t, err.Error(), "RESP3 type not supported yet")
	})
}

// Test constant definitions
func TestConstants(t *testing.T) {
	t.Run("Response constants", func(t *testing.T) {
		require.Equal(t, "+OK\r\n", string(resp.OK.MarshalRESP()))
		require.Equal(t, "$-1\r\n", string(resp.NULL.MarshalRESP()))
	})

	t.Run("Command constants", func(t *testing.T) {
		// Test some common command constants
		require.Equal(t, "$3\r\nSET\r\n", string(resp.SET.MarshalRESP()))
		require.Equal(t, "$3\r\nGET\r\n", string(resp.GET.MarshalRESP()))
		require.Equal(t, "$3\r\nDEL\r\n", string(resp.DEL.MarshalRESP()))
		require.Equal(t, "$2\r\nEX\r\n", string(resp.EX.MarshalRESP()))
		require.Equal(t, "$2\r\nPX\r\n", string(resp.PX.MarshalRESP()))
		require.Equal(t, "$2\r\nNX\r\n", string(resp.NX.MarshalRESP()))
		require.Equal(t, "$2\r\nXX\r\n", string(resp.XX.MarshalRESP()))
	})
}

// Test binary marshaling functionality
func TestBinaryMarshaling(t *testing.T) {
	t.Run("BulkString binary marshaling", func(t *testing.T) {
		b := resp.BulkString("hello world")
		result := b.MarshalBinary()
		require.Equal(t, []byte("hello world"), result)
	})

	t.Run("Integer binary marshaling", func(t *testing.T) {
		i := resp.Integer(12345)
		result := i.MarshalBinary()
		require.Equal(t, []byte("12345"), result)
	})
}

// Test BulkString special methods
func TestBulkStringMethods(t *testing.T) {
	t.Run("Upper method", func(t *testing.T) {
		b := resp.BulkString("hello world")
		upper := b.Upper()
		require.Equal(t, resp.BulkString("HELLO WORLD"), upper)
	})

	t.Run("Upper with mixed case", func(t *testing.T) {
		b := resp.BulkString("HeLLo WoRLd")
		upper := b.Upper()
		require.Equal(t, resp.BulkString("HELLO WORLD"), upper)
	})
}

// Test complete round-trip encoding and decoding
func TestRoundTrip(t *testing.T) {
	testCases := []struct {
		name     string
		data     interface{}
		expected interface{}
	}{
		{"SimpleString", resp.SimpleString("test"), resp.SimpleString("test")},
		{"Integer", resp.Integer(42), resp.Integer(42)},
		{"BulkString", resp.BulkString("test string"), resp.BulkString("test string")},
		{"Empty BulkString", resp.BulkString(""), resp.BulkString("")},
		{"Null", resp.Null{}, resp.BulkString("")}, // Null is parsed as empty BulkString
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			marshaler, ok := tc.data.(resp.Marshaler)
			require.True(t, ok, "Type should implement Marshaler interface")

			// Marshal
			encoded := marshaler.MarshalRESP()
			require.NotEmpty(t, encoded)

			// Unmarshal
			parser := resp.NewParser(bytes.NewReader(encoded))
			decoded, err := parser.Parse()
			require.NoError(t, err)
			require.Equal(t, tc.expected, decoded)
		})
	}
}

// Test Redis command scenarios
func TestRedisCommandScenarios(t *testing.T) {
	t.Run("SET command", func(t *testing.T) {
		command := resp.Array{
			resp.SET,
			resp.BulkString("mykey"),
			resp.BulkString("myvalue"),
		}

		encoded := command.MarshalRESP()
		parser := resp.NewParser(bytes.NewReader(encoded))
		decoded, err := parser.Parse()

		require.NoError(t, err)
		array, ok := decoded.(resp.Array)
		require.True(t, ok)
		require.Len(t, array, 3)
		require.Equal(t, resp.BulkString("SET"), array[0])
		require.Equal(t, resp.BulkString("mykey"), array[1])
		require.Equal(t, resp.BulkString("myvalue"), array[2])
	})

	t.Run("GET command", func(t *testing.T) {
		command := resp.Array{
			resp.GET,
			resp.BulkString("mykey"),
		}

		encoded := command.MarshalRESP()
		parser := resp.NewParser(bytes.NewReader(encoded))
		decoded, err := parser.Parse()

		require.NoError(t, err)
		array, ok := decoded.(resp.Array)
		require.True(t, ok)
		require.Len(t, array, 2)
		require.Equal(t, resp.BulkString("GET"), array[0])
		require.Equal(t, resp.BulkString("mykey"), array[1])
	})

	t.Run("DEL command with multiple keys", func(t *testing.T) {
		command := resp.Array{
			resp.DEL,
			resp.BulkString("key1"),
			resp.BulkString("key2"),
			resp.BulkString("key3"),
		}

		encoded := command.MarshalRESP()
		parser := resp.NewParser(bytes.NewReader(encoded))
		decoded, err := parser.Parse()

		require.NoError(t, err)
		array, ok := decoded.(resp.Array)
		require.True(t, ok)
		require.Len(t, array, 4)
		require.Equal(t, resp.BulkString("DEL"), array[0])
		require.Equal(t, resp.BulkString("key1"), array[1])
		require.Equal(t, resp.BulkString("key2"), array[2])
		require.Equal(t, resp.BulkString("key3"), array[3])
	})
}

// Test large data handling
func TestLargeDataHandling(t *testing.T) {
	t.Run("Large bulk string", func(t *testing.T) {
		// Create a large string
		largeData := strings.Repeat("x", 1024)
		bulkString := resp.BulkString(largeData)

		encoded := bulkString.MarshalRESP()
		parser := resp.NewParser(bytes.NewReader(encoded))
		decoded, err := parser.Parse()

		require.NoError(t, err)
		require.Equal(t, bulkString, decoded)
	})

	t.Run("Large array", func(t *testing.T) {
		// Create an array with multiple elements
		array := make(resp.Array, 100)
		for i := 0; i < 100; i++ {
			array[i] = resp.BulkString(strings.Repeat("data", i+1))
		}

		encoded := array.MarshalRESP()
		parser := resp.NewParser(bytes.NewReader(encoded))
		decoded, err := parser.Parse()

		require.NoError(t, err)
		decodedArray, ok := decoded.(resp.Array)
		require.True(t, ok)
		require.Len(t, decodedArray, 100)

		// Verify a few elements
		require.Equal(t, resp.BulkString("data"), decodedArray[0])
		require.Equal(t, resp.BulkString(strings.Repeat("data", 50)), decodedArray[49])
		require.Equal(t, resp.BulkString(strings.Repeat("data", 100)), decodedArray[99])
	})
}

// Test edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("Zero integer", func(t *testing.T) {
		i := resp.Integer(0)
		encoded := i.MarshalRESP()
		require.Equal(t, ":0\r\n", string(encoded))

		parser := resp.NewParser(bytes.NewReader(encoded))
		decoded, err := parser.Parse()
		require.NoError(t, err)
		require.Equal(t, i, decoded)
	})

	t.Run("Negative integer", func(t *testing.T) {
		i := resp.Integer(-123)
		encoded := i.MarshalRESP()
		require.Equal(t, ":-123\r\n", string(encoded))

		parser := resp.NewParser(bytes.NewReader(encoded))
		decoded, err := parser.Parse()
		require.NoError(t, err)
		require.Equal(t, i, decoded)
	})

	t.Run("BulkString with special characters", func(t *testing.T) {
		special := resp.BulkString("hello\r\nworld\ttab")
		encoded := special.MarshalRESP()

		parser := resp.NewParser(bytes.NewReader(encoded))
		decoded, err := parser.Parse()
		require.NoError(t, err)
		require.Equal(t, special, decoded)
	})

	t.Run("Empty simple string", func(t *testing.T) {
		s := resp.SimpleString("")
		encoded := s.MarshalRESP()
		require.Equal(t, "+\r\n", string(encoded))

		parser := resp.NewParser(bytes.NewReader(encoded))
		decoded, err := parser.Parse()
		require.NoError(t, err)
		require.Equal(t, s, decoded)
	})
}
