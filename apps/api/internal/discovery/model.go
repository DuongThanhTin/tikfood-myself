package discovery

import (
	"context"
	"errors"
)

var ErrVenueNotFound = errors.New("venue not found")

type Venue struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	Slug             string        `json:"slug"`
	ShortDescription string        `json:"short_description"`
	About            string        `json:"about"`
	Address          string        `json:"address"`
	City             string        `json:"city"`
	District         string        `json:"district"`
	Latitude         float64       `json:"latitude"`
	Longitude        float64       `json:"longitude"`
	Categories       []string      `json:"categories"`
	PriceLevel       int           `json:"price_level"`
	AvgPriceMinVND   int           `json:"avg_price_min_vnd"`
	AvgPriceMaxVND   int           `json:"avg_price_max_vnd"`
	Currency         string        `json:"currency"`
	SocialVideoCount int           `json:"social_video_count"`
	SocialVideos     []SocialVideo `json:"social_videos"`
	TrendScore       int           `json:"trend_score"`
	TrendingDishes   []string      `json:"trending_dishes"`
	Dishes           []VenueDish   `json:"dishes,omitempty"`
	OpeningHours     []OpeningHour `json:"opening_hours,omitempty"`
	AISummary        string        `json:"ai_summary"`
	DistanceMeters   *float64      `json:"distance_meters,omitempty"`
}

type SocialVideo struct {
	ID            string `json:"id"`
	Platform      string `json:"platform"`
	URL           string `json:"url"`
	CreatorHandle string `json:"creator_handle"`
	Caption       string `json:"caption"`
	ThumbnailURL  string `json:"thumbnail_url"`
	ViewCount     int64  `json:"view_count"`
	LikeCount     int64  `json:"like_count"`
	PublishedAt   string `json:"published_at,omitempty"`
}

type VenueDish struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Slug             string  `json:"slug"`
	ShortDescription string  `json:"short_description"`
	About            string  `json:"about"`
	Category         string  `json:"category"`
	Cuisine          string  `json:"cuisine"`
	PriceMinVND      int     `json:"price_min_vnd"`
	PriceMaxVND      int     `json:"price_max_vnd"`
	Currency         string  `json:"currency"`
	MentionCount     int     `json:"mention_count"`
	VideoCount       int     `json:"video_count"`
	ViewCount        int64   `json:"view_count"`
	TrendScore       float64 `json:"trend_score"`
}

type OpeningHour struct {
	DayOfWeek int    `json:"day_of_week"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
	IsClosed  bool   `json:"is_closed"`
}

type VenueSearch struct {
	Query       string
	City        string
	District    string
	Dish        string
	Tags        []string
	Platforms   []string
	Lat         *float64
	Lng         *float64
	RadiusM     int
	MinPriceVND int
	MaxPriceVND int
	OpenNow     bool
	Sort        VenueSort
	Limit       int
}

type VenueSort string

const (
	VenueSortTrending VenueSort = "trending"
	VenueSortVideos   VenueSort = "videos"
	VenueSortDistance VenueSort = "distance"
	VenueSortPrice    VenueSort = "price"
)

type VenueRepository interface {
	ListVenues(ctx context.Context, search VenueSearch) ([]Venue, error)
	GetVenueBySlug(ctx context.Context, slug string) (Venue, error)
}
