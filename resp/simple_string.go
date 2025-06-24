package resp

import "fmt"

type SimpleString string

func (s SimpleString) MarshalRESP() []byte {
	return fmt.Appendf(nil, "+%s\r\n", s)
}
