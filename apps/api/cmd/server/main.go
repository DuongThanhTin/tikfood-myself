package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/config"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
	apihttp "github.com/DuongThanhTin/tikfood-myself/apps/api/internal/http"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/storage/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.Load()
	venueService := discovery.NewVenueService()

	if cfg.DatabaseURL != "" {
		db, err := sql.Open("pgx", cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("open database: %v", err)
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			log.Fatalf("ping database: %v", err)
		}

		venueService = discovery.NewVenueServiceWithRepository(postgres.NewDiscoveryRepository(db))
		log.Printf("tikfood api using postgres discovery storage")
	} else {
		log.Printf("tikfood api using in-memory discovery storage")
	}

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: apihttp.NewRouter(venueService),
	}

	log.Printf("tikfood api listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("api server stopped: %v", err)
	}
}
