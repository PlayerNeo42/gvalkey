package log

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

func New(level string) *slog.Logger {
	var slogLevel slog.Level
	switch strings.ToUpper(level) {
	case "DEBUG":
		slogLevel = slog.LevelDebug
	case "INFO":
		slogLevel = slog.LevelInfo
	case "WARN":
		slogLevel = slog.LevelWarn
	case "ERROR":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	out := os.Stdout
	var handler slog.Handler
	if isatty.IsTerminal(out.Fd()) {
		handler = tint.NewHandler(colorable.NewColorable(out), &tint.Options{
			Level:      slogLevel,
			TimeFormat: time.DateTime,
		})
	} else {
		handler = slog.NewJSONHandler(out, &slog.HandlerOptions{
			Level: slogLevel,
		})
	}

	return slog.New(handler)
}
