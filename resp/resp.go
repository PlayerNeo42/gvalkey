// Package resp provides RESP (Redis Serialization Protocol) types and utilities.
package resp

import "io"

// Payload is used to marshal a value to a RESP-encoded byte slice.
type Payload interface {
	RESPReader() io.Reader
}

type Stringer interface {
	String() string
}

// Bytes() []byte
// String() string
