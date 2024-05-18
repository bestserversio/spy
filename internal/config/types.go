package config

type API struct {
	Host          string `json:"host"`
	Authorization string `json:"authorization"`
	Ssl           bool   `json:"ssl"`
	Timeout       int    `json:"timeout"`
}

type VMS struct {
	Enabled  bool   `json:"enabled"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
	ApiToken string `json:"api_token"`
	AppIds   []int  `json:"app_ids"`
}

type Scanner struct {
	MinWait int `json:"min_wait"`
	MaxWait int `json:"max_wait"`
}

type PlatformMapper struct {
	AppId      int `json:"appid"`
	PlatformId int `json:"platformid"`
}

type Config struct {
	Verbose int `json:"verbose"`

	Api          API              `json:"api"`
	Vms          VMS              `json:"vms"`
	Scanner      Scanner          `json:"scanner"`
	PlatformMaps []PlatformMapper `json:"platform_maps"`
}
