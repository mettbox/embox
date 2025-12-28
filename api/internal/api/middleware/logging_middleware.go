package middleware

import (
	"embox/internal/config"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

func createLogFile(cfg *config.RouterConfig) io.Writer {
	return &lumberjack.Logger{
		Filename:   cfg.LogOutput,
		MaxSize:    cfg.LogMaxSize,
		MaxBackups: cfg.LogMaxBackups,
		MaxAge:     cfg.LogMaxAge,
		Compress:   cfg.LogCompress,
	}
}

func LoggingMiddleware(cfg *config.RouterConfig) gin.HandlerFunc {
	var writer io.Writer = os.Stdout
	if cfg.LogOutput != "stdout" {
		writer = createLogFile(cfg)
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Only log errors and warnings (status >= 400)
		if c.Writer.Status() >= 400 {
			end := time.Now()
			latency := end.Sub(start)
			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()

			if raw != "" {
				path = path + "?" + raw
			}

			fmt.Fprintf(writer, "[GIN] %v | %3d | %13v | %15s | %-7s %#v\n%s",
				end.Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
				c.Errors.ByType(gin.ErrorTypePrivate).String(),
			)
		}
	}
}
