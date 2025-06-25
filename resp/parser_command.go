package resp

import (
	"errors"
	"fmt"
)

func ParseGetArgs(args Array) (BinaryMarshaler, error) {
	key, ok := args[1].(BinaryMarshaler)
	if !ok {
		return nil, errors.New("key is not a binary marshaler")
	}
	return key, nil
}

func ParseSetArgs(args Array) (*SetArgs, error) {
	key, ok := args[1].(BinaryMarshaler)
	if !ok {
		return nil, errors.New("key is not a binary marshaler")
	}

	value := args[2]

	parsedArgs := &SetArgs{
		Key:   key,
		Value: value,
	}

	length := len(args)

	// set without options
	if length <= 3 {
		return parsedArgs, nil
	}

	for i := 3; i < length; i++ {
		option, ok := args[i].(BulkString)
		if !ok {
			return nil, fmt.Errorf("option is not a bulk string: %T", args[i])
		}
		switch option.Upper() {
		case EX:
			ex, err := peekNextInteger(args, i)
			if err != nil {
				return nil, fmt.Errorf("syntax error: %w", err)
			}
			parsedArgs.EX = &ex
			// skip the next argument
			i++
		case PX:
			px, err := peekNextInteger(args, i)
			if err != nil {
				return nil, fmt.Errorf("syntax error: %w", err)
			}
			parsedArgs.PX = &px
			// skip the next argument
			i++
		case NX:
			parsedArgs.NX = true
		case XX:
			parsedArgs.XX = true
		case GET:
			parsedArgs.Get = true
		default:
			return nil, fmt.Errorf("syntax error: unsupported option '%s'", option.Upper())
		}
	}

	if parsedArgs.NX && parsedArgs.XX {
		return nil, errors.New("syntax error: NX and XX options cannot be used together")
	}

	if parsedArgs.EX != nil && parsedArgs.PX != nil {
		return nil, errors.New("syntax error: EX and PX options cannot be used together")
	}

	if parsedArgs.EX != nil && *parsedArgs.EX < 0 {
		return nil, errors.New("syntax error: EX option must be greater than 0")
	}

	if parsedArgs.PX != nil && *parsedArgs.PX < 0 {
		return nil, errors.New("syntax error: PX option must be greater than 0")
	}

	return parsedArgs, nil
}

func ParseDelArgs(args Array) ([]BinaryMarshaler, error) {
	keys := make([]BinaryMarshaler, len(args)-1)
	for i := 1; i < len(args); i++ {
		key, ok := args[i].(BinaryMarshaler)
		if !ok {
			return nil, errors.New("key is not a binary marshaler")
		}
		keys[i-1] = key
	}
	return keys, nil
}
