package cafelisting

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultAutocompleteLimit = 5
	maxAutocompleteLimit     = 10
)

type AddressSuggestion struct {
	Name         string  `json:"name"`
	Formatted    string  `json:"formatted"`
	AddressLine1 string  `json:"address_line1"`
	AddressLine2 string  `json:"address_line2"`
	City         string  `json:"city"`
	Postcode     string  `json:"postcode"`
	Country      string  `json:"country"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
}

type AddressAutocompleteResponse struct {
	Results []AddressSuggestion `json:"results"`
}

type AddressAutocompleteProvider interface {
	Autocomplete(ctx context.Context, text string, limit int) ([]AddressSuggestion, error)
}

type geoapifyResponse struct {
	Results []AddressSuggestion `json:"results"`
}

type GeoapifyClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewGeoapifyClient(apiKey string) *GeoapifyClient {
	return &GeoapifyClient{
		apiKey:  strings.TrimSpace(apiKey),
		baseURL: "https://api.geoapify.com/v1/geocode/autocomplete",
		httpClient: &http.Client{
			Timeout: 8 * time.Second,
		},
	}
}

func NewGeoapifyClientFromEnv() *GeoapifyClient {
	apiKey := strings.TrimSpace(os.Getenv("GEOAPIFY_API_KEY"))
	if apiKey == "" {
		return nil
	}
	return NewGeoapifyClient(apiKey)
}

func (c *GeoapifyClient) Autocomplete(ctx context.Context, text string, limit int) ([]AddressSuggestion, error) {
	if c == nil || c.apiKey == "" {
		return nil, fmt.Errorf("geoapify API key is not configured")
	}
	query := strings.TrimSpace(text)
	if query == "" {
		return []AddressSuggestion{}, nil
	}
	if limit <= 0 {
		limit = defaultAutocompleteLimit
	}
	if limit > maxAutocompleteLimit {
		limit = maxAutocompleteLimit
	}

	params := url.Values{}
	params.Set("text", query)
	params.Set("format", "json")
	params.Set("limit", strconv.Itoa(limit))
	params.Set("apiKey", c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("geoapify request failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var payload geoapifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Results, nil
}
