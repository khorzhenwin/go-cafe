package discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveCoordinatesUsesPropertiesFirst(t *testing.T) {
	feature := geoapifyFeature{
		Properties: geoapifyProperties{
			Lat: 1.3,
			Lon: 103.8,
		},
	}

	lat, lon := resolveCoordinates(feature)

	if assert.NotNil(t, lat) && assert.NotNil(t, lon) {
		assert.Equal(t, 1.3, *lat)
		assert.Equal(t, 103.8, *lon)
	}
}

func TestNormalizeGeoapifyFeature(t *testing.T) {
	place := normalizeGeoapifyFeature(geoapifyFeature{
		Properties: geoapifyProperties{
			Name:         "Daily Grind",
			PlaceID:      "abc123",
			Formatted:    "123 Main St, Singapore",
			City:         "Singapore",
			Suburb:       "Tiong Bahru",
			AddressLine1: "Daily Grind",
			AddressLine2: "123 Main St, Singapore",
			Categories:   []string{"catering.cafe", "catering.cafe.coffee_shop"},
			Catering: struct {
				Cuisine string `json:"cuisine"`
			}{
				Cuisine: "coffee_shop",
			},
		},
		Geometry: geoapifyGeometry{
			Type:        "Point",
			Coordinates: []float64{103.8, 1.3},
		},
	})

	assert.Equal(t, "abc123", place.ExternalPlaceID)
	assert.Equal(t, SourceProviderGeoapifyPlaces, place.SourceProvider)
	assert.Equal(t, "discover", place.VisitStatus)
	assert.Equal(t, "Tiong Bahru", place.Neighborhood)
	assert.Contains(t, place.Description, "coffee shop")
	assert.NotNil(t, place.Latitude)
	assert.NotNil(t, place.Longitude)
}

func TestFirstNonEmpty(t *testing.T) {
	assert.Equal(t, "Singapore", firstNonEmpty("", " ", "Singapore"))
	assert.Equal(t, "", firstNonEmpty("", " "))
}
