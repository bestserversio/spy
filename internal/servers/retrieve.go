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

func RetrieveServers(cfg *config.Config, platform_id *int, limit *int) ([]Server, error) {
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
		params.Add("platformId", strconv.Itoa(*platform_id))
	}

	if limit != nil {
		params.Add("limit", strconv.Itoa(*limit))
	}

	// Format URL.
	url := fmt.Sprintf("%s%s", cfg.Api.Host, "/api/servers/get")

	if len(params) > 0 {
		url = fmt.Sprintf("%s?%s", url, params.Encode())
	}

	utils.DebugMsg(4, cfg.Verbose, "[SERVERS] Sending API request to => '%s'.", url)

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

	// Assign servers to response.
	servers = retrieveResp.Servers

	// Verbose.
	if cfg.Verbose > 1 {
		fmt.Println("[GET] Found", retrieveResp.Count, "servers!")
		fmt.Println("[GET] Response message =>", retrieveResp.Message)
	}

	return servers, err
}
