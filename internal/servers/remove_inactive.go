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

type Resp struct {
	Count int `json:"count"`
}

func RemoveInactive(cfg *config.Config) {
	for {
		// Initialize error.
		var err error

		interval := time.Duration(cfg.RemoveInactive.Interval) * time.Second

		if interval < 1 {
			interval = 60 * time.Second
		}

		if !cfg.RemoveInactive.Enabled {
			time.Sleep(interval)

			continue
		}

		var inactive_time *int = nil

		if cfg.RemoveInactive.InactiveTime > 0 {
			inactive_time = new(int)
			*inactive_time = cfg.RemoveInactive.InactiveTime
		}

		params := url.Values{}

		if inactive_time != nil {
			params.Add("time", strconv.Itoa(*inactive_time))
		}

		// Format URL.
		url := fmt.Sprintf("%s%s", cfg.Api.Host, "/api/servers/removeinactive")

		if len(params) > 0 {
			url = fmt.Sprintf("%s?%s", url, params.Encode())
		}

		// Create new HTTP client.
		client := http.Client{
			Timeout: time.Duration(cfg.RemoveInactive.Timeout) * time.Second,
		}

		// Create new request and check for error.
		req, err := http.NewRequest("POST", url, nil)

		if err != nil {
			utils.DebugMsg(1, cfg, "[REM_INACTIVE] Request failed due to error :: %s", err.Error())

			time.Sleep(interval)

			continue
		}

		// Set authorization header.
		req.Header.Add("Authorization", cfg.Api.Authorization)

		// Set content type to JSON.
		req.Header.Add("Content-Type", "application/json")

		// Perform request and check for error.
		res, err := client.Do(req)

		if err != nil {
			utils.DebugMsg(1, cfg, "[REM_INACTIVE] Request failed due to error :: %s", err.Error())

			time.Sleep(interval)

			continue
		}

		if res.StatusCode != 200 {
			utils.DebugMsg(1, cfg, "[REM_INACTIVE] Request status code doesn't equal 200 (%d)", res.StatusCode)

			time.Sleep(interval)

			continue
		}

		// Read result.
		resBytes, err := io.ReadAll(res.Body)

		if err != nil {
			utils.DebugMsg(1, cfg, "[REM_INACTIVE] Failed to read response body due to error :: %s", err.Error())

			time.Sleep(interval)

			continue
		}

		// Close body.
		res.Body.Close()

		ret := Resp{}
		err = json.Unmarshal(resBytes, &ret)

		if err != nil {
			utils.DebugMsg(1, cfg, "[REM_INACTIVE] Failed to unmarshal response due to error :: %s", err.Error())
		} else {
			utils.DebugMsg(4, cfg, "[REM_INACTIVE] Removed %d inactive servers!", ret.Count)
		}

		time.Sleep(interval)
	}
}
