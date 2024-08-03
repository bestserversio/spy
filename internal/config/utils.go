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

	if len(cfg.Vms) > 0 {
		fmt.Println("\tVMS")

		for i, vms := range cfg.Vms {
			fmt.Println("\t\tVMS #" + strconv.Itoa(i+1))

			fmt.Println("\t\t\tMin Wait => " + strconv.Itoa(vms.MinWait))
			fmt.Println("\t\t\tMax Wait => " + strconv.Itoa(vms.MaxWait))
			fmt.Println("\t\t\tTimeout => " + strconv.Itoa(vms.Timeout))
			fmt.Println("\t\t\tAPI Token => " + vms.ApiToken)

			recv_only := "No"

			if vms.RecvOnly {
				recv_only = "Yes"
			}

			fmt.Println("\t\t\tReceive Only => " + recv_only)

			sub_bots := "No"

			if vms.SubBots {
				sub_bots = "Yes"
			}

			fmt.Println("\t\t\tSub Bots => " + sub_bots)

			exclude_empty := "No"

			if vms.ExcludeEmpty {
				exclude_empty = "Yes"
			}

			fmt.Println("\t\t\tExclude Empty => " + exclude_empty)

			only_empty := "No"

			if vms.OnlyEmpty {
				only_empty = "Yes"
			}

			fmt.Println("\t\t\tOnly Empty => " + only_empty)

			add_only := "No"

			if vms.AddOnly {
				add_only = "Yes"
			}

			fmt.Println("\t\t\tAdd Only => " + add_only)

			random_apps := "No"

			if vms.RandomApps {
				random_apps = "Yes"
			}

			fmt.Println("\t\t\tRandom Apps => " + random_apps)

			set_offline := "No"

			if vms.SetOffline {
				set_offline = "Yes"
			}

			fmt.Println("\t\t\tSet Offline => " + set_offline)

			appIdList := make([]string, len(vms.AppIds))

			for i, num := range vms.AppIds {
				appIdList[i] = strconv.Itoa(num)
			}

			fmt.Println("\t\t\tApp IDs => " + strings.Join(appIdList, ", "))
		}
	}

	if len(cfg.Scanners) > 0 {
		fmt.Println("\tScanners")

		for i, s := range cfg.Scanners {
			fmt.Println("\t\tScanner #" + strconv.Itoa(i+1))

			fmt.Println("\t\t\tProtocol => " + s.Protocol)
			fmt.Println("\t\t\tMin Wait => " + strconv.Itoa(s.MinWait))
			fmt.Println("\t\t\tMax Wait => " + strconv.Itoa(s.MaxWait))
			fmt.Println("\t\t\tRequest Delay (MS) => " + strconv.Itoa(s.RequestDelay))

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

	fmt.Println("\tRemove Timed Out")

	enabled = "No"

	if cfg.RemoveTimedOut.Enabled {
		enabled = "Yes"
	}

	fmt.Println("\t\tEnabled => " + enabled)
	fmt.Println("\t\tInterval => " + strconv.Itoa(cfg.RemoveTimedOut.Interval))

	platform_ids := "None"

	if len(cfg.RemoveTimedOut.PlatformIds) > 0 {
		str_ids := make([]string, len(cfg.RemoveTimedOut.PlatformIds))

		for i, id := range cfg.RemoveTimedOut.PlatformIds {
			str_ids[i] = strconv.Itoa(id)
		}

		platform_ids = strings.Join(str_ids, ", ")
	}

	fmt.Println("Platform IDs => " + platform_ids)

	fmt.Println("Timed Out Time => " + strconv.Itoa(cfg.RemoveTimedOut.TimedOutTime))
	fmt.Println("Timeout => " + strconv.Itoa(cfg.RemoveTimedOut.Timeout))

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
