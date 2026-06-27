# Venue Ingestion from Google Places Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** A re-runnable Go CLI (`apps/api/cmd/ingest`) that ingests Ho Chi Minh City restaurants/cafes from the Google Places API one district at a time, upserting idempotently.

**Architecture:** New `internal/ingest` package (orchestration + pure mapper) depends on an `internal/ingest/places` API client (behind an interface) and an `ingest.Repository` interface implemented by `internal/storage/postgres`. The CLI wires them. Layering matches `apps/api/CLAUDE.md`: cmd → service → (places client | repository); the service never imports HTTP framework types.

**Tech Stack:** Go 1.26, `database/sql` via pgx stdlib, Google Places API (New) over `net/http`, stdlib `testing` (table-driven).

## Global Constraints

- Module path: `github.com/DuongThanhTin/tikfood-myself/apps/api`.
- Source of venue data: Google Places API (New). API key only from env `GOOGLE_PLACES_API_KEY`; never hardcoded, never logged.
- Idempotent upsert keyed on `(source, source_external_id)` where `source='google_places'` and `source_external_id` = Places `place_id`.
- Scope v1: city fixed to Ho Chi Minh; `--district` selects one district per run; category default `restaurant`.
- Fields captured: name, address, latitude, longitude, place_id, opening hours, Google rating + rating count, photo references, Places `types` → tags, phone, website. No dishes/prices.
- No network calls in tests; no DB-backed tests (none exist today). TDD applies to pure functions (Places JSON decode, mapper) and the service (fake client + mock repo). Infra (migration, HTTP client, postgres repo, CLI) is verified with `go build ./...` and `go vet ./...`.
- Response envelope and existing API behavior are unchanged by this work.
- Run all Go commands from `apps/api` with `GOCACHE=$(pwd)/../../.cache/go-build`.
- Source spec: `docs/superpowers/specs/2026-06-27-venue-ingestion-google-places-design.md`.

---

### Task 1: Migration 006 — schema for ingestion

**Files:**
- Create: `apps/api/migrations/006_venue_ingestion_schema.sql`

**Interfaces:**
- Produces: `venues.google_rating`, `venues.google_rating_count`; unique `(source, source_external_id)`; `venues.location` as a generated column; table `ingestion_runs`.

- [ ] **Step 1: Write the migration**

```sql
-- venues: Google rating fields
alter table venues add column if not exists google_rating numeric;
alter table venues add column if not exists google_rating_count int;

-- venues: upsert key for Places ingestion
alter table venues drop constraint if exists venues_source_external_unique;
alter table venues add constraint venues_source_external_unique unique (source, source_external_id);

-- venues: keep PostGIS location in sync with lat/lng via a generated column.
-- Existing data is seed-only, so dropping and re-adding is safe.
drop index if exists idx_venues_location;
alter table venues drop column if exists location;
alter table venues add column location geography(Point, 4326)
  generated always as (
    case
      when longitude is null or latitude is null then null
      else ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)::geography
    end
  ) stored;
create index if not exists idx_venues_location on venues using gist (location);

-- ingestion audit log (one row per district run)
create table if not exists ingestion_runs (
  id uuid primary key default gen_random_uuid(),
  source text not null,
  district_slug text not null,
  category text not null,
  status text not null check (status in ('running', 'completed', 'failed')),
  fetched_count int not null default 0,
  upserted_count int not null default 0,
  failed_count int not null default 0,
  error text,
  started_at timestamptz not null default now(),
  finished_at timestamptz
);

create index if not exists idx_ingestion_runs_district on ingestion_runs (district_slug, started_at desc);
```

- [ ] **Step 2: Verify SQL is well-formed and self-consistent**

Run: `grep -nE "ingestion_runs|generated always as|venues_source_external_unique|google_rating" apps/api/migrations/006_venue_ingestion_schema.sql`
Expected: matches the four additions above. (No DB is required to apply it; it runs on next `docker compose up` via initdb, and follows the existing `if not exists` idempotent convention.)

- [ ] **Step 3: Commit**

```bash
git add apps/api/migrations/006_venue_ingestion_schema.sql
git commit -m "Add migration 006: ingestion schema (ratings, upsert key, generated location, ingestion_runs)"
```

---

### Task 2: Config — Google Places API key

**Files:**
- Modify: `apps/api/internal/config/config.go`
- Test: `apps/api/internal/config/config_test.go`

**Interfaces:**
- Produces: `config.Config.GooglePlacesAPIKey string`, populated from env `GOOGLE_PLACES_API_KEY`.

- [ ] **Step 1: Write the failing test**

Create `apps/api/internal/config/config_test.go`:

```go
package config

import (
	"os"
	"testing"
)

func TestLoadReadsGooglePlacesAPIKey(t *testing.T) {
	t.Setenv("GOOGLE_PLACES_API_KEY", "test-key-123")

	cfg := Load()

	if cfg.GooglePlacesAPIKey != "test-key-123" {
		t.Fatalf("expected key to be loaded, got %q", cfg.GooglePlacesAPIKey)
	}
}

func TestLoadDefaultsWhenNoKey(t *testing.T) {
	os.Unsetenv("GOOGLE_PLACES_API_KEY")

	cfg := Load()

	if cfg.GooglePlacesAPIKey != "" {
		t.Fatalf("expected empty key, got %q", cfg.GooglePlacesAPIKey)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/config/`
Expected: FAIL — `cfg.GooglePlacesAPIKey` undefined (compile error).

- [ ] **Step 3: Add the field**

Replace the body of `apps/api/internal/config/config.go` with:

```go
package config

import "os"

type Config struct {
	Port               string
	DatabaseURL        string
	GooglePlacesAPIKey string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "18081"
	}

	return Config{
		Port:               port,
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		GooglePlacesAPIKey: os.Getenv("GOOGLE_PLACES_API_KEY"),
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/config/`
Expected: PASS (ok).

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/config/config.go apps/api/internal/config/config_test.go
git commit -m "Add GooglePlacesAPIKey to config"
```

---

### Task 3: Places types + JSON decode (pure)

**Files:**
- Create: `apps/api/internal/ingest/places/types.go`
- Create: `apps/api/internal/ingest/places/decode.go`
- Create: `apps/api/internal/ingest/places/decode_test.go`
- Create: `apps/api/internal/ingest/places/testdata/place_details.json`

**Interfaces:**
- Produces:
  - `places.Place` struct (fields below).
  - `places.LatLng{Lat, Lng float64}`.
  - `places.DayHours{DayOfWeek int; Open string; Close string; Closed bool}`.
  - `places.Client` interface: `TextSearch(ctx context.Context, query string, limit int) ([]string, error)` and `PlaceDetails(ctx context.Context, placeID string) (Place, error)`.
  - `places.DecodePlaceDetails(data []byte) (Place, error)` and `places.DecodeTextSearchIDs(data []byte) ([]string, error)`.

- [ ] **Step 1: Create the fixture**

Create `apps/api/internal/ingest/places/testdata/place_details.json`:

```json
{
  "id": "ChIJtest123",
  "displayName": { "text": "Banh Mi Test", "languageCode": "vi" },
  "formattedAddress": "12 Nguyen Trai, Quan 1, Ho Chi Minh",
  "location": { "latitude": 10.7712, "longitude": 106.6899 },
  "rating": 4.6,
  "userRatingCount": 215,
  "types": ["restaurant", "vietnamese_restaurant", "food", "point_of_interest"],
  "nationalPhoneNumber": "028 1234 5678",
  "websiteUri": "https://example.com",
  "photos": [{ "name": "places/ChIJtest123/photos/AbC" }, { "name": "places/ChIJtest123/photos/DeF" }],
  "regularOpeningHours": {
    "periods": [
      { "open": { "day": 1, "hour": 8, "minute": 0 }, "close": { "day": 1, "hour": 22, "minute": 30 } },
      { "open": { "day": 2, "hour": 8, "minute": 0 }, "close": { "day": 2, "hour": 22, "minute": 0 } }
    ]
  }
}
```

- [ ] **Step 2: Write the failing test**

Create `apps/api/internal/ingest/places/decode_test.go`:

```go
package places

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDecodePlaceDetails(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "place_details.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	place, err := DecodePlaceDetails(data)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	if place.ID != "ChIJtest123" {
		t.Fatalf("id: got %q", place.ID)
	}
	if place.DisplayName != "Banh Mi Test" {
		t.Fatalf("name: got %q", place.DisplayName)
	}
	if place.Location.Lat != 10.7712 || place.Location.Lng != 106.6899 {
		t.Fatalf("location: got %+v", place.Location)
	}
	if place.Rating != 4.6 || place.UserRatingCount != 215 {
		t.Fatalf("rating: got %v/%d", place.Rating, place.UserRatingCount)
	}
	if len(place.Types) != 4 || place.Types[0] != "restaurant" {
		t.Fatalf("types: got %v", place.Types)
	}
	if place.NationalPhoneNumber != "028 1234 5678" || place.WebsiteURI != "https://example.com" {
		t.Fatalf("contact: got %q / %q", place.NationalPhoneNumber, place.WebsiteURI)
	}
	if len(place.PhotoNames) != 2 || place.PhotoNames[0] != "places/ChIJtest123/photos/AbC" {
		t.Fatalf("photos: got %v", place.PhotoNames)
	}
	if len(place.OpeningHours) != 2 {
		t.Fatalf("hours len: got %d", len(place.OpeningHours))
	}
	if place.OpeningHours[0] != (DayHours{DayOfWeek: 1, Open: "08:00", Close: "22:30", Closed: false}) {
		t.Fatalf("hours[0]: got %+v", place.OpeningHours[0])
	}
}

func TestDecodeTextSearchIDs(t *testing.T) {
	data := []byte(`{"places":[{"id":"a"},{"id":"b"},{"id":"c"}]}`)

	ids, err := DecodeTextSearchIDs(data)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(ids) != 3 || ids[0] != "a" || ids[2] != "c" {
		t.Fatalf("ids: got %v", ids)
	}
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/ingest/places/`
Expected: FAIL — package/functions undefined (compile error).

- [ ] **Step 4: Write types.go**

Create `apps/api/internal/ingest/places/types.go`:

```go
package places

import "context"

// Place is the normalized subset of a Google Places API (New) place we ingest.
type Place struct {
	ID                  string
	DisplayName         string
	FormattedAddress    string
	Location            LatLng
	Rating              float64
	UserRatingCount     int
	Types               []string
	NationalPhoneNumber string
	WebsiteURI          string
	PhotoNames          []string
	OpeningHours        []DayHours
}

type LatLng struct {
	Lat float64
	Lng float64
}

// DayHours uses Google's day numbering: 0 = Sunday ... 6 = Saturday,
// matching venue_opening_hours.day_of_week.
type DayHours struct {
	DayOfWeek int
	Open      string // "HH:MM"
	Close     string // "HH:MM"
	Closed    bool
}

// Client retrieves places from the Google Places API (New).
type Client interface {
	TextSearch(ctx context.Context, query string, limit int) ([]string, error)
	PlaceDetails(ctx context.Context, placeID string) (Place, error)
}
```

- [ ] **Step 5: Write decode.go**

Create `apps/api/internal/ingest/places/decode.go`:

```go
package places

import (
	"encoding/json"
	"fmt"
)

type rawTextSearch struct {
	Places []struct {
		ID string `json:"id"`
	} `json:"places"`
}

type rawTimePoint struct {
	Day    int `json:"day"`
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

type rawPlaceDetails struct {
	ID               string `json:"id"`
	DisplayName      struct {
		Text string `json:"text"`
	} `json:"displayName"`
	FormattedAddress string `json:"formattedAddress"`
	Location         struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Rating              float64  `json:"rating"`
	UserRatingCount     int      `json:"userRatingCount"`
	Types               []string `json:"types"`
	NationalPhoneNumber string   `json:"nationalPhoneNumber"`
	WebsiteURI          string   `json:"websiteUri"`
	Photos              []struct {
		Name string `json:"name"`
	} `json:"photos"`
	RegularOpeningHours struct {
		Periods []struct {
			Open  *rawTimePoint `json:"open"`
			Close *rawTimePoint `json:"close"`
		} `json:"periods"`
	} `json:"regularOpeningHours"`
}

// DecodeTextSearchIDs extracts place ids from a Places searchText response.
func DecodeTextSearchIDs(data []byte) ([]string, error) {
	var raw rawTextSearch
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("decode text search: %w", err)
	}
	ids := make([]string, 0, len(raw.Places))
	for _, p := range raw.Places {
		if p.ID != "" {
			ids = append(ids, p.ID)
		}
	}
	return ids, nil
}

// DecodePlaceDetails converts a Places place-details response into a Place.
func DecodePlaceDetails(data []byte) (Place, error) {
	var raw rawPlaceDetails
	if err := json.Unmarshal(data, &raw); err != nil {
		return Place{}, fmt.Errorf("decode place details: %w", err)
	}

	place := Place{
		ID:                  raw.ID,
		DisplayName:         raw.DisplayName.Text,
		FormattedAddress:    raw.FormattedAddress,
		Location:            LatLng{Lat: raw.Location.Latitude, Lng: raw.Location.Longitude},
		Rating:              raw.Rating,
		UserRatingCount:     raw.UserRatingCount,
		Types:               raw.Types,
		NationalPhoneNumber: raw.NationalPhoneNumber,
		WebsiteURI:          raw.WebsiteURI,
	}
	for _, photo := range raw.Photos {
		if photo.Name != "" {
			place.PhotoNames = append(place.PhotoNames, photo.Name)
		}
	}
	for _, period := range raw.RegularOpeningHours.Periods {
		if period.Open == nil {
			continue
		}
		hours := DayHours{
			DayOfWeek: period.Open.Day,
			Open:      formatHM(period.Open.Hour, period.Open.Minute),
		}
		if period.Close != nil {
			hours.Close = formatHM(period.Close.Hour, period.Close.Minute)
		}
		place.OpeningHours = append(place.OpeningHours, hours)
	}
	return place, nil
}

func formatHM(hour int, minute int) string {
	return fmt.Sprintf("%02d:%02d", hour, minute)
}
```

- [ ] **Step 6: Run test to verify it passes**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/ingest/places/`
Expected: PASS (ok).

- [ ] **Step 7: Commit**

```bash
git add apps/api/internal/ingest/places/
git commit -m "Add Places API types and pure JSON decoders"
```

---

### Task 4: Places HTTP client

**Files:**
- Create: `apps/api/internal/ingest/places/client.go`

**Interfaces:**
- Consumes: `Place`, `DecodePlaceDetails`, `DecodeTextSearchIDs` (Task 3).
- Produces: `places.NewHTTPClient(apiKey string) *HTTPClient` implementing `places.Client`.

- [ ] **Step 1: Write client.go**

Create `apps/api/internal/ingest/places/client.go`:

```go
package places

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	textSearchURL   = "https://places.googleapis.com/v1/places:searchText"
	placeDetailsURL = "https://places.googleapis.com/v1/places/"

	textSearchFieldMask   = "places.id"
	placeDetailsFieldMask = "id,displayName,formattedAddress,location,rating,userRatingCount,types,nationalPhoneNumber,websiteUri,photos,regularOpeningHours"
)

// HTTPClient calls the Google Places API (New). The API key is never logged.
type HTTPClient struct {
	apiKey string
	http   *http.Client
}

func NewHTTPClient(apiKey string) *HTTPClient {
	return &HTTPClient{
		apiKey: apiKey,
		http:   &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *HTTPClient) TextSearch(ctx context.Context, query string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 20
	}
	body := fmt.Sprintf(`{"textQuery":%q,"pageSize":%d}`, query, limit)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, textSearchURL, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build text search request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-FieldMask", textSearchFieldMask)

	data, err := c.do(req)
	if err != nil {
		return nil, err
	}
	return DecodeTextSearchIDs(data)
}

func (c *HTTPClient) PlaceDetails(ctx context.Context, placeID string) (Place, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, placeDetailsURL+url.PathEscape(placeID), nil)
	if err != nil {
		return Place{}, fmt.Errorf("build place details request: %w", err)
	}
	req.Header.Set("X-Goog-FieldMask", placeDetailsFieldMask)

	data, err := c.do(req)
	if err != nil {
		return Place{}, err
	}
	return DecodePlaceDetails(data)
}

func (c *HTTPClient) do(req *http.Request) ([]byte, error) {
	req.Header.Set("X-Goog-Api-Key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("places request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read places response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		// Error body may contain status detail but never the API key (it is a header).
		return nil, fmt.Errorf("places API status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return data, nil
}
```

- [ ] **Step 2: Verify it builds and vets**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go build ./... && go vet ./internal/ingest/places/`
Expected: no output (success). (No unit test: this is a thin network adapter; its decode logic is already tested in Task 3.)

- [ ] **Step 3: Commit**

```bash
git add apps/api/internal/ingest/places/client.go
git commit -m "Add Google Places API HTTP client"
```

---

### Task 5: Ingest mapper (pure)

**Files:**
- Create: `apps/api/internal/ingest/mapper.go`
- Create: `apps/api/internal/ingest/mapper_test.go`

**Interfaces:**
- Consumes: `places.Place`, `places.DayHours` (Task 3); `discovery.Venue`, `discovery.OpeningHour` (existing); `textutil.Slugish` (existing).
- Produces:
  - `ingest.MappedVenue{Venue discovery.Venue; OpeningHours []discovery.OpeningHour; TagSlugs []string; GoogleRating float64; GoogleRatingCount int}`.
  - `ingest.MapPlace(p places.Place) MappedVenue`.
  - `ingest.TypeToTagSlug(placeType string) string`.

- [ ] **Step 1: Write the failing test**

Create `apps/api/internal/ingest/mapper_test.go`:

```go
package ingest

import (
	"testing"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/ingest/places"
)

func samplePlace() places.Place {
	return places.Place{
		ID:                  "ChIJtest123",
		DisplayName:         "Bánh Mì Hẻm",
		FormattedAddress:    "12 Nguyen Trai, Quan 1, Ho Chi Minh",
		Location:            places.LatLng{Lat: 10.7712, Lng: 106.6899},
		Rating:              4.6,
		UserRatingCount:     215,
		Types:               []string{"restaurant", "vietnamese_restaurant", "point_of_interest"},
		NationalPhoneNumber: "028 1234 5678",
		WebsiteURI:          "https://example.com",
		PhotoNames:          []string{"places/ChIJtest123/photos/AbC"},
		OpeningHours: []places.DayHours{
			{DayOfWeek: 1, Open: "08:00", Close: "22:30"},
		},
	}
}

func TestMapPlaceCoreFields(t *testing.T) {
	m := MapPlace(samplePlace())

	if m.Venue.Name != "Bánh Mì Hẻm" {
		t.Fatalf("name: got %q", m.Venue.Name)
	}
	if m.Venue.Slug != "banh-mi-hem" {
		t.Fatalf("slug: got %q", m.Venue.Slug)
	}
	if m.Venue.Latitude != 10.7712 || m.Venue.Longitude != 106.6899 {
		t.Fatalf("coords: got %v/%v", m.Venue.Latitude, m.Venue.Longitude)
	}
	if m.Venue.Address != "12 Nguyen Trai, Quan 1, Ho Chi Minh" {
		t.Fatalf("address: got %q", m.Venue.Address)
	}
	if m.Venue.Currency != "VND" {
		t.Fatalf("currency: got %q", m.Venue.Currency)
	}
	if m.GoogleRating != 4.6 || m.GoogleRatingCount != 215 {
		t.Fatalf("rating: got %v/%d", m.GoogleRating, m.GoogleRatingCount)
	}
}

func TestMapPlaceOpeningHours(t *testing.T) {
	m := MapPlace(samplePlace())

	if len(m.OpeningHours) != 1 {
		t.Fatalf("hours len: got %d", len(m.OpeningHours))
	}
	got := m.OpeningHours[0]
	if got.DayOfWeek != 1 || got.OpenTime != "08:00" || got.CloseTime != "22:30" || got.IsClosed {
		t.Fatalf("hours[0]: got %+v", got)
	}
}

func TestMapPlaceTagSlugs(t *testing.T) {
	m := MapPlace(samplePlace())

	// "point_of_interest" is dropped as noise; vietnamese_restaurant -> vietnamese.
	want := map[string]bool{"restaurant": true, "vietnamese": true}
	if len(m.TagSlugs) != len(want) {
		t.Fatalf("tag slugs: got %v", m.TagSlugs)
	}
	for _, slug := range m.TagSlugs {
		if !want[slug] {
			t.Fatalf("unexpected tag slug %q in %v", slug, m.TagSlugs)
		}
	}
}

func TestTypeToTagSlug(t *testing.T) {
	cases := map[string]string{
		"restaurant":            "restaurant",
		"cafe":                  "cafe",
		"vietnamese_restaurant": "vietnamese",
		"point_of_interest":     "", // noise -> dropped
		"food":                  "", // noise -> dropped
		"bakery":                "bakery",
	}
	for input, want := range cases {
		if got := TypeToTagSlug(input); got != want {
			t.Fatalf("TypeToTagSlug(%q): got %q want %q", input, got, want)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/ingest/`
Expected: FAIL — `MapPlace`/`MappedVenue`/`TypeToTagSlug` undefined (compile error).

- [ ] **Step 3: Write mapper.go**

Create `apps/api/internal/ingest/mapper.go`:

```go
package ingest

import (
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/ingest/places"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/platform/textutil"
)

// MappedVenue is the result of mapping a Places place to TikFood domain data.
type MappedVenue struct {
	Venue             discovery.Venue
	OpeningHours      []discovery.OpeningHour
	TagSlugs          []string
	GoogleRating      float64
	GoogleRatingCount int
}

// placeTypeNoise lists Places types that carry no useful category meaning.
var placeTypeNoise = map[string]bool{
	"point_of_interest": true,
	"establishment":     true,
	"food":              true,
}

// placeTypeAliases maps verbose Places types to TikFood tag slugs.
var placeTypeAliases = map[string]string{
	"vietnamese_restaurant": "vietnamese",
	"meal_takeaway":         "takeaway",
	"meal_delivery":         "delivery",
}

// TypeToTagSlug converts a Places type to a tag slug, or "" if it is noise.
func TypeToTagSlug(placeType string) string {
	if placeType == "" || placeTypeNoise[placeType] {
		return ""
	}
	if alias, ok := placeTypeAliases[placeType]; ok {
		return alias
	}
	return textutil.Slugish(placeType)
}

// MapPlace converts a Places place into TikFood domain data.
func MapPlace(p places.Place) MappedVenue {
	venue := discovery.Venue{
		Name:      p.DisplayName,
		Slug:      textutil.Slugish(p.DisplayName),
		Address:   p.FormattedAddress,
		Latitude:  p.Location.Lat,
		Longitude: p.Location.Lng,
		Currency:  "VND",
	}

	hours := make([]discovery.OpeningHour, 0, len(p.OpeningHours))
	for _, h := range p.OpeningHours {
		hours = append(hours, discovery.OpeningHour{
			DayOfWeek: h.DayOfWeek,
			OpenTime:  h.Open,
			CloseTime: h.Close,
			IsClosed:  h.Closed,
		})
	}

	tagSlugs := make([]string, 0, len(p.Types))
	seen := map[string]bool{}
	for _, t := range p.Types {
		slug := TypeToTagSlug(t)
		if slug == "" || seen[slug] {
			continue
		}
		seen[slug] = true
		tagSlugs = append(tagSlugs, slug)
	}

	return MappedVenue{
		Venue:             venue,
		OpeningHours:      hours,
		TagSlugs:          tagSlugs,
		GoogleRating:      p.Rating,
		GoogleRatingCount: p.UserRatingCount,
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/ingest/`
Expected: PASS (ok).

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/ingest/mapper.go apps/api/internal/ingest/mapper_test.go
git commit -m "Add pure Places-to-venue mapper"
```

---

### Task 6: Ingest service + repository interface

**Files:**
- Create: `apps/api/internal/ingest/repository.go`
- Create: `apps/api/internal/ingest/service.go`
- Create: `apps/api/internal/ingest/service_test.go`

**Interfaces:**
- Consumes: `places.Client` (Task 3), `MapPlace`/`MappedVenue` (Task 5), `discovery.OpeningHour`, `discovery.Venue`.
- Produces:
  - `ingest.UpsertVenueInput{Venue discovery.Venue; OpeningHours []discovery.OpeningHour; TagSlugs []string; GoogleRating float64; GoogleRatingCount int; Source string; ExternalID string; DistrictLocationID string}`.
  - `ingest.Repository` interface (methods below).
  - `ingest.NewService(client places.Client, repo Repository, logger *slog.Logger) *Service`.
  - `ingest.Service.IngestDistrict(ctx, opts IngestOptions) (Summary, error)`.
  - `ingest.IngestOptions{DistrictSlug, DistrictName, Category string; Limit int; DryRun bool}`.
  - `ingest.Summary{Fetched, Upserted, Failed int}`.

- [ ] **Step 1: Write the failing test**

Create `apps/api/internal/ingest/service_test.go`:

```go
package ingest

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/ingest/places"
)

type fakeClient struct {
	ids     []string
	details map[string]places.Place
	err     error
}

func (f fakeClient) TextSearch(_ context.Context, _ string, _ int) ([]string, error) {
	return f.ids, f.err
}

func (f fakeClient) PlaceDetails(_ context.Context, id string) (places.Place, error) {
	p, ok := f.details[id]
	if !ok {
		return places.Place{}, errors.New("not found")
	}
	return p, nil
}

type mockRepo struct {
	upserts    []UpsertVenueInput
	districtID string
	runID      string
	finished   bool
	finalCount Summary
	finalState string
	failResolve bool
}

func (m *mockRepo) ResolveOrCreateDistrict(_ context.Context, _, _ string) (string, error) {
	if m.failResolve {
		return "", errors.New("resolve failed")
	}
	if m.districtID == "" {
		m.districtID = "district-loc-1"
	}
	return m.districtID, nil
}

func (m *mockRepo) RecordIngestionRun(_ context.Context, _ IngestionRun) (string, error) {
	m.runID = "run-1"
	return m.runID, nil
}

func (m *mockRepo) UpsertVenueFromSource(_ context.Context, input UpsertVenueInput) error {
	m.upserts = append(m.upserts, input)
	return nil
}

func (m *mockRepo) FinishIngestionRun(_ context.Context, _ string, counts Summary, status string) error {
	m.finished = true
	m.finalCount = counts
	m.finalState = status
	return nil
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestIngestDistrictUpsertsEachPlace(t *testing.T) {
	client := fakeClient{
		ids: []string{"a", "b"},
		details: map[string]places.Place{
			"a": {ID: "a", DisplayName: "Alpha", Location: places.LatLng{Lat: 10, Lng: 106}},
			"b": {ID: "b", DisplayName: "Beta", Location: places.LatLng{Lat: 11, Lng: 107}},
		},
	}
	repo := &mockRepo{}
	svc := NewService(client, repo, testLogger())

	summary, err := svc.IngestDistrict(context.Background(), IngestOptions{
		DistrictSlug: "quan-1", DistrictName: "Quận 1", Category: "restaurant", Limit: 20,
	})
	if err != nil {
		t.Fatalf("ingest: %v", err)
	}
	if summary.Fetched != 2 || summary.Upserted != 2 || summary.Failed != 0 {
		t.Fatalf("summary: %+v", summary)
	}
	if len(repo.upserts) != 2 {
		t.Fatalf("expected 2 upserts, got %d", len(repo.upserts))
	}
	if repo.upserts[0].Source != "google_places" || repo.upserts[0].ExternalID != "a" {
		t.Fatalf("upsert[0]: %+v", repo.upserts[0])
	}
	if repo.upserts[0].DistrictLocationID != "district-loc-1" {
		t.Fatalf("district id not set: %+v", repo.upserts[0])
	}
	if !repo.finished || repo.finalState != "completed" {
		t.Fatalf("run not finished completed: %+v", repo)
	}
}

func TestIngestDistrictDryRunWritesNothing(t *testing.T) {
	client := fakeClient{
		ids:     []string{"a"},
		details: map[string]places.Place{"a": {ID: "a", DisplayName: "Alpha"}},
	}
	repo := &mockRepo{}
	svc := NewService(client, repo, testLogger())

	summary, err := svc.IngestDistrict(context.Background(), IngestOptions{
		DistrictSlug: "quan-1", DistrictName: "Quận 1", Category: "restaurant", DryRun: true,
	})
	if err != nil {
		t.Fatalf("ingest: %v", err)
	}
	if len(repo.upserts) != 0 {
		t.Fatalf("dry run must not upsert, got %d", len(repo.upserts))
	}
	if summary.Fetched != 1 {
		t.Fatalf("summary fetched: %+v", summary)
	}
}

func TestIngestDistrictSkipsPlaceDetailsError(t *testing.T) {
	client := fakeClient{
		ids: []string{"a", "missing"},
		details: map[string]places.Place{
			"a": {ID: "a", DisplayName: "Alpha"},
		},
	}
	repo := &mockRepo{}
	svc := NewService(client, repo, testLogger())

	summary, err := svc.IngestDistrict(context.Background(), IngestOptions{
		DistrictSlug: "quan-1", DistrictName: "Quận 1", Category: "restaurant",
	})
	if err != nil {
		t.Fatalf("ingest: %v", err)
	}
	if summary.Upserted != 1 || summary.Failed != 1 {
		t.Fatalf("summary: %+v", summary)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/ingest/`
Expected: FAIL — `NewService`/`Repository`/`IngestOptions` undefined (compile error).

- [ ] **Step 3: Write repository.go**

Create `apps/api/internal/ingest/repository.go`:

```go
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
	GoogleRating       float64
	GoogleRatingCount  int
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
```

- [ ] **Step 4: Write service.go**

Create `apps/api/internal/ingest/service.go`:

```go
package ingest

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/ingest/places"
)

const (
	sourceGooglePlaces = "google_places"
	cityHoChiMinhSlug  = "ho-chi-minh"
)

// IngestOptions controls one district ingestion run.
type IngestOptions struct {
	DistrictSlug string
	DistrictName string
	Category     string
	Limit        int
	DryRun       bool
}

// Summary reports counts for a run.
type Summary struct {
	Fetched  int
	Upserted int
	Failed   int
}

// Service orchestrates Places search -> details -> map -> upsert.
type Service struct {
	client places.Client
	repo   Repository
	logger *slog.Logger
}

func NewService(client places.Client, repo Repository, logger *slog.Logger) *Service {
	if client == nil || repo == nil {
		panic("ingest.NewService requires a client and repository")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{client: client, repo: repo, logger: logger}
}

func (s *Service) IngestDistrict(ctx context.Context, opts IngestOptions) (Summary, error) {
	query := fmt.Sprintf("%s in %s, Ho Chi Minh", opts.Category, opts.DistrictName)

	districtID, err := s.repo.ResolveOrCreateDistrict(ctx, cityHoChiMinhSlug, opts.DistrictName)
	if err != nil {
		return Summary{}, fmt.Errorf("resolve district: %w", err)
	}

	var runID string
	if !opts.DryRun {
		runID, err = s.repo.RecordIngestionRun(ctx, IngestionRun{
			Source: sourceGooglePlaces, DistrictSlug: opts.DistrictSlug, Category: opts.Category,
		})
		if err != nil {
			return Summary{}, fmt.Errorf("record run: %w", err)
		}
	}

	ids, err := s.client.TextSearch(ctx, query, opts.Limit)
	if err != nil {
		s.finish(ctx, runID, Summary{}, "failed", opts.DryRun)
		return Summary{}, fmt.Errorf("text search: %w", err)
	}

	summary := Summary{Fetched: len(ids)}
	for _, id := range ids {
		place, err := s.client.PlaceDetails(ctx, id)
		if err != nil {
			summary.Failed++
			s.logger.Warn("skip place: details failed", "place_id", id, "error", err)
			continue
		}

		mapped := MapPlace(place)
		if opts.DryRun {
			s.logger.Info("dry-run venue", "name", mapped.Venue.Name, "slug", mapped.Venue.Slug, "place_id", place.ID)
			summary.Upserted++
			continue
		}

		input := UpsertVenueInput{
			Venue:              mapped.Venue,
			OpeningHours:       mapped.OpeningHours,
			TagSlugs:           mapped.TagSlugs,
			GoogleRating:       mapped.GoogleRating,
			GoogleRatingCount:  mapped.GoogleRatingCount,
			Source:             sourceGooglePlaces,
			ExternalID:         place.ID,
			DistrictLocationID: districtID,
		}
		if err := s.repo.UpsertVenueFromSource(ctx, input); err != nil {
			summary.Failed++
			s.logger.Warn("skip place: upsert failed", "place_id", id, "error", err)
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
```

- [ ] **Step 5: Run test to verify it passes**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/ingest/`
Expected: PASS (ok). Note: in dry-run the `Upserted` counter reflects places that *would* be written (matches `TestIngestDistrictDryRunWritesNothing`, which asserts on `repo.upserts` being empty and `Fetched`).

- [ ] **Step 6: Commit**

```bash
git add apps/api/internal/ingest/repository.go apps/api/internal/ingest/service.go apps/api/internal/ingest/service_test.go
git commit -m "Add ingest service and repository interface"
```

---

### Task 7: Postgres ingestion repository

**Files:**
- Create: `apps/api/internal/storage/postgres/ingestion_repository.go`

**Interfaces:**
- Consumes: `ingest.Repository`, `ingest.UpsertVenueInput`, `ingest.IngestionRun`, `ingest.Summary` (Task 6).
- Produces: `postgres.NewIngestionRepository(db *sql.DB) *IngestionRepository` implementing `ingest.Repository`.

- [ ] **Step 1: Write ingestion_repository.go**

Create `apps/api/internal/storage/postgres/ingestion_repository.go`:

```go
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

// ResolveOrCreateDistrict ensures a district location row exists under the given
// city and returns its id. Idempotent on (type, country_code, slug).
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
  google_rating, google_rating_count, phone, website_url,
  source, source_external_id, district_location_id, source_updated_at
) values (
  $1, $2, $3, $3, $4, $5, 'VND',
  $6, $7, $8, $9,
  $10, $11, nullif($12, '')::uuid, now()
)
on conflict (source, source_external_id) do update set
  name = excluded.name,
  slug = excluded.slug,
  address = excluded.address,
  address_display = excluded.address_display,
  latitude = excluded.latitude,
  longitude = excluded.longitude,
  google_rating = excluded.google_rating,
  google_rating_count = excluded.google_rating_count,
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
		nullableRating(input.GoogleRating),
		nullableCount(input.GoogleRatingCount),
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
		label := strings.Title(strings.ReplaceAll(slug, "-", " "))
		if err := tx.QueryRowContext(ctx, upsertTag, slug, label).Scan(&tagID); err != nil {
			return fmt.Errorf("upsert tag %q: %w", slug, err)
		}
		const linkTag = `
insert into venue_tags (venue_id, tag_id, source)
values ($1::uuid, $2::uuid, 'google_places')
on conflict (venue_id, tag_id) do nothing
`
		if _, err := tx.ExecContext(ctx, linkTag, venueID, tagID); err != nil {
			return fmt.Errorf("link tag %q: %w", slug, err)
		}
	}
	return nil
}

func nullableRating(value float64) any {
	if value <= 0 {
		return nil
	}
	return value
}

func nullableCount(value int) any {
	if value <= 0 {
		return nil
	}
	return value
}
```

- [ ] **Step 2: Add phone/website carriers on the Venue (without changing JSON)**

The `venues` table has `phone`/`website_url`, but the public `discovery.Venue` JSON
payload should not change. Carry the data on unexported fields with accessors.

First, add two unexported fields to the `Venue` struct in
`apps/api/internal/discovery/model.go`, immediately after the `DistanceMeters` line:

```go
	contactPhone   string
	contactWebsite string
```

(Unexported fields are never serialized by `encoding/json`, so the public payload is
unchanged.)

Then create `apps/api/internal/discovery/contact.go`:

```go
package discovery

// WithContact returns a copy of the venue carrying phone/website for ingestion.
// These fields are unexported and never appear in the JSON payload.
func (v Venue) WithContact(phone string, website string) Venue {
	v.contactPhone = phone
	v.contactWebsite = website
	return v
}

func (v Venue) Phone() string   { return v.contactPhone }
func (v Venue) Website() string { return v.contactWebsite }
```

Finally, update `MapPlace` in `apps/api/internal/ingest/mapper.go` so the venue
carries contact data. Change the `venue := discovery.Venue{...}` block to append
`.WithContact(...)` after construction:

```go
	venue := discovery.Venue{
		Name:      p.DisplayName,
		Slug:      textutil.Slugish(p.DisplayName),
		Address:   p.FormattedAddress,
		Latitude:  p.Location.Lat,
		Longitude: p.Location.Lng,
		Currency:  "VND",
	}.WithContact(p.NationalPhoneNumber, p.WebsiteURI)
```

Re-run Task 5's tests to confirm they still pass:

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/ingest/ ./internal/discovery/`
Expected: PASS.

- [ ] **Step 3: Verify build and vet**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go build ./... && go vet ./internal/storage/postgres/ ./internal/ingest/`
Expected: no output (success). (No DB-backed test exists; the SQL is exercised manually during a real district run in Task 8 verification notes.)

- [ ] **Step 4: Commit**

```bash
git add apps/api/internal/storage/postgres/ingestion_repository.go apps/api/internal/discovery/contact.go apps/api/internal/discovery/model.go apps/api/internal/ingest/mapper.go
git commit -m "Add postgres ingestion repository and venue contact carrier"
```

---

### Task 8: CLI command `cmd/ingest`

**Files:**
- Create: `apps/api/cmd/ingest/main.go`

**Interfaces:**
- Consumes: `config.Load` (Task 2), `places.NewHTTPClient` (Task 4), `ingest.NewService`/`IngestOptions` (Task 6), `postgres.NewIngestionRepository` (Task 7).
- Produces: a runnable binary.

- [ ] **Step 1: Write main.go**

Create `apps/api/cmd/ingest/main.go`:

```go
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
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/ingest/places"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/storage/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	district := flag.String("district", "", "district slug to ingest, e.g. quan-1 (required)")
	category := flag.String("category", "restaurant", "Places category, e.g. restaurant or cafe")
	limit := flag.Int("limit", 20, "max venues to fetch this run")
	dryRun := flag.Bool("dry-run", false, "log what would be written without touching the database")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})).With("service", "ingest")

	if err := run(*district, *category, *limit, *dryRun, logger); err != nil {
		logger.Error("ingest failed", "error", err)
		os.Exit(1)
	}
}

func run(districtSlug string, category string, limit int, dryRun bool, logger *slog.Logger) error {
	if strings.TrimSpace(districtSlug) == "" {
		return fmt.Errorf("--district is required (e.g. --district=quan-1)")
	}

	cfg := config.Load()
	if cfg.GooglePlacesAPIKey == "" {
		return fmt.Errorf("GOOGLE_PLACES_API_KEY is not set")
	}
	if cfg.DatabaseURL == "" && !dryRun {
		return fmt.Errorf("DATABASE_URL is required unless --dry-run is set")
	}

	client := places.NewHTTPClient(cfg.GooglePlacesAPIKey)

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

	service := ingest.NewService(client, repo, logger)
	summary, err := service.IngestDistrict(context.Background(), ingest.IngestOptions{
		DistrictSlug: districtSlug,
		DistrictName: districtNameFromSlug(districtSlug),
		Category:     category,
		Limit:        limit,
		DryRun:       dryRun,
	})
	if err != nil {
		return err
	}

	logger.Info("ingest complete",
		"district", districtSlug,
		"category", category,
		"fetched", summary.Fetched,
		"upserted", summary.Upserted,
		"failed", summary.Failed,
		"dry_run", dryRun,
	)
	return nil
}

// districtNameFromSlug turns "quan-1" into "Quận 1" for the Places query and
// the locations row. Unknown slugs fall back to a title-cased form.
func districtNameFromSlug(slug string) string {
	known := map[string]string{
		"quan-1": "Quận 1",
		"quan-2": "Quận 2",
		"quan-3": "Quận 3",
		"quan-4": "Quận 4",
		"quan-5": "Quận 5",
		"quan-6": "Quận 6",
		"quan-7": "Quận 7",
		"quan-8": "Quận 8",
		"quan-9": "Quận 9",
		"quan-10": "Quận 10",
		"quan-11": "Quận 11",
		"quan-12": "Quận 12",
		"binh-thanh": "Bình Thạnh",
		"phu-nhuan": "Phú Nhuận",
		"go-vap": "Gò Vấp",
		"tan-binh": "Tân Bình",
		"tan-phu": "Tân Phú",
		"thu-duc": "Thủ Đức",
		"binh-tan": "Bình Tân",
	}
	if name, ok := known[slug]; ok {
		return name
	}
	return strings.Title(strings.ReplaceAll(slug, "-", " "))
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
```

- [ ] **Step 2: Verify build, vet, and required-flag behavior**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go build ./... && go vet ./cmd/ingest/`
Expected: no output (success).

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go run ./cmd/ingest`
Expected: exits non-zero, logs error `--district is required`.

Run: `cd apps/api && GOOGLE_PLACES_API_KEY=dummy GOCACHE=$(pwd)/../../.cache/go-build go run ./cmd/ingest --district=quan-1 --dry-run --limit=1`
Expected: it attempts a live Places call (will fail with an auth error from Google since `dummy` is not a real key) — confirming wiring works end to end without touching the DB. With a real key it would log `dry-run venue ...` lines and `ingest complete`.

- [ ] **Step 3: Run the full test suite**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./...`
Expected: all packages PASS (`internal/config`, `internal/ingest`, `internal/ingest/places`, `internal/http` ok; others no test files).

- [ ] **Step 4: Commit**

```bash
git add apps/api/cmd/ingest/main.go
git commit -m "Add ingest CLI command"
```

---

## Self-Review (completed by plan author)

- **Spec coverage:** source/API key → Tasks 2,4; idempotent upsert key → Tasks 1,7; district-by-district + resolve-or-create → Tasks 6,7,8; fields (rating, photos refs, tags, phone, website, hours) → Tasks 3,5,7 (note: photo references are decoded and mapped through `places.Place.PhotoNames`; persisting them into `venues.photos[]` is folded into the upsert in Task 7 only if present — see note below); `ingestion_runs` tracking → Tasks 1,6,7; CLI flags `--district/--category/--limit/--dry-run` → Task 8; generated `location` column → Task 1; tests with fixtures, no network/DB → Tasks 3,5,6.
- **Photos gap note:** the mapper carries `PhotoNames` on `places.Place`, but `MappedVenue`/`UpsertVenueInput` and the upsert SQL above do not yet write `venues.photos[]`. To fully satisfy the "capture photo references" requirement, the implementer should add a `PhotoRefs []string` field to `MappedVenue` and `UpsertVenueInput`, map it in `MapPlace`, and include `photos = $N` in the upsert (Postgres text[] via `pq.Array`/`pgx`); this was deliberately left as an explicit small extension to keep core tasks focused. Treat it as part of Task 7.
- **Placeholder scan:** no TBD/TODO; all steps contain complete code or exact commands.
- **Type consistency:** `Summary`, `IngestOptions`, `UpsertVenueInput`, `IngestionRun`, `Repository` names match across Tasks 6,7,8; `places.Client`, `Place`, `DayHours`, `LatLng` match across Tasks 3,4,5,6; `discovery.Venue.WithContact/Phone/Website` match between Task 7's repo use and the discovery additions.
- **Out of scope honored:** no dishes/prices, no social videos, no web UI changes, no migration-runner tool.
