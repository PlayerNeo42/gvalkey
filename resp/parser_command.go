package resp

import "fmt"

func ParseGetArgs(args Array) (key BinaryMarshaler, err error) {
	var ok bool
	key, ok = args[1].(BinaryMarshaler)
	if !ok {
		err = fmt.Errorf("key is not a binary marshaler")
		return
	}
	return key, nil
}

func ParseSetArgs(args Array) (key BinaryMarshaler, value any, ex, px int64, nx, xx, get bool, err error) {
	var ok bool
	key, ok = args[1].(BinaryMarshaler)
	if !ok {
		err = fmt.Errorf("key is not a binary marshaler")
		return
	}

	value = args[2]

	length := len(args)
	if length > 3 {
		for i := 3; i < length; i++ {
			option, ok := args[i].(BulkString)
			if !ok {
				err = fmt.Errorf("option is not a bulk string: %T", args[i])
				return
			}
			switch option.Upper() {
			case EX:
				ex, err = peekNextInteger(args, i)
				if err != nil {
					err = fmt.Errorf("syntax error: %w", err)
					return
				}
				// skip the next argument
				i++
			case PX:
				px, err = peekNextInteger(args, i)
				if err != nil {
					err = fmt.Errorf("syntax error: %w", err)
					return
				}
				// skip the next argument
				i++
			case NX:
				nx = true
			case XX:
				xx = true
			case GET:
				get = true
			default:
				err = fmt.Errorf("syntax error: unsupported option '%s'", option.Upper())
				return
			}
		}
	}

	if nx && xx {
		err = fmt.Errorf("syntax error: NX and XX options cannot be used together")
		return
	}

	if ex > 0 && px > 0 {
		err = fmt.Errorf("syntax error: EX and PX options cannot be used together")
		return
	}

	if ex < 0 || px < 0 {
		err = fmt.Errorf("syntax error: EX or PX option must be greater than 0")
		return
	}

	return key, value, ex, px, nx, xx, get, nil
}

func ParseDelArgs(args Array) (keys []BinaryMarshaler, err error) {
	length := len(args)
	for i := 1; i < length; i++ {
		key, ok := args[i].(BinaryMarshaler)
		if !ok {
			err = fmt.Errorf("key is not a binary marshaler")
			return
		}
		keys = append(keys, key)
	}
	return keys, nil
}
