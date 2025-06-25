package main

import (
	"fmt"
	"os"

	"github.com/PlayerNeo42/gvalkey/internal/config"
	"github.com/PlayerNeo42/gvalkey/internal/log"
	"github.com/PlayerNeo42/gvalkey/server"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config:", err)
		os.Exit(1)
	}

	logger := log.New(conf.LogLevel)

	tcpServer := server.NewServer(fmt.Sprintf("%s:%d", conf.Host, conf.Port), server.WithLogger(logger))
	if err := tcpServer.ListenAndServe(); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
