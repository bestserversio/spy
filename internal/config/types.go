package config

type API struct {
	Host          string `json:"host"`
	Authorization string `json:"authorization"`
	Timeout       int    `json:"timeout"`
}

type VMS struct {
	Enabled      bool   `json:"enabled"`
	Timeout      int    `json:"timeout"`
	ApiToken     string `json:"api_token"`
	AppIds       []int  `json:"app_ids"`
	RecvOnly     bool   `json:"recv_only"`
	MinWait      int    `json:"min_wait"`
	MaxWait      int    `json:"max_wait"`
	Limit        int    `json:"limit"`
	ExcludeEmpty bool   `json:"hide_empty"`
	SubBots      bool   `json:"sub_bots"`
}

type Scanner struct {
	Protocol    string `json:"protocol"`
	PlatformIds []int  `json:"platform_ids"`
	MinWait     int    `json:"min_wait"`
	MaxWait     int    `json:"max_wait"`
	Limit       int    `json:"limit"`
	RecvOnly    bool   `json:"recv_only"`
	SubBots     bool   `json:"sub_bots"`
}

type PlatformMapper struct {
	AppId      int `json:"app_id"`
	PlatformId int `json:"platform_id"`
}

type Config struct {
	Verbose int `json:"verbose"`

	Api          API              `json:"api"`
	Vms          VMS              `json:"vms"`
	Scanners     []Scanner        `json:"scanners"`
	PlatformMaps []PlatformMapper `json:"platform_maps"`
	BadNames     []string         `json:"bad_names"`
	BadIps       []string         `json:"bad_ips"`
	BadAsns      []uint           `json:"bad_asns"`
}
