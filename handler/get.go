package handler

import (
	"fmt"

	"github.com/PlayerNeo42/gvalkey/resp"
)

func (h *Handler) handleGet(args resp.Array) (resp.Payload, error) {
	key, err := resp.ParseGetArgs(args)
	if err != nil {
		return nil, err
	}

	value, ok := h.store.Get(key.String())
	if !ok {
		return resp.NULL, nil
	}

	if value == nil {
		return resp.NULL, nil
	}

	if val, ok := value.(resp.Payload); ok {
		return val, nil
	}

	return nil, fmt.Errorf("value is not a valid type: %T", value)
}
