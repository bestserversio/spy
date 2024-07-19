package config

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func (cfg *Config) LoadFromFs(path string) error {
	f, err := os.Open(path)

	if err != nil {
		return err
	}

	defer f.Close()

	stat, err := f.Stat()

	if err != nil {
		return err
	}

	data := make([]byte, stat.Size())

	_, err = f.Read(data)

	if err != nil {
		return err
	}

	err = cfg.Load(string(data))

	return err
}

func (cfg *Config) LoadFromWeb() (string, error) {
	var err error

	if !cfg.WebApi.Enabled {
		return "", err
	}

	client := http.Client{
		Timeout: time.Duration(cfg.WebApi.Timeout) * time.Second,
	}

	url := fmt.Sprintf("%s%s", cfg.WebApi.Host, cfg.WebApi.Endpoint)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", cfg.WebApi.Authorization)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	err = cfg.Load(string(b))

	return string(b), err
}
