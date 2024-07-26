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

func Respin(cfg *config.Config) {
	utils.RandomWait(cfg.Vms.MinWait, cfg.Vms.MaxWait)
}

func DoVms(cfg *config.Config) {
	utils.DebugMsg(1, cfg, "[VMS] Starting....")

	next_app := 0

	for {
		if !cfg.Vms.Enabled {
			utils.DebugMsg(5, cfg, "[VMS] Found VMS disabled. Aborting DoVms().")

			Respin(cfg)

			continue
		}

		if cfg.Vms.RandomApps {
			// Retrieve random app ID.
			rand.Seed(time.Now().UnixNano())

			next_app = rand.Intn(len(cfg.Vms.AppIds))
		} else {
			utils.GetNextIndex(&next_app, len(cfg.Vms.AppIds))
		}

		appId := cfg.Vms.AppIds[next_app]

		// Get platform ID.
		platform_id := utils.AppIdToPlatformId(cfg, appId)

		utils.DebugMsg(4, cfg, "[VMS] Using (random) app ID '%d'.", appId)

		// Retrieve servers.
		allServers, err := RetrieveServers(cfg, appId)

		if err != nil {
			utils.DebugMsg(1, cfg, "[VMS] Failed to retrieve servers for app ID '%d' due to error.", appId)
			utils.DebugMsg(1, cfg, err.Error())

			Respin(cfg)

			continue
		}

		utils.DebugMsg(3, cfg, "[VMS] Retrieved %d servers from app ID '%d'.", len(allServers), appId)

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

			if cfg.Vms.SetOffline {
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
			if cfg.Vms.SubBots && (*newSrv.CurUsers-*newSrv.Bots) > -1 {
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
					utils.DebugMsg(1, cfg, "")

					Respin(cfg)

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
				utils.DebugMsg(1, cfg, "[VMS] Failed to filter server '%s:%d' due to error. Setting to invisible just in case.", *newSrv.Ip, *newSrv.Port)
				utils.DebugMsg(1, cfg, err.Error())

				*newSrv.Visible = false
			}

			if filtered {
				utils.DebugMsg(5, cfg, "[VMS] Setting '%s:%d' to invisible due to being filtered.", *newSrv.Ip, *newSrv.Port)

				*newSrv.Visible = false
			}

			// Append to servers to update array.
			serversToUpdate = append(serversToUpdate, newSrv)

			utils.DebugMsg(5, cfg, "[VMS] Found and adding/updating server '%s:%d'. Host Name => '%s'. Players => %d. Max Players => %d. Map Name => '%s'.", *newSrv.Ip, *newSrv.Port, *newSrv.Name, *newSrv.CurUsers, *newSrv.MaxUsers, *newSrv.MapName)
		}

		if len(serversToUpdate) < 1 {
			utils.DebugMsg(3, cfg, "[VMS] Found no servers to update for app ID %d.", appId)

			Respin(cfg)

			continue
		}

		// Add/update servers.
		if !cfg.Vms.RecvOnly {
			cnt, err := servers.AddServers(cfg, serversToUpdate, cfg.Vms.AddOnly)

			if err != nil {
				utils.DebugMsg(1, cfg, "[VMS] Failed to add/update servers for app ID %d due to error.", appId)
				utils.DebugMsg(1, cfg, err.Error())

				Respin(cfg)

				continue
			}

			utils.DebugMsg(2, cfg, "[VMS] Added/Updated %d servers!", cnt)
		}

		Respin(cfg)
	}
}
