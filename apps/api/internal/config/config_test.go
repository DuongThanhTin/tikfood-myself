package config

import (
	"testing"
)

func TestLoadOverpassEndpointDefault(t *testing.T) {
	t.Setenv("OVERPASS_ENDPOINT", "")

	cfg := Load()

	if cfg.OverpassEndpoint != "https://overpass-api.de/api/interpreter" {
		t.Fatalf("expected default endpoint, got %q", cfg.OverpassEndpoint)
	}
}

func TestLoadOverpassEndpointOverride(t *testing.T) {
	t.Setenv("OVERPASS_ENDPOINT", "https://overpass.kumi.systems/api/interpreter")

	cfg := Load()

	if cfg.OverpassEndpoint != "https://overpass.kumi.systems/api/interpreter" {
		t.Fatalf("expected override endpoint, got %q", cfg.OverpassEndpoint)
	}
}
