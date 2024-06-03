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
		rand_num := rand.Intn(len(scanner.AppIds))

		app_id := scanner.AppIds[rand_num]

		// Retrieve servers from API.
		allServers, err := servers.RetrieveServers(cfg, &app_id, &scanner.Limit)

		if err != nil {
			utils.DebugMsg(1, cfg.Verbose, "[SCANNER] Failed to retrieve servers using app ID '%d' due to error.", app_id)
			utils.DebugMsg(1, cfg.Verbose, err.Error())

			Respin(scanner)
		}

		utils.DebugMsg(4, cfg.Verbose, "[SCANNER] Found %d servers to update from API for app ID '%d'!", len(allServers), app_id)

		// Loop through each server and update.
		for i := 0; i < len(allServers); i++ {
			srv := &allServers[i]

			if srv.Online == nil {
				srv.Online = new(bool)
			}

			// Set where clause.
			if srv.Where.Id == nil {
				srv.Where.Id = new(int)
			}

			*srv.Where.Id = *srv.Id

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
					utils.DebugMsg(1, cfg.Verbose, "[SCANNER] Missing IP/port for server. Skipping.")

					continue
				}

				// Try querying server with A2S and check for error..
				err = QueryA2s(srv)

				if err != nil {
					utils.DebugMsg(3, cfg.Verbose, "[SCANNER] Failed to query A2S server '%s:%d' due to error.", *srv.Ip, *srv.Port)
					utils.DebugMsg(3, cfg.Verbose, err.Error())

					if srv.MaxUsers == nil {
						srv.MaxUsers = new(int)
					}

					*srv.MaxUsers = 0

					*srv.Online = false
				}
			}

			utils.DebugMsg(4, cfg.Verbose, "[SCANNER] Updating server '%s:%d' for app ID '%d'. Players => %d. Max players => %d. Map => %s.", *srv.Ip, *srv.Port, app_id, *srv.CurUsers, *srv.MaxUsers, *srv.MapName)
		}

		if !scanner.RecvOnly {
			// Update servers.
			cnt, err := servers.AddServers(cfg, allServers, false)

			if err != nil {
				utils.DebugMsg(1, cfg.Verbose, "[SCANNER] Failed to update servers for app ID '%d' due to error.", app_id)
				utils.DebugMsg(1, cfg.Verbose, err.Error())
			} else {
				utils.DebugMsg(3, cfg.Verbose, "[SCANNER] Updated %d servers in database for app ID '%d'!", cnt, app_id)
			}
		}

		Respin(scanner)
	}
}
