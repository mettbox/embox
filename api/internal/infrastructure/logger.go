package infrastructure

import (
	"log/slog"
	"os"
)

// InitLogger configures the global slog logger.
// "release" mode uses JSON output (production); all other modes use human-readable text with debug level.
func InitLogger(mode string) *slog.Logger {
	var handler slog.Handler
	if mode == "release" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
