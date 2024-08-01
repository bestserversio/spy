package servers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bestserversio/spy/internal/config"
	"github.com/bestserversio/spy/internal/utils"
)

type TimedOutResp struct {
	TimedOut int `json:"timedOut"`
}

func RemoveTimedOut(cfg *config.Config) {
	for {
		// Initialize error.
		var err error

		interval := time.Duration(cfg.RemoveTimedOut.Interval) * time.Second

		if interval < 1 {
			interval = 60 * time.Second
		}

		if !cfg.RemoveTimedOut.Enabled {
			time.Sleep(interval)

			continue
		}

		params := url.Values{}

		platform_ids := cfg.RemoveTimedOut.PlatformIds

		if len(platform_ids) > 0 {
			str_ids := make([]string, len(platform_ids))

			for i, id := range platform_ids {
				str_ids[i] = strconv.Itoa(id)
			}

			params.Add("platformIds", strings.Join(str_ids, ","))
		}

		timed_out_time := cfg.RemoveTimedOut.TimedOutTime

		if timed_out_time < 1 {
			timed_out_time = 3600
		}

		params.Add("timeout", strconv.Itoa(timed_out_time))

		// Format URL.
		url := fmt.Sprintf("%s%s?%s", cfg.Api.Host, "/api/servers/removetimedout", params.Encode())

		// Create new HTTP client.
		client := http.Client{
			Timeout: time.Duration(cfg.RemoveTimedOut.Timeout) * time.Second,
		}

		// Create new request and check for error.
		req, err := http.NewRequest("POST", url, nil)

		if err != nil {
			utils.DebugMsg(1, cfg, "[REM_TIMED_OUT] Request failed due to error :: %s", err.Error())

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
			utils.DebugMsg(1, cfg, "[REM_TIMED_OUT] Request failed due to error :: %s", err.Error())

			time.Sleep(interval)

			continue
		}

		if res.StatusCode != 200 {
			utils.DebugMsg(1, cfg, "[REM_TIMED_OUT] Request status code doesn't equal 200 (%d)", res.StatusCode)

			time.Sleep(interval)

			continue
		}

		// Read result.
		resBytes, err := io.ReadAll(res.Body)

		if err != nil {
			utils.DebugMsg(1, cfg, "[REM_TIMED_OUT] Failed to read response body due to error :: %s", err.Error())

			time.Sleep(interval)

			continue
		}

		// Close body.
		res.Body.Close()

		ret := TimedOutResp{}
		err = json.Unmarshal(resBytes, &ret)

		if err != nil {
			utils.DebugMsg(1, cfg, "[REM_TIMED_OUT] Failed to unmarshal response due to error :: %s", err.Error())
		} else {
			utils.DebugMsg(4, cfg, "[REM_TIMED_OUT] Removed %d timed out servers!", ret.TimedOut)
		}

		time.Sleep(interval)
	}
}
