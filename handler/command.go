package handler

import "github.com/PlayerNeo42/gvalkey/resp"

func (h *Handler) handleCommand(command resp.Array) (resp.Payload, error) {
	return resp.OK, nil
}
