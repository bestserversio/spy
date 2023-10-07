package servers

type Server struct {
	Id      int    `json:"id"`
	Visible bool   `json:"visible"`
	Url     string `json:"url"`

	Ip       string `json:"ip"`
	Ip6      string `json:"ip6"`
	Port     int    `json:"port"`
	HostName string `json:"hostName"`

	PlatformId int `json:"platformId"`
	CategoryId int `json:"categoryId"`

	Name             string `json:"name"`
	DescriptionShort string `json:"descriptionShort"`
	Description      string `json:"description"`
	Features         string `json:"features"`
	Rules            string `json:"rules"`

	Online   bool   `json:"online"`
	CurUsers int    `json:"curUsers"`
	MaxUsers int    `json:"maxUsers"`
	Bots     int    `json:"bots"`
	MapName  string `json:"mapName"`
	AvgUsers int    `json:"avgUsers"`

	Region      string  `json:"region"`
	LocationLat float64 `json:"locationLat"`
	LocationLon float64 `json:"locationLon"`

	LastQueried string `json:"lastQueried"`
}

type ServerClaimKey struct {
	ServerId int    `json:"serverId"`
	UserId   string `json:"userId"`

	Key     string `json:"key"`
	Expires string `json:"expires"`
}
