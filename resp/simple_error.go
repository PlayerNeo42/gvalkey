package resp

import "fmt"

type SimpleError struct {
	message string
}

func NewSimpleError(message string) SimpleError {
	return SimpleError{message: message}
}

func (e SimpleError) MarshalRESP() []byte {
	return fmt.Appendf(nil, "-ERR %s\r\n", e.message)
}
