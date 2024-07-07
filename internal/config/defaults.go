package config

func (cfg *Config) LoadDefaults() {
	cfg.Api.Host = "http://localhost"
	cfg.Api.Authorization = ""
	cfg.Api.Timeout = 5

	cfg.Vms.Enabled = false
	cfg.Vms.MinWait = 60
	cfg.Vms.MaxWait = 180
	cfg.Vms.Timeout = 5
	cfg.Vms.Limit = 100
	cfg.Vms.ExcludeEmpty = true
	cfg.Vms.SubBots = true
}
