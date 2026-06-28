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
