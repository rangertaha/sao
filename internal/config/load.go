package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// EnsureConfig creates a config file from the embedded default if missing.
func EnsureConfig(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return fmt.Errorf("config path cannot be empty")
	}

	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	if err := os.WriteFile(path, DefaultConfigBytes(), 0o644); err != nil {
		return fmt.Errorf("write default config: %w", err)
	}

	return nil
}

// Load reads and parses the HCL config file at path.
func Load(path string) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	parser := hclparse.NewParser()
	hclFile, diags := parser.ParseHCL(content, path)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parse config: %s", diags.Error())
	}

	var cfg Config
	diags = gohcl.DecodeBody(hclFile.Body, nil, &cfg)
	if diags.HasErrors() {
		return nil, fmt.Errorf("decode config: %s", diags.Error())
	}

	applyDefaults(&cfg)
	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// EnsureAndLoad ensures a config file exists, then loads it.
func EnsureAndLoad(path string) (*Config, error) {
	if strings.TrimSpace(path) == "" {
		path = DefaultPath
	}

	if err := EnsureConfig(path); err != nil {
		return nil, err
	}
	return Load(path)
}

func validate(cfg *Config) error {
	if strings.TrimSpace(cfg.Server.Address) == "" {
		return fmt.Errorf("server address cannot be empty")
	}
	if strings.TrimSpace(cfg.UI.Address) == "" {
		return fmt.Errorf("ui address cannot be empty")
	}
	if strings.TrimSpace(cfg.NATS.Host) == "" {
		return fmt.Errorf("nats host cannot be empty")
	}
	if cfg.NATS.Port <= 0 || cfg.NATS.Port > 65535 {
		return fmt.Errorf("nats port must be between 1 and 65535")
	}
	return nil
}
