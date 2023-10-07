package config

type Config struct {
	Host          string `json:"host"`
	EndPoint      string `json:"endpoint"`
	Authorization string `json:"authorization"`
	Ssl           bool   `json:"ssl"`
}
