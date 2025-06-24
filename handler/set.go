package handler

import (
	"github.com/PlayerNeo42/gvalkey/resp"
)

func (h *Handler) handleSet(args resp.Array) (resp.Marshaler, error) {
	key, value, ex, px, nx, xx, get, err := resp.ParseSetArgs(args)
	if err != nil {
		return nil, err
	}

	oldValue, success := h.store.Set(string(key.MarshalBinary()), value, ex, px, nx, xx, get)

	if get {
		if success && oldValue != nil {
			if val, ok := oldValue.(resp.Marshaler); ok {
				return val, nil
			}
			// 如果不是 Marshaler，包装成 BulkString
			return resp.BulkString(oldValue.(string)), nil
		}
		return resp.NULL, nil
	}

	// 对于 NX 和 XX 选项，如果操作失败，返回 NULL
	if (nx || xx) && !success {
		return resp.NULL, nil
	}

	return resp.OK, nil
}
