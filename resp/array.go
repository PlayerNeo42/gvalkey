package resp

import "fmt"

type Array []any

func (a Array) MarshalRESP() []byte {
	buf := make([]byte, 0, 1024)

	length := len(a)
	buf = fmt.Appendf(buf, "*%d\r\n", length)

	for _, v := range a {
		marshaler, ok := v.(Marshaler)
		if !ok {
			return nil
		}

		buf = append(buf, marshaler.MarshalRESP()...)
	}

	return buf
}
