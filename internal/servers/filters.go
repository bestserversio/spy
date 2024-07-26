package servers

import (
	"github.com/bestserversio/spy/internal/config"
	"github.com/bestserversio/spy/internal/utils"
)

func (srv *Server) FilterServer(cfg *config.Config) (bool, error) {
	var err error

	// Scan for bad names.
	if len(cfg.BadNames) > 0 && srv.Name != nil {
		if utils.ContainsStrings(*srv.Name, cfg.BadNames) {
			return true, err
		}
	}

	// Scan for bad IP ranges.
	if len(cfg.BadIps) > 0 && (srv.Ip != nil || srv.Ip6 != nil) {
		if srv.Ip != nil {
			filtered, err := utils.ContainsIps(*srv.Ip, cfg.BadIps)

			if filtered || err != nil {
				return filtered, err
			}
		}

		if srv.Ip6 != nil {
			filtered, err := utils.ContainsIps(*srv.Ip6, cfg.BadIps)

			if filtered || err != nil {
				return filtered, err
			}
		}
	}

	// Scan for platform-specific filters.
	p_filters := srv.GetPlatformFilters(cfg)

	// Check max current users.
	if p_filters.MaxCurUsers != nil && srv.CurUsers != nil && *srv.CurUsers > *p_filters.MaxCurUsers {
		return true, err
	}

	// Check max users.
	if p_filters.MaxUsers != nil && srv.MaxUsers != nil && *srv.MaxUsers > *p_filters.MaxUsers {
		return true, err
	}

	// Check user overflow.
	if p_filters.AllowUserOverflow != nil && srv.CurUsers != nil && srv.MaxUsers != nil && *srv.CurUsers > *srv.MaxUsers {
		return true, err
	}

	return false, err
}

func (srv *Server) GetPlatformFilters(cfg *config.Config) config.PlatformFilter {
	ret := config.PlatformFilter{}

	if srv.PlatformId == nil {
		return ret
	}

	for _, filter := range cfg.PlatformFilters {
		if filter.Id == *srv.PlatformId {
			// Check for max current users.
			if filter.MaxCurUsers != nil {
				ret.MaxCurUsers = new(int)
				*ret.MaxCurUsers = *filter.MaxCurUsers
			}

			// Check for max users.
			if filter.MaxUsers != nil {
				ret.MaxUsers = new(int)
				*ret.MaxUsers = *filter.MaxUsers
			}

			if filter.AllowUserOverflow != nil {
				ret.AllowUserOverflow = new(bool)
				*ret.AllowUserOverflow = *filter.AllowUserOverflow
			}

			break
		}
	}

	return ret
}
