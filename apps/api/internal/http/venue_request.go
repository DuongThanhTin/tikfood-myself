package http

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
	"github.com/gin-gonic/gin"
)

func parseVenueSearch(c *gin.Context) (discovery.VenueSearch, *errorResponse) {
	search := discovery.VenueSearch{
		Query:    strings.TrimSpace(c.Query("q")),
		District: strings.TrimSpace(c.Query("district")),
		Dish:     strings.TrimSpace(c.Query("dish")),
		Tags:     splitCSV(c.Query("tags")),
		Sort:     discovery.VenueSort(strings.TrimSpace(c.DefaultQuery("sort", string(discovery.VenueSortTrending)))),
		Limit:    50,
	}

	if search.Query != "" && len(search.Query) > 120 {
		return search, invalidQuery("q", "Query must be 120 characters or fewer.")
	}

	if c.Query("lat") != "" || c.Query("lng") != "" {
		lat, err := parseFloatParam(c, "lat")
		if err != nil {
			return search, invalidQuery("lat", err.Error())
		}
		lng, err := parseFloatParam(c, "lng")
		if err != nil {
			return search, invalidQuery("lng", err.Error())
		}
		if lat < -90 || lat > 90 {
			return search, invalidQuery("lat", "Latitude must be between -90 and 90.")
		}
		if lng < -180 || lng > 180 {
			return search, invalidQuery("lng", "Longitude must be between -180 and 180.")
		}
		search.Lat = &lat
		search.Lng = &lng
	}

	if c.Query("radius_m") != "" {
		radius, err := parseIntParam(c, "radius_m")
		if err != nil {
			return search, invalidQuery("radius_m", err.Error())
		}
		if radius < 0 || radius > 50000 {
			return search, invalidQuery("radius_m", "Radius must be between 0 and 50000 meters.")
		}
		search.RadiusM = radius
	}

	if c.Query("max_price_vnd") != "" {
		maxPrice, err := parseIntParam(c, "max_price_vnd")
		if err != nil {
			return search, invalidQuery("max_price_vnd", err.Error())
		}
		if maxPrice < 0 {
			return search, invalidQuery("max_price_vnd", "Max price must be greater than or equal to 0.")
		}
		search.MaxPriceVND = maxPrice
	}

	if c.Query("limit") != "" {
		limit, err := parseIntParam(c, "limit")
		if err != nil {
			return search, invalidQuery("limit", err.Error())
		}
		if limit < 1 || limit > 100 {
			return search, invalidQuery("limit", "Limit must be between 1 and 100.")
		}
		search.Limit = limit
	}

	if c.Query("open_now") != "" {
		openNow, err := strconv.ParseBool(c.Query("open_now"))
		if err != nil {
			return search, invalidQuery("open_now", "Open now must be true or false.")
		}
		search.OpenNow = openNow
	}

	switch search.Sort {
	case discovery.VenueSortTrending, discovery.VenueSortVideos, discovery.VenueSortDistance, discovery.VenueSortPrice:
	default:
		return search, invalidQuery("sort", "Sort must be one of trending, videos, distance, or price.")
	}
	if search.Sort == discovery.VenueSortDistance && (search.Lat == nil || search.Lng == nil) {
		return search, invalidQuery("sort", "Distance sort requires lat and lng.")
	}

	return search, nil
}

func splitCSV(value string) []string {
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	results := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			results = append(results, trimmed)
		}
	}
	return results
}

func parseFloatParam(c *gin.Context, name string) (float64, error) {
	value := strings.TrimSpace(c.Query(name))
	if value == "" {
		return 0, fmt.Errorf("%s is required", name)
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number", name)
	}
	return parsed, nil
}

func parseIntParam(c *gin.Context, name string) (int, error) {
	value := strings.TrimSpace(c.Query(name))
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", name)
	}
	return parsed, nil
}

func invalidQuery(field string, message string) *errorResponse {
	return &errorResponse{
		Code:    "invalid_request",
		Message: message,
		Details: map[string]any{
			"field": field,
		},
	}
}
