package resp

import "fmt"

func peekNextInteger(args Array, index int) (int64, error) {
	nextIndex := index + 1
	if nextIndex >= len(args) {
		return 0, fmt.Errorf("argument required")
	}
	next, ok := args[nextIndex].(Integer)
	if !ok {
		return 0, fmt.Errorf("value is not an integer: %T", args[nextIndex])
	}
	return int64(next), nil
}
