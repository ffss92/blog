package logging

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func NewLogger(level slog.Leveler, dev bool) *slog.Logger {
	var h slog.Handler
	if dev {
		h = tint.NewHandler(os.Stderr, &tint.Options{Level: level})
	} else {
		h = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	}
	return slog.New(h)
}
