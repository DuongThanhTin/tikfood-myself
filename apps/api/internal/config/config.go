package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultOverpassEndpoint = "https://overpass-api.de/api/interpreter"
	defaultAccessTokenTTL   = 15 * time.Minute
	defaultRefreshTokenTTL  = 720 * time.Hour
	defaultVerificationTTL  = 24 * time.Hour
	defaultWebOrigin        = "http://localhost:3000"
)

type Config struct {
	Port             string
	DatabaseURL      string
	OverpassEndpoint string

	// Authentication (see ADR-0007). Secrets come from the environment only;
	// they are never hardcoded and never logged.
	// AuthEnabled gates the whole auth surface for staged rollout: when false, no auth
	// routes are registered and discovery stays public (default true).
	AuthEnabled          bool
	JWTSecret            string
	AccessTokenTTL       time.Duration
	RefreshTokenTTL      time.Duration
	VerificationTokenTTL time.Duration
	CookieSecure         bool

	// Google OAuth SSO.
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	// Browser origins allowed to make credentialed requests.
	WebOrigin      string
	AllowedOrigins []string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "18081"
	}

	endpoint := os.Getenv("OVERPASS_ENDPOINT")
	if endpoint == "" {
		endpoint = defaultOverpassEndpoint
	}

	webOrigin := strings.TrimSpace(os.Getenv("WEB_ORIGIN"))
	if webOrigin == "" {
		webOrigin = defaultWebOrigin
	}

	allowedOrigins := splitAndTrim(os.Getenv("ALLOWED_ORIGINS"))
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{webOrigin}
	}

	return Config{
		Port:             port,
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		OverpassEndpoint: endpoint,

		AuthEnabled:          boolOrDefault("AUTH_ENABLED", true),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		AccessTokenTTL:       durationOrDefault("ACCESS_TOKEN_TTL", defaultAccessTokenTTL),
		RefreshTokenTTL:      durationOrDefault("REFRESH_TOKEN_TTL", defaultRefreshTokenTTL),
		VerificationTokenTTL: durationOrDefault("EMAIL_VERIFICATION_TTL", defaultVerificationTTL),
		CookieSecure:         boolOrDefault("COOKIE_SECURE", true),

		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),

		WebOrigin:      webOrigin,
		AllowedOrigins: allowedOrigins,
	}
}

func durationOrDefault(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func boolOrDefault(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func splitAndTrim(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	results := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			results = append(results, trimmed)
		}
	}
	return results
}
