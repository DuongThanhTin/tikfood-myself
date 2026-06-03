package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/config"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
	apihttp "github.com/DuongThanhTin/tikfood-myself/apps/api/internal/http"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/storage/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})).With("service", "api")
	venueService := discovery.NewVenueService()

	if cfg.DatabaseURL != "" {
		db, err := sql.Open("pgx", cfg.DatabaseURL)
		if err != nil {
			logger.Error("open database", "error", err)
			os.Exit(1)
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			logger.Error("ping database", "error", err)
			os.Exit(1)
		}

		venueService = discovery.NewVenueServiceWithRepository(postgres.NewDiscoveryRepository(db))
		logger.Info("using postgres discovery storage")
	} else {
		logger.Info("using in-memory discovery storage")
	}

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: apihttp.NewRouterWithLogger(venueService, logger),
	}

	logger.Info("api listening", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("api server stopped", "error", err)
		os.Exit(1)
	}
}
