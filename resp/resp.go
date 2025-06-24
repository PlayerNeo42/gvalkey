package resp

type Marshaler interface {
	MarshalRESP() []byte
}

type BinaryMarshaler interface {
	MarshalBinary() []byte
}
