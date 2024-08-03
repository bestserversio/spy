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
	Timeout      int       `json:"timeout"`
	ApiToken     string    `json:"api_token"`
	AppIds       []int     `json:"app_ids"`
	RecvOnly     bool      `json:"recv_only"`
	MinWait      int       `json:"min_wait"`
	MaxWait      int       `json:"max_wait"`
	Limit        int       `json:"limit"`
	ExcludeEmpty bool      `json:"exclude_empty"`
	OnlyEmpty    bool      `json:"only_empty"`
	SubBots      bool      `json:"sub_bots"`
	AddOnly      bool      `json:"add_only"`
	RandomApps   bool      `json:"random_apps"`
	SetOffline   bool      `json:"set_offline"`
	Channel      chan bool `json:"-"`
}

type Scanner struct {
	Protocol         string    `json:"protocol"`
	PlatformIds      []uint    `json:"platform_ids"`
	MinWait          int       `json:"min_wait"`
	MaxWait          int       `json:"max_wait"`
	Limit            int       `json:"limit"`
	RecvOnly         bool      `json:"recv_only"`
	SubBots          bool      `json:"sub_bots"`
	QueryTimeout     int       `json:"query_timeout"`
	A2sPlayer        bool      `json:"a2s_player"`
	RandomPlatforms  bool      `json:"random_platforms"`
	VisibleSkipCount int       `json:"visible_skip_count"`
	RequestDelay     int       `json:"request_delay"`
	Channel          chan bool `json:"-"`
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

type RemoveDups struct {
	Enabled    bool `json:"enabled"`
	Interval   int  `json:"interval"`
	Limit      int  `json:"limit"`
	MaxServers int  `json:"max_servers"`
	Timeout    int  `json:"timeout"`
}

type RemoveTimedOut struct {
	Enabled      bool  `json:"enabled"`
	PlatformIds  []int `json:"platform_ids"`
	TimedOutTime int   `json:"timed_out_time"`
	Interval     int   `json:"interval"`
	Timeout      int   `json:"timeout"`
}

type Config struct {
	Verbose      int     `json:"verbose"`
	LogDirectory *string `json:"log_directory"`

	Api             BsAPI            `json:"api"`
	WebApi          WebApi           `json:"web_api"`
	Vms             []VMS            `json:"vms"`
	Scanners        []Scanner        `json:"scanners"`
	PlatformMaps    []PlatformMapper `json:"platform_maps"`
	BadNames        []string         `json:"bad_names"`
	BadIps          []string         `json:"bad_ips"`
	BadAsns         []uint           `json:"bad_asns"`
	GoodIps         []string         `json:"good_ips"`
	RemoveInactive  RemoveInactive   `json:"remove_inactive"`
	PlatformFilters []PlatformFilter `json:"platform_filters"`
	RemoveDups      RemoveDups       `json:"remove_dups"`
	RemoveTimedOut  RemoveTimedOut   `json:"remove_timed_out"`
}
