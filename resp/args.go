package resp

type SetArgs struct {
	Key   BinaryMarshaler
	Value any
	EX    *int64
	PX    *int64
	NX    bool
	XX    bool
	Get   bool
}
