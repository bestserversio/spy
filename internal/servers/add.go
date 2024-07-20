package servers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/bestserversio/spy/internal/config"
	"github.com/bestserversio/spy/internal/utils"
)

type UpdateResp struct {
	ServerCount int      `json:"serverCount"`
	Servers     []Server `json:"servers"`
	ErrorCount  int      `json:"errorCount"`
	Errors      []string `json:"errors"`
	Message     string   `json:"message"`
}

type WhereClause struct {
	Id   int    `json:"id"`
	Url  string `json:"url"`
	Ip   string `json:"ip"`
	Ip6  string `json:"ip6"`
	Port int    `json:"port"`
}

type ServerWithWhere struct {
	Where WhereClause `json:"where"`
	Server
}

type UpdateBody struct {
	Servers []Server `json:"servers"`
}

func AddServers(cfg *config.Config, servers []Server, addOnly bool) (int, error) {
	// Initialize error.
	var err error
	cnt := 0

	params := url.Values{}

	if addOnly {
		params.Add("addonly", "1")
	}

	// Format URL.
	url := fmt.Sprintf("%s%s", cfg.Api.Host, "/api/servers/add")

	if len(params) > 0 {
		url = fmt.Sprintf("%s?%s", url, params.Encode())
	}

	utils.DebugMsg(4, cfg.Verbose, "[SRV_AS] Using API url '%s'.", url)

	// Create new HTTP client.
	client := http.Client{
		Timeout: time.Duration(cfg.Api.Timeout) * time.Second,
	}

	// Create request body.
	upBody := UpdateBody{
		Servers: servers,
	}

	// Convert update body to JSON string and check for error.
	upBytes, err := json.Marshal(&upBody)

	if err != nil {
		return cnt, err
	}

	utils.DebugMsg(6, cfg.Verbose, "[SERVERS] Using request body: %s", string(upBytes))

	// Create new request and check for error.
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(upBytes))

	if err != nil {
		return cnt, err
	}

	// Set authorization header.
	req.Header.Add("Authorization", cfg.Api.Authorization)

	// Set content type to JSON.
	req.Header.Add("Content-Type", "application/json")

	// Perform request and check for error.
	res, err := client.Do(req)

	if err != nil {
		return cnt, err
	}

	// Close body at end.
	defer res.Body.Close()

	// Check for status code.
	if res.StatusCode != 200 {
		err = fmt.Errorf("[SERVERS] Failed to update/add servers due to invalid status code %d", res.StatusCode)

		return cnt, err
	}

	// Read result.
	resBytes, err := io.ReadAll(res.Body)

	if err != nil {
		return cnt, err
	}

	// Unmarshal as update response and check for error.
	upResp := UpdateResp{}
	err = json.Unmarshal(resBytes, &upResp)

	if err != nil {
		return cnt, err
	}

	// Set count to server count.
	cnt = upResp.ServerCount

	return cnt, err
}
