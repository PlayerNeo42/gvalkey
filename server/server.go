package server

import (
	"io"
	"log/slog"
	"net"

	"github.com/PlayerNeo42/gvalkey/handler"
	"github.com/PlayerNeo42/gvalkey/store"
)

type Server struct {
	addr    string
	logger  *slog.Logger
	store   *store.Store
	handler *handler.Handler
}

func NewServer(addr string, opts ...Option) *Server {
	s := &Server{
		addr:  addr,
		store: store.NewStore(),
	}
	s.handler = handler.New(s.logger, s.store)
	for _, opt := range opts {
		opt(s)
	}

	// disable logging if not set
	if s.logger == nil {
		s.logger = slog.New(slog.NewTextHandler(io.Discard, nil))
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
		go s.handler.Serve(conn)
	}
}
