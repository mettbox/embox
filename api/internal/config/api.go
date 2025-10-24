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
		Csrf:    LoadCsrfConfig(server.Domain),
		Auth:    LoadAuthConfig(server.Domain),
		Storage: LoadStorageConfig(),
		Email:   LoadEmailConfig(),
	}
}
