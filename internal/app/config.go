package app

import (
	"encoding/json"
	"fmt"
	"os"
)

type ServerConfig struct {
	ListenAddress string       `json:"listen_address"`
	StaticPath    string       `json:"static_path"`
	ProxyRoutes   []ProxyRoute `json:"proxy_routes"`
}

type ProxyRoute struct {
	Prefix   string   `json:"prefix"`
	Backends []string `json:"backends"`
}

func LoadConfig() (*ServerConfig, error) {
	path := "config/server.json"
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open config file: %w", err)
	}
	defer file.Close()

	var cfg ServerConfig
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}
