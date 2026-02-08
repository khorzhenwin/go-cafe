package auth

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/khorzhenwin/go-cafe/backend/internal/config"
)

// LoginFinder is implemented by user service for login (returns id and password hash for verification).
type LoginFinder interface {
	GetByEmailForAuth(email string) (id uint, passwordHash string, err error)
}

// RegisterCreator is implemented by user service for registration.
type RegisterCreator interface {
	CreateWithPassword(email, name, password string) (id uint, err error)
}

type Handler struct {
	AuthCfg *config.AuthConfig
	Finder  LoginFinder
	Creator RegisterCreator
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.LoginHandler)
		r.Post("/register", h.RegisterHandler)
	})
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}
	id, passwordHash, err := h.Finder.GetByEmailForAuth(req.Email)
	if err != nil || passwordHash == "" {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	if !CheckPassword(req.Password, passwordHash) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	token, expiresAt, err := h.issueToken(id)
	if err != nil {
		http.Error(w, "Failed to issue token", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(TokenResponse{Token: token, ExpiresAt: expiresAt.Format(time.RFC3339)})
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}
	id, err := h.Creator.CreateWithPassword(req.Email, req.Name, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, expiresAt, err := h.issueToken(id)
	if err != nil {
		http.Error(w, "Failed to issue token", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(TokenResponse{Token: token, ExpiresAt: expiresAt.Format(time.RFC3339)})
}

func (h *Handler) issueToken(userID uint) (string, time.Time, error) {
	expiresAt := time.Now().Add(h.AuthCfg.JWTExpiry)
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(userID), 10),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(h.AuthCfg.JWTSecret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}
