package scanners

import (
	"math/rand"
	"time"

	"github.com/bestserversio/spy/internal/config"
	"github.com/bestserversio/spy/internal/servers"
	"github.com/bestserversio/spy/internal/utils"
)

func Respin(scanner *config.Scanner) {
	utils.RandomWait(scanner.MinWait, scanner.MaxWait)
}

func DoScanner(cfg *config.Config, scanner *config.Scanner, idx int) {
	// Defaults to A2S
	query_type := 0

	utils.DebugMsg(1, cfg.Verbose, "[SCANNER %d] Starting scanner with protocol '%s'!", idx, scanner.Protocol)

	for {
		// Reseed.
		rand.Seed(time.Now().UnixNano())

		// We need to pick a random app ID.
		rand_num := rand.Intn(len(scanner.PlatformIds))

		platform_id := scanner.PlatformIds[rand_num]

		utils.DebugMsg(4, cfg.Verbose, "[SCANNER %d] Using platform ID %d.", idx, platform_id)

		// Retrieve servers from API.
		allServers, err := servers.RetrieveServers(cfg, &platform_id, &scanner.Limit)

		if err != nil {
			utils.DebugMsg(1, cfg.Verbose, "[SCANNER %d] Failed to retrieve servers using platform ID '%d' due to error.", idx, platform_id)
			utils.DebugMsg(1, cfg.Verbose, err.Error())

			Respin(scanner)
		}

		utils.DebugMsg(4, cfg.Verbose, "[SCANNER %d] Found %d servers to update from API for platform ID '%d'!", idx, len(allServers), platform_id)

		// Make sure we have servers.
		if len(allServers) < 1 {
			Respin(scanner)
		}

		// Loop through each server and update.
		for i := 0; i < len(allServers); i++ {
			srv := &allServers[i]

			// Allocate visibility if needed.
			if srv.Visible == nil {
				srv.Visible = new(bool)
			}

			*srv.Visible = true

			// Allocate online if needed.
			if srv.Online == nil {
				srv.Online = new(bool)
			}

			// Set where clause.
			if srv.Where.Id == nil {
				srv.Where.Id = new(int)
			}

			*srv.Where.Id = *srv.Id

			// Make original ID nil.
			srv.Id = nil

			// Update last queried to now.
			if srv.LastQueried == nil {
				srv.LastQueried = new(string)
			}

			now := time.Now().UTC()
			isoDate := now.Format("2006-01-02T15:04:05Z")

			*srv.LastQueried = isoDate

			switch query_type {
			case 0:
				if srv.Ip == nil || srv.Port == nil {
					utils.DebugMsg(1, cfg.Verbose, "[SCANNER %d] Missing IP/port for server. Skipping.", idx)

					continue
				}

				// Try querying server with A2S and check for error..
				err = QueryA2s(srv)

				if err != nil {
					utils.DebugMsg(4, cfg.Verbose, "[SCANNER %d] Failed to query A2S server '%s:%d' due to error.", idx, *srv.Ip, *srv.Port)
					utils.DebugMsg(4, cfg.Verbose, err.Error())

					*srv.Online = false
				}
			}

			// Check for filters.
			filtered, err := srv.FilterServer(cfg)

			if err != nil {
				utils.DebugMsg(1, cfg.Verbose, "[SCANNER %d] Failed to filter server '%s:%d' due to error. Setting to invisible.", idx, *srv.Ip, *srv.Port)
				utils.DebugMsg(1, cfg.Verbose, err.Error())

				*srv.Visible = false
			}

			if filtered {
				utils.DebugMsg(3, cfg.Verbose, "[SCANNER %d] Found '%s:%d' filtered. Setting to invisible.", idx, *srv.Ip, *srv.Port)

				*srv.Visible = false
			}

			utils.DebugMsg(4, cfg.Verbose, "[SCANNER %d] Updating server '%s:%d' for platform ID '%d'. Players => %d. Max players => %d. Map => %s.", idx, *srv.Ip, *srv.Port, platform_id, *srv.CurUsers, *srv.MaxUsers, *srv.MapName)
		}

		if !scanner.RecvOnly {
			// Update servers.
			cnt, err := servers.AddServers(cfg, allServers, false)

			if err != nil {
				utils.DebugMsg(1, cfg.Verbose, "[SCANNER %d] Failed to update servers for platform ID '%d' due to error.", idx, platform_id)
				utils.DebugMsg(1, cfg.Verbose, err.Error())
			} else {
				utils.DebugMsg(3, cfg.Verbose, "[SCANNER %d] Updated %d servers in database for platform ID '%d'!", idx, cnt, platform_id)
			}
		}

		Respin(scanner)
	}
}
