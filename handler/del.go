package handler

import (
	"github.com/PlayerNeo42/gvalkey/resp"
)

func (h *Handler) handleDel(args resp.Array) (resp.Marshaler, error) {
	keys, err := resp.ParseDelArgs(args)
	if err != nil {
		return nil, err
	}

	count := 0
	for _, key := range keys {
		if h.store.Del(string(key.MarshalBinary())) {
			count++
		}
	}
	return resp.Integer(count), nil
}
