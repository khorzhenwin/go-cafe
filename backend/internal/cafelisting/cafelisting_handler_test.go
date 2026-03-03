package cafelisting

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockAutocompleteProvider struct {
	results []AddressSuggestion
	err     error
}

func (m mockAutocompleteProvider) Autocomplete(ctx context.Context, text string, limit int) ([]AddressSuggestion, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.results, nil
}

func TestAddressAutocompleteHandler_Success(t *testing.T) {
	h := &Handler{
		Autocomplete: mockAutocompleteProvider{
			results: []AddressSuggestion{
				{
					Name:      "Cafe ABC",
					Formatted: "123 Main St, Singapore",
					City:      "Singapore",
					Country:   "Singapore",
				},
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/cafes/autocomplete?text=cafe&limit=3", nil)
	rec := httptest.NewRecorder()
	h.AddressAutocompleteHandler(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var payload AddressAutocompleteResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload))
	require.Len(t, payload.Results, 1)
	assert.Equal(t, "Cafe ABC", payload.Results[0].Name)
	assert.Equal(t, "123 Main St, Singapore", payload.Results[0].Formatted)
}

func TestAddressAutocompleteHandler_MissingText(t *testing.T) {
	h := &Handler{Autocomplete: mockAutocompleteProvider{}}
	req := httptest.NewRequest(http.MethodGet, "/cafes/autocomplete?limit=3", nil)
	rec := httptest.NewRecorder()

	h.AddressAutocompleteHandler(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAddressAutocompleteHandler_InvalidLimit(t *testing.T) {
	h := &Handler{Autocomplete: mockAutocompleteProvider{}}
	req := httptest.NewRequest(http.MethodGet, "/cafes/autocomplete?text=cafe&limit=bad", nil)
	rec := httptest.NewRecorder()

	h.AddressAutocompleteHandler(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAddressAutocompleteHandler_NotConfigured(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodGet, "/cafes/autocomplete?text=cafe", nil)
	rec := httptest.NewRecorder()

	h.AddressAutocompleteHandler(rec, req)
	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
}

func TestAddressAutocompleteHandler_ProviderError(t *testing.T) {
	h := &Handler{
		Autocomplete: mockAutocompleteProvider{err: errors.New("upstream failed")},
	}
	req := httptest.NewRequest(http.MethodGet, "/cafes/autocomplete?text=cafe", nil)
	rec := httptest.NewRecorder()

	h.AddressAutocompleteHandler(rec, req)
	require.Equal(t, http.StatusBadGateway, rec.Code)
}
