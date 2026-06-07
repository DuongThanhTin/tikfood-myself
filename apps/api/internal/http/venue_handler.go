package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
	"github.com/gin-gonic/gin"
)

type VenueHandler struct {
	venues *discovery.VenueService
}

func NewVenueHandler(venues *discovery.VenueService) *VenueHandler {
	if venues == nil {
		panic("http.NewVenueHandler requires a venue service")
	}
	return &VenueHandler{venues: venues}
}

func (handler *VenueHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	v1.GET("/map/venues", handler.Search)
	v1.GET("/discovery/venues", handler.Search)
	v1.GET("/discovery/venues/:slug", handler.Detail)
}

func (handler *VenueHandler) Search(c *gin.Context) {
	search, parseError := parseVenueSearch(c)
	if parseError != nil {
		respondWithBadRequest(c, parseError)
		return
	}

	result, err := handler.venues.List(c.Request.Context(), search)
	if err != nil {
		respondWithInternalServerError(c, MessageFailedLoadVenues)
		return
	}
	respondWithData(c, result)
}

func (handler *VenueHandler) Detail(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	if slug == "" || len(slug) > 160 {
		respondWithBadRequest(c, invalidQuery("slug", MessageVenueSlugInvalid))
		return
	}

	result, err := handler.venues.GetBySlug(c.Request.Context(), slug)
	if errors.Is(err, discovery.ErrVenueNotFound) {
		respondWithError(c, http.StatusNotFound, ErrorCodeNotFound, MessageVenueNotFound)
		return
	}
	if err != nil {
		respondWithInternalServerError(c, MessageFailedLoadVenue)
		return
	}

	respondWithData(c, result)
}
