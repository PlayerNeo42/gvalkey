package resp

import "fmt"

func NewErrorMessage(message string) []byte {
	return fmt.Appendf(ERR, "%s\r\n", message)
}
