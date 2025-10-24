package config

import (
	"embox/pkg/env"
)

type DbConfig struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPass       string
	DBName       string
	DBSystemUser string
}

func LoadDbConfig() *DbConfig {
	return &DbConfig{
		DBHost:       env.GetEnv("DB_HOST", "localhost"),
		DBPort:       env.GetEnv("DB_PORT", "3306"),
		DBUser:       env.GetEnv("DB_USER", "embox"),
		DBPass:       env.GetEnv("DB_PASSWORD", "embox"),
		DBName:       env.GetEnv("DB_NAME", "embox"),
		DBSystemUser: env.GetEnv("DB_SYSTEM_USER", ""),
	}
}
