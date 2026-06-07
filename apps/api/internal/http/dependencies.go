package http

import "github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"

type HandlerDependencies struct {
	Venues *discovery.VenueService
}

func DefaultRouteRegistrars(deps HandlerDependencies) []RouteRegistrar {
	return []RouteRegistrar{
		NewVenueHandler(deps.Venues),
	}
}
