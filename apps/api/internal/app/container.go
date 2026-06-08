package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/config"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
	apihttp "github.com/DuongThanhTin/tikfood-myself/apps/api/internal/http"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/storage/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Container struct {
	Logger          *slog.Logger
	VenueRepository discovery.VenueRepository
	VenueService    *discovery.VenueService
	RouteRegistrars []apihttp.RouteRegistrar
	close           func() error
}

func NewContainer(cfg config.Config) (*Container, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})).With("service", "api")

	venueRepo, closeRepo, err := buildVenueRepository(cfg, logger)
	if err != nil {
		return nil, err
	}

	venueService := discovery.NewVenueService(venueRepo)
	routeRegistrars := apihttp.DefaultRouteRegistrars(apihttp.HandlerDependencies{
		Venues: venueService,
	})

	return &Container{
		Logger:          logger,
		VenueRepository: venueRepo,
		VenueService:    venueService,
		RouteRegistrars: routeRegistrars,
		close:           closeRepo,
	}, nil
}

func (container *Container) Close() error {
	if container.close == nil {
		return nil
	}
	return container.close()
}

func buildVenueRepository(cfg config.Config, logger *slog.Logger) (discovery.VenueRepository, func() error, error) {
	if cfg.DatabaseURL == "" {
		logger.Info("using in-memory discovery storage")
		return discovery.NewFallbackVenueRepository(), nil, nil
	}

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("ping database: %w", err)
	}

	logger.Info("using postgres discovery storage")
	return postgres.NewDiscoveryRepository(db), db.Close, nil
}
