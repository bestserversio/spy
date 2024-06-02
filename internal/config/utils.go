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

	fmt.Println("\tScanners")
	for i, s := range cfg.Scanners {
		fmt.Println("\t\tScanner #" + strconv.Itoa(i+1))

		fmt.Println("\t\t\tProtocol => " + s.Protocol)
		fmt.Println("\t\t\tMin Wait => " + strconv.Itoa(s.MinWait))
		fmt.Println("\t\t\tMax Wait => " + strconv.Itoa(s.MaxWait))

		ids := "None"

		if len(s.AppIds) > 0 {
			var ids_s []string

			for _, a := range s.AppIds {
				ids_s = append(ids_s, strconv.Itoa(a))
			}

			ids = strings.Join(ids_s, ", ")
		}

		fmt.Println("\t\t\tApp IDs => " + ids)
	}

	fmt.Println("\tPlatform Mappings (App ID => Platform ID)")
	for _, v := range cfg.PlatformMaps {
		fmt.Println("\t\t- " + strconv.Itoa(v.AppId) + " => " + strconv.Itoa(v.PlatformId))
	}

}
