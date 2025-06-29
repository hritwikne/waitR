package app

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ServerConfig struct {
	ListenAddress string        `json:"listen_address"`
	StaticPath    string        `json:"static_path"`
	RootFile      string        `json:"root_file"`
	Timeouts      TimeoutConfig `json:"timeouts"`
	ProxyRoutes   []ProxyRoute  `json:"proxy_routes"`
}

type ProxyRoute struct {
	Prefix   string   `json:"prefix"`
	Backends []string `json:"backends"`
}

type TimeoutConfig struct {
	ReadHeaderTimeout int `json:"read_header_timeout_ms"`
	ReadTimeout       int `json:"read_timeout_ms"`
	WriteTimeout      int `json:"write_timeout_ms"`
	IdleTimeout       int `json:"idle_timeout_ms"`
}

func LoadConfig(path string) (*ServerConfig, error) {
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

func validateListenAddress(addr string) error {
	if strings.TrimSpace(addr) == "" {
		return fmt.Errorf("listen_address is required and cannot be empty")
	}

	// Check if address can be resolved
	if _, err := net.ResolveTCPAddr("tcp", addr); err != nil {
		return fmt.Errorf("invalid listen_address: %v", err)
	}

	// Split host and port
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("invalid listen_address format: %w", err)
	}

	// Ensure port is numeric
	if _, err := strconv.Atoi(portStr); err != nil {
		return fmt.Errorf("port must be numeric in listen_address: %s", addr)
	}

	return nil
}

func validateStaticPath(staticPath, rootFile string) error {
	trimmed := strings.TrimSpace(staticPath)
	if trimmed == "" {
		return fmt.Errorf("static_path is required and cannot be blank")
	}

	info, err := os.Stat(trimmed)
	if err != nil {
		return fmt.Errorf("static_path does not exist: %v", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("static_path must be a directory")
	}

	// If rootFile is empty, fallback to common defaults
	rootCandidates := []string{"index.html", "index.htm"}
	if rootFile != "" {
		rootCandidates = []string{rootFile}
	}

	found := false
	for _, name := range rootCandidates {
		if _, err := os.Stat(filepath.Join(trimmed, name)); err == nil {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("none of the root files (%v) exist in static_path", rootCandidates)
	}

	return nil
}

func ValidateConfig(cfg *ServerConfig) error {
	if err := validateListenAddress(cfg.ListenAddress); err != nil {
		return err
	}

	if err := validateStaticPath(cfg.StaticPath, cfg.RootFile); err != nil {
		return err
	}

	if cfg.Timeouts.ReadTimeout <= 0 {
		return fmt.Errorf("read_timeout_ms must be greater than 0")
	}
	if cfg.Timeouts.WriteTimeout <= 0 {
		return fmt.Errorf("write_timeout_ms must be greater than 0")
	}
	if cfg.Timeouts.IdleTimeout <= 0 {
		return fmt.Errorf("idle_timeout_ms must be greater than 0")
	}
	if cfg.Timeouts.ReadHeaderTimeout <= 0 {
		return fmt.Errorf("read_header_timeout_ms must be greater than 0")
	}

	for _, route := range cfg.ProxyRoutes {
		if route.Prefix == "" {
			return fmt.Errorf("each proxy_route must have a prefix")
		}
		if len(route.Backends) == 0 {
			return fmt.Errorf("proxy_route for prefix %s must have at least one backend", route.Prefix)
		}
	}

	return nil
}
