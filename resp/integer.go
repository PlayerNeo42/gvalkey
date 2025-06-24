package resp

import "fmt"

type Integer int64

func (i Integer) MarshalRESP() []byte {
	return fmt.Appendf(nil, ":%d\r\n", i)
}

func (i Integer) MarshalBinary() []byte {
	return fmt.Appendf(nil, "%d", i)
}
