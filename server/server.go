package server

import (
	"log/slog"
	"net"

	"github.com/PlayerNeo42/gvalkey/handler"
	"github.com/PlayerNeo42/gvalkey/store"
	"github.com/PlayerNeo42/gvalkey/store/naive"
)

type Server struct {
	addr    string
	logger  *slog.Logger
	storage store.Store
	handler *handler.Handler
}

func NewServer(addr string, opts ...Option) *Server {
	// storage := eventloop.NewEventloopStore()
	storage := naive.NewNaiveStore()

	s := &Server{
		addr:    addr,
		storage: storage,
		logger:  slog.New(slog.DiscardHandler),
	}
	for _, opt := range opts {
		opt(s)
	}

	s.handler = handler.New(s.logger, s.storage)

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
