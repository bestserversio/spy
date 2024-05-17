package config

func (cfg *Config) LoadDefaults() {
	cfg.Api.Host = "localhost"
	cfg.Api.EndPoint = "/api"
	cfg.Api.Authorization = ""
	cfg.Api.Ssl = false
}
