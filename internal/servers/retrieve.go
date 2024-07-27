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

type RetrieveResp struct {
	Count   int      `json:"count"`
	Servers []Server `json:"servers"`
	Message string   `json:"message"`
}

type RetrieveCountResp struct {
	Count int `json:"count"`
}

func RetrieveServers(cfg *config.Config, platform_id *uint, limit *int, visibleOnly bool) ([]Server, error) {
	// Initiailize what we're returning.
	servers := []Server{}
	var err error

	// Create new client with specific timeout.
	client := http.Client{
		Timeout: time.Duration(cfg.Api.Timeout) * time.Second,
	}

	// Create query parameters.
	params := url.Values{}

	if platform_id != nil {
		params.Add("platformId", strconv.Itoa(int(*platform_id)))
	}

	if visibleOnly {
		params.Add("visibleOnly", "1")
	}

	if limit != nil {
		params.Add("limit", strconv.Itoa(*limit))
	}

	// Format URL.
	url := fmt.Sprintf("%s%s", cfg.Api.Host, "/api/servers/get")

	if len(params) > 0 {
		url = fmt.Sprintf("%s?%s", url, params.Encode())
	}

	utils.DebugMsg(4, cfg, "[SERVERS] Sending API request to => '%s'.", url)

	// Create a new request and check for error.
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return servers, err
	}

	// Set authorization header token.
	req.Header.Add("Authorization", cfg.Api.Authorization)

	// Set content type to JSON.
	req.Header.Add("Content-Type", "application/json")

	// Send response and check for error.
	res, err := client.Do(req)

	if err != nil {
		return servers, err
	}

	if res.StatusCode != 200 {
		return servers, fmt.Errorf("status code did not return 200 (%d)", res.StatusCode)
	}

	// Make sure we close body at end.
	defer res.Body.Close()

	// Read response into byte array and check for error.
	b, err := io.ReadAll(res.Body)

	if err != nil {
		return servers, err
	}

	// Initialize retrieve response.
	retrieveResp := RetrieveResp{}

	// Unmarshal byte array into servers structure and return result.
	err = json.Unmarshal(b, &retrieveResp)

	if err != nil {
		return servers, err
	}

	// Assign servers to response.
	servers = retrieveResp.Servers

	return servers, err
}

func RetrieveServerCount(cfg *config.Config, ip string) (int, error) {
	var err error
	cnt := 0

	// Create new client with specific timeout.
	client := http.Client{
		Timeout: time.Duration(cfg.Api.Timeout) * time.Second,
	}

	// Create query parameters.
	params := url.Values{}

	params.Add("ip", ip)
	params.Add("countOnly", "1")

	// Format URL.
	url := fmt.Sprintf("%s%s?%s", cfg.Api.Host, "/api/servers/get", params.Encode())

	utils.DebugMsg(4, cfg, "[SERVERS] Sending server count API request to => '%s'.", url)

	// Create a new request and check for error.
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return cnt, err
	}

	// Set authorization header token.
	req.Header.Add("Authorization", cfg.Api.Authorization)

	// Set content type to JSON.
	req.Header.Add("Content-Type", "application/json")

	// Send response and check for error.
	res, err := client.Do(req)

	if err != nil {
		return cnt, err
	}

	if res.StatusCode != 200 {
		return cnt, fmt.Errorf("status code did not return 200 (%d)", res.StatusCode)
	}

	// Make sure we close body at end.
	defer res.Body.Close()

	// Read response into byte array and check for error.
	b, err := io.ReadAll(res.Body)

	if err != nil {
		return cnt, err
	}

	// Initialize retrieve response.
	retrieveResp := RetrieveCountResp{}

	// Unmarshal byte array into servers structure and return result.
	err = json.Unmarshal(b, &retrieveResp)

	if err != nil {
		return cnt, err
	}

	// Assign servers to response.
	cnt = retrieveResp.Count

	return cnt, err
}
