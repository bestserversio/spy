package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bestserversio/spy/internal/config"
	"github.com/bestserversio/spy/internal/scanners"
	"github.com/bestserversio/spy/internal/servers"
	"github.com/bestserversio/spy/internal/utils"
	"github.com/bestserversio/spy/internal/vms"
)

const VERSION = "1.0.0"
const HELPMENU = `
./spy [OPTIONS]\n\n
-l --list => Prints config and exits.\n
-v --version => Prints version and exits.\n
-h --help => Prints help menu and exits.\n
-c --cfg => Path to config file.\n
`

func Respin(min int, max int) {
	wait := rand.Intn(max-min+1) + min

	time.Sleep(time.Second * time.Duration(wait))
}

func VmsRespin(cfg *config.Config) {
	Respin(cfg.Vms.MinWait, cfg.Vms.MaxWait)
}

func DoVms(cfg *config.Config) {
	for {
		if !cfg.Vms.Enabled {
			VmsRespin(cfg)

			continue
		}

		// Retrieve random app ID.
		rand.Seed(time.Now().UnixNano())

		randIndex := rand.Intn(len(cfg.Vms.AppIds))

		appId := cfg.Vms.AppIds[randIndex]

		// Get platform ID.
		platform_id := utils.AppIdToPlatformId(cfg, appId)

		utils.DebugMsg(4, cfg.Verbose, "[VMS] Using (random) app ID '%d'.", appId)

		// Retrieve servers.
		allServers, err := vms.RetrieveServers(cfg, appId)

		if err != nil {
			utils.DebugMsg(1, cfg.Verbose, "[VMS] Failed to retrieve servers for app ID '%d' due to error.", appId)
			utils.DebugMsg(1, cfg.Verbose, err.Error())

			VmsRespin(cfg)

			continue
		}

		utils.DebugMsg(3, cfg.Verbose, "[VMS] Retrieved %d servers from app ID '%d'.", len(allServers), appId)

		var serversToUpdate []servers.Server

		// Loop through all servers from VMS result.
		for _, srv := range allServers {
			// Create new servers object from servers package and assign basic values.
			newSrv := servers.Server{
				Online:      new(bool),
				HostName:    new(string),
				CurUsers:    new(int),
				MaxUsers:    new(int),
				MapName:     new(string),
				Ip:          new(string),
				Port:        new(int),
				PlatformId:  new(int),
				Os:          new(string),
				LastQueried: new(string),
				Region:      new(string),
			}
			*newSrv.Online = true
			*newSrv.HostName = srv.HostName
			*newSrv.CurUsers = srv.Players
			*newSrv.MaxUsers = srv.MaxPlayers
			*newSrv.MapName = srv.Map
			*newSrv.PlatformId = platform_id
			*newSrv.Os = srv.Os
			*newSrv.Region = utils.GetRegion(srv.Region)

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
					utils.DebugMsg(1, cfg.Verbose, "")

					VmsRespin(cfg)

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

			// Append to servers to update array.
			serversToUpdate = append(serversToUpdate, newSrv)

			utils.DebugMsg(4, cfg.Verbose, "[VMS] Adding server '%s:%d'. Host Name => '%s'. Players => %d. Max Players => %d. Map Name => %s", *newSrv.Ip, *newSrv.Port, *newSrv.HostName, *newSrv.CurUsers, *newSrv.MaxUsers, *newSrv.MapName)
		}

		if len(serversToUpdate) < 1 {
			utils.DebugMsg(3, cfg.Verbose, "[VMS] Found no servers for app ID %d.", appId)

			VmsRespin(cfg)

			continue
		}

		// Add/update servers.
		if !cfg.Vms.RecvOnly {
			cnt, err := servers.AddServers(cfg, serversToUpdate, false)

			if err != nil {
				utils.DebugMsg(1, cfg.Verbose, "[VMS] Failed to add/update servers for app ID %d due to error.", appId)
				utils.DebugMsg(1, cfg.Verbose, err.Error())

				VmsRespin(cfg)

				continue
			}

			utils.DebugMsg(2, cfg.Verbose, "[VMS] Added/Updated %d servers!", cnt)
		}

		VmsRespin(cfg)
	}
}

func ScannerRespin(scanner *config.Scanner) {
	Respin(scanner.MinWait, scanner.MaxWait)
}

func HandleScanner(cfg *config.Config, scanner *config.Scanner, idx int) {
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

			ScannerRespin(scanner)
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
				err = scanners.QueryA2s(srv)

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

		ScannerRespin(scanner)
	}
}

func main() {
	// Command line options and parse command line.
	var list bool
	var version bool
	var help bool

	flag.BoolVar(&list, "l", false, "Prints config settings and exits.")
	flag.BoolVar(&list, "list", false, "Prints config settings and exits.")

	flag.BoolVar(&version, "v", false, "Prints version and exits.")
	flag.BoolVar(&version, "version", false, "Prints number and exits.")

	flag.BoolVar(&help, "h", false, "Prints help menu and exits.")
	flag.BoolVar(&help, "help", false, "Prints help menu and exits.")

	cfgPath := flag.String("cfg", "/etc/bestservers/spy.json", "Path to config file.")

	flag.Parse()

	// Check for version.
	if version {
		fmt.Println(VERSION)

		os.Exit(0)
	}

	// Check for help menu.
	if version {
		fmt.Print(HELPMENU)

		os.Exit(0)
	}

	// Initialize config.
	cfg := config.Config{}

	// Load defaults.
	cfg.LoadDefaults()

	// Attempt to load config.
	err := cfg.LoadFromFs(*cfgPath)

	if err != nil {
		fmt.Println("Error loading config file. Resorting to defaults...")
		fmt.Println(err)
	}

	utils.DebugMsg(4, cfg.Verbose, "[MAIN] Config loaded...")

	// Check if we want to print our config settings.
	if list {
		cfg.PrintConfig()

		os.Exit(0)
	}

	// Create VMS.
	go DoVms(&cfg)

	// Create scanners.
	for i, s := range cfg.Scanners {
		go HandleScanner(&cfg, &s, i)
	}

	// Make a signal.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

	<-sigc

	utils.DebugMsg(0, cfg.Verbose, "[MAIN] Exiting...")
	os.Exit(0)
}
