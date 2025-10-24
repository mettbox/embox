package middleware

import (
	"embox/internal/config"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/gin-gonic/gin"
)

func createLogFile(cfg *config.RouterConfig) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   cfg.LogOutput,
		MaxSize:    cfg.LogMaxSize,
		MaxBackups: cfg.LogMaxBackups,
		MaxAge:     cfg.LogMaxAge,
		Compress:   cfg.LogCompress,
	}
}

func LoggingMiddleware(cfg *config.RouterConfig) gin.HandlerFunc {
	if cfg.LogOutput == "stdout" {
		return gin.Logger()
	}
	logFile := createLogFile(cfg)
	return gin.LoggerWithWriter(logFile)
}
