package app

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/ingest"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/ingest/overpass"
)

// osmDemoArea mirrors the district centers used by cmd/ingest so the in-memory
// preview hits the same Overpass area as a real run.
type osmDemoArea struct {
	name    string
	lat     float64
	lng     float64
	radiusM int
}

var osmDemoAreas = map[string]osmDemoArea{
	"quan-1": {"Quận 1", 10.7756, 106.7019, 1300},
	"quan-3": {"Quận 3", 10.7860, 106.6840, 1200},
	"quan-4": {"Quận 4", 10.7578, 106.7050, 1200},
	"quan-5": {"Quận 5", 10.7540, 106.6630, 1300},
}

// seedFallbackFromOSM ingests live OpenStreetMap venues for one district and
// appends them to the in-memory repository so they render on the UI. This runs
// the real Overpass client + mapper pipeline; only persistence differs.
func seedFallbackFromOSM(repo *discovery.FallbackVenueRepository, endpoint string, districtSlug string, limit int, logger *slog.Logger) error {
	area, ok := osmDemoAreas[districtSlug]
	if !ok {
		return fmt.Errorf("unknown district %q for INGEST_ON_START", districtSlug)
	}

	client := overpass.NewClient(endpoint)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	places, err := client.Search(ctx, ingest.SearchParams{
		Lat: area.lat, Lng: area.lng, RadiusM: area.radiusM,
		Amenities: []string{"restaurant", "cafe", "fast_food"},
	})
	if err != nil {
		return fmt.Errorf("overpass search: %w", err)
	}

	// Dedupe by name so chain branches (KFC, Lotteria...) show once, giving a
	// more varied preview, then cap at the requested limit.
	seenName := map[string]bool{}
	venues := make([]discovery.Venue, 0, limit)
	for _, place := range places {
		if limit > 0 && len(venues) >= limit {
			break
		}
		nameKey := strings.ToLower(strings.TrimSpace(place.Name))
		if seenName[nameKey] {
			continue
		}
		seenName[nameKey] = true

		mapped := ingest.MapPlace(place)
		v := mapped.Venue

		v.ID = "osm-" + place.ExternalID
		v.City = "Thành phố Hồ Chí Minh"
		v.District = area.name
		v.Categories = nonNilStrings(mapped.TagSlugs)
		v.OpeningHours = mapped.OpeningHours
		v.SocialVideos = []discovery.SocialVideo{}
		v.TrendingDishes = []string{}
		v.ShortDescription = "Quán nhập từ OpenStreetMap."
		if v.Address != "" {
			v.About = "Địa điểm tại " + v.Address + ", " + area.name + " — dữ liệu từ OpenStreetMap."
		} else {
			v.About = "Địa điểm tại " + area.name + " — dữ liệu từ OpenStreetMap."
		}
		v.AISummary = "Nhập từ OpenStreetMap (© OpenStreetMap contributors)."

		venues = append(venues, v)
	}

	repo.AddVenues(venues)
	logger.Info("seeded in-memory storage from OpenStreetMap",
		"district", districtSlug, "fetched", len(places), "added", len(venues))
	return nil
}

func nonNilStrings(in []string) []string {
	if in == nil {
		return []string{}
	}
	return in
}
