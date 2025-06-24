package main

import (
	"log/slog"

	"github.com/PlayerNeo42/gvalkey/internal/log"
	"github.com/PlayerNeo42/gvalkey/server"
)

func main() {
	logger := log.New(slog.LevelInfo)

	tcpServer := server.NewServer(":6379", server.WithLogger(logger))
	panic(tcpServer.ListenAndServe())
}
