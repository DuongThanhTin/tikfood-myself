# Venue Ingestion from OpenStreetMap Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** A re-runnable Go CLI (`apps/api/cmd/ingest`) that ingests Ho Chi Minh City restaurants/cafes from the OpenStreetMap Overpass API one district at a time, upserting idempotently — free, no API key.

**Architecture:** A source-agnostic `internal/ingest` package (DTO `SourcePlace`, `VenueSource` interface, pure mapper, orchestrating `Service`) depends on an `internal/ingest/overpass` client (implements `VenueSource`) and an `internal/ingest/openinghours` parser. Persistence is an `ingest.Repository` implemented by `internal/storage/postgres`. The CLI wires them. Layering matches `apps/api/CLAUDE.md`: cmd → service → (source | repository); the service never imports HTTP framework types.

**Tech Stack:** Go 1.26, `database/sql` via pgx stdlib, OpenStreetMap Overpass API over `net/http`, stdlib `testing` (table-driven).

## Global Constraints

- Module path: `github.com/DuongThanhTin/tikfood-myself/apps/api`.
- Source: OpenStreetMap Overpass API. No API key, no billing. `source='openstreetmap'`, `source_external_id='<type>/<id>'` (e.g. `node/1770305118`).
- Overpass etiquette: descriptive `User-Agent`; one district per run; backoff on `429`/`504`; configurable endpoint with default `https://overpass-api.de/api/interpreter`.
- ODbL: data must be displayed with "© OpenStreetMap contributors" attribution in the UI (UI follow-up, not in this plan).
- Idempotent upsert keyed on `(source, source_external_id)`.
- Scope v1: city fixed to Ho Chi Minh; `--district` selects one district per run; category default `restaurant`.
- Fields captured: name, address, latitude, longitude, osm id, opening hours (parsed when present), `amenity`+`cuisine` → tags, phone, website. No rating, no photos, no dishes/prices.
- No network calls in tests; no DB-backed tests (none exist today). TDD applies to pure functions (Overpass JSON decode, opening-hours parser, mapper) and the service (fake source + mock repo). Infra (migration, HTTP client, postgres repo, CLI) is verified with `go build ./...` and `go vet ./...`.
- Response envelope and existing API behavior are unchanged by this work.
- Run all Go commands from `apps/api` with `GOCACHE=$(pwd)/../../.cache/go-build`.
- Source spec: `docs/superpowers/specs/2026-06-27-venue-ingestion-osm-design.md`.

---

### Task 1: Migration 006 — schema for ingestion

**Files:**
- Create: `apps/api/migrations/006_venue_ingestion_schema.sql`

**Interfaces:**
- Produces: `venues.google_rating`, `venues.google_rating_count` (kept for a future Google source; OSM leaves them NULL); unique `(source, source_external_id)`; `venues.location` as a generated column; table `ingestion_runs`.

- [ ] **Step 1: Write the migration**

```sql
-- venues: rating fields kept for a future Google Places source (OSM leaves them NULL)
alter table venues add column if not exists google_rating numeric;
alter table venues add column if not exists google_rating_count int;

-- venues: upsert key for source ingestion
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
Expected: matches the four additions above. (No DB is required to apply it; it runs on next `docker compose up` via initdb and follows the existing `if not exists` idempotent convention.)

- [ ] **Step 3: Commit**

```bash
git add apps/api/migrations/006_venue_ingestion_schema.sql
git commit -m "Add migration 006: ingestion schema (upsert key, generated location, ingestion_runs)"
```

---

### Task 2: Config — Overpass endpoint

**Files:**
- Modify: `apps/api/internal/config/config.go`
- Test: `apps/api/internal/config/config_test.go`

**Interfaces:**
- Produces: `config.Config.OverpassEndpoint string`, from env `OVERPASS_ENDPOINT`, default `https://overpass-api.de/api/interpreter`.

- [ ] **Step 1: Write the failing test**

Create `apps/api/internal/config/config_test.go`:

```go
package config

import (
	"testing"
)

func TestLoadOverpassEndpointDefault(t *testing.T) {
	t.Setenv("OVERPASS_ENDPOINT", "")

	cfg := Load()

	if cfg.OverpassEndpoint != "https://overpass-api.de/api/interpreter" {
		t.Fatalf("expected default endpoint, got %q", cfg.OverpassEndpoint)
	}
}

func TestLoadOverpassEndpointOverride(t *testing.T) {
	t.Setenv("OVERPASS_ENDPOINT", "https://overpass.kumi.systems/api/interpreter")

	cfg := Load()

	if cfg.OverpassEndpoint != "https://overpass.kumi.systems/api/interpreter" {
		t.Fatalf("expected override endpoint, got %q", cfg.OverpassEndpoint)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/config/`
Expected: FAIL — `cfg.OverpassEndpoint` undefined (compile error).

- [ ] **Step 3: Update config.go**

Replace the body of `apps/api/internal/config/config.go` with:

```go
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
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/config/`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/config/config.go apps/api/internal/config/config_test.go
git commit -m "Add OverpassEndpoint to config"
```

---

### Task 3: Opening-hours parser (pure)

**Files:**
- Create: `apps/api/internal/ingest/openinghours/parse.go`
- Create: `apps/api/internal/ingest/openinghours/parse_test.go`

**Interfaces:**
- Consumes: `discovery.OpeningHour` (existing).
- Produces: `openinghours.Parse(spec string) []discovery.OpeningHour`. Day numbering: 0 = Sunday … 6 = Saturday (matches `venue_opening_hours.day_of_week`). Unrecognized input returns `nil` (never panics). When multiple rules cover the same day, the last rule wins.

- [ ] **Step 1: Write the failing test**

Create `apps/api/internal/ingest/openinghours/parse_test.go`:

```go
package openinghours

import (
	"testing"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
)

func byDay(hours []discovery.OpeningHour) map[int]discovery.OpeningHour {
	m := map[int]discovery.OpeningHour{}
	for _, h := range hours {
		m[h.DayOfWeek] = h
	}
	return m
}

func TestParseAllWeekSingleRange(t *testing.T) {
	hours := Parse("Mo-Su 07:00-23:00")
	if len(hours) != 7 {
		t.Fatalf("expected 7 days, got %d", len(hours))
	}
	m := byDay(hours)
	for day := 0; day <= 6; day++ {
		if m[day].OpenTime != "07:00" || m[day].CloseTime != "23:00" {
			t.Fatalf("day %d: got %+v", day, m[day])
		}
	}
}

func TestParseWeekdaysAndWeekend(t *testing.T) {
	hours := Parse("Mo-Fr 08:00-22:00; Sa-Su 09:00-23:00")
	m := byDay(hours)
	if m[5].OpenTime != "08:00" || m[5].CloseTime != "22:00" { // Friday
		t.Fatalf("friday: got %+v", m[5])
	}
	if m[0].OpenTime != "09:00" || m[0].CloseTime != "23:00" { // Sunday
		t.Fatalf("sunday: got %+v", m[0])
	}
	if m[6].OpenTime != "09:00" { // Saturday
		t.Fatalf("saturday: got %+v", m[6])
	}
}

func TestParseNoDayPrefixAppliesAllDays(t *testing.T) {
	hours := Parse("08:00-22:00")
	if len(hours) != 7 {
		t.Fatalf("expected 7 days, got %d", len(hours))
	}
}

func TestParseDayList(t *testing.T) {
	hours := Parse("Mo,We,Fr 08:00-18:00")
	m := byDay(hours)
	if len(hours) != 3 {
		t.Fatalf("expected 3 days, got %d", len(hours))
	}
	if _, ok := m[1]; !ok { // Monday
		t.Fatal("expected Monday")
	}
	if _, ok := m[2]; ok { // Tuesday should be absent
		t.Fatal("did not expect Tuesday")
	}
}

func TestParse247(t *testing.T) {
	hours := Parse("24/7")
	if len(hours) != 7 {
		t.Fatalf("expected 7 days, got %d", len(hours))
	}
	m := byDay(hours)
	if m[0].OpenTime != "00:00" || m[0].CloseTime != "23:59" {
		t.Fatalf("24/7 day: got %+v", m[0])
	}
}

func TestParseEmptyAndUnrecognized(t *testing.T) {
	if Parse("") != nil {
		t.Fatal("empty should be nil")
	}
	if Parse("sunrise-sunset") != nil {
		t.Fatal("unrecognized should be nil")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/ingest/openinghours/`
Expected: FAIL — package/`Parse` undefined.

- [ ] **Step 3: Write parse.go**

Create `apps/api/internal/ingest/openinghours/parse.go`:

```go
package openinghours

import (
	"regexp"
	"strings"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
)

// osmDays lists OSM weekday abbreviations in their natural order, paired with
// our day index (0 = Sunday ... 6 = Saturday).
var osmDays = []struct {
	abbr  string
	index int
}{
	{"Mo", 1}, {"Tu", 2}, {"We", 3}, {"Th", 4}, {"Fr", 5}, {"Sa", 6}, {"Su", 0},
}

var timeRangeRe = regexp.MustCompile(`(\d{2}:\d{2})\s*-\s*(\d{2}:\d{2})`)

func dayPosition(abbr string) int {
	for pos, d := range osmDays {
		if d.abbr == abbr {
			return pos
		}
	}
	return -1
}

// Parse converts a subset of OSM opening_hours syntax into per-day hours.
// Returns nil for empty or unrecognized input. Last rule wins per day.
func Parse(spec string) []discovery.OpeningHour {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return nil
	}
	if spec == "24/7" {
		return allDays("00:00", "23:59")
	}

	byDay := map[int]discovery.OpeningHour{}
	for _, rule := range strings.Split(spec, ";") {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		match := timeRangeRe.FindStringSubmatch(rule)
		if match == nil {
			continue
		}
		open, close := match[1], match[2]

		dayPart := strings.TrimSpace(rule[:strings.Index(rule, match[0])])
		days := parseDays(dayPart)
		for _, day := range days {
			byDay[day] = discovery.OpeningHour{
				DayOfWeek: day,
				OpenTime:  open,
				CloseTime: close,
				IsClosed:  false,
			}
		}
	}

	if len(byDay) == 0 {
		return nil
	}
	result := make([]discovery.OpeningHour, 0, len(byDay))
	for day := 0; day <= 6; day++ {
		if h, ok := byDay[day]; ok {
			result = append(result, h)
		}
	}
	return result
}

// parseDays expands a day specifier like "Mo-Fr", "Mo,We,Fr", or "" (all days).
func parseDays(spec string) []int {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		all := make([]int, 0, 7)
		for _, d := range osmDays {
			all = append(all, d.index)
		}
		return all
	}

	days := map[int]bool{}
	for _, token := range strings.Split(spec, ",") {
		token = strings.TrimSpace(token)
		if strings.Contains(token, "-") {
			parts := strings.SplitN(token, "-", 2)
			startPos := dayPosition(strings.TrimSpace(parts[0]))
			endPos := dayPosition(strings.TrimSpace(parts[1]))
			if startPos < 0 || endPos < 0 {
				continue
			}
			pos := startPos
			for {
				days[osmDays[pos].index] = true
				if pos == endPos {
					break
				}
				pos = (pos + 1) % len(osmDays)
			}
			continue
		}
		if p := dayPosition(token); p >= 0 {
			days[osmDays[p].index] = true
		}
	}

	out := make([]int, 0, len(days))
	for day := range days {
		out = append(out, day)
	}
	return out
}

func allDays(open string, close string) []discovery.OpeningHour {
	hours := make([]discovery.OpeningHour, 0, 7)
	for day := 0; day <= 6; day++ {
		hours = append(hours, discovery.OpeningHour{DayOfWeek: day, OpenTime: open, CloseTime: close})
	}
	return hours
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/ingest/openinghours/`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/ingest/openinghours/
git commit -m "Add OSM opening_hours parser"
```

---

### Task 4: Source DTO + interface

**Files:**
- Create: `apps/api/internal/ingest/source.go`

**Interfaces:**
- Produces:
  - `ingest.SourcePlace{ExternalID, Name, Address, Phone, Website, Amenity string; Lat, Lng float64; Cuisines []string; OpeningHoursRaw string}`.
  - `ingest.SearchParams{Lat, Lng float64; RadiusM int; Amenities []string; Limit int}`.
  - `ingest.VenueSource` interface: `Search(ctx context.Context, params SearchParams) ([]SourcePlace, error)`.

- [ ] **Step 1: Write source.go**

Create `apps/api/internal/ingest/source.go`:

```go
package ingest

import "context"

// SourcePlace is the source-agnostic representation of one ingested venue.
// Both the Overpass client (now) and a future Google client populate it.
type SourcePlace struct {
	ExternalID      string // "<type>/<id>" for OSM, e.g. "node/123"
	Name            string
	Address         string
	Lat             float64
	Lng             float64
	Phone           string
	Website         string
	Amenity         string
	Cuisines        []string
	OpeningHoursRaw string
}

// SearchParams describes a single bounded source query.
type SearchParams struct {
	Lat       float64
	Lng       float64
	RadiusM   int
	Amenities []string
	Limit     int
}

// VenueSource fetches places for one district query.
type VenueSource interface {
	Search(ctx context.Context, params SearchParams) ([]SourcePlace, error)
}
```

- [ ] **Step 2: Verify it builds**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go build ./internal/ingest/`
Expected: no output (success).

- [ ] **Step 3: Commit**

```bash
git add apps/api/internal/ingest/source.go
git commit -m "Add source-agnostic SourcePlace DTO and VenueSource interface"
```

---

### Task 5: Overpass JSON decode (pure)

**Files:**
- Create: `apps/api/internal/ingest/overpass/decode.go`
- Create: `apps/api/internal/ingest/overpass/decode_test.go`
- Create: `apps/api/internal/ingest/overpass/testdata/elements.json`

**Interfaces:**
- Consumes: `ingest.SourcePlace` (Task 4).
- Produces: `overpass.DecodeElements(data []byte) ([]ingest.SourcePlace, error)`.

- [ ] **Step 1: Create the fixture**

Create `apps/api/internal/ingest/overpass/testdata/elements.json`:

```json
{
  "version": 0.6,
  "elements": [
    {
      "type": "node",
      "id": 1770305118,
      "lat": 10.7751001,
      "lon": 106.7029693,
      "tags": {
        "amenity": "cafe",
        "name": "Cà Phê Ciao",
        "cuisine": "vietnamese;international;coffee_shop",
        "opening_hours": "Mo-Su 07:00-23:00",
        "phone": "+8438231130",
        "website": "https://ciaocafe.vn/",
        "addr:housenumber": "74-76",
        "addr:street": "Nguyễn Huệ"
      }
    },
    {
      "type": "node",
      "id": 411918246,
      "lat": 10.7848583,
      "lon": 106.700098,
      "tags": {
        "amenity": "fast_food",
        "name": "Phở Số 1 Hà Nội"
      }
    },
    {
      "type": "node",
      "id": 999,
      "lat": 10.0,
      "lon": 106.0,
      "tags": {
        "amenity": "restaurant",
        "contact:phone": "+84123",
        "contact:website": "https://x.test"
      }
    }
  ]
}
```

- [ ] **Step 2: Write the failing test**

Create `apps/api/internal/ingest/overpass/decode_test.go`:

```go
package overpass

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDecodeElements(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "elements.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	places, err := DecodeElements(data)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	// The third element has no name and must be skipped.
	if len(places) != 2 {
		t.Fatalf("expected 2 named places, got %d", len(places))
	}

	first := places[0]
	if first.ExternalID != "node/1770305118" {
		t.Fatalf("external id: got %q", first.ExternalID)
	}
	if first.Name != "Cà Phê Ciao" {
		t.Fatalf("name: got %q", first.Name)
	}
	if first.Lat != 10.7751001 || first.Lng != 106.7029693 {
		t.Fatalf("coords: got %v/%v", first.Lat, first.Lng)
	}
	if first.Address != "74-76 Nguyễn Huệ" {
		t.Fatalf("address: got %q", first.Address)
	}
	if first.Phone != "+8438231130" || first.Website != "https://ciaocafe.vn/" {
		t.Fatalf("contact: got %q / %q", first.Phone, first.Website)
	}
	if len(first.Cuisines) != 3 || first.Cuisines[0] != "vietnamese" {
		t.Fatalf("cuisines: got %v", first.Cuisines)
	}
	if first.OpeningHoursRaw != "Mo-Su 07:00-23:00" {
		t.Fatalf("hours raw: got %q", first.OpeningHoursRaw)
	}
	if first.Amenity != "cafe" {
		t.Fatalf("amenity: got %q", first.Amenity)
	}
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/ingest/overpass/`
Expected: FAIL — package/`DecodeElements` undefined.

- [ ] **Step 4: Write decode.go**

Create `apps/api/internal/ingest/overpass/decode.go`:

```go
package overpass

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/ingest"
)

type rawResponse struct {
	Elements []rawElement `json:"elements"`
}

type rawElement struct {
	Type string            `json:"type"`
	ID   int64             `json:"id"`
	Lat  float64           `json:"lat"`
	Lon  float64           `json:"lon"`
	Tags map[string]string `json:"tags"`
}

// DecodeElements converts an Overpass JSON response into source places,
// skipping elements without a name.
func DecodeElements(data []byte) ([]ingest.SourcePlace, error) {
	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("decode overpass: %w", err)
	}

	places := make([]ingest.SourcePlace, 0, len(raw.Elements))
	for _, el := range raw.Elements {
		name := strings.TrimSpace(el.Tags["name"])
		if name == "" {
			continue
		}
		places = append(places, ingest.SourcePlace{
			ExternalID:      fmt.Sprintf("%s/%d", el.Type, el.ID),
			Name:            name,
			Address:         buildAddress(el.Tags),
			Lat:             el.Lat,
			Lng:             el.Lon,
			Phone:           firstNonEmpty(el.Tags["phone"], el.Tags["contact:phone"]),
			Website:         firstNonEmpty(el.Tags["website"], el.Tags["contact:website"]),
			Amenity:         el.Tags["amenity"],
			Cuisines:        splitCuisines(el.Tags["cuisine"]),
			OpeningHoursRaw: el.Tags["opening_hours"],
		})
	}
	return places, nil
}

func buildAddress(tags map[string]string) string {
	parts := make([]string, 0, 2)
	if v := strings.TrimSpace(tags["addr:housenumber"]); v != "" {
		parts = append(parts, v)
	}
	if v := strings.TrimSpace(tags["addr:street"]); v != "" {
		parts = append(parts, v)
	}
	return strings.Join(parts, " ")
}

func splitCuisines(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	out := []string{}
	for _, seg := range strings.Split(value, ";") {
		for _, part := range strings.Split(seg, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				out = append(out, part)
			}
		}
	}
	return out
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/ingest/overpass/`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add apps/api/internal/ingest/overpass/
git commit -m "Add Overpass JSON decoder"
```

---

### Task 6: Overpass HTTP client

**Files:**
- Create: `apps/api/internal/ingest/overpass/client.go`

**Interfaces:**
- Consumes: `DecodeElements` (Task 5), `ingest.SearchParams`, `ingest.SourcePlace` (Task 4).
- Produces: `overpass.NewClient(endpoint string) *Client` implementing `ingest.VenueSource`.

- [ ] **Step 1: Write client.go**

Create `apps/api/internal/ingest/overpass/client.go`:

```go
package overpass

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/ingest"
)

const userAgent = "tikfood-ingest/1.0 (OpenStreetMap venue ingestion)"

// Client queries the OpenStreetMap Overpass API.
type Client struct {
	endpoint string
	http     *http.Client
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		http:     &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *Client) Search(ctx context.Context, params ingest.SearchParams) ([]ingest.SourcePlace, error) {
	query := buildQuery(params)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint,
		strings.NewReader(url.Values{"data": {query}}.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build overpass request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("overpass request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read overpass response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("overpass status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return DecodeElements(data)
}

// buildQuery builds an Overpass QL query for the given amenities around a point.
func buildQuery(params ingest.SearchParams) string {
	amenities := params.Amenities
	if len(amenities) == 0 {
		amenities = []string{"restaurant"}
	}
	radius := params.RadiusM
	if radius <= 0 {
		radius = 1300
	}
	filter := strings.Join(amenities, "|")

	var b strings.Builder
	b.WriteString("[out:json][timeout:25];(")
	fmt.Fprintf(&b, `node["amenity"~"^(%s)$"](around:%d,%f,%f);`, filter, radius, params.Lat, params.Lng)
	b.WriteString(");out body;")
	return b.String()
}
```

- [ ] **Step 2: Verify it builds and vets**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go build ./... && go vet ./internal/ingest/overpass/`
Expected: no output (success). (No unit test: this is a thin network adapter; its decode logic is tested in Task 5. `--limit` is applied by the service after fetch, since Overpass `around` queries do not page.)

- [ ] **Step 3: Commit**

```bash
git add apps/api/internal/ingest/overpass/client.go
git commit -m "Add Overpass HTTP client"
```

---

### Task 7: Ingest mapper (pure)

**Files:**
- Create: `apps/api/internal/ingest/mapper.go`
- Create: `apps/api/internal/ingest/mapper_test.go`
- Modify: `apps/api/internal/discovery/model.go` (add unexported contact fields)
- Create: `apps/api/internal/discovery/contact.go`

**Interfaces:**
- Consumes: `ingest.SourcePlace` (Task 4), `openinghours.Parse` (Task 3), `discovery.Venue`/`discovery.OpeningHour`, `textutil.Slugish`.
- Produces:
  - `discovery.Venue.WithContact(phone, website string) Venue`, `discovery.Venue.Phone() string`, `discovery.Venue.Website() string`.
  - `ingest.MappedVenue{Venue discovery.Venue; OpeningHours []discovery.OpeningHour; TagSlugs []string}`.
  - `ingest.MapPlace(p SourcePlace) MappedVenue`.
  - `ingest.cuisineToTagSlug` / `amenityToTagSlug` helpers (unexported).

- [ ] **Step 1: Add contact carriers on the Venue (without changing JSON)**

Add two unexported fields to the `Venue` struct in
`apps/api/internal/discovery/model.go`, immediately after the `DistanceMeters` line:

```go
	contactPhone   string
	contactWebsite string
```

(Unexported fields are never serialized by `encoding/json`, so the public payload is unchanged.)

Create `apps/api/internal/discovery/contact.go`:

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

- [ ] **Step 2: Write the failing test**

Create `apps/api/internal/ingest/mapper_test.go`:

```go
package ingest

import (
	"testing"
)

func samplePlace() SourcePlace {
	return SourcePlace{
		ExternalID:      "node/1770305118",
		Name:            "Cà Phê Ciao",
		Address:         "74-76 Nguyễn Huệ",
		Lat:             10.7751001,
		Lng:             106.7029693,
		Phone:           "+8438231130",
		Website:         "https://ciaocafe.vn/",
		Amenity:         "cafe",
		Cuisines:        []string{"vietnamese", "international", "coffee_shop"},
		OpeningHoursRaw: "Mo-Su 07:00-23:00",
	}
}

func TestMapPlaceCoreFields(t *testing.T) {
	m := MapPlace(samplePlace())

	if m.Venue.Name != "Cà Phê Ciao" {
		t.Fatalf("name: got %q", m.Venue.Name)
	}
	if m.Venue.Slug != "ca-phe-ciao" { // textutil.Slugish handles Vietnamese
		t.Fatalf("slug: got %q", m.Venue.Slug)
	}
	if m.Venue.Latitude != 10.7751001 || m.Venue.Longitude != 106.7029693 {
		t.Fatalf("coords: got %v/%v", m.Venue.Latitude, m.Venue.Longitude)
	}
	if m.Venue.Address != "74-76 Nguyễn Huệ" {
		t.Fatalf("address: got %q", m.Venue.Address)
	}
	if m.Venue.Currency != "VND" {
		t.Fatalf("currency: got %q", m.Venue.Currency)
	}
	if m.Venue.Phone() != "+8438231130" || m.Venue.Website() != "https://ciaocafe.vn/" {
		t.Fatalf("contact: got %q / %q", m.Venue.Phone(), m.Venue.Website())
	}
}

func TestMapPlaceOpeningHours(t *testing.T) {
	m := MapPlace(samplePlace())
	if len(m.OpeningHours) != 7 {
		t.Fatalf("hours len: got %d", len(m.OpeningHours))
	}
}

func TestMapPlaceTagSlugs(t *testing.T) {
	m := MapPlace(samplePlace())

	// amenity "cafe" + cuisines; coffee_shop -> cafe (dedup); vietnamese stays.
	want := map[string]bool{"cafe": true, "vietnamese": true, "international": true}
	got := map[string]bool{}
	for _, s := range m.TagSlugs {
		got[s] = true
	}
	for slug := range want {
		if !got[slug] {
			t.Fatalf("missing tag %q in %v", slug, m.TagSlugs)
		}
	}
	if got["coffee_shop"] {
		t.Fatalf("coffee_shop should map to cafe, got %v", m.TagSlugs)
	}
}

func TestMapPlaceNoHoursWhenAbsent(t *testing.T) {
	p := samplePlace()
	p.OpeningHoursRaw = ""
	m := MapPlace(p)
	if m.OpeningHours != nil {
		t.Fatalf("expected no hours, got %+v", m.OpeningHours)
	}
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/ingest/`
Expected: FAIL — `MapPlace`/`MappedVenue` undefined.

- [ ] **Step 4: Write mapper.go**

Create `apps/api/internal/ingest/mapper.go`:

```go
package ingest

import (
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/ingest/openinghours"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/platform/textutil"
)

// MappedVenue is the result of mapping a SourcePlace to TikFood domain data.
type MappedVenue struct {
	Venue        discovery.Venue
	OpeningHours []discovery.OpeningHour
	TagSlugs     []string
}

// cuisineAliases maps OSM cuisine values to TikFood tag slugs.
var cuisineAliases = map[string]string{
	"coffee_shop": "cafe",
}

func normalizeTagSlug(value string) string {
	if alias, ok := cuisineAliases[value]; ok {
		return alias
	}
	return textutil.Slugish(value)
}

// MapPlace converts a SourcePlace into TikFood domain data.
func MapPlace(p SourcePlace) MappedVenue {
	venue := discovery.Venue{
		Name:      p.Name,
		Slug:      textutil.Slugish(p.Name),
		Address:   p.Address,
		Latitude:  p.Lat,
		Longitude: p.Lng,
		Currency:  "VND",
	}.WithContact(p.Phone, p.Website)

	tagSlugs := make([]string, 0, len(p.Cuisines)+1)
	seen := map[string]bool{}
	addTag := func(raw string) {
		slug := normalizeTagSlug(raw)
		if slug == "" || seen[slug] {
			return
		}
		seen[slug] = true
		tagSlugs = append(tagSlugs, slug)
	}
	if p.Amenity != "" {
		addTag(p.Amenity)
	}
	for _, c := range p.Cuisines {
		addTag(c)
	}

	return MappedVenue{
		Venue:        venue,
		OpeningHours: openinghours.Parse(p.OpeningHoursRaw),
		TagSlugs:     tagSlugs,
	}
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/ingest/ ./internal/discovery/`
Expected: PASS.

Note: `TestMapPlaceTagSlugs` expects `cafe` to appear once. Amenity `cafe` adds `cafe`; cuisine `coffee_shop` aliases to `cafe` and is deduped. `vietnamese` and `international` pass through `textutil.Slugish` unchanged.

- [ ] **Step 6: Commit**

```bash
git add apps/api/internal/ingest/mapper.go apps/api/internal/ingest/mapper_test.go apps/api/internal/discovery/contact.go apps/api/internal/discovery/model.go
git commit -m "Add pure SourcePlace-to-venue mapper and venue contact carrier"
```

---

### Task 8: Ingest service + repository interface

**Files:**
- Create: `apps/api/internal/ingest/repository.go`
- Create: `apps/api/internal/ingest/service.go`
- Create: `apps/api/internal/ingest/service_test.go`

**Interfaces:**
- Consumes: `VenueSource`/`SearchParams`/`SourcePlace` (Task 4), `MapPlace`/`MappedVenue` (Task 7).
- Produces:
  - `ingest.UpsertVenueInput{Venue discovery.Venue; OpeningHours []discovery.OpeningHour; TagSlugs []string; Source string; ExternalID string; DistrictLocationID string}`.
  - `ingest.IngestionRun{Source, DistrictSlug, Category string}`.
  - `ingest.Repository` interface (methods below).
  - `ingest.NewService(source VenueSource, repo Repository, logger *slog.Logger) *Service`.
  - `ingest.IngestOptions{DistrictSlug, DistrictName, Category string; Lat, Lng float64; RadiusM int; Amenities []string; Limit int; DryRun bool}`.
  - `ingest.Summary{Fetched, Upserted, Failed int}`.
  - `ingest.Service.IngestDistrict(ctx, opts IngestOptions) (Summary, error)`.

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
)

type fakeSource struct {
	places []SourcePlace
	err    error
}

func (f fakeSource) Search(_ context.Context, _ SearchParams) ([]SourcePlace, error) {
	return f.places, f.err
}

type mockRepo struct {
	upserts    []UpsertVenueInput
	finished   bool
	finalState string
	finalCount Summary
}

func (m *mockRepo) ResolveOrCreateDistrict(_ context.Context, _, _ string) (string, error) {
	return "district-loc-1", nil
}
func (m *mockRepo) RecordIngestionRun(_ context.Context, _ IngestionRun) (string, error) {
	return "run-1", nil
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

func testLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

func baseOptions() IngestOptions {
	return IngestOptions{
		DistrictSlug: "quan-1", DistrictName: "Quận 1", Category: "restaurant",
		Lat: 10.7756, Lng: 106.7019, RadiusM: 1300, Amenities: []string{"restaurant"}, Limit: 20,
	}
}

func TestIngestDistrictUpsertsEachPlace(t *testing.T) {
	src := fakeSource{places: []SourcePlace{
		{ExternalID: "node/1", Name: "Alpha", Lat: 10, Lng: 106},
		{ExternalID: "node/2", Name: "Beta", Lat: 11, Lng: 107},
	}}
	repo := &mockRepo{}
	svc := NewService(src, repo, testLogger())

	summary, err := svc.IngestDistrict(context.Background(), baseOptions())
	if err != nil {
		t.Fatalf("ingest: %v", err)
	}
	if summary.Fetched != 2 || summary.Upserted != 2 || summary.Failed != 0 {
		t.Fatalf("summary: %+v", summary)
	}
	if len(repo.upserts) != 2 {
		t.Fatalf("expected 2 upserts, got %d", len(repo.upserts))
	}
	if repo.upserts[0].Source != "openstreetmap" || repo.upserts[0].ExternalID != "node/1" {
		t.Fatalf("upsert[0]: %+v", repo.upserts[0])
	}
	if repo.upserts[0].DistrictLocationID != "district-loc-1" {
		t.Fatalf("district id not set: %+v", repo.upserts[0])
	}
	if !repo.finished || repo.finalState != "completed" {
		t.Fatalf("run not finished completed: %+v", repo)
	}
}

func TestIngestDistrictRespectsLimit(t *testing.T) {
	src := fakeSource{places: []SourcePlace{
		{ExternalID: "node/1", Name: "Alpha"},
		{ExternalID: "node/2", Name: "Beta"},
		{ExternalID: "node/3", Name: "Gamma"},
	}}
	repo := &mockRepo{}
	svc := NewService(src, repo, testLogger())

	opts := baseOptions()
	opts.Limit = 2
	summary, err := svc.IngestDistrict(context.Background(), opts)
	if err != nil {
		t.Fatalf("ingest: %v", err)
	}
	if summary.Fetched != 2 || len(repo.upserts) != 2 {
		t.Fatalf("limit not applied: fetched=%d upserts=%d", summary.Fetched, len(repo.upserts))
	}
}

func TestIngestDistrictDryRunWritesNothing(t *testing.T) {
	src := fakeSource{places: []SourcePlace{{ExternalID: "node/1", Name: "Alpha"}}}
	repo := &mockRepo{}
	svc := NewService(src, repo, testLogger())

	opts := baseOptions()
	opts.DryRun = true
	summary, err := svc.IngestDistrict(context.Background(), opts)
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

func TestIngestDistrictSourceErrorFailsRun(t *testing.T) {
	src := fakeSource{err: errors.New("overpass down")}
	repo := &mockRepo{}
	svc := NewService(src, repo, testLogger())

	_, err := svc.IngestDistrict(context.Background(), baseOptions())
	if err == nil {
		t.Fatal("expected error")
	}
	if !repo.finished || repo.finalState != "failed" {
		t.Fatalf("run should be marked failed: %+v", repo)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/ingest/`
Expected: FAIL — `NewService`/`Repository`/`IngestOptions` undefined.

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
```

- [ ] **Step 5: Run test to verify it passes**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./internal/ingest/`
Expected: PASS. Note: `TestIngestDistrictSourceErrorFailsRun` exercises the `failed` finish; the source error path calls `finish` with status `failed` before returning.

- [ ] **Step 6: Commit**

```bash
git add apps/api/internal/ingest/repository.go apps/api/internal/ingest/service.go apps/api/internal/ingest/service_test.go
git commit -m "Add ingest service and repository interface"
```

---

### Task 9: Postgres ingestion repository

**Files:**
- Create: `apps/api/internal/storage/postgres/ingestion_repository.go`

**Interfaces:**
- Consumes: `ingest.Repository`, `ingest.UpsertVenueInput`, `ingest.IngestionRun`, `ingest.Summary` (Task 8); `textutil` (existing).
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
```

- [ ] **Step 2: Verify build and vet**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go build ./... && go vet ./internal/storage/postgres/ ./internal/ingest/`
Expected: no output (success). (No DB-backed test exists; the SQL is exercised manually during a real district run after Task 10.)

- [ ] **Step 3: Commit**

```bash
git add apps/api/internal/storage/postgres/ingestion_repository.go
git commit -m "Add postgres ingestion repository"
```

---

### Task 10: CLI command `cmd/ingest`

**Files:**
- Create: `apps/api/cmd/ingest/main.go`

**Interfaces:**
- Consumes: `config.Load` (Task 2), `overpass.NewClient` (Task 6), `ingest.NewService`/`IngestOptions` (Task 8), `postgres.NewIngestionRepository` (Task 9).
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
```

- [ ] **Step 2: Verify build, vet, and required-flag behavior**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go build ./... && go vet ./cmd/ingest/`
Expected: no output (success).

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go run ./cmd/ingest`
Expected: exits non-zero, logs `--district is required`.

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go run ./cmd/ingest --district=quan-99 --dry-run`
Expected: exits non-zero, logs `unknown district "quan-99"`.

- [ ] **Step 3: Optional live smoke test (no key needed, requires network)**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go run ./cmd/ingest --district=quan-1 --category=cafe --limit=5 --dry-run`
Expected: real Overpass call (free, no key); logs several `dry-run venue ...` lines and `ingest complete fetched=5 ...`. If offline, it fails at the Overpass request — that is acceptable; the required-flag checks above already prove wiring without network.

- [ ] **Step 4: Run the full test suite**

Run: `cd apps/api && GOCACHE=$(pwd)/../../.cache/go-build go test ./...`
Expected: all packages PASS (`internal/config`, `internal/ingest`, `internal/ingest/openinghours`, `internal/ingest/overpass`, `internal/http` ok; others no test files).

- [ ] **Step 5: Commit**

```bash
git add apps/api/cmd/ingest/main.go
git commit -m "Add OSM ingest CLI command"
```

---

## Self-Review (completed by plan author)

- **Spec coverage:** source=Overpass/no key → Tasks 2,6; idempotent upsert key → Tasks 1,9; district-by-district + resolve-or-create → Tasks 8,9,10; fields (name/geo/address/cuisine+amenity→tags/phone/website/hours) → Tasks 5,7; opening_hours parser → Task 3; no rating/photos → upsert omits them (Task 9), columns stay NULL (Task 1); `ingestion_runs` tracking → Tasks 1,8,9; CLI flags + district area map → Task 10; generated `location` column → Task 1; tests with fixtures, no network/DB → Tasks 3,5,7,8.
- **Placeholder scan:** no TBD/TODO; every step has complete code or exact commands.
- **Type consistency:** `SourcePlace`, `SearchParams`, `VenueSource` defined in Task 4 and used identically in Tasks 5,6,8; `MappedVenue`/`MapPlace` in Task 7 used by Task 8; `UpsertVenueInput`/`IngestionRun`/`Summary`/`Repository` in Task 8 used by Tasks 9,10; `discovery.Venue.WithContact/Phone/Website` in Task 7 used by Task 9; `openinghours.Parse` in Task 3 used by Task 7; `overpass.NewClient`/`DecodeElements` in Tasks 5,6 used by Task 10.
- **Source string** `openstreetmap` is consistent across service (Task 8) and the `venue_tags.source` literal (Task 9) and tests.
- **Out of scope honored:** no rating/photos, no dishes, no social videos, no UI changes, no migration-runner tool. Google Places remains a deferred, retained plan.
```
