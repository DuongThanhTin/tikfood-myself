package http

import (
	"encoding/json"
	"net/http"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
)

type response struct {
	Data  any    `json:"data"`
	Error string `json:"error,omitempty"`
}

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	venues := discovery.NewVenueService()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, response{Data: map[string]bool{"ok": true}})
	})

	mux.HandleFunc("GET /api/v1/map/venues", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		result := venues.List(discovery.VenueSearch{
			District: query.Get("district"),
			Dish:     query.Get("dish"),
		})
		writeJSON(w, http.StatusOK, response{Data: result})
	})

	return mux
}

func writeJSON(w http.ResponseWriter, status int, body response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
