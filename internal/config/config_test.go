package config

import "testing"

func TestLoadUsesDefaultPort(t *testing.T) {
	t.Setenv("PORT", "")
	cfg := Load()

	if cfg.Address != ":8080" {
		t.Errorf("expected address to be :8080, got %q", cfg.Address)
	}

}
func TestLoadUsesEnvironmentPort(t *testing.T) {

	t.Setenv("PORT", "9090")

	cfg := Load()

	if cfg.Address != ":9090" {
		t.Errorf("expected address to be :9090, got %q", cfg.Address)
	}

}
