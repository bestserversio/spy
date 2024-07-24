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
)

type Resp struct {
	Count int `json:"count"`
}

func RemoveInactive(cfg *config.Config, inactive_time *int) (int, error) {
	// Initialize error.
	var err error

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
		return 0, err
	}

	// Set authorization header.
	req.Header.Add("Authorization", cfg.Api.Authorization)

	// Set content type to JSON.
	req.Header.Add("Content-Type", "application/json")

	// Perform request and check for error.
	res, err := client.Do(req)

	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	// Read result.
	resBytes, err := io.ReadAll(res.Body)

	if err != nil {
		return 0, err
	}

	ret := Resp{}
	err = json.Unmarshal(resBytes, &ret)

	if err != nil {
		return 0, err
	}

	return ret.Count, err
}
