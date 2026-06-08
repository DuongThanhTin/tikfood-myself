package discovery

import (
	"context"
	"fmt"
	"strings"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/platform/textutil"
)

type FallbackVenueRepository struct {
	venues []Venue
}

func NewFallbackVenueRepository() *FallbackVenueRepository {
	return &FallbackVenueRepository{venues: fallbackVenues()}
}

func (repo *FallbackVenueRepository) ListVenues(_ context.Context, search VenueSearch) ([]Venue, error) {
	results := make([]Venue, 0, len(repo.venues))

	for _, venue := range repo.venues {
		if search.Query != "" && !matchesQuery(venue, search.Query) {
			continue
		}
		if search.City != "" && !matchesCity(venue.City, search.City) {
			continue
		}
		if search.District != "" && !matchesDistrict(venue.District, search.District) {
			continue
		}
		if search.Dish != "" && !hasDish(venue.TrendingDishes, search.Dish) {
			continue
		}
		if search.MaxPriceVND > 0 && venue.AvgPriceMaxVND > search.MaxPriceVND {
			continue
		}
		if search.MinPriceVND > 0 && venue.AvgPriceMaxVND < search.MinPriceVND {
			continue
		}
		if len(search.Tags) > 0 && !hasAnyTag(venue.Categories, search.Tags) {
			continue
		}
		if len(search.Platforms) > 0 && !hasAnyPlatform(venue.SocialVideos, search.Platforms) {
			continue
		}
		results = append(results, venue)
	}

	if search.Limit > 0 && len(results) > search.Limit {
		results = results[:search.Limit]
	}

	return results, nil
}

func (repo *FallbackVenueRepository) GetVenueBySlug(_ context.Context, slug string) (Venue, error) {
	for _, venue := range repo.venues {
		if venue.Slug == slug {
			venue.Dishes = fallbackVenueDishes(venue)
			venue.OpeningHours = fallbackOpeningHours(venue)
			return venue, nil
		}
	}

	return Venue{}, ErrVenueNotFound
}

func fallbackVenueDishes(venue Venue) []VenueDish {
	dishes := make([]VenueDish, 0, len(venue.TrendingDishes))
	for index, dish := range venue.TrendingDishes {
		dishes = append(dishes, VenueDish{
			ID:               fmt.Sprintf("%s-dish-%d", venue.ID, index+1),
			Name:             dish,
			Slug:             textutil.Slugish(dish),
			ShortDescription: "Popular dish mentioned in social videos.",
			About:            "A venue-specific dish signal used by TikFood discovery.",
			Category:         firstCategory(venue.Categories),
			Cuisine:          "vietnamese",
			PriceMinVND:      venue.AvgPriceMinVND,
			PriceMaxVND:      venue.AvgPriceMaxVND,
			Currency:         venue.Currency,
			MentionCount:     1,
			VideoCount:       1,
			ViewCount:        0,
			TrendScore:       float64(venue.TrendScore),
		})
	}
	return dishes
}

func fallbackOpeningHours(venue Venue) []OpeningHour {
	openTime := "08:00"
	closeTime := "23:30"
	if strings.Contains(strings.ToLower(venue.Name), "pho") {
		openTime = "06:00"
		closeTime = "14:00"
	}

	hours := make([]OpeningHour, 0, 7)
	for day := 0; day <= 6; day++ {
		hours = append(hours, OpeningHour{
			DayOfWeek: day,
			OpenTime:  openTime,
			CloseTime: closeTime,
			IsClosed:  false,
		})
	}
	return hours
}

func firstCategory(categories []string) string {
	if len(categories) == 0 {
		return ""
	}
	return categories[0]
}
