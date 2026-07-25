package config

import "testing"

func TestLoadUsesDefaultPort(t *testing.T) {
	t.Setenv("PORT", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Address != ":8080" {
		t.Errorf("expected address to be :8080, got %q", cfg.Address)
	}

}
func TestLoadUsesEnvironmentPort(t *testing.T) {

	t.Setenv("PORT", "9090")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Address != ":9090" {
		t.Errorf("expected address to be :9090, got %q", cfg.Address)
	}

}

func TestLoadRejectsNonNumericalPort(t *testing.T) {
	t.Setenv("PORT", "banana")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLoadRejectsOutOfRangePort(t *testing.T) {
	t.Setenv("PORT", "70000")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
