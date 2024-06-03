package utils

import "github.com/bestserversio/spy/internal/config"

func AppIdToPlatformId(cfg *config.Config, app_id int) int {
	ret := 0

	for _, t := range cfg.PlatformMaps {
		if t.AppId == app_id {
			ret = t.PlatformId
		}
	}

	return ret
}
