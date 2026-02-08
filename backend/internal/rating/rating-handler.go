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
		http.Error(w, "Failed to create rating", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(rating)
}

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
