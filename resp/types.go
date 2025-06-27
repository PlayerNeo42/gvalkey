package resp

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Array
type Array []any

func (a Array) RESPReader() io.Reader {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	buf.WriteByte('*')
	buf.WriteString(strconv.FormatInt(int64(len(a)), 10))
	buf.WriteByte('\r')
	buf.WriteByte('\n')

	for _, v := range a {
		marshaler, ok := v.(Payload)
		if !ok {
			return nil
		}

		if _, err := io.Copy(buf, marshaler.RESPReader()); err != nil {
			return nil
		}
	}

	return buf
}

func (a Array) Bytes() []byte {
	return nil
}

func (a Array) String() string {
	return ""
}

// SimpleString
type SimpleString string

func (s SimpleString) RESPReader() io.Reader {
	return bytes.NewReader(fmt.Appendf(nil, "+%s\r\n", s))
}

func (s SimpleString) Bytes() []byte {
	return []byte(s)
}

func (s SimpleString) String() string {
	return string(s)
}

// SimpleError
type SimpleError struct {
	message string
}

func NewSimpleError(message string) SimpleError {
	return SimpleError{message: message}
}

func (e SimpleError) RESPReader() io.Reader {
	return bytes.NewReader(fmt.Appendf(nil, "-ERR %s\r\n", e.message))
}

func (e SimpleError) Bytes() []byte {
	return []byte(e.message)
}

func (e SimpleError) String() string {
	return e.message
}

// BulkString
type BulkString string

func (b BulkString) RESPReader() io.Reader {
	return bytes.NewReader(fmt.Appendf(nil, "$%d\r\n%s\r\n", len(b), b))
}

func (b BulkString) Bytes() []byte {
	return []byte(b)
}

func (b BulkString) String() string {
	return string(b)
}

func (b BulkString) Upper() BulkString {
	return BulkString(strings.ToUpper(string(b)))
}

// Integer
type Integer int64

func (i Integer) RESPReader() io.Reader {
	return bytes.NewReader(fmt.Appendf(nil, ":%d\r\n", i))
}

func (i Integer) Bytes() []byte {
	return []byte(strconv.FormatInt(int64(i), 10))
}

func (i Integer) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// Null
type Null struct{}

func (n Null) RESPReader() io.Reader {
	return bytes.NewReader([]byte("$-1\r\n"))
}

func (n Null) Bytes() []byte {
	return nil
}

func (n Null) String() string {
	return ""
}
