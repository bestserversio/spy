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

	return false, err
}
