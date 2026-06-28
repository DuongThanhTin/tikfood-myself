package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/ingest"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/platform/textutil"
)

type IngestionRepository struct {
	db *sql.DB
}

func NewIngestionRepository(db *sql.DB) *IngestionRepository {
	return &IngestionRepository{db: db}
}

// ResolveOrCreateDistrict ensures a district location exists under the given city
// and returns its id. Idempotent on (type, country_code, slug).
func (repo *IngestionRepository) ResolveOrCreateDistrict(ctx context.Context, citySlug string, districtName string) (string, error) {
	districtSlug := textutil.Slugish(districtName)
	normalized := textutil.Normalize(districtName)

	const query = `
with city as (
  select id from locations where type = 'city' and slug = $1 limit 1
)
insert into locations (type, parent_id, official_name, normalized_name, slug, country_code)
select 'district', city.id, $2, $3, $4, 'VN' from city
on conflict (type, country_code, slug) do update set official_name = excluded.official_name, updated_at = now()
returning id::text
`
	var id string
	if err := repo.db.QueryRowContext(ctx, query, citySlug, districtName, normalized, districtSlug).Scan(&id); err != nil {
		return "", fmt.Errorf("resolve district %q: %w", districtName, err)
	}
	return id, nil
}

func (repo *IngestionRepository) RecordIngestionRun(ctx context.Context, run ingest.IngestionRun) (string, error) {
	const query = `
insert into ingestion_runs (source, district_slug, category, status)
values ($1, $2, $3, 'running')
returning id::text
`
	var id string
	if err := repo.db.QueryRowContext(ctx, query, run.Source, run.DistrictSlug, run.Category).Scan(&id); err != nil {
		return "", fmt.Errorf("record ingestion run: %w", err)
	}
	return id, nil
}

func (repo *IngestionRepository) FinishIngestionRun(ctx context.Context, runID string, counts ingest.Summary, status string) error {
	const query = `
update ingestion_runs
set status = $2, fetched_count = $3, upserted_count = $4, failed_count = $5, finished_at = now()
where id = $1::uuid
`
	if _, err := repo.db.ExecContext(ctx, query, runID, status, counts.Fetched, counts.Upserted, counts.Failed); err != nil {
		return fmt.Errorf("finish ingestion run: %w", err)
	}
	return nil
}

// UpsertVenueFromSource upserts a venue keyed on (source, source_external_id),
// then replaces its opening hours and links its tags. latitude/longitude drive
// the generated `location` column, so geo stays in sync automatically.
func (repo *IngestionRepository) UpsertVenueFromSource(ctx context.Context, input ingest.UpsertVenueInput) error {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	venueID, err := upsertVenueRow(ctx, tx, input)
	if err != nil {
		return err
	}
	if err := replaceOpeningHours(ctx, tx, venueID, input); err != nil {
		return err
	}
	if err := linkTags(ctx, tx, venueID, input.TagSlugs); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit venue upsert: %w", err)
	}
	return nil
}

func upsertVenueRow(ctx context.Context, tx *sql.Tx, input ingest.UpsertVenueInput) (string, error) {
	const query = `
insert into venues (
  name, slug, address, address_display, latitude, longitude, currency,
  phone, website_url, source, source_external_id, district_location_id, source_updated_at
) values (
  $1, $2, $3, $3, $4, $5, 'VND',
  $6, $7, $8, $9, nullif($10, '')::uuid, now()
)
on conflict (source, source_external_id) do update set
  name = excluded.name,
  slug = excluded.slug,
  address = excluded.address,
  address_display = excluded.address_display,
  latitude = excluded.latitude,
  longitude = excluded.longitude,
  phone = excluded.phone,
  website_url = excluded.website_url,
  district_location_id = excluded.district_location_id,
  source_updated_at = now(),
  updated_at = now()
returning id::text
`
	var venueID string
	err := tx.QueryRowContext(ctx, query,
		input.Venue.Name,
		input.Venue.Slug,
		input.Venue.Address,
		input.Venue.Latitude,
		input.Venue.Longitude,
		input.Venue.Phone(),
		input.Venue.Website(),
		input.Source,
		input.ExternalID,
		input.DistrictLocationID,
	).Scan(&venueID)
	if err != nil {
		return "", fmt.Errorf("upsert venue %q: %w", input.ExternalID, err)
	}
	return venueID, nil
}

func replaceOpeningHours(ctx context.Context, tx *sql.Tx, venueID string, input ingest.UpsertVenueInput) error {
	if _, err := tx.ExecContext(ctx, `delete from venue_opening_hours where venue_id = $1::uuid`, venueID); err != nil {
		return fmt.Errorf("clear opening hours: %w", err)
	}
	for _, h := range input.OpeningHours {
		const query = `
insert into venue_opening_hours (venue_id, day_of_week, open_time, close_time, is_closed)
values ($1::uuid, $2, nullif($3, '')::time, nullif($4, '')::time, $5)
`
		if _, err := tx.ExecContext(ctx, query, venueID, h.DayOfWeek, h.OpenTime, h.CloseTime, h.IsClosed); err != nil {
			return fmt.Errorf("insert opening hour: %w", err)
		}
	}
	return nil
}

func linkTags(ctx context.Context, tx *sql.Tx, venueID string, tagSlugs []string) error {
	for _, slug := range tagSlugs {
		var tagID string
		const upsertTag = `
insert into tags (slug, label, type)
values ($1, $2, 'category')
on conflict (slug) do update set label = excluded.label
returning id::text
`
		label := tagLabel(slug)
		if err := tx.QueryRowContext(ctx, upsertTag, slug, label).Scan(&tagID); err != nil {
			return fmt.Errorf("upsert tag %q: %w", slug, err)
		}
		const linkTag = `
insert into venue_tags (venue_id, tag_id, source)
values ($1::uuid, $2::uuid, 'openstreetmap')
on conflict (venue_id, tag_id) do nothing
`
		if _, err := tx.ExecContext(ctx, linkTag, venueID, tagID); err != nil {
			return fmt.Errorf("link tag %q: %w", slug, err)
		}
	}
	return nil
}

func tagLabel(slug string) string {
	words := strings.Split(slug, "-")
	for i, w := range words {
		if w == "" {
			continue
		}
		words[i] = strings.ToUpper(w[:1]) + w[1:]
	}
	return strings.Join(words, " ")
}
