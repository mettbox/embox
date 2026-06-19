package config

type ApiConfig struct {
	Server  *ServerConfig
	Router  *RouterConfig
	Csrf    *CsrfConfig
	Auth    *AuthConfig
	Storage *StorageConfig
	Email   *EmailConfig
}

func LoadApiConfig() *ApiConfig {
	server := LoadServerConfig()
	return &ApiConfig{
		Server:  server,
		Router:  LoadRouterConfig(),
		Csrf:    LoadCsrfConfig(server.Domain, server.IsSecure),
		Auth:    LoadAuthConfig(server.Domain, server.IsSecure),
		Storage: LoadStorageConfig(),
		Email:   LoadEmailConfig(),
	}
}
