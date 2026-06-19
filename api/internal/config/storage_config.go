package config

import (
	"embox/pkg/env"
)

type StorageConfig struct {
	Adapter  string // "luckycloud" | "local"
	LocalDir string
	Url      string
	Username string
	Password string
	RepoID   string
}

func LoadStorageConfig() *StorageConfig {
	return &StorageConfig{
		Adapter:  env.GetEnv("STORAGE_ADAPTER", "luckycloud"),
		LocalDir: env.GetEnv("STORAGE_LOCAL_DIR", "./dev-storage"),
		Url:      env.GetEnv("STORAGE_URL", "https://sync.luckycloud.de/api2"),
		Username: env.GetEnv("STORAGE_USERNAME", ""),
		Password: env.GetEnv("STORAGE_PASSWORD", ""),
		RepoID:   env.GetEnv("STORAGE_REPO_ID", ""),
	}
}
