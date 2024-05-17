package config

import (
	"fmt"
	"strconv"
)

func (cfg *Config) PrintConfig() {
	fmt.Println("Config Settings")

	fmt.Println("\tHost => " + cfg.Api.Host)
	fmt.Println("\tAuthorization => " + cfg.Api.Authorization)
	fmt.Println("\tSSL => " + strconv.FormatBool(cfg.Api.Ssl))
}
