package scanners

import (
	"math/rand"
	"sync"
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

	utils.DebugMsg(1, cfg, "[SCANNER %d] Starting scanner with protocol '%s'!", idx, scanner.Protocol)

	next_platform := 0

	visible_cnt := make(map[int]int)

	for {
		select {
		case <-scanner.Channel:
			utils.DebugMsg(1, cfg, "[SCANNER %d] Found close request. Closing scanner.", idx)

			return
		default:
			if scanner.RandomPlatforms {
				// Reseed.
				rand.Seed(time.Now().UnixNano())

				// We need to pick a random app ID.
				next_platform = rand.Intn(len(scanner.PlatformIds))
			} else {
				utils.GetNextIndex(&next_platform, len(scanner.PlatformIds))
			}

			platform_id := scanner.PlatformIds[next_platform]

			// Check if we should only scan visible servers.
			visibleOnly := false

			if scanner.VisibleSkipCount > 0 {
				if val, exists := visible_cnt[int(platform_id)]; exists {
					visible_cnt[int(platform_id)] = val + 1
				} else {
					visible_cnt[int(platform_id)] = 1
				}

				if visible_cnt[int(platform_id)] < scanner.VisibleSkipCount {
					visibleOnly = true
				} else {
					visible_cnt[int(platform_id)] = 1
				}
			}

			visibleOnlyStr := "No"

			if visibleOnly {
				visibleOnlyStr = "Yes"
			}

			utils.DebugMsg(4, cfg, "[SCANNER %d] Using platform ID %d. Visible Only => %s", idx, platform_id, visibleOnlyStr)

			// Retrieve servers from API.
			allServers, err := servers.RetrieveServers(cfg, &platform_id, &scanner.Limit, visibleOnly)

			if err != nil {
				utils.DebugMsg(1, cfg, "[SCANNER %d] Failed to retrieve servers using platform ID '%d' due to error.", idx, platform_id)
				utils.DebugMsg(1, cfg, err.Error())

				Respin(scanner)

				continue
			}

			utils.DebugMsg(4, cfg, "[SCANNER %d] Found %d servers to update from API for platform ID '%d'!", idx, len(allServers), platform_id)

			// Make sure we have servers.
			if len(allServers) < 1 {
				Respin(scanner)

				continue
			}

			onlineCnt := 0
			offlineCnt := 0

			var wg sync.WaitGroup

			// Loop through each server and update.
			for i := 0; i < len(allServers); i++ {
				wg.Add(1)

				srv := &allServers[i]

				go func(srv *servers.Server, i int) {
					defer func() {
						if r := recover(); r != nil {
							utils.DebugMsg(1, cfg, "[SCANNER %d] Found panic when scanning '%s:%d'.", idx, *srv.Ip, *srv.Port)
						}
					}()

					var err error
					defer wg.Done()

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
							utils.DebugMsg(1, cfg, "[SCANNER %d] Missing IP/port for server. Skipping.", idx)

							return
						}

						// Retrieve timeout.
						timeout := scanner.QueryTimeout

						if timeout < 1 {
							timeout = 3
						}

						// Try querying server with A2S and check for error..
						err = QueryA2s(srv, timeout, scanner.A2sPlayer)

						if err != nil {
							utils.DebugMsg(5, cfg, "[SCANNER %d] Failed to query A2S server '%s:%d' due to error :: %s.", idx, *srv.Ip, *srv.Port, err.Error())

							*srv.Online = false
						}
					}

					// Check for filters.
					filtered, err := srv.FilterServer(cfg)

					if err != nil {
						utils.DebugMsg(1, cfg, "[SCANNER %d] Failed to filter server '%s:%d' due to error. Setting to invisible.", idx, *srv.Ip, *srv.Port)
						utils.DebugMsg(1, cfg, err.Error())

						*srv.Visible = false
					}

					if filtered {
						utils.DebugMsg(5, cfg, "[SCANNER %d] Found '%s:%d' filtered. Setting to invisible.", idx, *srv.Ip, *srv.Port)

						*srv.Visible = false
					}

					// Check for sub bots.
					if scanner.SubBots && (srv.CurUsers != nil && srv.Bots != nil) && ((*srv.CurUsers - *srv.Bots) > -1) {
						*srv.CurUsers -= *srv.Bots
					}

					// If we're online, set last online.
					if *srv.Online {
						// Update last queried to now.
						if srv.LastOnline == nil {
							srv.LastOnline = new(string)
						}

						now := time.Now().UTC()
						isoDate := now.Format("2006-01-02T15:04:05Z")

						*srv.LastOnline = isoDate

						onlineCnt++
					} else {
						offlineCnt++
					}

					utils.DebugMsg(5, cfg, "[SCANNER %d] Updating server '%s:%d' for platform ID '%d'. Players => %d. Max players => %d. Map => %s.", idx, *srv.Ip, *srv.Port, platform_id, *srv.CurUsers, *srv.MaxUsers, *srv.MapName)
				}(srv, i)

				// Check for request delay.
				if scanner.RequestDelay > 0 {
					time.Sleep(time.Millisecond * time.Duration(scanner.RequestDelay))
				}
			}

			wg.Wait()

			if !scanner.RecvOnly {
				// Update servers.
				cnt, err := servers.AddServers(cfg, allServers, false)

				if err != nil {
					utils.DebugMsg(1, cfg, "[SCANNER %d] Failed to update servers for platform ID '%d' due to error.", idx, platform_id)
					utils.DebugMsg(1, cfg, err.Error())
				} else {
					utils.DebugMsg(3, cfg, "[SCANNER %d] Updated %d servers in database for platform ID '%d'! Online => %d. Offline => %d", idx, cnt, platform_id, onlineCnt, offlineCnt)
				}
			}
		}

		Respin(scanner)
	}
}

func SetupScanners(cfg *config.Config) {
	for i := 0; i < len(cfg.Scanners); i++ {
		s := &cfg.Scanners[i]

		// Check if channel if valid here and if so, cancel existing go routine.
		if s.Channel != nil {
			s.Channel <- true
		}

		s.Channel = make(chan bool)

		go DoScanner(cfg, s, i+1)
	}
}
