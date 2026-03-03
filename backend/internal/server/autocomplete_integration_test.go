package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/khorzhenwin/go-cafe/backend/internal/cafelisting"
	appconfig "github.com/khorzhenwin/go-cafe/backend/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockAddressAutocompleteProvider struct {
	results []cafelisting.AddressSuggestion
	err     error
}

func (m mockAddressAutocompleteProvider) Autocomplete(ctx context.Context, text string, limit int) ([]cafelisting.AddressSuggestion, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.results, nil
}

func testServerConfig() Config {
	return Config{
		BasePath:     "/api/v1",
		Address:      ":0",
		WriteTimeout: 0,
		ReadTimeout:  0,
	}
}

func testAuthConfig() *appconfig.AuthConfig {
	return &appconfig.AuthConfig{
		JWTSecret: []byte("test-secret"),
	}
}

func TestIntegration_AutocompleteRoute_Success(t *testing.T) {
	handler := NewWithDependencies(
		nil,
		testAuthConfig(),
		testServerConfig(),
		mockAddressAutocompleteProvider{
			results: []cafelisting.AddressSuggestion{
				{
					Name:      "Cafe Route Test",
					Formatted: "10 Test Street, Singapore",
				},
			},
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cafes/autocomplete?text=cafe&limit=5", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var payload cafelisting.AddressAutocompleteResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload))
	require.Len(t, payload.Results, 1)
	assert.Equal(t, "Cafe Route Test", payload.Results[0].Name)
}

func TestIntegration_AutocompleteRoute_ProviderFailure(t *testing.T) {
	handler := NewWithDependencies(
		nil,
		testAuthConfig(),
		testServerConfig(),
		mockAddressAutocompleteProvider{err: errors.New("upstream failure")},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cafes/autocomplete?text=cafe", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadGateway, rec.Code)
}

func TestIntegration_AutocompleteRoute_BadRequest(t *testing.T) {
	handler := NewWithDependencies(
		nil,
		testAuthConfig(),
		testServerConfig(),
		mockAddressAutocompleteProvider{},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cafes/autocomplete", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
