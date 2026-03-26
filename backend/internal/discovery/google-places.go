package discovery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	SourceProviderGooglePlaces = "google_places"
	defaultSearchLimit         = 12
	maxSearchLimit             = 20
)

type Place struct {
	ID              string   `json:"id"`
	ExternalPlaceID string   `json:"external_place_id"`
	SourceProvider  string   `json:"source_provider"`
	Name            string   `json:"name"`
	Address         string   `json:"address"`
	City            string   `json:"city,omitempty"`
	Neighborhood    string   `json:"neighborhood,omitempty"`
	Description     string   `json:"description,omitempty"`
	ImageURL        string   `json:"image_url,omitempty"`
	Latitude        *float64 `json:"latitude,omitempty"`
	Longitude       *float64 `json:"longitude,omitempty"`
	AvgRating       float64  `json:"avg_rating"`
	ReviewCount     int64    `json:"review_count"`
	GoogleMapsURL   string   `json:"google_maps_url,omitempty"`
	VisitStatus     string   `json:"visit_status,omitempty"`
}

type SearchFilter struct {
	Query string
	City  string
	Limit int
}

type Provider interface {
	Search(ctx context.Context, filter SearchFilter) ([]Place, error)
	GetByID(ctx context.Context, placeID string) (*Place, error)
}

type googlePlacesSearchRequest struct {
	TextQuery string `json:"textQuery"`
	PageSize  int    `json:"pageSize"`
}

type googlePlacesSearchResponse struct {
	Places []googlePlace `json:"places"`
}

type googlePlace struct {
	ID          string `json:"id"`
	DisplayName struct {
		Text string `json:"text"`
	} `json:"displayName"`
	FormattedAddress      string `json:"formattedAddress"`
	ShortFormattedAddress string `json:"shortFormattedAddress"`
	Location              struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Rating                 float64 `json:"rating"`
	UserRatingCount        int64   `json:"userRatingCount"`
	GoogleMapsUri          string  `json:"googleMapsUri"`
	PrimaryTypeDisplayName struct {
		Text string `json:"text"`
	} `json:"primaryTypeDisplayName"`
}

type GooglePlacesClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewGooglePlacesClient(apiKey string) *GooglePlacesClient {
	return &GooglePlacesClient{
		apiKey:  strings.TrimSpace(apiKey),
		baseURL: "https://places.googleapis.com/v1",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func NewGooglePlacesClientFromEnv() *GooglePlacesClient {
	apiKey := strings.TrimSpace(os.Getenv("GOOGLE_PLACES_API_KEY"))
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("GOOGLE_MAPS_API_KEY"))
	}
	if apiKey == "" {
		return nil
	}
	return NewGooglePlacesClient(apiKey)
}

func (c *GooglePlacesClient) Search(ctx context.Context, filter SearchFilter) ([]Place, error) {
	if c == nil || c.apiKey == "" {
		return nil, fmt.Errorf("google places API key is not configured")
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = defaultSearchLimit
	}
	if limit > maxSearchLimit {
		limit = maxSearchLimit
	}

	textQuery := buildTextQuery(filter.Query, filter.City)
	body, err := json.Marshal(googlePlacesSearchRequest{
		TextQuery: textQuery,
		PageSize:  limit,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/places:searchText", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", c.apiKey)
	req.Header.Set("X-Goog-FieldMask", "places.id,places.displayName,places.formattedAddress,places.shortFormattedAddress,places.location,places.rating,places.userRatingCount,places.googleMapsUri,places.primaryTypeDisplayName")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("google places search failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var payload googlePlacesSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	places := make([]Place, 0, len(payload.Places))
	for _, place := range payload.Places {
		places = append(places, normalizeGooglePlace(place))
	}

	return places, nil
}

func (c *GooglePlacesClient) GetByID(ctx context.Context, placeID string) (*Place, error) {
	if c == nil || c.apiKey == "" {
		return nil, fmt.Errorf("google places API key is not configured")
	}

	cleanPlaceID := strings.TrimSpace(placeID)
	if cleanPlaceID == "" {
		return nil, fmt.Errorf("place ID is required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/places/"+url.PathEscape(cleanPlaceID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Goog-Api-Key", c.apiKey)
	req.Header.Set("X-Goog-FieldMask", "id,displayName,formattedAddress,shortFormattedAddress,location,rating,userRatingCount,googleMapsUri,primaryTypeDisplayName")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("google place detail failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var place googlePlace
	if err := json.NewDecoder(resp.Body).Decode(&place); err != nil {
		return nil, err
	}

	normalized := normalizeGooglePlace(place)
	return &normalized, nil
}

func buildTextQuery(query, city string) string {
	trimmedQuery := strings.TrimSpace(query)
	trimmedCity := strings.TrimSpace(city)

	base := "coffee shops"
	if trimmedQuery != "" {
		base = trimmedQuery
	}
	if trimmedCity != "" {
		base += " in " + trimmedCity
	}
	return base
}

func normalizeGooglePlace(place googlePlace) Place {
	lat := place.Location.Latitude
	lng := place.Location.Longitude

	description := strings.TrimSpace(place.PrimaryTypeDisplayName.Text)
	if description == "" {
		description = "Discovered from Google Places."
	}

	return Place{
		ID:              place.ID,
		ExternalPlaceID: place.ID,
		SourceProvider:  SourceProviderGooglePlaces,
		Name:            strings.TrimSpace(place.DisplayName.Text),
		Address:         strings.TrimSpace(place.FormattedAddress),
		Neighborhood:    strings.TrimSpace(place.ShortFormattedAddress),
		Description:     description,
		Latitude:        &lat,
		Longitude:       &lng,
		AvgRating:       place.Rating,
		ReviewCount:     place.UserRatingCount,
		GoogleMapsURL:   strings.TrimSpace(place.GoogleMapsUri),
		VisitStatus:     "discover",
	}
}
