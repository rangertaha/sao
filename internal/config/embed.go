package config

import (
	_ "embed"
)

//go:embed default.hcl
var defaultConfigHCL []byte

// DefaultConfigBytes returns the embedded default config bytes.
func DefaultConfigBytes() []byte {
	return append([]byte(nil), defaultConfigHCL...)
}
