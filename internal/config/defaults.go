package config

func (cfg *Config) LoadDefaults() {
	cfg.Api.Host = "http://localhost"
	cfg.Api.Authorization = ""
	cfg.Api.Timeout = 5

	cfg.WebApi.Host = "http://localhost"
	cfg.WebApi.Endpoint = "/api/spy/get"
	cfg.WebApi.Timeout = 5
	cfg.WebApi.Interval = 120

	cfg.Vms.Enabled = false
	cfg.Vms.MinWait = 60
	cfg.Vms.MaxWait = 180
	cfg.Vms.Timeout = 5
	cfg.Vms.Limit = 100
	cfg.Vms.ExcludeEmpty = true
	cfg.Vms.SubBots = true
	cfg.Vms.AddOnly = true
	cfg.Vms.SetOffline = true

	cfg.RemoveInactive.Enabled = false
	cfg.RemoveInactive.InactiveTime = 2592000
	cfg.RemoveInactive.Interval = 86400
	cfg.RemoveInactive.Timeout = 5
}
