package config

type Config struct {
	Verbose       int    `json:"verbose"`
	Host          string `json:"host"`
	EndPoint      string `json:"endpoint"`
	Authorization string `json:"authorization"`
	Ssl           bool   `json:"ssl"`
	Timeout       uint   `json:"timeout"`
}
