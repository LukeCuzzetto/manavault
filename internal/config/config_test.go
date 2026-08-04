package config

import "testing"

const testDatabaseURL = "postgres://deckengine:test@localhost:5432/deckengine"

func TestLoadUsesDefaultPort(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", testDatabaseURL)

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
	t.Setenv("DATABASE_URL", testDatabaseURL)

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
	t.Setenv("DATABASE_URL", testDatabaseURL)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLoadRejectsOutOfRangePort(t *testing.T) {
	t.Setenv("PORT", "70000")
	t.Setenv("DATABASE_URL", testDatabaseURL)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLoadRejectsMissingDatabaseURL(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("DATABASE_URL", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
