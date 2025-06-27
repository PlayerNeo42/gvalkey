package resp

import (
	"time"
)

type SetArgs struct {
	Key      Stringer
	Value    any
	ExpireAt time.Time
	NX       bool
	XX       bool
	Get      bool
}
