package cafelisting

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/khorzhenwin/go-cafe/backend/internal/auth"
	"github.com/khorzhenwin/go-cafe/backend/internal/models"
	"gorm.io/gorm"
)

type Handler struct {
	Service      *Service
	Autocomplete AddressAutocompleteProvider
}

// RegisterRoutes registers cafe listing routes. Pass authMiddleware for protected routes (required for create/update/delete and /me).
func RegisterRoutes(
	r chi.Router,
	service *Service,
	authMiddleware func(http.Handler) http.Handler,
	autocompleteProvider AddressAutocompleteProvider,
) {
	if autocompleteProvider == nil {
		autocompleteProvider = NewGeoapifyClientFromEnv()
	}
	h := &Handler{
		Service:      service,
		Autocomplete: autocompleteProvider,
	}
	// Authenticated "me" routes - user ID from JWT context
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/me/cafes", h.ListMyHandler)
		r.Post("/me/cafes", h.CreateMyHandler)
	})
	// Legacy user-scoped routes - require auth and path userId must match JWT
	r.Route("/users/{userId}/cafes", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/", h.ListByUserHandler)
		r.Post("/", h.CreateHandler)
	})
	r.Route("/cafes", func(r chi.Router) {
		r.Get("/", h.ListDiscoveryHandler)
		r.Get("/autocomplete", h.AddressAutocompleteHandler)
		r.Get("/{id}", h.GetByIDHandler)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Put("/{id}", h.UpdateHandler)
			r.Delete("/{id}", h.DeleteHandler)
		})
	})
}

// ListDiscoveryHandler godoc
// @Summary Discover cafes
// @Description Returns public community-submitted cafes for discovery surfaces.
// @Tags cafes
// @Produce json
// @Param query query string false "Search query"
// @Param city query string false "City filter"
// @Param sort query string false "Sort: rating_desc|newest|name_asc"
// @Param limit query int false "Result limit (1-60)"
// @Success 200 {array} models.CafeListing
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /cafes [get]
func (h *Handler) ListDiscoveryHandler(w http.ResponseWriter, r *http.Request) {
	limit := 18
	if limitStr := strings.TrimSpace(r.URL.Query().Get("limit")); limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		limit = parsed
	}

	listings, err := h.Service.ListDiscovery(
		r.URL.Query().Get("query"),
		r.URL.Query().Get("city"),
		r.URL.Query().Get("sort"),
		limit,
	)
	if err != nil {
		http.Error(w, "Failed to retrieve cafes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(listings)
}

// AddressAutocompleteHandler godoc
// @Summary Address autocomplete
// @Description Returns autocomplete suggestions from Geoapify.
// @Tags cafes
// @Produce json
// @Param text query string true "Partial address or place text"
// @Param limit query int false "Max suggestions (1-10)"
// @Success 200 {object} AddressAutocompleteResponse
// @Failure 400 {string} string
// @Failure 503 {string} string
// @Failure 502 {string} string
// @Router /cafes/autocomplete [get]
func (h *Handler) AddressAutocompleteHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("text"))
	if query == "" {
		http.Error(w, "Missing query parameter: text", http.StatusBadRequest)
		return
	}

	if h.Autocomplete == nil {
		http.Error(w, "Address autocomplete is not configured", http.StatusServiceUnavailable)
		return
	}

	limit := defaultAutocompleteLimit
	if limitStr := strings.TrimSpace(r.URL.Query().Get("limit")); limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		limit = parsed
	}

	results, err := h.Autocomplete.Autocomplete(r.Context(), query, limit)
	if err != nil {
		http.Error(w, "Failed to fetch address suggestions", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(AddressAutocompleteResponse{Results: results})
}

// GetByIDHandler godoc
// @Summary Get cafe by ID
// @Description Returns a cafe listing by ID.
// @Tags cafes
// @Produce json
// @Param id path int true "Cafe ID"
// @Success 200 {object} models.CafeListing
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /cafes/{id} [get]
func (h *Handler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	listing, err := h.Service.GetByID(uint(id))
	if err != nil {
		http.Error(w, "Failed to retrieve cafe listing", http.StatusInternalServerError)
		return
	}
	if listing == nil {
		http.Error(w, "Cafe listing not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(listing)
}

// ListMyHandler godoc
// @Summary List my cafes
// @Description Returns cafe listings owned by the authenticated user.
// @Tags cafes
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status (to_visit|visited)"
// @Param sort query string false "Sort: updated_desc|created_desc|name_asc|name_desc|status_asc|status_desc"
// @Success 200 {array} models.CafeListing
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /me/cafes [get]
func (h *Handler) ListMyHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	visitStatus := r.URL.Query().Get("status")
	sort := r.URL.Query().Get("sort")
	listings, err := h.Service.GetByUserIDFiltered(userID, visitStatus, sort)
	if err != nil {
		if errors.Is(err, ErrInvalidVisitStatus) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to retrieve cafe listings", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(listings)
}

// CreateMyHandler godoc
// @Summary Create my cafe
// @Description Creates a cafe listing for the authenticated user.
// @Tags cafes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.CafeListing true "Create cafe payload"
// @Success 201 {object} models.CafeListing
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /me/cafes [post]
func (h *Handler) CreateMyHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var listing models.CafeListing
	if err := json.NewDecoder(r.Body).Decode(&listing); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	listing.UserID = userID
	if err := h.Service.CreateListing(&listing); err != nil {
		if errors.Is(err, ErrInvalidVisitStatus) || errors.Is(err, ErrInvalidCafeName) || errors.Is(err, ErrInvalidCoordinates) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Source cafe not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to create cafe listing", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(listing)
}

// ListByUserHandler godoc
// @Summary List cafes by user
// @Description Returns cafes for a user path parameter; must match authenticated user.
// @Tags cafes
// @Produce json
// @Security BearerAuth
// @Param userId path int true "User ID"
// @Param status query string false "Filter by status (to_visit|visited)"
// @Param sort query string false "Sort: updated_desc|created_desc|name_asc|name_desc|status_asc|status_desc"
// @Success 200 {array} models.CafeListing
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /users/{userId}/cafes/ [get]
func (h *Handler) ListByUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	pathUserIDStr := chi.URLParam(r, "userId")
	pathUserID, err := strconv.Atoi(pathUserIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	if uint(pathUserID) != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	visitStatus := r.URL.Query().Get("status")
	sort := r.URL.Query().Get("sort")
	listings, err := h.Service.GetByUserIDFiltered(userID, visitStatus, sort)
	if err != nil {
		if errors.Is(err, ErrInvalidVisitStatus) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to retrieve cafe listings", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(listings)
}

// CreateHandler godoc
// @Summary Create cafe by user route
// @Description Creates a cafe via legacy user-scoped route; path userId must match authenticated user.
// @Tags cafes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path int true "User ID"
// @Param body body models.CafeListing true "Create cafe payload"
// @Success 201 {object} models.CafeListing
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /users/{userId}/cafes/ [post]
func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	pathUserIDStr := chi.URLParam(r, "userId")
	pathUserID, err := strconv.Atoi(pathUserIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	if uint(pathUserID) != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var listing models.CafeListing
	if err := json.NewDecoder(r.Body).Decode(&listing); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	listing.UserID = userID
	if err := h.Service.CreateListing(&listing); err != nil {
		if errors.Is(err, ErrInvalidVisitStatus) || errors.Is(err, ErrInvalidCafeName) || errors.Is(err, ErrInvalidCoordinates) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Source cafe not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to create cafe listing", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(listing)
}

// UpdateHandler godoc
// @Summary Update cafe
// @Description Updates a cafe listing by ID.
// @Tags cafes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cafe ID"
// @Param body body models.CafeListing true "Update cafe payload"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /cafes/{id} [put]
func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var listing models.CafeListing
	if err := json.NewDecoder(r.Body).Decode(&listing); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := h.Service.UpdateListing(uint(id), userID, listing); err != nil {
		if errors.Is(err, ErrNotOwner) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if errors.Is(err, ErrInvalidVisitStatus) || errors.Is(err, ErrInvalidCafeName) || errors.Is(err, ErrInvalidCoordinates) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Cafe listing not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update cafe listing", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "updated"})
}

// DeleteHandler godoc
// @Summary Delete cafe
// @Description Deletes a cafe listing by ID.
// @Tags cafes
// @Security BearerAuth
// @Param id path int true "Cafe ID"
// @Success 204 {string} string
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /cafes/{id} [delete]
func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if err := h.Service.DeleteListing(uint(id), userID); err != nil {
		if errors.Is(err, ErrNotOwner) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Cafe listing not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete cafe listing", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
