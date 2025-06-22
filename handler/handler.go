package handler

import (
	"io"
	"log/slog"
	"net"

	"github.com/PlayerNeo42/gvalkey/resp"
)

type Handler struct {
	conn   net.Conn
	logger *slog.Logger
}

func New(conn net.Conn, logger *slog.Logger) *Handler {
	return &Handler{conn: conn, logger: logger}
}

func (h *Handler) Serve() {
	defer h.conn.Close()

	parser := resp.NewParser(h.conn)

	for {
		value, err := parser.Parse()
		if err != nil {
			if err == io.EOF {
				h.logger.Info("client closed connection", "remote_addr", h.conn.RemoteAddr().String())
				return
			}
			h.logger.Error("parse command failed", "error", err)
			_, err = h.conn.Write(resp.NewErrorMessage(err.Error()))
			if err != nil {
				h.logger.Error("write error message to client failed", "error", err)
			}
			return
		}

		command, ok := value.([]string)
		if !ok {
			h.logger.Error("command is not a string array")
			_, err = h.conn.Write(resp.NewErrorMessage("command must be an array"))
			if err != nil {
				h.logger.Error("write error message to client failed", "error", err)
			}
			continue
		}

		h.logger.Debug("received command", "remote_addr", h.conn.RemoteAddr().String(), "command", command)

		// TODO: 暂时先简单回复一个 OK
		_, err = h.conn.Write(resp.OK)
		if err != nil {
			h.logger.Error("write ok message to client failed", "error", err)
		}
	}
}
