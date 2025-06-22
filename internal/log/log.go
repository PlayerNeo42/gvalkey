package log

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

func New(level slog.Level) *slog.Logger {
	out := os.Stdout
	var handler slog.Handler
	if isatty.IsTerminal(out.Fd()) {
		handler = tint.NewHandler(colorable.NewColorable(out), &tint.Options{
			Level:      level,
			TimeFormat: time.DateTime,
		})
	} else {
		handler = slog.NewJSONHandler(out, &slog.HandlerOptions{
			Level: level,
		})
	}

	return slog.New(handler)
}
