package discovery

type Venue struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Address        string   `json:"address"`
	District       string   `json:"district"`
	Latitude       float64  `json:"latitude"`
	Longitude      float64  `json:"longitude"`
	Categories     []string `json:"categories"`
	TrendScore     int      `json:"trend_score"`
	TrendingDishes []string `json:"trending_dishes"`
	AISummary      string   `json:"ai_summary"`
}

type VenueSearch struct {
	District string
	Dish     string
}

type VenueService struct {
	venues []Venue
}

func NewVenueService() *VenueService {
	return &VenueService{
		venues: []Venue{
			{
				ID:             "venue_001",
				Name:           "Banh Mi Hem",
				Address:        "12 Nguyen Trai",
				District:       "District 1",
				Latitude:       10.7712,
				Longitude:      106.6899,
				Categories:     []string{"street_food", "banh_mi"},
				TrendScore:     92,
				TrendingDishes: []string{"banh mi thit nuong", "banh mi pate"},
				AISummary:      "Trending for late-night banh mi clips with strong local social proof.",
			},
			{
				ID:             "venue_002",
				Name:           "Pho Bo Nguyen",
				Address:        "88 Le Van Sy",
				District:       "District 3",
				Latitude:       10.7864,
				Longitude:      106.6767,
				Categories:     []string{"noodle", "pho"},
				TrendScore:     87,
				TrendingDishes: []string{"pho bo tai", "pho bo vien"},
				AISummary:      "Popular for breakfast pho videos and consistent creator mentions.",
			},
		},
	}
}

func (service *VenueService) List(search VenueSearch) []Venue {
	results := make([]Venue, 0, len(service.venues))

	for _, venue := range service.venues {
		if search.District != "" && venue.District != search.District {
			continue
		}
		if search.Dish != "" && !hasDish(venue.TrendingDishes, search.Dish) {
			continue
		}
		results = append(results, venue)
	}

	return results
}

func hasDish(dishes []string, target string) bool {
	for _, dish := range dishes {
		if dish == target {
			return true
		}
	}
	return false
}
