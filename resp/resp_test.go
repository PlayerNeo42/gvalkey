package resp_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/PlayerNeo42/gvalkey/resp"

	"github.com/stretchr/testify/require"
)

// 测试RESP编组功能的集成测试
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

// 测试解析器的集成功能
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

// 测试Parser错误处理
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

// 测试常量定义
func TestConstants(t *testing.T) {
	t.Run("Response constants", func(t *testing.T) {
		require.Equal(t, "+OK\r\n", string(resp.OK.MarshalRESP()))
		require.Equal(t, "$-1\r\n", string(resp.NULL.MarshalRESP()))
	})

	t.Run("Command constants", func(t *testing.T) {
		// 测试一些常用命令常量
		require.Equal(t, "$3\r\nSET\r\n", string(resp.SET.MarshalRESP()))
		require.Equal(t, "$3\r\nGET\r\n", string(resp.GET.MarshalRESP()))
		require.Equal(t, "$3\r\nDEL\r\n", string(resp.DEL.MarshalRESP()))
		require.Equal(t, "$2\r\nEX\r\n", string(resp.EX.MarshalRESP()))
		require.Equal(t, "$2\r\nPX\r\n", string(resp.PX.MarshalRESP()))
		require.Equal(t, "$2\r\nNX\r\n", string(resp.NX.MarshalRESP()))
		require.Equal(t, "$2\r\nXX\r\n", string(resp.XX.MarshalRESP()))
	})
}

// 测试二进制编组功能
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

// 测试BulkString特殊方法
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

// 测试完整的往返编解码
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
		{"Null", resp.Null{}, resp.BulkString("")}, // Null被解析为空的BulkString
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			marshaler, ok := tc.data.(resp.Marshaler)
			require.True(t, ok, "Type should implement Marshaler interface")

			// 编组
			encoded := marshaler.MarshalRESP()
			require.NotEmpty(t, encoded)

			// 解组
			parser := resp.NewParser(bytes.NewReader(encoded))
			decoded, err := parser.Parse()
			require.NoError(t, err)
			require.Equal(t, tc.expected, decoded)
		})
	}
}

// 测试Redis命令场景
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

// 测试大数据处理
func TestLargeDataHandling(t *testing.T) {
	t.Run("Large bulk string", func(t *testing.T) {
		// 创建一个较大的字符串
		largeData := strings.Repeat("x", 1024)
		bulkString := resp.BulkString(largeData)

		encoded := bulkString.MarshalRESP()
		parser := resp.NewParser(bytes.NewReader(encoded))
		decoded, err := parser.Parse()

		require.NoError(t, err)
		require.Equal(t, bulkString, decoded)
	})

	t.Run("Large array", func(t *testing.T) {
		// 创建一个包含多个元素的数组
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

		// 验证几个元素
		require.Equal(t, resp.BulkString("data"), decodedArray[0])
		require.Equal(t, resp.BulkString(strings.Repeat("data", 50)), decodedArray[49])
		require.Equal(t, resp.BulkString(strings.Repeat("data", 100)), decodedArray[99])
	})
}

// 测试边界情况
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
