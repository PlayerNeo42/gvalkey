package handler

import (
	"errors"
	"fmt"

	"github.com/PlayerNeo42/gvalkey/resp"
)

func (h *Handler) dispatch(args resp.Array) (resp.Marshaler, error) {
	val, ok := args[0].(resp.BulkString)
	if !ok {
		return resp.NULL, errors.New("command must be a bulk string")
	}

	var cmd *Command

	cmd, ok = h.commandTable.Get(val.Upper())
	if !ok {
		return nil, errors.New("unsupported command")
	}

	// positive value means fixed number of arguments
	// negative value means at least that number of arguments
	if (cmd.Args > 0 && len(args) != cmd.Args) ||
		(cmd.Args < 0 && len(args) < -cmd.Args) {
		return nil, fmt.Errorf("ERR wrong number of arguments for '%s' command", cmd.Name)
	}

	return cmd.Handler(args)
}
