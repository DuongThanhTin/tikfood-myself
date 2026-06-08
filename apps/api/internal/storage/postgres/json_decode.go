package postgres

import (
	"encoding/json"
	"fmt"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
)

func decodeJSONField[T any](raw string, target *T, field string) error {
	if err := json.Unmarshal([]byte(raw), target); err != nil {
		return fmt.Errorf("decode %s: %w", field, err)
	}
	return nil
}

func decodeVenueListFields(venue *discovery.Venue, categoriesJSON string, dishesJSON string, socialVideosJSON string) error {
	if err := decodeJSONField(categoriesJSON, &venue.Categories, "venue categories"); err != nil {
		return err
	}
	if err := decodeJSONField(dishesJSON, &venue.TrendingDishes, "venue dishes"); err != nil {
		return err
	}
	if err := decodeJSONField(socialVideosJSON, &venue.SocialVideos, "venue social videos"); err != nil {
		return err
	}
	return nil
}

func decodeVenueDetailFields(
	venue *discovery.Venue,
	categoriesJSON string,
	socialVideosJSON string,
	trendingDishesJSON string,
	dishesJSON string,
	openingHoursJSON string,
) error {
	if err := decodeJSONField(categoriesJSON, &venue.Categories, "venue categories"); err != nil {
		return err
	}
	if err := decodeJSONField(socialVideosJSON, &venue.SocialVideos, "venue social videos"); err != nil {
		return err
	}
	if err := decodeJSONField(trendingDishesJSON, &venue.TrendingDishes, "venue trending dishes"); err != nil {
		return err
	}
	if err := decodeJSONField(dishesJSON, &venue.Dishes, "venue dishes"); err != nil {
		return err
	}
	if err := decodeJSONField(openingHoursJSON, &venue.OpeningHours, "venue opening hours"); err != nil {
		return err
	}
	return nil
}
