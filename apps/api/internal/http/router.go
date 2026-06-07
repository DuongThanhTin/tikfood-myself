package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
	"github.com/gin-gonic/gin"
)

type response struct {
	Data  any            `json:"data"`
	Error *errorResponse `json:"error,omitempty"`
}

type errorResponse struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func NewRouter(venueServices ...*discovery.VenueService) http.Handler {
	venues := discovery.NewVenueService()
	if len(venueServices) > 0 {
		venues = venueServices[0]
	}
	return NewRouterWithLogger(venues, slog.Default())
}

func NewRouterWithLogger(venues *discovery.VenueService, logger *slog.Logger) http.Handler {
	if venues == nil {
		venues = discovery.NewVenueService()
	}
	if logger == nil {
		logger = slog.Default()
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(recoveryMiddleware(logger), requestIDMiddleware(), requestLoggingMiddleware(logger))
	_ = router.SetTrustedProxies(nil)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, response{Data: map[string]bool{"ok": true}})
	})

	v1 := router.Group("/api/v1")
	v1.GET("/map/venues", venueSearchHandler(venues))
	v1.GET("/discovery/venues", venueSearchHandler(venues))
	v1.GET("/discovery/venues/:slug", venueDetailHandler(venues))

	return router
}

func venueSearchHandler(venues *discovery.VenueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		search, parseError := parseVenueSearch(c)
		if parseError != nil {
			c.JSON(http.StatusBadRequest, response{Data: nil, Error: parseError})
			return
		}

		result, err := venues.List(c.Request.Context(), search)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response{
				Data: nil,
				Error: &errorResponse{
					Code:    "internal_error",
					Message: "Failed to load venues.",
				},
			})
			return
		}
		c.JSON(http.StatusOK, response{Data: result})
	}
}

func venueDetailHandler(venues *discovery.VenueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := strings.TrimSpace(c.Param("slug"))
		if slug == "" || len(slug) > 160 {
			c.JSON(http.StatusBadRequest, response{
				Data: nil,
				Error: &errorResponse{
					Code:    "invalid_request",
					Message: "Venue slug is required and must be 160 characters or fewer.",
					Details: map[string]any{
						"field": "slug",
					},
				},
			})
			return
		}

		result, err := venues.GetBySlug(c.Request.Context(), slug)
		if errors.Is(err, discovery.ErrVenueNotFound) {
			c.JSON(http.StatusNotFound, response{
				Data: nil,
				Error: &errorResponse{
					Code:    "not_found",
					Message: "Venue was not found.",
				},
			})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, response{
				Data: nil,
				Error: &errorResponse{
					Code:    "internal_error",
					Message: "Failed to load venue.",
				},
			})
			return
		}

		c.JSON(http.StatusOK, response{Data: result})
	}
}
