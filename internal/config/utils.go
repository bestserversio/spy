package config

import (
	"fmt"
	"strconv"
	"strings"
)

func (cfg *Config) PrintConfig() {
	fmt.Println("Config Settings")

	fmt.Println("\tVerbose => " + strconv.Itoa(cfg.Verbose))

	fmt.Println("\tAPI")
	fmt.Println("\t\tHost => " + cfg.Api.Host)
	fmt.Println("\t\tAuthorization => " + cfg.Api.Authorization)
	fmt.Println("\t\tSSL => " + strconv.FormatBool(cfg.Api.Ssl))
	fmt.Println("\t\tTimeout => " + strconv.Itoa(cfg.Api.Timeout))

	fmt.Println("\tVMS")
	fmt.Println("\t\tEnabled => " + strconv.FormatBool(cfg.Vms.Enabled))
	fmt.Println("\t\tInterval => " + strconv.Itoa(cfg.Vms.Interval))
	fmt.Println("\t\tTimeout => " + strconv.Itoa(cfg.Vms.Timeout))
	fmt.Println("\t\tAPI Token => " + cfg.Vms.ApiToken)

	appIdList := make([]string, len(cfg.Vms.AppIds))
	for i, num := range cfg.Vms.AppIds {
		appIdList[i] = strconv.Itoa(num)
	}

	fmt.Println("\t\tApp IDs => " + strings.Join(appIdList, ", "))

	fmt.Println("\tScanner")
	fmt.Println("\t\tMin Wait => " + strconv.Itoa(cfg.Scanner.MinWait))
	fmt.Println("\t\tMax Wait => " + strconv.Itoa(cfg.Scanner.MaxWait))

	fmt.Println("\tPlatform Mappings (App ID => Platform ID)")
	for _, v := range cfg.PlatformMaps {
		fmt.Println("\t\t- " + strconv.Itoa(v.AppId) + " => " + strconv.Itoa(v.PlatformId))
	}

}
