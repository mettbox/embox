package config

import (
	"embox/pkg/env"
)

type StorageConfig struct {
	Url      string
	Username string
	Password string
	RepoID   string
}

func LoadStorageConfig() *StorageConfig {
	return &StorageConfig{
		Url:      env.GetEnv("STORAGE_URL", "https://sync.luckycloud.de/api2"),
		Username: env.GetEnv("STORAGE_USERNAME", ""),
		Password: env.GetEnv("STORAGE_PASSWORD", ""),
		RepoID:   env.GetEnv("STORAGE_REPO_ID", ""),
	}
}
