package config

import (
	"embox/pkg/env"
	"strings"
)

type AuthConfig struct {
	Domain               string
	AccessExpiration     int
	RefreshExpiration    int
	AccessSecret         string
	RefreshSecret        string
	IsSecure             bool
	LoginTokenExpiration int
	LoginEmailSubject    string
	LoginEmailTemplate   string
}

func LoadAuthConfig(domain string) *AuthConfig {
	isSecure := strings.HasPrefix(domain, "https")

	return &AuthConfig{
		Domain:               domain,
		AccessExpiration:     env.GetEnvAsInt("AUTH_ACCESS_COOKIE_EXPIRATION", 60*60),       // 1 hour
		RefreshExpiration:    env.GetEnvAsInt("AUTH_REFRESH_COOKIE_EXPIRATION", 5*60*60*24), // 5 days
		AccessSecret:         env.GetEnv("AUTH_ACCESS_JWT_SECRET", "emboxAccessSecret"),
		RefreshSecret:        env.GetEnv("AUTH_REFRESH_JWT_SECRET", "emboxRefreshSecret"),
		IsSecure:             isSecure,
		LoginTokenExpiration: env.GetEnvAsInt("AUTH_LOGIN_TOKEN_EXPIRATION", 10*60), // 10 minutes
		LoginEmailSubject:    env.GetEnv("TOKEN_EMAIL_SUBJECT", "Your Login Token"),
		LoginEmailTemplate:   env.GetEnv("TOKEN_EMAIL_TEMPLATE", "<p>Hello %s,</p><p>here is your login token:</p><p><b>%s</b></p><p>The token is valid for %d minutes.</p>"),
	}
}
