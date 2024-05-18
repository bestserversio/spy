package servers

type Server struct {
	Id      *int    `json:"id,omitempty"`
	Visible *bool   `json:"visible,omitempty"`
	Url     *string `json:"url,omitempty"`

	Ip       *string `json:"ip,omitempty"`
	Ip6      *string `json:"ip6,omitempty"`
	Port     *int    `json:"port,omitempty"`
	HostName *string `json:"hostName,omitempty"`

	PlatformId *int `json:"platformId,omitempty"`
	CategoryId *int `json:"categoryId,omitempty"`

	Name             *string `json:"name,omitempty"`
	DescriptionShort *string `json:"descriptionShort,omitempty"`
	Description      *string `json:"description,omitempty"`
	Features         *string `json:"features,omitempty"`
	Rules            *string `json:"rules,omitempty"`

	Online   *bool   `json:"online,omitempty"`
	CurUsers *int    `json:"curUsers,omitempty"`
	MaxUsers *int    `json:"maxUsers,omitempty"`
	Bots     *int    `json:"bots,omitempty"`
	MapName  *string `json:"mapName,omitempty"`
	AvgUsers *int    `json:"avgUsers,omitempty"`

	Region      *string  `json:"region,omitempty"`
	LocationLat *float64 `json:"locationLat,omitempty"`
	LocationLon *float64 `json:"locationLon,omitempty"`

	LastQueried *string `json:"lastQueried,omitempty"`
}

type ServerClaimKey struct {
	ServerId int    `json:"serverId"`
	UserId   string `json:"userId"`

	Key     string `json:"key"`
	Expires string `json:"expires"`
}
