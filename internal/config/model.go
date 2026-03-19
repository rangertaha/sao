package config

// DefaultPath is the expected SAO config path on Linux systems.
const DefaultPath = "/etc/sao/config.hcl"

// Config is the root SAO configuration.
type Config struct {
	Server ServerConfig `hcl:"server,block"`
	UI     UIConfig     `hcl:"ui,block"`
	NATS   NATSConfig   `hcl:"nats,block"`
}

// ServerConfig controls TOC server runtime settings.
type ServerConfig struct {
	Address string `hcl:"address,optional"`
}

// UIConfig controls the embedded UI HTTP service.
type UIConfig struct {
	Address string `hcl:"address,optional"`
}

// NATSConfig controls embedded NATS startup.
type NATSConfig struct {
	Host string `hcl:"host,optional"`
	Port int    `hcl:"port,optional"`
}
