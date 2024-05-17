package platforms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bestserversio/spy/internal/config"
)

type RetrieveResp struct {
	Count     int        `json:"count"`
	Platforms []Platform `json:"platforms"`
	Message   string     `json:"message"`
}

func RetrievePlatforms(cfg config.Config) ([]Platform, error) {
	// Initiailize what we're returning.
	platforms := []Platform{}
	var err error

	// Create new client with specific timeout.
	client := http.Client{
		Timeout: time.Duration(cfg.Api.Timeout) * time.Second,
	}

	// Format URL.
	url := fmt.Sprintf("%s%s", cfg.Api.Host, "/api/platforms/get")

	// Create a new request and check for error.
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return platforms, err
	}

	// Set authorization header token.
	req.Header.Add("Authorization", cfg.Api.Authorization)

	// Set content type to JSON.
	req.Header.Add("Content-Type", "application/json")

	// Send response and check for error.
	res, err := client.Do(req)

	if err != nil {
		return platforms, err
	}

	// Make sure we close body at end.
	defer res.Body.Close()

	// Read response into byte array and check for error.
	b, err := io.ReadAll(res.Body)

	if err != nil {
		return platforms, err
	}

	// Initialize retrieve response.
	retrieveResp := RetrieveResp{}

	// Unmarshal byte array into platforms structure and return result.
	err = json.Unmarshal(b, &retrieveResp)

	// Assign platforms to response.
	platforms = retrieveResp.Platforms

	// Verbose.
	if cfg.Verbose > 1 {
		fmt.Println("[GET] Found", retrieveResp.Count, "platforms!")
		fmt.Println("[GET] Response message =>", retrieveResp.Message)
	}

	return platforms, err
}
