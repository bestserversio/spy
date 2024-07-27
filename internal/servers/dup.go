package servers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/bestserversio/spy/internal/config"
	"github.com/bestserversio/spy/internal/utils"
)

type DupResp struct {
	Filtered int `json:"filtered"`
}

func CheckForDups(cfg *config.Config) {
	for {
		var err error
		cnt := 0

		interval := time.Duration(cfg.RemoveDups.Interval) * time.Second

		if interval < 1 {
			interval = 60 * time.Second
		}

		if !cfg.RemoveDups.Enabled {
			time.Sleep(interval)

			continue
		}

		// Create new client with specific timeout.
		client := http.Client{
			Timeout: time.Duration(cfg.RemoveInactive.Timeout) * time.Second,
		}

		// Create query parameters.
		params := url.Values{}

		params.Add("limit", strconv.Itoa(cfg.RemoveDups.Limit))
		params.Add("maxServers", strconv.Itoa(cfg.RemoveDups.MaxServers))

		// Format URL.
		url := fmt.Sprintf("%s%s?%s", cfg.Api.Host, "/api/servers/removedups", params.Encode())

		utils.DebugMsg(4, cfg, "[DUPS] Sending duplicate API request to => '%s'.", url)

		// Create a new request and check for error.
		req, err := http.NewRequest("POST", url, nil)

		if err != nil {
			utils.DebugMsg(1, cfg, "[DUPS] Request failed due to error :: %s", err.Error())

			time.Sleep(interval)

			continue
		}

		// Set authorization header token.
		req.Header.Add("Authorization", cfg.Api.Authorization)

		// Set content type to JSON.
		req.Header.Add("Content-Type", "application/json")

		// Send response and check for error.
		res, err := client.Do(req)

		if err != nil {
			utils.DebugMsg(1, cfg, "[DUPS] Request failed due to error :: %s", err.Error())

			time.Sleep(interval)

			continue
		}

		if res.StatusCode != 200 {
			utils.DebugMsg(1, cfg, "[DUPS] Request failed due to error :: status code did not return 200 (%d)", res.StatusCode)

			time.Sleep(interval)

			continue
		}

		// Read response into byte array and check for error.
		b, err := io.ReadAll(res.Body)

		if err != nil {
			utils.DebugMsg(1, cfg, "[DUPS] Failed to read response body :: %s", err.Error())

			time.Sleep(interval)

			continue
		}

		// Make sure we close body at end.
		res.Body.Close()

		// Initialize retrieve response.
		retrieveResp := DupResp{}

		// Unmarshal byte array into servers structure and return result.
		err = json.Unmarshal(b, &retrieveResp)

		if err != nil {
			utils.DebugMsg(1, cfg, "[DUPS] Failed to unmarshal JSON due to error :: %s", err.Error())

			time.Sleep(interval)

			continue
		}

		// Assign servers to response.
		cnt = retrieveResp.Filtered

		utils.DebugMsg(4, cfg, "[DUPS] Filtered %d servers.", cnt)

		time.Sleep(interval)
	}
}
