package discovery

import (
	"context"
	"errors"
	"strings"
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
				SocialVideos: []SocialVideo{
					{
						ID:            "33333333-3333-4333-8333-333333333331",
						Platform:      "tiktok",
						URL:           "https://www.tiktok.com/@tikfood/video/banh-mi-hem-1",
						CreatorHandle: "@tikfood",
						Caption:       "Late-night banh mi with grilled pork near Nguyen Trai.",
						ThumbnailURL:  "",
						ViewCount:     120000,
						LikeCount:     8200,
						PublishedAt:   "2026-05-20T10:00:00Z",
					},
					{
						ID:            "33333333-3333-4333-8333-333333333332",
						Platform:      "instagram",
						URL:           "https://www.instagram.com/reel/banh-mi-hem-2/",
						CreatorHandle: "@saigonbites",
						Caption:       "Crispy banh mi pate and grilled pork combo.",
						ThumbnailURL:  "",
						ViewCount:     80000,
						LikeCount:     5100,
						PublishedAt:   "2026-05-22T10:00:00Z",
					},
				},
				TrendScore:     92,
				TrendingDishes: []string{"banh mi thit nuong", "banh mi pate"},
				AISummary:      "Trending for late-night banh mi clips with strong local social proof.",
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
				SocialVideos: []SocialVideo{
					{
						ID:            "44444444-4444-4444-8444-444444444441",
						Platform:      "tiktok",
						URL:           "https://www.tiktok.com/@tikfood/video/pho-bo-nguyen-1",
						CreatorHandle: "@tikfood",
						Caption:       "Breakfast pho with clear broth and rare beef.",
						ThumbnailURL:  "",
						ViewCount:     95000,
						LikeCount:     6900,
						PublishedAt:   "2026-05-19T23:00:00Z",
					},
				},
				TrendScore:     87,
				TrendingDishes: []string{"pho bo tai", "pho bo vien"},
				AISummary:      "Popular for breakfast pho videos and consistent creator mentions.",
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

func (service *VenueService) GetBySlug(ctx context.Context, slug string) (Venue, error) {
	slug = strings.TrimSpace(slug)
	if service.repo != nil {
		return service.repo.GetVenueBySlug(ctx, slug)
	}

	for _, venue := range service.venues {
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
			ID:               venue.ID + "-dish-" + string(rune('1'+index)),
			Name:             dish,
			Slug:             strings.ReplaceAll(strings.ToLower(dish), " ", "-"),
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

func normalizeSearch(search VenueSearch) VenueSearch {
	search.Query = strings.ToLower(strings.TrimSpace(search.Query))
	search.City = strings.TrimSpace(search.City)
	search.District = strings.TrimSpace(search.District)
	search.Dish = strings.ToLower(strings.TrimSpace(search.Dish))
	search.Tags = normalizeTags(search.Tags)
	search.Platforms = normalizeTags(search.Platforms)

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

func matchesCity(value string, target string) bool {
	return normalizeCityAlias(value) == normalizeCityAlias(target)
}

func normalizeCityAlias(value string) string {
	normalized := normalizeVietnameseText(value)
	normalized = strings.ReplaceAll(normalized, ".", "")
	normalized = strings.ReplaceAll(normalized, " ", "")
	normalized = strings.ReplaceAll(normalized, "thanhpho", "tp")
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
	normalized := normalizeVietnameseText(value)
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

func normalizeVietnameseText(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	replacer := strings.NewReplacer(
		"à", "a", "á", "a", "ạ", "a", "ả", "a", "ã", "a", "â", "a", "ầ", "a", "ấ", "a", "ậ", "a", "ẩ", "a", "ẫ", "a", "ă", "a", "ằ", "a", "ắ", "a", "ặ", "a", "ẳ", "a", "ẵ", "a",
		"è", "e", "é", "e", "ẹ", "e", "ẻ", "e", "ẽ", "e", "ê", "e", "ề", "e", "ế", "e", "ệ", "e", "ể", "e", "ễ", "e",
		"ì", "i", "í", "i", "ị", "i", "ỉ", "i", "ĩ", "i",
		"ò", "o", "ó", "o", "ọ", "o", "ỏ", "o", "õ", "o", "ô", "o", "ồ", "o", "ố", "o", "ộ", "o", "ổ", "o", "ỗ", "o", "ơ", "o", "ờ", "o", "ớ", "o", "ợ", "o", "ở", "o", "ỡ", "o",
		"ù", "u", "ú", "u", "ụ", "u", "ủ", "u", "ũ", "u", "ư", "u", "ừ", "u", "ứ", "u", "ự", "u", "ử", "u", "ữ", "u",
		"ỳ", "y", "ý", "y", "ỵ", "y", "ỷ", "y", "ỹ", "y",
		"đ", "d",
	)
	return replacer.Replace(normalized)
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
