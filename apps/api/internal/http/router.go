package http

import (
	"log/slog"
	"net/http"

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
	v1.GET("/map/venues", func(c *gin.Context) {
		result, err := venues.List(c.Request.Context(), discovery.VenueSearch{
			District: c.Query("district"),
			Dish:     c.Query("dish"),
		})
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
	})

	return router
}
