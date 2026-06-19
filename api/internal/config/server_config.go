package config

import (
	"embox/pkg/env"
	"log"
	"time"
)

type ServerConfig struct {
	Host        string
	Port        string
	Domain      string
	CorsOrigins []string
	CorsMaxAge  time.Duration
}

func LoadServerConfig() *ServerConfig {
	maxAgeSeconds := env.GetEnvAsInt("SERVER_CORS_MAXAGE", 43200) // 12h = 43200s

	corsOrigins := env.GetEnvSlice("SERVER_CORS_ORIGINS", []string{})
	if len(corsOrigins) == 0 {
		log.Println("WARNING: SERVER_CORS_ORIGINS is not set — all cross-origin requests will be rejected")
	}

	return &ServerConfig{
		Host:        env.GetEnv("SERVER_HOST", "0.0.0.0"),
		Port:        env.GetEnv("SERVER_PORT", "2705"),
		Domain:      env.GetEnv("SERVER_DOMAIN", "localhost"),
		CorsOrigins: corsOrigins,
		CorsMaxAge:  time.Duration(maxAgeSeconds) * time.Second,
	}
}
