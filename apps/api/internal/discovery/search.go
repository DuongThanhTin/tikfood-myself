package discovery

import (
	"strings"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/platform/textutil"
)

func normalizeSearch(search VenueSearch) VenueSearch {
	search.Query = strings.ToLower(strings.TrimSpace(search.Query))
	search.City = strings.TrimSpace(search.City)
	search.District = strings.TrimSpace(search.District)
	search.Dish = strings.ToLower(strings.TrimSpace(search.Dish))
	search.Tags = normalizeTokens(search.Tags)
	search.Platforms = normalizeTokens(search.Platforms)

	if search.Sort == "" {
		search.Sort = VenueSortTrending
	}
	if search.Limit <= 0 {
		search.Limit = 50
	}
	if search.Limit > 100 {
		search.Limit = 100
	}

	return search
}

func normalizeTokens(values []string) []string {
	results := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, value := range values {
		normalized := textutil.Slugish(value)
		if normalized == "" || seen[normalized] {
			continue
		}
		seen[normalized] = true
		results = append(results, normalized)
	}
	return results
}

func matchesCity(value string, target string) bool {
	return normalizeCityAlias(value) == normalizeCityAlias(target)
}

func normalizeCityAlias(value string) string {
	normalized := strings.ReplaceAll(textutil.Compact(value), "thanhpho", "tp")
	switch normalized {
	case "hochiminhcity", "hochiminh", "hcm", "tphcm", "tphochiminh":
		return "ho-chi-minh"
	default:
		return normalized
	}
}

func matchesDistrict(value string, target string) bool {
	return normalizeDistrictAlias(value) == normalizeDistrictAlias(target)
}

func normalizeDistrictAlias(value string) string {
	switch textutil.Compact(value) {
	case "quan1", "district1", "q1":
		return "quan-1"
	case "quan3", "district3", "q3":
		return "quan-3"
	default:
		return textutil.Compact(value)
	}
}

func hasDish(dishes []string, target string) bool {
	for _, dish := range dishes {
		if strings.EqualFold(dish, target) {
			return true
		}
	}
	return false
}

func hasAnyPlatform(videos []SocialVideo, platforms []string) bool {
	for _, video := range videos {
		for _, platform := range platforms {
			if strings.EqualFold(video.Platform, platform) {
				return true
			}
		}
	}
	return false
}

func matchesQuery(venue Venue, query string) bool {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return true
	}

	fields := []string{
		venue.Name,
		venue.ShortDescription,
		venue.About,
		venue.Address,
		venue.District,
	}
	fields = append(fields, venue.Categories...)
	fields = append(fields, venue.TrendingDishes...)

	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), query) {
			return true
		}
	}
	return false
}

func hasAnyTag(categories []string, tags []string) bool {
	for _, category := range categories {
		for _, tag := range tags {
			if strings.EqualFold(category, tag) {
				return true
			}
		}
	}
	return false
}
