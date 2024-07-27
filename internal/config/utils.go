package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func (cfg *Config) Load(data string) error {
	err := json.Unmarshal([]byte(data), cfg)

	return err
}

func (cfg *Config) PrintConfig() {
	fmt.Println("Config Settings")

	fmt.Println("\tVerbose => " + strconv.Itoa(cfg.Verbose))

	logDir := "NULL"

	if cfg.LogDirectory != nil {
		logDir = *cfg.LogDirectory
	}

	fmt.Println("\tLog Directory => " + logDir)

	fmt.Println("\tBest Servers API")
	fmt.Println("\t\tHost => " + cfg.Api.Host)
	fmt.Println("\t\tAuthorization => " + cfg.Api.Authorization)
	fmt.Println("\t\tTimeout => " + strconv.Itoa(cfg.Api.Timeout))

	fmt.Println("\tWeb API")
	fmt.Println("\t\tHost => " + cfg.WebApi.Host)
	fmt.Println("\t\tEndpoint => " + cfg.WebApi.Endpoint)
	fmt.Println("\t\tAuthorization => " + cfg.WebApi.Authorization)
	fmt.Println("\t\tInterval => " + strconv.Itoa(cfg.WebApi.Interval))
	fmt.Println("\t\tTimeout => " + strconv.Itoa(cfg.WebApi.Timeout))
	saveToFs := "No"

	if cfg.WebApi.SaveToFs {
		saveToFs = "Yes"
	}

	fmt.Println("\t\tSave To FS => " + saveToFs)

	fmt.Println("\tVMS")
	fmt.Println("\t\tEnabled => " + strconv.FormatBool(cfg.Vms.Enabled))
	fmt.Println("\t\tMin Wait => " + strconv.Itoa(cfg.Vms.MinWait))
	fmt.Println("\t\tMax Wait => " + strconv.Itoa(cfg.Vms.MaxWait))
	fmt.Println("\t\tTimeout => " + strconv.Itoa(cfg.Vms.Timeout))
	fmt.Println("\t\tAPI Token => " + cfg.Vms.ApiToken)

	recv_only := "No"

	if cfg.Vms.RecvOnly {
		recv_only = "Yes"
	}

	fmt.Println("\t\tReceive Only => " + recv_only)

	sub_bots := "No"

	if cfg.Vms.SubBots {
		sub_bots = "Yes"
	}

	fmt.Println("\t\tSub Bots => " + sub_bots)

	exclude_empty := "No"

	if cfg.Vms.ExcludeEmpty {
		exclude_empty = "Yes"
	}

	fmt.Println("\t\tExclude Empty => " + exclude_empty)

	add_only := "No"

	if cfg.Vms.AddOnly {
		add_only = "Yes"
	}

	fmt.Println("\t\tAdd Only => " + add_only)

	random_apps := "No"

	if cfg.Vms.RandomApps {
		random_apps = "Yes"
	}

	fmt.Println("\t\tRandom Apps => " + random_apps)

	set_offline := "No"

	if cfg.Vms.SetOffline {
		set_offline = "Yes"
	}

	fmt.Println("\t\tSet Offline => " + set_offline)

	appIdList := make([]string, len(cfg.Vms.AppIds))
	for i, num := range cfg.Vms.AppIds {
		appIdList[i] = strconv.Itoa(num)
	}

	fmt.Println("\t\tApp IDs => " + strings.Join(appIdList, ", "))

	if len(cfg.Scanners) > 0 {
		fmt.Println("\tScanners")

		for i, s := range cfg.Scanners {
			fmt.Println("\t\tScanner #" + strconv.Itoa(i+1))

			fmt.Println("\t\t\tProtocol => " + s.Protocol)
			fmt.Println("\t\t\tMin Wait => " + strconv.Itoa(s.MinWait))
			fmt.Println("\t\t\tMax Wait => " + strconv.Itoa(s.MaxWait))
			sub_bots := "No"

			if s.SubBots {
				sub_bots = "Yes"
			}

			fmt.Println("\t\t\tSub Bots => " + sub_bots)
			fmt.Println("\t\t\tQuery Timeout => " + strconv.Itoa(s.QueryTimeout))

			a2s_player := "No"

			if s.A2sPlayer {
				a2s_player = "Yes"
			}

			fmt.Println("\t\t\tA2S Player => " + a2s_player)

			random_platforms := "No"

			if s.RandomPlatforms {
				random_platforms = "Yes"
			}

			fmt.Println("\t\t\tRandom Platforms => " + random_platforms)

			fmt.Println("\t\t\tVisible Skip Count => " + strconv.Itoa(s.VisibleSkipCount))

			ids := "None"

			if len(s.PlatformIds) > 0 {
				var ids_s []string

				for _, a := range s.PlatformIds {
					ids_s = append(ids_s, strconv.Itoa(int(a)))
				}

				ids = strings.Join(ids_s, ", ")
			}

			fmt.Println("\t\t\tPlatform IDs => " + ids)
		}
	}

	if len(cfg.PlatformMaps) > 0 {
		fmt.Println("\tPlatform Mappings (App ID => Platform ID)")

		for _, v := range cfg.PlatformMaps {
			fmt.Println("\t\t- " + strconv.Itoa(v.AppId) + " => " + strconv.Itoa(v.PlatformId))
		}
	}

	fmt.Println("\tRemove Inactive")

	enabled := "No"

	if cfg.RemoveInactive.Enabled {
		enabled = "Yes"
	}

	fmt.Println("\t\tEnabled => " + enabled)
	fmt.Println("\t\tInterval => " + strconv.Itoa(cfg.RemoveInactive.Interval))
	fmt.Println("\t\tInactive Time => " + strconv.Itoa(cfg.RemoveInactive.InactiveTime))
	fmt.Println("\t\tTimeout => " + strconv.Itoa(cfg.RemoveInactive.Timeout))

	if len(cfg.PlatformFilters) > 0 {
		fmt.Println("\tPlatform Filters")

		for _, f := range cfg.PlatformFilters {
			fmt.Println("\t\tPlatform " + strconv.Itoa(f.Id))

			if f.MaxCurUsers != nil {
				fmt.Println("\t\t\tMax Current Users => " + strconv.Itoa(*f.MaxCurUsers))
			}

			if f.MaxUsers != nil {
				fmt.Println("\t\t\tMax Users => " + strconv.Itoa(*f.MaxUsers))
			}

			if f.AllowUserOverflow != nil {
				allow_user_overflow := "No"

				if *f.AllowUserOverflow {
					allow_user_overflow = "Yes"
				}

				fmt.Println("\t\t\tAllow User Overflow => " + allow_user_overflow)
			}
		}
	}

	fmt.Println("\tRemove Duplicates")

	enabled = "No"

	if cfg.RemoveDups.Enabled {
		enabled = "Yes"
	}

	fmt.Println("\t\tEnabled => " + enabled)
	fmt.Println("\t\tInterval => " + strconv.Itoa(cfg.RemoveDups.Interval))
	fmt.Println("\t\tLimit => " + strconv.Itoa(cfg.RemoveDups.Limit))
	fmt.Println("\t\tMax Servers => " + strconv.Itoa(cfg.RemoveDups.MaxServers))
	fmt.Println("\t\tTimeout => " + strconv.Itoa(cfg.RemoveDups.Timeout))

	if len(cfg.BadNames) > 0 {
		fmt.Println("\tBad Names")

		for _, s := range cfg.BadNames {
			fmt.Println("\t\t- ", s)
		}
	}

	if len(cfg.BadIps) > 0 {
		fmt.Println("\tBad IPs")

		for _, s := range cfg.BadIps {
			fmt.Println("\t\t- ", s)
		}
	}

	if len(cfg.BadAsns) > 0 {
		fmt.Println("\tBad ASNs")

		for _, s := range cfg.BadAsns {
			fmt.Println("\t\t- ", s)
		}
	}

	if len(cfg.GoodIps) > 0 {
		fmt.Println("\tGood IPs")

		for _, s := range cfg.GoodIps {
			fmt.Println("\t\t- ", s)
		}
	}
}
