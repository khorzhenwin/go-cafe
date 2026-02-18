package rating

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/khorzhenwin/go-cafe/backend/internal/auth"
	"github.com/khorzhenwin/go-cafe/backend/internal/models"
	"gorm.io/gorm"
)

type Handler struct {
	Service *Service
}

// RegisterRoutes registers rating routes. authMiddleware is required for create/update/delete and /me.
func RegisterRoutes(r chi.Router, service *Service, authMiddleware func(http.Handler) http.Handler) {
	h := &Handler{Service: service}
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/me/ratings", h.ListMyHandler)
	})
	r.Route("/cafes/{id}/ratings", func(r chi.Router) {
		r.Get("/", h.ListByCafeHandler)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Post("/", h.CreateHandler)
		})
	})
	r.Route("/users/{userId}/ratings", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/", h.ListByUserHandler)
	})
	r.Route("/ratings", func(r chi.Router) {
		r.Get("/{id}", h.GetByIDHandler)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Put("/{id}", h.UpdateHandler)
			r.Delete("/{id}", h.DeleteHandler)
		})
	})
}

// GetByIDHandler godoc
// @Summary Get rating by ID
// @Description Returns a rating by ID.
// @Tags ratings
// @Produce json
// @Param id path int true "Rating ID"
// @Success 200 {object} models.Rating
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /ratings/{id} [get]
func (h *Handler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	rating, err := h.Service.GetByID(uint(id))
	if err != nil {
		http.Error(w, "Failed to retrieve rating", http.StatusInternalServerError)
		return
	}
	if rating == nil {
		http.Error(w, "Rating not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rating)
}

// ListByCafeHandler godoc
// @Summary List ratings by cafe
// @Description Returns ratings for a cafe listing.
// @Tags ratings
// @Produce json
// @Param id path int true "Cafe ID"
// @Success 200 {array} models.Rating
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /cafes/{id}/ratings/ [get]
func (h *Handler) ListByCafeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	cafeID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid cafe ID", http.StatusBadRequest)
		return
	}
	ratings, err := h.Service.GetByCafeListingID(uint(cafeID))
	if err != nil {
		http.Error(w, "Failed to retrieve ratings", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ratings)
}

// ListMyHandler godoc
// @Summary List my ratings
// @Description Returns ratings created by the authenticated user.
// @Tags ratings
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Rating
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /me/ratings [get]
func (h *Handler) ListMyHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ratings, err := h.Service.GetByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve ratings", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ratings)
}

// ListByUserHandler godoc
// @Summary List ratings by user
// @Description Returns ratings for a user path parameter; must match authenticated user.
// @Tags ratings
// @Produce json
// @Security BearerAuth
// @Param userId path int true "User ID"
// @Success 200 {array} models.Rating
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Router /users/{userId}/ratings/ [get]
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
	ratings, err := h.Service.GetByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve ratings", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ratings)
}

// CreateHandler godoc
// @Summary Create rating
// @Description Creates a rating for a cafe listing.
// @Tags ratings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cafe ID"
// @Param body body models.Rating true "Create rating payload"
// @Success 201 {object} models.Rating
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /cafes/{id}/ratings/ [post]
func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	cafeIDStr := chi.URLParam(r, "id")
	cafeID, err := strconv.Atoi(cafeIDStr)
	if err != nil {
		http.Error(w, "Invalid cafe ID", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var rating models.Rating
	if err := json.NewDecoder(r.Body).Decode(&rating); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	rating.CafeListingID = uint(cafeID)
	rating.UserID = userID
	if err := h.Service.CreateRating(&rating); err != nil {
		if errors.Is(err, ErrCafeNotVisited) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Cafe listing not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to create rating", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(rating)
}

// UpdateHandler godoc
// @Summary Update rating
// @Description Updates a rating by ID.
// @Tags ratings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Rating ID"
// @Param body body models.Rating true "Update rating payload"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /ratings/{id} [put]
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
	var rating models.Rating
	if err := json.NewDecoder(r.Body).Decode(&rating); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := h.Service.UpdateRating(uint(id), userID, rating); err != nil {
		if errors.Is(err, ErrNotOwner) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Rating not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update rating", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "updated"})
}

// DeleteHandler godoc
// @Summary Delete rating
// @Description Deletes a rating by ID.
// @Tags ratings
// @Security BearerAuth
// @Param id path int true "Rating ID"
// @Success 204 {string} string
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /ratings/{id} [delete]
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
	if err := h.Service.DeleteRating(uint(id), userID); err != nil {
		if errors.Is(err, ErrNotOwner) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Rating not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete rating", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
