package vms

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bestserversio/spy/internal/config"
	"github.com/bestserversio/spy/internal/servers"
	"github.com/bestserversio/spy/internal/utils"
)

func Respin(vms *config.VMS) {
	utils.RandomWait(vms.MinWait, vms.MaxWait)
}

func DoVms(cfg *config.Config, vms *config.VMS, idx int) {
	utils.DebugMsg(1, cfg, "[VMS %d] Starting...", idx)

	next_app := 0

	for {
		select {
		case <-vms.Channel:
			utils.DebugMsg(1, cfg, "[VMS %d] Found close request. Closing VMS.", idx)

			return
		default:
			if vms.RandomApps {
				// Retrieve random app ID.
				rand.Seed(time.Now().UnixNano())

				next_app = rand.Intn(len(vms.AppIds))
			} else {
				utils.GetNextIndex(&next_app, len(vms.AppIds))
			}

			appId := vms.AppIds[next_app]

			// Get platform ID.
			platform_id := utils.AppIdToPlatformId(cfg, appId)

			utils.DebugMsg(4, cfg, "[VMS %d] Using app ID '%d'.", idx, appId)

			// Retrieve servers.
			allServers, err := RetrieveServers(vms, appId)

			if err != nil {
				utils.DebugMsg(1, cfg, "[VMS %d] Failed to retrieve servers for app ID '%d' due to error.", idx, appId)
				utils.DebugMsg(1, cfg, err.Error())

				Respin(vms)

				continue
			}

			utils.DebugMsg(3, cfg, "[VMS %d] Retrieved %d servers from app ID '%d'.", idx, len(allServers), appId)

			var serversToUpdate []servers.Server

			// Loop through all servers from VMS result.
			for _, srv := range allServers {
				// Create new servers object from servers package and assign basic values.
				newSrv := servers.Server{
					Visible:     new(bool),
					Online:      new(bool),
					Name:        new(string),
					CurUsers:    new(int),
					MaxUsers:    new(int),
					Bots:        new(int),
					MapName:     new(string),
					Ip:          new(string),
					Port:        new(int),
					PlatformId:  new(int),
					Os:          new(string),
					LastQueried: new(string),
					Region:      new(string),
				}

				initialOnline := true

				if vms.SetOffline {
					initialOnline = false
				}

				*newSrv.Visible = true
				*newSrv.Online = initialOnline
				*newSrv.Name = srv.HostName
				*newSrv.CurUsers = srv.Players
				*newSrv.MaxUsers = srv.MaxPlayers
				*newSrv.Bots = srv.Bots
				*newSrv.MapName = srv.Map
				*newSrv.PlatformId = platform_id
				*newSrv.Os = srv.Os
				*newSrv.Region = utils.GetRegion(srv.Region)

				// Check for sub bots.
				if vms.SubBots && (*newSrv.CurUsers-*newSrv.Bots) > -1 {
					*newSrv.CurUsers -= *newSrv.Bots
				}

				// Set last queries.
				now := time.Now().UTC()
				isoDate := now.Format("2006-01-02T15:04:05Z")

				*newSrv.LastQueried = isoDate

				// We need to split IP and port from address.
				split := strings.Split(srv.Address, ":")

				if len(split) > 1 {
					*newSrv.Ip = split[0]

					portNum, err := strconv.Atoi(split[1])

					if err != nil {
						utils.DebugMsg(1, cfg, "[VMS %d] Failed to convert port string to number due to error :: %s", idx, err.Error())

						Respin(vms)

						continue
					}

					*newSrv.Port = portNum
				}

				// Assign IP and port as where.
				if newSrv.Where.Ip == nil {
					newSrv.Where.Ip = new(string)
				}

				if newSrv.Where.Port == nil {
					newSrv.Where.Port = new(int)
				}

				*newSrv.Where.Ip = *newSrv.Ip
				*newSrv.Where.Port = *newSrv.Port

				// Before adding, do filter check.
				filtered, err := newSrv.FilterServer(cfg)

				if err != nil {
					utils.DebugMsg(1, cfg, "[VMS %d] Failed to filter server '%s:%d' due to error. Setting to invisible just in case.", idx, *newSrv.Ip, *newSrv.Port)
					utils.DebugMsg(1, cfg, err.Error())

					*newSrv.Visible = false
				}

				if filtered {
					utils.DebugMsg(5, cfg, "[VMS %d] Setting '%s:%d' to invisible due to being filtered.", idx, *newSrv.Ip, *newSrv.Port)

					*newSrv.Visible = false
				}

				// Append to servers to update array.
				serversToUpdate = append(serversToUpdate, newSrv)

				utils.DebugMsg(5, cfg, "[VMS %d] Found and adding/updating server '%s:%d'. Host Name => '%s'. Players => %d. Max Players => %d. Map Name => '%s'.", idx, *newSrv.Ip, *newSrv.Port, *newSrv.Name, *newSrv.CurUsers, *newSrv.MaxUsers, *newSrv.MapName)
			}

			if len(serversToUpdate) < 1 {
				utils.DebugMsg(3, cfg, "[VMS %d] Found no servers to update for app ID %d.", idx, appId)

				Respin(vms)

				continue
			}

			// Add/update servers.
			if !vms.RecvOnly {
				cnt, err := servers.AddServers(cfg, serversToUpdate, vms.AddOnly)

				if err != nil {
					utils.DebugMsg(1, cfg, "[VMS %d] Failed to add/update servers for app ID %d due to error.", idx, appId)
					utils.DebugMsg(1, cfg, err.Error())

					Respin(vms)

					continue
				}

				utils.DebugMsg(2, cfg, "[VMS %d] Added/Updated %d servers!", idx, cnt)
			}
		}

		Respin(vms)
	}
}

func SetupVms(cfg *config.Config) {
	for i := 0; i < len(cfg.Vms); i++ {
		vms := &cfg.Vms[i]

		// Check channel.
		if vms.Channel != nil {
			vms.Channel <- true
		}

		vms.Channel = make(chan bool)

		go DoVms(cfg, vms, i+1)
	}
}
