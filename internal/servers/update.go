package servers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bestserversio/spy/internal/config"
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

type ServersWithWhere struct {
	Where WhereClause `json:"where"`
	Server
}

type UpdateBody struct {
	Servers []ServersWithWhere `json:"servers"`
}

func UpdateServers(cfg config.Config, servers []ServersWithWhere) (int, error) {
	// Initialize error.
	var err error
	cnt := 0

	// Format URL.
	url := fmt.Sprintf("%s%s", cfg.Api.Host, "/api/servers/update")

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

	// Read result.
	resBytes, err := io.ReadAll(res.Body)

	// Unmarshal as update response and check for error.
	upResp := UpdateResp{}
	err = json.Unmarshal(resBytes, &upResp)

	if err != nil {
		return cnt, err
	}

	// Set count to server count.
	cnt = upResp.ServerCount

	// Verbose.
	if cfg.Verbose > 0 {
		// Print every server.
		if cfg.Verbose > 2 {
			for i, srv := range upResp.Servers {
				fmt.Printf("[3][%d] Updating server %s:%d (%s). URL => %s. Platform ID => %d.\n", i, srv.Ip, srv.Port, srv.Ip6, srv.Url, srv.PlatformId)
			}
		}

		// Print generic details on response.
		if cfg.Verbose > 1 {
			fmt.Println("[2][PUT] Found", upResp.ServerCount, "servers!")
			fmt.Println("[2][PUT] Message =>", upResp.Message)
		}

		// If we have errors, display them.
		if upResp.ErrorCount > 0 {
			fmt.Println("[1][PUT] Found", upResp.ErrorCount, "errors. Listing...")

			for i, msg := range upResp.Errors {
				fmt.Printf("[1][Err %d] %s\n", i, msg)
			}
		}
	}

	return cnt, err
}
