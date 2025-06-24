package handler

import (
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/PlayerNeo42/gvalkey/resp"
	"github.com/PlayerNeo42/gvalkey/store"
)

type Handler struct {
	logger       *slog.Logger
	store        *store.Store
	commandTable *CommandTable
}

func New(logger *slog.Logger, store *store.Store) *Handler {
	commandTable := NewCommandTable()
	h := &Handler{logger: logger, store: store, commandTable: commandTable}

	commandTable.MustRegister(&Command{resp.GET, 2, h.handleGet})
	commandTable.MustRegister(&Command{resp.SET, -3, h.handleSet})
	commandTable.MustRegister(&Command{resp.DEL, -2, h.handleDel})
	commandTable.MustRegister(&Command{resp.COMMAND, -1, h.handleCommand})

	return h
}

func (h *Handler) Serve(conn net.Conn) {
	defer conn.Close()

	parser := resp.NewParser(conn)

	for {
		value, err := parser.Parse()
		if err != nil {
			if err == io.EOF {
				h.logger.Info("client closed connection", "remote_addr", conn.RemoteAddr().String())
				return
			}
			h.logger.Error("parse command failed", "error", err)
			_, err = conn.Write(resp.NewSimpleError(err.Error()).MarshalRESP())
			if err != nil {
				h.logger.Error("write error message to client failed", "error", err)
			}
			return
		}

		var response resp.Marshaler
		var commandErr error

		switch v := value.(type) {
		case resp.Array:
			h.logger.Debug("received array command", "remote_addr", conn.RemoteAddr().String(), "command", v)
			response, commandErr = h.dispatch(v)
		default:
			h.logger.Error("unsupported command type", "remote_addr", conn.RemoteAddr().String(), "command", v)
			commandErr = fmt.Errorf("command must be an array")
		}

		if commandErr != nil {
			response = resp.NewSimpleError(commandErr.Error())
		}

		h.logger.Debug("writing response", "remote_addr", conn.RemoteAddr().String(), "response", response, "payload", fmt.Sprintf("%q", response.MarshalRESP()))
		_, err = conn.Write(response.MarshalRESP())
		if err != nil {
			h.logger.Error("write ok message to client failed", "error", err)
		}
	}
}
