package resp

import "time"

type SetArgs struct {
	Key      BinaryMarshaler
	Value    any
	ExpireAt time.Time
	NX       bool
	XX       bool
	Get      bool
}
