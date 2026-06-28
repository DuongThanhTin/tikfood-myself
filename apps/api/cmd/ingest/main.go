package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/config"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/ingest"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/ingest/overpass"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/storage/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// districtArea maps a district slug to its display name and Overpass search center.
type districtArea struct {
	name    string
	lat     float64
	lng     float64
	radiusM int
}

var districtAreas = map[string]districtArea{
	"quan-1": {"Quận 1", 10.7756, 106.7019, 1300},
	"quan-3": {"Quận 3", 10.7860, 106.6840, 1200},
	"quan-4": {"Quận 4", 10.7578, 106.7050, 1200},
	"quan-5": {"Quận 5", 10.7540, 106.6630, 1300},
}

func main() {
	district := flag.String("district", "", "district slug to ingest, e.g. quan-1 (required)")
	category := flag.String("category", "restaurant", "category: restaurant, cafe, fast_food, or all")
	limit := flag.Int("limit", 50, "max venues to upsert this run")
	dryRun := flag.Bool("dry-run", false, "log what would be written without touching the database")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})).With("service", "ingest")

	if err := run(*district, *category, *limit, *dryRun, logger); err != nil {
		logger.Error("ingest failed", "error", err)
		os.Exit(1)
	}
}

func run(districtSlug string, category string, limit int, dryRun bool, logger *slog.Logger) error {
	districtSlug = strings.TrimSpace(districtSlug)
	if districtSlug == "" {
		return fmt.Errorf("--district is required (e.g. --district=quan-1)")
	}
	area, ok := districtAreas[districtSlug]
	if !ok {
		return fmt.Errorf("unknown district %q; known: %s", districtSlug, knownDistricts())
	}
	amenities, err := categoryToAmenities(category)
	if err != nil {
		return err
	}

	cfg := config.Load()
	if cfg.DatabaseURL == "" && !dryRun {
		return fmt.Errorf("DATABASE_URL is required unless --dry-run is set")
	}

	source := overpass.NewClient(cfg.OverpassEndpoint)

	var repo ingest.Repository
	var closeDB func() error
	if dryRun {
		repo = noopRepository{}
	} else {
		db, err := sql.Open("pgx", cfg.DatabaseURL)
		if err != nil {
			return fmt.Errorf("open database: %w", err)
		}
		if err := db.Ping(); err != nil {
			_ = db.Close()
			return fmt.Errorf("ping database: %w", err)
		}
		closeDB = db.Close
		repo = postgres.NewIngestionRepository(db)
	}
	if closeDB != nil {
		defer func() { _ = closeDB() }()
	}

	service := ingest.NewService(source, repo, logger)
	summary, err := service.IngestDistrict(context.Background(), ingest.IngestOptions{
		DistrictSlug: districtSlug,
		DistrictName: area.name,
		Category:     category,
		Lat:          area.lat,
		Lng:          area.lng,
		RadiusM:      area.radiusM,
		Amenities:    amenities,
		Limit:        limit,
		DryRun:       dryRun,
	})
	if err != nil {
		return err
	}

	logger.Info("ingest complete",
		"district", districtSlug, "category", category,
		"fetched", summary.Fetched, "upserted", summary.Upserted, "failed", summary.Failed,
		"dry_run", dryRun,
	)
	return nil
}

func categoryToAmenities(category string) ([]string, error) {
	switch category {
	case "restaurant", "cafe", "fast_food":
		return []string{category}, nil
	case "all":
		return []string{"restaurant", "cafe", "fast_food"}, nil
	default:
		return nil, fmt.Errorf("unknown category %q; use restaurant, cafe, fast_food, or all", category)
	}
}

func knownDistricts() string {
	keys := make([]string, 0, len(districtAreas))
	for k := range districtAreas {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

// noopRepository satisfies ingest.Repository for --dry-run (no DB connection).
type noopRepository struct{}

func (noopRepository) ResolveOrCreateDistrict(context.Context, string, string) (string, error) {
	return "", nil
}
func (noopRepository) RecordIngestionRun(context.Context, ingest.IngestionRun) (string, error) {
	return "", nil
}
func (noopRepository) UpsertVenueFromSource(context.Context, ingest.UpsertVenueInput) error {
	return nil
}
func (noopRepository) FinishIngestionRun(context.Context, string, ingest.Summary, string) error {
	return nil
}
