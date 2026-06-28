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
