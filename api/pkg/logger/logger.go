package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func Print(c *gin.Context) *logrus.Entry {
	logger.Formatter = &logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
		DisableColors:   false,
		ForceColors:     true,
	}

	if c == nil {
		return logger.WithFields(logrus.Fields{})
	}

	return logger.WithFields(logrus.Fields{
		"remote_ip": c.Request.RemoteAddr,
		"method":    c.Request.Method,
		"uri":       c.Request.URL.String(),
	})
}
