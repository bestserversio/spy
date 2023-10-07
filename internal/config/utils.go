package config

import (
	"fmt"
	"strconv"
)

func (cfg *Config) PrintConfig() {
	fmt.Println("Config Settings")

	fmt.Println("\tHost => " + cfg.Host)
	fmt.Println("\tEndPoint => " + cfg.EndPoint)
	fmt.Println("\tAuthorization => " + cfg.Authorization)
	fmt.Println("\tSSL => " + strconv.FormatBool(cfg.Ssl))
}
