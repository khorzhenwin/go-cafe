package discovery

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
	SourceProviderGeoapifyPlaces = "geoapify_places"
	defaultSearchLimit           = 12
	maxSearchLimit               = 20
	defaultDiscoveryLat          = 1.3521
	defaultDiscoveryLon          = 103.8198
	defaultDiscoveryRadiusMeters = 12000
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

type geoapifyFeatureCollection struct {
	Features []geoapifyFeature `json:"features"`
}

type geoapifyFeature struct {
	Properties geoapifyProperties `json:"properties"`
	Geometry   geoapifyGeometry   `json:"geometry"`
}

type geoapifyGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type geoapifyProperties struct {
	Name         string   `json:"name"`
	City         string   `json:"city"`
	District     string   `json:"district"`
	Suburb       string   `json:"suburb"`
	Quarter      string   `json:"quarter"`
	Formatted    string   `json:"formatted"`
	AddressLine1 string   `json:"address_line1"`
	AddressLine2 string   `json:"address_line2"`
	Lat          float64  `json:"lat"`
	Lon          float64  `json:"lon"`
	PlaceID      string   `json:"place_id"`
	Categories   []string `json:"categories"`
	Distance     int64    `json:"distance"`
	Catering     struct {
		Cuisine string `json:"cuisine"`
	} `json:"catering"`
}

type geoapifyGeocodeResponse struct {
	Results []geoapifyGeocodeResult `json:"results"`
}

type geoapifyGeocodeResult struct {
	PlaceID   string  `json:"place_id"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	City      string  `json:"city"`
	Formatted string  `json:"formatted"`
}

type searchCenter struct {
	Lat           float64
	Lon           float64
	RadiusMeters  int
	ResolvedLabel string
}

type GeoapifyPlacesClient struct {
	apiKey           string
	httpClient       *http.Client
	placesURL        string
	placeDetailsURL  string
	geocodeSearchURL string
}

func NewGeoapifyPlacesClient(apiKey string) *GeoapifyPlacesClient {
	return &GeoapifyPlacesClient{
		apiKey:           strings.TrimSpace(apiKey),
		placesURL:        "https://api.geoapify.com/v2/places",
		placeDetailsURL:  "https://api.geoapify.com/v2/place-details",
		geocodeSearchURL: "https://api.geoapify.com/v1/geocode/search",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func NewGeoapifyPlacesClientFromEnv() *GeoapifyPlacesClient {
	apiKey := strings.TrimSpace(os.Getenv("GEOAPIFY_API_KEY"))
	if apiKey == "" {
		return nil
	}
	return NewGeoapifyPlacesClient(apiKey)
}

func (c *GeoapifyPlacesClient) Search(ctx context.Context, filter SearchFilter) ([]Place, error) {
	if c == nil || c.apiKey == "" {
		return nil, fmt.Errorf("geoapify API key is not configured")
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = defaultSearchLimit
	}
	if limit > maxSearchLimit {
		limit = maxSearchLimit
	}

	center, err := c.resolveSearchCenter(ctx, filter.City)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("categories", "catering.cafe.coffee,catering.cafe.coffee_shop")
	params.Set("conditions", "named")
	params.Set("filter", fmt.Sprintf("circle:%f,%f,%d", center.Lon, center.Lat, center.RadiusMeters))
	params.Set("bias", fmt.Sprintf("proximity:%f,%f", center.Lon, center.Lat))
	params.Set("limit", strconv.Itoa(limit))
	params.Set("lang", "en")
	params.Set("apiKey", c.apiKey)

	if query := strings.TrimSpace(filter.Query); query != "" {
		params.Set("name", query)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.placesURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("geoapify places search failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var payload geoapifyFeatureCollection
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	places := make([]Place, 0, len(payload.Features))
	for _, feature := range payload.Features {
		place := normalizeGeoapifyFeature(feature)
		if place.Name == "" {
			continue
		}
		places = append(places, place)
	}

	return places, nil
}

func (c *GeoapifyPlacesClient) GetByID(ctx context.Context, placeID string) (*Place, error) {
	if c == nil || c.apiKey == "" {
		return nil, fmt.Errorf("geoapify API key is not configured")
	}

	cleanPlaceID := strings.TrimSpace(placeID)
	if cleanPlaceID == "" {
		return nil, fmt.Errorf("place ID is required")
	}

	params := url.Values{}
	params.Set("id", cleanPlaceID)
	params.Set("features", "details")
	params.Set("lang", "en")
	params.Set("apiKey", c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.placeDetailsURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

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
		return nil, fmt.Errorf("geoapify place details failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var payload geoapifyFeatureCollection
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if len(payload.Features) == 0 {
		return nil, nil
	}

	place := normalizeGeoapifyFeature(payload.Features[0])
	return &place, nil
}

func (c *GeoapifyPlacesClient) resolveSearchCenter(ctx context.Context, city string) (searchCenter, error) {
	trimmedCity := strings.TrimSpace(city)
	if trimmedCity == "" {
		return searchCenter{
			Lat:           defaultDiscoveryLat,
			Lon:           defaultDiscoveryLon,
			RadiusMeters:  defaultDiscoveryRadiusMeters,
			ResolvedLabel: "Singapore",
		}, nil
	}

	params := url.Values{}
	params.Set("text", trimmedCity)
	params.Set("type", "city")
	params.Set("format", "json")
	params.Set("limit", "1")
	params.Set("lang", "en")
	params.Set("apiKey", c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.geocodeSearchURL+"?"+params.Encode(), nil)
	if err != nil {
		return searchCenter{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return searchCenter{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return searchCenter{}, fmt.Errorf("geoapify geocode search failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var payload geoapifyGeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return searchCenter{}, err
	}
	if len(payload.Results) == 0 {
		return searchCenter{}, fmt.Errorf("could not resolve city %q in Geoapify", trimmedCity)
	}

	result := payload.Results[0]
	return searchCenter{
		Lat:           result.Lat,
		Lon:           result.Lon,
		RadiusMeters:  defaultDiscoveryRadiusMeters,
		ResolvedLabel: firstNonEmpty(result.City, result.Formatted, trimmedCity),
	}, nil
}

func normalizeGeoapifyFeature(feature geoapifyFeature) Place {
	lat, lon := resolveCoordinates(feature)
	description := buildGeoapifyDescription(feature.Properties)

	return Place{
		ID:              feature.Properties.PlaceID,
		ExternalPlaceID: feature.Properties.PlaceID,
		SourceProvider:  SourceProviderGeoapifyPlaces,
		Name:            firstNonEmpty(feature.Properties.Name, feature.Properties.AddressLine1, "Cafe"),
		Address:         firstNonEmpty(feature.Properties.Formatted, feature.Properties.AddressLine2, feature.Properties.AddressLine1),
		City:            strings.TrimSpace(feature.Properties.City),
		Neighborhood:    firstNonEmpty(feature.Properties.Suburb, feature.Properties.District, feature.Properties.Quarter),
		Description:     description,
		Latitude:        lat,
		Longitude:       lon,
		AvgRating:       0,
		ReviewCount:     0,
		VisitStatus:     "discover",
	}
}

func resolveCoordinates(feature geoapifyFeature) (*float64, *float64) {
	if feature.Properties.Lat != 0 || feature.Properties.Lon != 0 {
		lat := feature.Properties.Lat
		lon := feature.Properties.Lon
		return &lat, &lon
	}

	if len(feature.Geometry.Coordinates) >= 2 {
		lon := feature.Geometry.Coordinates[0]
		lat := feature.Geometry.Coordinates[1]
		return &lat, &lon
	}

	return nil, nil
}

func buildGeoapifyDescription(properties geoapifyProperties) string {
	if cuisine := strings.TrimSpace(properties.Catering.Cuisine); cuisine != "" {
		return fmt.Sprintf("Geoapify discovery: %s cafe.", strings.ReplaceAll(cuisine, "_", " "))
	}

	if len(properties.Categories) > 0 {
		return fmt.Sprintf("Geoapify discovery: %s.", strings.ReplaceAll(properties.Categories[len(properties.Categories)-1], "_", " "))
	}

	return "Geoapify discovery result."
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
