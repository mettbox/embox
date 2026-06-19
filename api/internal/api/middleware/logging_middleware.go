package middleware

import (
	"embox/internal/config"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

func LoggingMiddleware(cfg *config.RouterConfig) gin.HandlerFunc {
	var out io.Writer = os.Stdout
	if cfg.LogOutput != "stdout" {
		out = &lumberjack.Logger{
			Filename:   cfg.LogOutput,
			MaxSize:    cfg.LogMaxSize,
			MaxBackups: cfg.LogMaxBackups,
			MaxAge:     cfg.LogMaxAge,
			Compress:   cfg.LogCompress,
		}
	}
	logger := slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelError}))

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if c.Writer.Status() >= 400 {
			if raw != "" {
				path = path + "?" + raw
			}
			logger.Error("request failed",
				"status", c.Writer.Status(),
				"method", c.Request.Method,
				"path", path,
				"latency", time.Since(start),
				"ip", c.ClientIP(),
			)
		}
	}
}
