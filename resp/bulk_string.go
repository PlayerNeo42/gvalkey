package resp

import (
	"bytes"
	"fmt"
)

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
