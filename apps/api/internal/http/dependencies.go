package http

import (
	"time"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/auth"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
)

type HandlerDependencies struct {
	Venues *discovery.VenueService

	// Auth is optional; when nil, the auth routes are not registered (e.g. tests that
	// exercise only discovery).
	Auth         *auth.AuthService
	Google       auth.GoogleAuthenticator // optional; nil disables Google login
	RefreshTTL   time.Duration
	CookieSecure bool
	WebOrigin    string
}

func DefaultRouteRegistrars(deps HandlerDependencies) []RouteRegistrar {
	registrars := []RouteRegistrar{
		NewVenueHandler(deps.Venues),
	}
	if deps.Auth != nil {
		registrars = append(registrars, NewAuthHandler(AuthHandlerConfig{
			Service:      deps.Auth,
			Google:       deps.Google,
			RefreshTTL:   deps.RefreshTTL,
			CookieSecure: deps.CookieSecure,
			WebOrigin:    deps.WebOrigin,
		}))
	}
	return registrars
}
