package vms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bestserversio/spy/internal/config"
	"github.com/bestserversio/spy/internal/utils"
)

const STEAM_API_URL = "https://api.steampowered.com/IGameServersService/GetServerList/v1/"

func RetrieveServers(cfg config.Config, appId int) ([]Server, error) {
	var servers []Server
	var err error = nil

	if !cfg.Vms.Enabled {
		return servers, err
	}

	// Create HTTP client with timeout.
	client := http.Client{
		Timeout: time.Duration(cfg.Vms.Timeout) * time.Second,
	}

	// Compile URL.
	url := fmt.Sprintf("%s?key=%s&filter=\\appid\\%d", STEAM_API_URL, cfg.Vms.ApiToken, appId)

	// Create response and check.
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return servers, err
	}

	// Only accept JSON.
	req.Header.Add("Content-Type", "application/json")

	// Send response and check.
	res, err := client.Do(req)

	if err != nil {
		return servers, nil
	}

	defer res.Body.Close()

	// Read response.
	b, err := io.ReadAll(res.Body)

	if err != nil {
		return servers, err
	}

	retrieveResp := Response{}

	err = json.Unmarshal(b, &retrieveResp)

	servers = retrieveResp.Response.Servers

	// Some debug.
	utils.DebugMsg(2, cfg.Verbose, "[VMS] Retrieved %d servers from app ID %d.", len(servers), appId)

	return servers, err
}
