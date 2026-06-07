package discovery

import (
	"context"
	"strings"
)

type Venue struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Slug             string   `json:"slug"`
	ShortDescription string   `json:"short_description"`
	About            string   `json:"about"`
	Address          string   `json:"address"`
	City             string   `json:"city"`
	District         string   `json:"district"`
	Latitude         float64  `json:"latitude"`
	Longitude        float64  `json:"longitude"`
	Categories       []string `json:"categories"`
	PriceLevel       int      `json:"price_level"`
	AvgPriceMinVND   int      `json:"avg_price_min_vnd"`
	AvgPriceMaxVND   int      `json:"avg_price_max_vnd"`
	Currency         string   `json:"currency"`
	SocialVideoCount int      `json:"social_video_count"`
	TrendScore       int      `json:"trend_score"`
	TrendingDishes   []string `json:"trending_dishes"`
	AISummary        string   `json:"ai_summary"`
	DistanceMeters   *float64 `json:"distance_meters,omitempty"`
}

type VenueSearch struct {
	Query       string
	District    string
	Dish        string
	Tags        []string
	Lat         *float64
	Lng         *float64
	RadiusM     int
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
}

type VenueService struct {
	venues []Venue
	repo   VenueRepository
}

func NewVenueService() *VenueService {
	return &VenueService{
		venues: []Venue{
			{
				ID:               "11111111-1111-4111-8111-111111111111",
				Name:             "Banh Mi Hem",
				Slug:             "banh-mi-hem-nguyen-trai-district-1",
				ShortDescription: "Late-night banh mi spot trending on social video.",
				About:            "A compact street-food venue known for grilled pork banh mi, quick service, and strong late-night local buzz near Nguyen Trai.",
				Address:          "12 Nguyen Trai",
				City:             "Thành phố Hồ Chí Minh",
				District:         "Quận 1",
				Latitude:         10.7712,
				Longitude:        106.6899,
				Categories:       []string{"street_food", "banh_mi"},
				PriceLevel:       1,
				AvgPriceMinVND:   30000,
				AvgPriceMaxVND:   80000,
				Currency:         "VND",
				SocialVideoCount: 42,
				TrendScore:       92,
				TrendingDishes:   []string{"banh mi thit nuong", "banh mi pate"},
				AISummary:        "Trending for late-night banh mi clips with strong local social proof.",
			},
			{
				ID:               "22222222-2222-4222-8222-222222222222",
				Name:             "Pho Bo Nguyen",
				Slug:             "pho-bo-nguyen-le-van-sy-district-3",
				ShortDescription: "Breakfast pho shop with consistent creator mentions.",
				About:            "A neighborhood pho venue known for clear broth, beef toppings, and steady breakfast traffic from local regulars and food creators.",
				Address:          "88 Le Van Sy",
				City:             "Thành phố Hồ Chí Minh",
				District:         "Quận 3",
				Latitude:         10.7864,
				Longitude:        106.6767,
				Categories:       []string{"noodle", "pho"},
				PriceLevel:       1,
				AvgPriceMinVND:   50000,
				AvgPriceMaxVND:   120000,
				Currency:         "VND",
				SocialVideoCount: 35,
				TrendScore:       87,
				TrendingDishes:   []string{"pho bo tai", "pho bo vien"},
				AISummary:        "Popular for breakfast pho videos and consistent creator mentions.",
			},
		},
	}
}

func NewVenueServiceWithRepository(repo VenueRepository) *VenueService {
	service := NewVenueService()
	service.repo = repo
	return service
}

func (service *VenueService) List(ctx context.Context, search VenueSearch) ([]Venue, error) {
	search = normalizeSearch(search)

	if service.repo != nil {
		return service.repo.ListVenues(ctx, search)
	}

	results := make([]Venue, 0, len(service.venues))

	for _, venue := range service.venues {
		if search.Query != "" && !matchesQuery(venue, search.Query) {
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
		if len(search.Tags) > 0 && !hasAnyTag(venue.Categories, search.Tags) {
			continue
		}
		results = append(results, venue)
	}

	if search.Limit > 0 && len(results) > search.Limit {
		results = results[:search.Limit]
	}

	return results, nil
}

func normalizeSearch(search VenueSearch) VenueSearch {
	search.Query = strings.ToLower(strings.TrimSpace(search.Query))
	search.District = strings.TrimSpace(search.District)
	search.Dish = strings.ToLower(strings.TrimSpace(search.Dish))
	search.Tags = normalizeTags(search.Tags)

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

func normalizeTags(tags []string) []string {
	results := make([]string, 0, len(tags))
	seen := map[string]bool{}
	for _, tag := range tags {
		normalized := strings.ToLower(strings.TrimSpace(tag))
		normalized = strings.ReplaceAll(normalized, "_", "-")
		if normalized == "" || seen[normalized] {
			continue
		}
		seen[normalized] = true
		results = append(results, normalized)
	}
	return results
}

func matchesDistrict(value string, target string) bool {
	return normalizeDistrictAlias(value) == normalizeDistrictAlias(target)
}

func normalizeDistrictAlias(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.ReplaceAll(normalized, "quận", "quan")
	normalized = strings.ReplaceAll(normalized, ".", "")
	normalized = strings.ReplaceAll(normalized, " ", "")
	switch normalized {
	case "quan1", "district1", "q1":
		return "quan-1"
	case "quan3", "district3", "q3":
		return "quan-3"
	default:
		return normalized
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
