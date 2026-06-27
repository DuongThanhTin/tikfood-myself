package config

import "os"

const defaultOverpassEndpoint = "https://overpass-api.de/api/interpreter"

type Config struct {
	Port             string
	DatabaseURL      string
	OverpassEndpoint string
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

	return Config{
		Port:             port,
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		OverpassEndpoint: endpoint,
	}
}
