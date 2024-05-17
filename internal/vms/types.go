package vms

type Server struct {
	Address    string `json:"addr"`
	Port       int    `json:"gameport"`
	SteamId    string `json:"steamid"`
	HostName   string `json:"name"`
	AppId      int    `json:"appid"`
	Gamedir    string `json:"gamedir"`
	Version    string `json:"version"`
	Product    string `json:"product"`
	Region     int    `json:"region"`
	Players    int    `json:"players"`
	MaxPlayers int    `json:"max_players"`
	Bots       int    `json:"bots"`
	Map        string `json:"map"`
	Secure     bool   `json:"secure"`
	Dedicated  bool   `json:"dedicated"`
	Os         string `json:"os"`
	GameType   string `json:"gametype"`
}

type DataResponse struct {
	Servers []Server `json:"servers"`
}

type Response struct {
	Response DataResponse `json:"response"`
}
