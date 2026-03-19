package config

import "strings"

func applyDefaults(cfg *Config) {
	if strings.TrimSpace(cfg.Server.Address) == "" {
		cfg.Server.Address = ":8080"
	}
	if strings.TrimSpace(cfg.UI.Address) == "" {
		cfg.UI.Address = ":8081"
	}
	if strings.TrimSpace(cfg.NATS.Host) == "" {
		cfg.NATS.Host = "127.0.0.1"
	}
	if cfg.NATS.Port == 0 {
		cfg.NATS.Port = 4222
	}
}
