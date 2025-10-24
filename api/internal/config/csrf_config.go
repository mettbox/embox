package config

import (
	"embox/pkg/env"
	"strings"
)

type CsrfConfig struct {
	Secret   string
	MaxAge   int
	Domain   string
	IsSecure bool
}

func LoadCsrfConfig(domain string) *CsrfConfig {
	isSecure := strings.HasPrefix(domain, "https")

	return &CsrfConfig{
		Secret:   env.GetEnv("CSRF_SECRET", "csrfSecret"),
		MaxAge:   env.GetEnvAsInt("CSRF_MAXAGE", 60*60),
		Domain:   domain,
		IsSecure: isSecure,
	}
}
