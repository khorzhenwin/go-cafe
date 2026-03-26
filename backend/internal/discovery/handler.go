package discovery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type RatingsFinder interface {
	GetByExternalPlaceID(externalPlaceID string) ([]map[string]any, error)
}

type Handler struct {
	Provider Provider
}

func RegisterRoutes(r chi.Router, provider Provider) {
	if provider == nil {
		provider = NewGooglePlacesClientFromEnv()
	}

	h := &Handler{Provider: provider}
	r.Route("/discovery/cafes", func(r chi.Router) {
		r.Get("/", h.ListHandler)
		r.Get("/{placeId}", h.GetByIDHandler)
	})
}

// ListHandler godoc
// @Summary Discover cafes from Google Places
// @Description Returns discovery cafes sourced from Google Places instead of the shared application database.
// @Tags discovery
// @Produce json
// @Param query query string false "Search query"
// @Param city query string false "City filter"
// @Param limit query int false "Result limit"
// @Success 200 {array} Place
// @Failure 500 {string} string
// @Failure 503 {string} string
// @Router /discovery/cafes/ [get]
func (h *Handler) ListHandler(w http.ResponseWriter, r *http.Request) {
	if h.Provider == nil {
		http.Error(w, "Google Places discovery is not configured", http.StatusServiceUnavailable)
		return
	}

	limit := defaultSearchLimit
	if limitStr := strings.TrimSpace(r.URL.Query().Get("limit")); limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		limit = parsed
	}

	places, err := h.Provider.Search(r.Context(), SearchFilter{
		Query: r.URL.Query().Get("query"),
		City:  r.URL.Query().Get("city"),
		Limit: limit,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(places)
}

// GetByIDHandler godoc
// @Summary Get discovery cafe by external place ID
// @Description Returns a Google Places cafe detail payload by place ID.
// @Tags discovery
// @Produce json
// @Param placeId path string true "External place ID"
// @Success 200 {object} Place
// @Failure 404 {string} string
// @Failure 503 {string} string
// @Router /discovery/cafes/{placeId} [get]
func (h *Handler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	if h.Provider == nil {
		http.Error(w, "Google Places discovery is not configured", http.StatusServiceUnavailable)
		return
	}

	placeID := strings.TrimSpace(chi.URLParam(r, "placeId"))
	if placeID == "" {
		http.Error(w, "Invalid place ID", http.StatusBadRequest)
		return
	}

	place, err := h.Provider.GetByID(r.Context(), placeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	if place == nil {
		http.Error(w, "Place not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(place)
}
