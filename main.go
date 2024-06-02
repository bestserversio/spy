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

func DoVms(cfg *config.Config) {
	if !cfg.Vms.Enabled {
		return
	}

	// Retrieve random app ID.
	rand.Seed(time.Now().UnixNano())

	randIndex := rand.Intn(len(cfg.Vms.AppIds))

	appId := cfg.Vms.AppIds[randIndex]

	utils.DebugMsg(4, cfg.Verbose, "[VMS] Using (random) app ID %d.", appId)

	// Retrieve servers.
	allServers, err := vms.RetrieveServers(cfg, appId)

	if err != nil {
		utils.DebugMsg(1, cfg.Verbose, "[VMS] Failed to retrieve servers for app %d due to error.", appId)
		utils.DebugMsg(1, cfg.Verbose, err.Error())

		return
	}

	utils.DebugMsg(4, cfg.Verbose, "[VMS] Retrieved %d servers from app ID %d.", len(allServers), appId)

	var serversToUpdate []servers.Server

	// Loop through all servers from VMS result.
	for _, srv := range allServers {
		// Create new servers object from servers package and assign basic values.
		newSrv := servers.Server{
			HostName: new(string),
			CurUsers: new(int),
			MaxUsers: new(int),
			MapName:  new(string),
			Ip:       new(string),
			Port:     new(int),
		}
		*newSrv.HostName = srv.HostName
		*newSrv.CurUsers = srv.Players
		*newSrv.MaxUsers = srv.MaxPlayers
		*newSrv.MapName = srv.Map

		// We need to split IP and port from address.
		split := strings.Split(srv.Address, ":")

		if len(split) > 1 {
			*newSrv.Ip = split[0]

			portNum, err := strconv.Atoi(split[1])

			if err != nil {
				utils.DebugMsg(1, cfg.Verbose, "")

				return
			}

			*newSrv.Port = portNum
		}

		serversToUpdate = append(serversToUpdate, newSrv)
	}

	if len(serversToUpdate) < 1 {
		utils.DebugMsg(3, cfg.Verbose, "[VMS] Found no servers for app ID %d.", appId)

		return
	}

	// Add/update servers.
	cnt, err := servers.AddServers(cfg, serversToUpdate, true)

	if err != nil {
		utils.DebugMsg(1, cfg.Verbose, "[VMS] Failed to add/update servers for app ID %d due to error.", appId)
		utils.DebugMsg(1, cfg.Verbose, err.Error())

		return
	}

	utils.DebugMsg(2, cfg.Verbose, "[VMS] Added/Updated %d servers!", cnt)
}

func ScannerRespin(scanner *config.Scanner) {
	min := scanner.MinWait
	max := scanner.MaxWait

	wait := rand.Intn(max-min+1) + min

	time.Sleep(time.Second * time.Duration(wait))
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
		allServers, err := servers.RetrieveServers(cfg, &app_id)

		if err != nil {
			utils.DebugMsg(1, cfg.Verbose, "[SCANNER] Failed to retrieve servers using app ID '%d' due to error.", app_id)
			utils.DebugMsg(1, cfg.Verbose, err.Error())

			ScannerRespin(scanner)
		}

		utils.DebugMsg(4, cfg.Verbose, "[SCANNER] Found %d to update from API for app ID '%d'!", len(allServers), app_id)

		// Loop through each server and update.
		for i := 0; i < len(allServers); i++ {
			srv := &allServers[i]

			switch query_type {
			case 0:
				if srv.Ip == nil || srv.Port == nil {
					utils.DebugMsg(1, cfg.Verbose, "[SCANNER] Missing IP/port for server. Skipping.")

					continue
				}

				// Try querying server with A2S and check for error..
				err = scanners.QueryA2s(srv)

				if err != nil {
					utils.DebugMsg(1, cfg.Verbose, "[SCANNER] Failed to query A2S server '%s:%d' due to error.", *srv.Ip, *srv.Port)
					utils.DebugMsg(1, cfg.Verbose, err.Error())
				} else {
					// Set maxplayers to 0 to indicate offline.
					if srv.MaxUsers == nil {
						srv.MaxUsers = new(int)
					}

					*srv.MaxUsers = 0
				}
			}
		}

		// Update servers.
		cnt, err := servers.AddServers(cfg, allServers, false)

		if err != nil {
			utils.DebugMsg(1, cfg.Verbose, "[SCANNER] Failed to update servers for app ID '"+strconv.Itoa(app_id)+"' due to error.")
			utils.DebugMsg(1, cfg.Verbose, err.Error())
		} else {
			utils.DebugMsg(3, cfg.Verbose, "[SCANNER] Updated "+strconv.Itoa(cnt)+" servers for app ID '"+strconv.Itoa(app_id)+"'!")
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

	// Create VMS timer.
	utils.DebugMsg(4, cfg.Verbose, "[MAIN] Creating VMS timer...")
	vmsTicker := time.NewTicker(time.Second * time.Duration(cfg.Vms.Interval))
	vmsQuit := make(chan struct{})

	go func() {
		for {
			select {
			case <-vmsTicker.C:
				DoVms(&cfg)
			case <-vmsQuit:
				vmsTicker.Stop()

				return
			}
		}
	}()

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
