package config

func (cfg *Config) LoadDefaults() {
	cfg.Host = "localhost"
	cfg.EndPoint = "/api"
	cfg.Authorization = ""
	cfg.Ssl = false
}
