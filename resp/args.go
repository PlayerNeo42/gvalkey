package resp

type SetArgs struct {
	Key    BinaryMarshaler
	Value  any
	Expire int64
	NX     bool
	XX     bool
	Get    bool
}
