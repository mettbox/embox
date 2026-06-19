package config

import (
	"embox/pkg/env"
)

type CsrfConfig struct {
	Secret   string
	MaxAge   int
	Domain   string
	IsSecure bool
}

func LoadCsrfConfig(domain string, isSecure bool) *CsrfConfig {
	return &CsrfConfig{
		Secret:   env.GetEnv("CSRF_SECRET", "csrfSecret"),
		MaxAge:   env.GetEnvAsInt("CSRF_MAXAGE", 60*60),
		Domain:   domain,
		IsSecure: isSecure,
	}
}
