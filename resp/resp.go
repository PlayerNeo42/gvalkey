// Package resp provides RESP (Redis Serialization Protocol) types and utilities.
package resp

// Marshaler is used to marshal a value to a RESP-encoded byte slice.
type Marshaler interface {
	MarshalRESP() []byte
}

// BinaryMarshaler is used to marshal a value to a binary byte slice.
type BinaryMarshaler interface {
	MarshalBinary() []byte
}
