package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/khorzhenwin/go-cafe/backend/internal/models"
	"gorm.io/gorm"
)

type Handler struct {
	Service *Service
}

func RegisterRoutes(r chi.Router, service *Service) {
	h := &Handler{Service: service}
	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.GetAllHandler)
		r.Post("/", h.CreateHandler)
		r.Get("/{id}", h.GetByIDHandler)
		r.Put("/{id}", h.UpdateHandler)
		r.Delete("/{id}", h.DeleteHandler)
	})
}

// GetAllHandler godoc
// @Summary List users
// @Description Returns all users.
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {string} string
// @Router /users/ [get]
func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.FindAll()
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(users)
}

// GetByIDHandler godoc
// @Summary Get user by ID
// @Description Returns a single user by ID.
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /users/{id} [get]
func (h *Handler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	user, err := h.Service.GetByID(uint(id))
	if err != nil {
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// CreateHandler godoc
// @Summary Create user
// @Description Creates a new user.
// @Tags users
// @Accept json
// @Produce json
// @Param body body CreateUserRequest true "Create user payload"
// @Success 201 {object} models.User
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /users/ [post]
func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}
	id, err := h.Service.CreateWithPassword(req.Email, req.Name, req.Password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	u, _ := h.Service.GetByID(id)
	w.WriteHeader(http.StatusCreated)
	if u != nil {
		_ = json.NewEncoder(w).Encode(u)
	} else {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	}
}

// UpdateHandler godoc
// @Summary Update user
// @Description Updates an existing user by ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body models.User true "Update user payload"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /users/{id} [put]
func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := h.Service.UpdateUser(uint(id), u); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "updated"})
}

// DeleteHandler godoc
// @Summary Delete user
// @Description Deletes a user by ID.
// @Tags users
// @Param id path int true "User ID"
// @Success 204 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /users/{id} [delete]
func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if err := h.Service.DeleteUser(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
