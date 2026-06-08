package discovery

import (
	"context"
	"strings"
)

type VenueService struct {
	repo VenueRepository
}

func NewVenueService(repo VenueRepository) *VenueService {
	if repo == nil {
		panic("discovery.NewVenueService requires a venue repository")
	}
	return &VenueService{repo: repo}
}

func (service *VenueService) List(ctx context.Context, search VenueSearch) ([]Venue, error) {
	return service.repo.ListVenues(ctx, normalizeSearch(search))
}

func (service *VenueService) GetBySlug(ctx context.Context, slug string) (Venue, error) {
	return service.repo.GetVenueBySlug(ctx, strings.TrimSpace(slug))
}
