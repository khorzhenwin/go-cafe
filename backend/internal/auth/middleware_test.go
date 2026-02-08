package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/khorzhenwin/go-cafe/backend/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddleware_NoAuthHeader(t *testing.T) {
	cfg := &config.AuthConfig{JWTSecret: []byte("secret"), JWTExpiry: time.Hour}
	mw := Middleware(cfg)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	mw(next).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Authorization")
}

func TestMiddleware_InvalidPrefix(t *testing.T) {
	cfg := &config.AuthConfig{JWTSecret: []byte("secret"), JWTExpiry: time.Hour}
	mw := Middleware(cfg)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic xyz")
	rec := httptest.NewRecorder()
	mw(next).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestMiddleware_InvalidToken(t *testing.T) {
	cfg := &config.AuthConfig{JWTSecret: []byte("secret"), JWTExpiry: time.Hour}
	mw := Middleware(cfg)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid.jwt.here")
	rec := httptest.NewRecorder()
	mw(next).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestMiddleware_ValidToken(t *testing.T) {
	secret := []byte("test-secret")
	cfg := &config.AuthConfig{JWTSecret: secret, JWTExpiry: time.Hour}
	mw := Middleware(cfg)
	var capturedID uint
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := UserIDFromContext(r.Context())
		require.True(t, ok)
		capturedID = id
		w.WriteHeader(200)
	})
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   "7",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	tokenStr, err := token.SignedString(secret)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()
	mw(next).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, uint(7), capturedID)
}
