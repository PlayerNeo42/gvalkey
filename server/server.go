package server

import (
	"log/slog"
	"net"

	"github.com/PlayerNeo42/gvalkey/handler"
)

type Server struct {
	addr   string
	logger *slog.Logger
}

func NewServer(addr string, opts ...Option) *Server {
	s := &Server{
		addr:   addr,
		logger: slog.Default(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) ListenAndServe() error {
	s.logger.Info("server started", "addr", s.addr)
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error("accept connection failed", "error", err)
			continue
		}

		s.logger.Info("new connection", "remote_addr", conn.RemoteAddr().String())
		handler := handler.New(conn, s.logger)
		go handler.Serve()
	}
}
