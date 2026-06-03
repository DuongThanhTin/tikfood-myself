package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
)

func TestHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	NewRouter().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
	if response.Header().Get("X-Request-ID") == "" {
		t.Fatal("expected request id header")
	}
}

func TestMapVenues(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/map/venues?district=District%201", nil)
	response := httptest.NewRecorder()

	NewRouter().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}

	var body struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if len(body.Data) != 1 {
		t.Fatalf("expected one venue, got %d", len(body.Data))
	}
	if body.Data[0]["slug"] == "" {
		t.Fatal("expected venue slug")
	}
	if body.Data[0]["about"] == "" {
		t.Fatal("expected venue about")
	}
}

func TestDiscoveryVenuesSearch(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/discovery/venues?q=pho&max_price_vnd=120000&tags=pho&limit=10", nil)
	response := httptest.NewRecorder()

	NewRouter().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}

	var body struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if len(body.Data) != 1 {
		t.Fatalf("expected one venue, got %d", len(body.Data))
	}
	if body.Data[0]["slug"] != "pho-bo-nguyen-le-van-sy-district-3" {
		t.Fatalf("expected pho venue, got %q", body.Data[0]["slug"])
	}
}

func TestDiscoveryVenuesInvalidSort(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/discovery/venues?sort=random", nil)
	response := httptest.NewRecorder()

	NewRouter().ServeHTTP(response, request)

	assertInvalidRequest(t, response, "sort")
}

func TestDiscoveryVenuesDistanceSortRequiresLocation(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/discovery/venues?sort=distance", nil)
	response := httptest.NewRecorder()

	NewRouter().ServeHTTP(response, request)

	assertInvalidRequest(t, response, "sort")
}

func TestMapVenuesStorageError(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/map/venues", nil)
	response := httptest.NewRecorder()

	venueService := discovery.NewVenueServiceWithRepository(failingVenueRepository{})
	NewRouter(venueService).ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", response.Code)
	}

	var body struct {
		Data  any `json:"data"`
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if body.Error.Code != "internal_error" {
		t.Fatalf("expected internal_error, got %q", body.Error.Code)
	}
}

type failingVenueRepository struct{}

func (failingVenueRepository) ListVenues(context.Context, discovery.VenueSearch) ([]discovery.Venue, error) {
	return nil, errors.New("storage failed")
}

func assertInvalidRequest(t *testing.T, response *httptest.ResponseRecorder, field string) {
	t.Helper()

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}

	var body struct {
		Data  any `json:"data"`
		Error struct {
			Code    string         `json:"code"`
			Message string         `json:"message"`
			Details map[string]any `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if body.Error.Code != "invalid_request" {
		t.Fatalf("expected invalid_request, got %q", body.Error.Code)
	}
	if body.Error.Details["field"] != field {
		t.Fatalf("expected field %q, got %q", field, body.Error.Details["field"])
	}
}
