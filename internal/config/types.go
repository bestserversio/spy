package config

type BsAPI struct {
	Host          string `json:"host"`
	Authorization string `json:"authorization"`
	Timeout       int    `json:"timeout"`
}

type WebApi struct {
	Enabled       bool   `json:"enabled"`
	Host          string `json:"host"`
	Endpoint      string `json:"endpoint"`
	Authorization string `json:"authorization"`
	Timeout       int    `json:"timeout"`
	Interval      int    `json:"interval"`
	SaveToFs      bool   `json:"save_to_fs"`
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
	ExcludeEmpty bool   `json:"exclude_empty"`
	SubBots      bool   `json:"sub_bots"`
	AddOnly      bool   `json:"add_only"`
	RandomApps   bool   `json:"random_apps"`
}

type Scanner struct {
	Protocol        string `json:"protocol"`
	PlatformIds     []uint `json:"platform_ids"`
	MinWait         int    `json:"min_wait"`
	MaxWait         int    `json:"max_wait"`
	Limit           int    `json:"limit"`
	RecvOnly        bool   `json:"recv_only"`
	SubBots         bool   `json:"sub_bots"`
	QueryTimeout    int    `json:"query_timeout"`
	A2sPlayer       bool   `json:"a2s_player"`
	RandomPlatforms bool   `json:"random_platforms"`
	Channel         chan bool
}

type PlatformMapper struct {
	AppId      int `json:"app_id"`
	PlatformId int `json:"platform_id"`
}

type RemoveInactive struct {
	Enabled      bool `json:"enabled"`
	Interval     int  `json:"interval"`
	InactiveTime int  `json:"inactive_time"`
	Timeout      int  `json:"timeout"`
}

type PlatformFilter struct {
	Id                int   `json:"id"`
	MaxUsers          *int  `json:"max_users"`
	MaxCurUsers       *int  `json:"max_cur_users"`
	AllowUserOverflow *bool `json:"allow_user_overflow"`
}

type Config struct {
	Verbose int `json:"verbose"`

	Api             BsAPI            `json:"api"`
	WebApi          WebApi           `json:"web_api"`
	Vms             VMS              `json:"vms"`
	Scanners        []Scanner        `json:"scanners"`
	PlatformMaps    []PlatformMapper `json:"platform_maps"`
	BadNames        []string         `json:"bad_names"`
	BadIps          []string         `json:"bad_ips"`
	BadAsns         []uint           `json:"bad_asns"`
	RemoveInactive  RemoveInactive   `json:"remove_inactive"`
	PlatformFilters []PlatformFilter `json:"platform_filters"`
}
