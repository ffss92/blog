package logging

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func NewLogger(dev bool) *slog.Logger {
	var h slog.Handler
	if dev {
		h = tint.NewHandler(os.Stderr, &tint.Options{Level: slog.LevelDebug})
	} else {
		h = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	return slog.New(h)
}
