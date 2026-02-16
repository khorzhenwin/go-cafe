//go:build integration
// +build integration

package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	appconfig "github.com/khorzhenwin/go-cafe/backend/internal/config"
	"github.com/khorzhenwin/go-cafe/backend/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_RegisterLoginAndCafeFlow(t *testing.T) {
	_ = godotenv.Load()
	dbCfg, err := appconfig.LoadAWSConfig()
	if err != nil {
		t.Skipf("skip integration: DB not configured: %v", err)
		return
	}
	authCfg, err := appconfig.LoadAuthConfig()
	if err != nil {
		t.Skipf("skip integration: auth not configured: %v", err)
		return
	}

	conn, err := db.NewAWSClient(dbCfg)
	require.NoError(t, err)

	// Run migrations (use same DB; migrations are idempotent with IF NOT EXISTS)
	migrationsPath, _ := filepath.Abs("../../migrations")
	sourceURL := "file://" + filepath.ToSlash(migrationsPath)
	m, err := migrate.New(sourceURL, dbCfg.GetMigrationDSN())
	require.NoError(t, err)
	defer m.Close()
	_ = m.Up() // ignore ErrNoChange

	cfg := Config{
		BasePath:     "/api/v1",
		Address:      ":0",
		WriteTimeout: 0,
		ReadTimeout:  0,
	}
	handler := New(conn, authCfg, cfg)

	base := "/api/v1"
	testEmail := "inttest+" + strconv.FormatInt(time.Now().UnixNano(), 10) + "@example.com"

	// 1. Register
	regBody := map[string]string{"email": testEmail, "name": "Int Test", "password": "secret123"}
	regJSON, _ := json.Marshal(regBody)
	req := httptest.NewRequest(http.MethodPost, base+"/auth/register", bytes.NewReader(regJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code, "register: %s", rec.Body.String())
	var regResp struct {
		Token string `json:"token"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&regResp))
	require.NotEmpty(t, regResp.Token)

	token := regResp.Token

	// 2. Login
	loginBody := map[string]string{"email": testEmail, "password": "secret123"}
	loginJSON, _ := json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, base+"/auth/login", bytes.NewReader(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code, "login: %s", rec.Body.String())

	// 3. Create cafe (with token)
	cafeBody := map[string]string{"name": "Test Cafe", "address": "123 Main"}
	cafeJSON, _ := json.Marshal(cafeBody)
	req = httptest.NewRequest(http.MethodPost, base+"/me/cafes", bytes.NewReader(cafeJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code, "create cafe: %s", rec.Body.String())
	var cafeResp struct {
		ID uint `json:"id"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&cafeResp))
	require.NotZero(t, cafeResp.ID)

	// 4. List my cafes
	req = httptest.NewRequest(http.MethodGet, base+"/me/cafes", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func init() {
	// Try common working-directory locations first.
	_ = godotenv.Load()
	_ = godotenv.Load(".env")

	// Resolve backend root from this test file path so loading works
	// no matter where `go test` is executed from.
	if _, file, _, ok := runtime.Caller(0); ok {
		backendRoot := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
		_ = godotenv.Load(filepath.Join(backendRoot, ".env"))
	}

	if os.Getenv("DB_HOST") == "" {
		_ = godotenv.Load("backend/.env")
	}
}
