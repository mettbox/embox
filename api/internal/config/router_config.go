package config

import (
	"embox/pkg/env"
)

type RouterConfig struct {
	ReleaseMode   string
	LogOutput     string
	LogMaxSize    int
	LogMaxBackups int
	LogMaxAge     int
	LogCompress   bool
	RateLimit     int
}

func LoadRouterConfig() *RouterConfig {
	return &RouterConfig{
		ReleaseMode:   env.GetEnv("ROUTER_RUNTIME", "release"),
		LogOutput:     env.GetEnv("ROUTER_LOG_OUTPUT", "stdout"),
		LogMaxSize:    env.GetEnvAsInt("ROUTER_LOG_MAX_SIZE", 100),
		LogMaxBackups: env.GetEnvAsInt("ROUTER_LOG_MAX_BACKUPS", 7),
		LogMaxAge:     env.GetEnvAsInt("ROUTER_LOG_MAX_AGE", 30),
		LogCompress:   env.GetEnvAsBool("ROUTER_LOG_COMPRESS", true),
		RateLimit:     env.GetEnvAsInt("ROUTER_AUTH_RATE_LIMIT", 10),
	}
}
