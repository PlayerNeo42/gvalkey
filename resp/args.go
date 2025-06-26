package resp

type SetArgs struct {
	Key    BinaryMarshaler
	Value  any
	Expire int64 // in milliseconds
	NX     bool
	XX     bool
	Get    bool
}
