package resp

import (
	"errors"
	"fmt"
	"strconv"
)

func peekNextInteger(args Array, index int) (int64, error) {
	nextIndex := index + 1
	if nextIndex >= len(args) {
		return 0, errors.New("argument required")
	}
	next, ok := args[nextIndex].(BulkString)
	if !ok {
		return 0, fmt.Errorf("value is not an integer: %T", args[nextIndex])
	}
	val, err := strconv.ParseInt(string(next), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("value is not an integer: %w", err)
	}
	return val, nil
}
