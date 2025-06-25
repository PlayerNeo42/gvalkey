package resp

import (
	"bytes"
	"fmt"
	"strconv"
)

// Array
type Array []any

func (a Array) MarshalRESP() []byte {
	buf := make([]byte, 0, 1024)

	buf = append(buf, '*')
	buf = strconv.AppendInt(buf, int64(len(a)), 10)
	buf = append(buf, '\r', '\n')

	for _, v := range a {
		marshaler, ok := v.(Marshaler)
		if !ok {
			return nil
		}

		buf = append(buf, marshaler.MarshalRESP()...)
	}

	return buf
}

// SimpleString
type SimpleString string

func (s SimpleString) MarshalRESP() []byte {
	return fmt.Appendf(nil, "+%s\r\n", s)
}

// SimpleError
type SimpleError struct {
	message string
}

func NewSimpleError(message string) SimpleError {
	return SimpleError{message: message}
}

func (e SimpleError) MarshalRESP() []byte {
	return fmt.Appendf(nil, "-ERR %s\r\n", e.message)
}

// BulkString
type BulkString string

func (b BulkString) MarshalRESP() []byte {
	return fmt.Appendf(nil, "$%d\r\n%s\r\n", len(b), b)
}

func (b BulkString) MarshalBinary() []byte {
	return []byte(b)
}

func (b BulkString) Upper() BulkString {
	return BulkString(bytes.ToUpper([]byte(b)))
}

// Integer
type Integer int64

func (i Integer) MarshalRESP() []byte {
	return fmt.Appendf(nil, ":%d\r\n", i)
}

func (i Integer) MarshalBinary() []byte {
	return fmt.Appendf(nil, "%d", i)
}

// Null
type Null struct{}

func (n Null) MarshalRESP() []byte {
	return []byte("$-1\r\n")
}
