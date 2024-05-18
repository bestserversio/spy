package config

func (cfg *Config) LoadDefaults() {
	cfg.Api.Host = "http://localhost"
	cfg.Api.Authorization = ""
	cfg.Api.Ssl = false
	cfg.Api.Timeout = 5

	cfg.Vms.Enabled = false
	cfg.Vms.Interval = 120
	cfg.Vms.Timeout = 5

	cfg.Scanner.MinWait = 10
	cfg.Scanner.MaxWait = 15
}
