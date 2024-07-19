package vms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bestserversio/spy/internal/config"
)

const STEAM_API_URL = "https://api.steampowered.com/IGameServersService/GetServerList/v1/"

func RetrieveServers(cfg *config.Config, appId int) ([]Server, error) {
	var servers []Server
	var err error = nil

	if !cfg.Vms.Enabled {
		return servers, err
	}

	// Create HTTP client with timeout.
	client := http.Client{
		Timeout: time.Duration(cfg.Vms.Timeout) * time.Second,
	}

	// Start building filters string.
	filters := fmt.Sprintf("\\appid\\%d", appId)

	// Add empty if exclude empty is set.
	if cfg.Vms.ExcludeEmpty {
		filters = fmt.Sprintf("%s\\empty\\1", filters)
	}

	// Compile URL.
	url := fmt.Sprintf("%s?key=%s&filter=%s&limit=%d", STEAM_API_URL, cfg.Vms.ApiToken, filters, cfg.Vms.Limit)

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
		return servers, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return servers, fmt.Errorf("status code returned isn't 200 (%d)", res.StatusCode)
	}

	// Read response.
	b, err := io.ReadAll(res.Body)

	if err != nil {
		return servers, err
	}

	retrieveResp := Response{}

	err = json.Unmarshal(b, &retrieveResp)

	servers = retrieveResp.Response.Servers

	return servers, err
}
