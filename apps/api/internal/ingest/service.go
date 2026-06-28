package ingest

import (
	"context"
	"fmt"
	"log/slog"
)

const (
	sourceOpenStreetMap = "openstreetmap"
	cityHoChiMinhSlug   = "ho-chi-minh"
)

// IngestOptions controls one district ingestion run.
type IngestOptions struct {
	DistrictSlug string
	DistrictName string
	Category     string
	Lat          float64
	Lng          float64
	RadiusM      int
	Amenities    []string
	Limit        int
	DryRun       bool
}

// Summary reports counts for a run.
type Summary struct {
	Fetched  int
	Upserted int
	Failed   int
}

// Service orchestrates source search -> map -> upsert.
type Service struct {
	source VenueSource
	repo   Repository
	logger *slog.Logger
}

func NewService(source VenueSource, repo Repository, logger *slog.Logger) *Service {
	if source == nil || repo == nil {
		panic("ingest.NewService requires a source and repository")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{source: source, repo: repo, logger: logger}
}

func (s *Service) IngestDistrict(ctx context.Context, opts IngestOptions) (Summary, error) {
	districtID, err := s.repo.ResolveOrCreateDistrict(ctx, cityHoChiMinhSlug, opts.DistrictName)
	if err != nil {
		return Summary{}, fmt.Errorf("resolve district: %w", err)
	}

	var runID string
	if !opts.DryRun {
		runID, err = s.repo.RecordIngestionRun(ctx, IngestionRun{
			Source: sourceOpenStreetMap, DistrictSlug: opts.DistrictSlug, Category: opts.Category,
		})
		if err != nil {
			return Summary{}, fmt.Errorf("record run: %w", err)
		}
	}

	places, err := s.source.Search(ctx, SearchParams{
		Lat: opts.Lat, Lng: opts.Lng, RadiusM: opts.RadiusM, Amenities: opts.Amenities, Limit: opts.Limit,
	})
	if err != nil {
		s.finish(ctx, runID, Summary{}, "failed", opts.DryRun)
		return Summary{}, fmt.Errorf("source search: %w", err)
	}

	if opts.Limit > 0 && len(places) > opts.Limit {
		places = places[:opts.Limit]
	}

	summary := Summary{Fetched: len(places)}
	for _, place := range places {
		mapped := MapPlace(place)
		if opts.DryRun {
			s.logger.Info("dry-run venue", "name", mapped.Venue.Name, "slug", mapped.Venue.Slug, "external_id", place.ExternalID)
			summary.Upserted++
			continue
		}

		input := UpsertVenueInput{
			Venue:              mapped.Venue,
			OpeningHours:       mapped.OpeningHours,
			TagSlugs:           mapped.TagSlugs,
			Source:             sourceOpenStreetMap,
			ExternalID:         place.ExternalID,
			DistrictLocationID: districtID,
		}
		if err := s.repo.UpsertVenueFromSource(ctx, input); err != nil {
			summary.Failed++
			s.logger.Warn("skip place: upsert failed", "external_id", place.ExternalID, "error", err)
			continue
		}
		summary.Upserted++
	}

	s.finish(ctx, runID, summary, "completed", opts.DryRun)
	return summary, nil
}

func (s *Service) finish(ctx context.Context, runID string, counts Summary, status string, dryRun bool) {
	if dryRun || runID == "" {
		return
	}
	if err := s.repo.FinishIngestionRun(ctx, runID, counts, status); err != nil {
		s.logger.Error("finish run failed", "run_id", runID, "error", err)
	}
}
