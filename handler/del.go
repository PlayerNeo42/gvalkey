package handler

import (
	"github.com/PlayerNeo42/gvalkey/resp"
)

func (h *Handler) handleDel(args resp.Array) (resp.Payload, error) {
	keys, err := resp.ParseDelArgs(args)
	if err != nil {
		return nil, err
	}

	count := 0
	for _, key := range keys {
		if h.store.Del(key.String()) {
			count++
		}
	}
	return resp.Integer(count), nil
}
