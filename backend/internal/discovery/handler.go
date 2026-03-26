package discovery

import (
	"encoding/json"
	"errors"
	"io"
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
		provider = NewGeoapifyPlacesClientFromEnv()
	}

	h := &Handler{Provider: provider}
	r.Route("/discovery/cafes", func(r chi.Router) {
		r.Get("/", h.ListHandler)
		r.Get("/static-map", h.StaticMapHandler)
		r.Get("/{placeId}", h.GetByIDHandler)
	})
}

// ListHandler godoc
// @Summary Discover cafes from Geoapify Places
// @Description Returns discovery cafes sourced from Geoapify Places instead of the shared application database.
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
		http.Error(w, "Geoapify discovery is not configured", http.StatusServiceUnavailable)
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
// @Description Returns a Geoapify cafe detail payload by place ID.
// @Tags discovery
// @Produce json
// @Param placeId path string true "External place ID"
// @Success 200 {object} Place
// @Failure 404 {string} string
// @Failure 503 {string} string
// @Router /discovery/cafes/{placeId} [get]
func (h *Handler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	if h.Provider == nil {
		http.Error(w, "Geoapify discovery is not configured", http.StatusServiceUnavailable)
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

// StaticMapHandler godoc
// @Summary Get Geoapify static map image
// @Description Returns a Geoapify Static Maps image for discovery coordinates without exposing the API key to the browser.
// @Tags discovery
// @Produce image/png
// @Param point query []string false "Map points in lat,lon format"
// @Param selected query string false "Selected point in lat,lon format"
// @Param width query int false "Image width"
// @Param height query int false "Image height"
// @Success 200 {file} binary
// @Failure 400 {string} string
// @Failure 502 {string} string
// @Failure 503 {string} string
// @Router /discovery/cafes/static-map [get]
func (h *Handler) StaticMapHandler(w http.ResponseWriter, r *http.Request) {
	staticMapClient := NewStaticMapClientFromEnv()
	if staticMapClient == nil {
		http.Error(w, "Geoapify static maps is not configured", http.StatusServiceUnavailable)
		return
	}

	width, err := parseStaticMapDimension(r.URL.Query().Get("width"), defaultStaticMapWidth)
	if err != nil {
		http.Error(w, "Invalid width", http.StatusBadRequest)
		return
	}

	height, err := parseStaticMapDimension(r.URL.Query().Get("height"), defaultStaticMapHeight)
	if err != nil {
		http.Error(w, "Invalid height", http.StatusBadRequest)
		return
	}

	points, err := parseStaticMapPoints(r.URL.Query()["point"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	selected, err := parseOptionalStaticMapPoint(r.URL.Query().Get("selected"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := staticMapClient.GetMap(r.Context(), points, selected, width, height)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, resp.Body)
}

func parseStaticMapDimension(raw string, fallback int) (int, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func parseStaticMapPoints(values []string) ([]StaticMapPoint, error) {
	points := make([]StaticMapPoint, 0, len(values))
	for _, value := range values {
		point, err := parseStaticMapPoint(value)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}
	return points, nil
}

func parseOptionalStaticMapPoint(raw string) (*StaticMapPoint, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil, nil
	}

	point, err := parseStaticMapPoint(value)
	if err != nil {
		return nil, err
	}
	return &point, nil
}

func parseStaticMapPoint(raw string) (StaticMapPoint, error) {
	parts := strings.Split(strings.TrimSpace(raw), ",")
	if len(parts) != 2 {
		return StaticMapPoint{}, errors.New("Invalid map point")
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return StaticMapPoint{}, errors.New("Invalid map point")
	}
	lon, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return StaticMapPoint{}, errors.New("Invalid map point")
	}
	return StaticMapPoint{Lat: lat, Lon: lon}, nil
}
