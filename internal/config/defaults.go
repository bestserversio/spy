package config

func (cfg *Config) LoadDefaults() {
	cfg.Api.Host = "localhost"
	cfg.Api.Authorization = ""
	cfg.Api.Ssl = false
	cfg.Api.Timeout = 5

	cfg.Vms.Enabled = false
	cfg.Vms.Timeout = 5
}
