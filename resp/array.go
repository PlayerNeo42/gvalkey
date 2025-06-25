package resp

import (
	"strconv"
)

type Array []any

func (a Array) MarshalRESP() []byte {
	buf := make([]byte, 0, 1024)

	buf = append(buf, '*')
	buf = strconv.AppendInt(buf, int64(len(a)), 10)
	buf = append(buf, '\r', '\n')

	for _, v := range a {
		marshaler, ok := v.(Marshaler)
		if !ok {
			return nil
		}

		buf = append(buf, marshaler.MarshalRESP()...)
	}

	return buf
}
