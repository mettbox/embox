package config

import (
	"embox/pkg/env"
	"log/slog"
	"time"
)

type ServerConfig struct {
	Host        string
	Port        string
	Domain      string
	IsSecure    bool
	CorsOrigins []string
	CorsMaxAge  time.Duration
}

func LoadServerConfig() *ServerConfig {
	maxAgeSeconds := env.GetEnvAsInt("SERVER_CORS_MAXAGE", 43200) // 12h = 43200s

	corsOrigins := env.GetEnvSlice("SERVER_CORS_ORIGINS", []string{})
	if len(corsOrigins) == 0 {
		slog.Warn("SERVER_CORS_ORIGINS not set — cross-origin requests will be rejected")
	}

	return &ServerConfig{
		Host:        env.GetEnv("SERVER_HOST", "0.0.0.0"),
		Port:        env.GetEnv("SERVER_PORT", "2705"),
		Domain:      env.GetEnv("SERVER_DOMAIN", "localhost"),
		IsSecure:    env.GetEnvAsBool("SERVER_SECURE", false),
		CorsOrigins: corsOrigins,
		CorsMaxAge:  time.Duration(maxAgeSeconds) * time.Second,
	}
}
