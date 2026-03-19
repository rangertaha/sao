package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureAndLoad(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "sao", "config.hcl")

	cfg, err := EnsureAndLoad(path)
	if err != nil {
		t.Fatalf("EnsureAndLoad() error: %v", err)
	}

	if cfg.Server.Address == "" || cfg.UI.Address == "" {
		t.Fatalf("expected default addresses to be populated")
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected config file to exist: %v", err)
	}
}
