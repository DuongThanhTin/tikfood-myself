package config

import (
	"testing"
	"time"
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

func TestConfigDefaultsAuth(t *testing.T) {
	t.Setenv("ACCESS_TOKEN_TTL", "")
	t.Setenv("REFRESH_TOKEN_TTL", "")
	t.Setenv("WEB_ORIGIN", "")
	t.Setenv("COOKIE_SECURE", "")

	cfg := Load()

	if cfg.AccessTokenTTL != 15*time.Minute {
		t.Fatalf("expected access TTL 15m, got %s", cfg.AccessTokenTTL)
	}
	if cfg.RefreshTokenTTL != 720*time.Hour {
		t.Fatalf("expected refresh TTL 720h, got %s", cfg.RefreshTokenTTL)
	}
	if cfg.WebOrigin != "http://localhost:3000" {
		t.Fatalf("expected default web origin, got %q", cfg.WebOrigin)
	}
	if len(cfg.AllowedOrigins) != 1 || cfg.AllowedOrigins[0] != "http://localhost:3000" {
		t.Fatalf("expected allowed origins to default to the web origin, got %v", cfg.AllowedOrigins)
	}
	if !cfg.CookieSecure {
		t.Fatal("expected cookie secure to default to true")
	}
}

func TestConfigAuthOverrides(t *testing.T) {
	t.Setenv("ACCESS_TOKEN_TTL", "5m")
	t.Setenv("REFRESH_TOKEN_TTL", "168h")
	t.Setenv("WEB_ORIGIN", "https://app.tikfood.local")
	t.Setenv("ALLOWED_ORIGINS", "https://app.tikfood.local, https://admin.tikfood.local")
	t.Setenv("COOKIE_SECURE", "false")

	cfg := Load()

	if cfg.AccessTokenTTL != 5*time.Minute {
		t.Fatalf("expected access TTL 5m, got %s", cfg.AccessTokenTTL)
	}
	if cfg.RefreshTokenTTL != 168*time.Hour {
		t.Fatalf("expected refresh TTL 168h, got %s", cfg.RefreshTokenTTL)
	}
	if cfg.WebOrigin != "https://app.tikfood.local" {
		t.Fatalf("unexpected web origin %q", cfg.WebOrigin)
	}
	if len(cfg.AllowedOrigins) != 2 {
		t.Fatalf("expected two allowed origins, got %v", cfg.AllowedOrigins)
	}
	if cfg.CookieSecure {
		t.Fatal("expected cookie secure to be false when overridden")
	}
}
