package config

import (
	"embox/pkg/env"
)

type EmailConfig struct {
	Host     string
	Port     int
	From     string
	Password string
}

func LoadEmailConfig() *EmailConfig {
	return &EmailConfig{
		Host:     env.GetEnv("GMAIL_HOST", "smtp.gmail.com"),
		Port:     env.GetEnvAsInt("GMAIL_PORT", 587),
		From:     env.GetEnv("GMAIL_FROM", ""),
		Password: env.GetEnv("GMAIL_APP_PASSWORD", ""),
	}
}
