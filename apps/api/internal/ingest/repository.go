package ingest

import (
	"context"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
)

// UpsertVenueInput is everything needed to persist one ingested venue.
type UpsertVenueInput struct {
	Venue              discovery.Venue
	OpeningHours       []discovery.OpeningHour
	TagSlugs           []string
	Source             string
	ExternalID         string
	DistrictLocationID string
}

// IngestionRun describes a district run for the audit log.
type IngestionRun struct {
	Source       string
	DistrictSlug string
	Category     string
}

// Repository persists ingested venues and run metadata.
type Repository interface {
	ResolveOrCreateDistrict(ctx context.Context, citySlug string, districtName string) (string, error)
	RecordIngestionRun(ctx context.Context, run IngestionRun) (string, error)
	UpsertVenueFromSource(ctx context.Context, input UpsertVenueInput) error
	FinishIngestionRun(ctx context.Context, runID string, counts Summary, status string) error
}
