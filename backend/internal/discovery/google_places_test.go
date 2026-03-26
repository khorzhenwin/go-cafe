package discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildTextQuery(t *testing.T) {
	assert.Equal(t, "coffee shops in Singapore", buildTextQuery("", "Singapore"))
	assert.Equal(t, "flat white in Singapore", buildTextQuery("flat white", "Singapore"))
	assert.Equal(t, "espresso", buildTextQuery("espresso", ""))
}

func TestNormalizeGooglePlace(t *testing.T) {
	place := normalizeGooglePlace(googlePlace{
		ID:                    "abc123",
		FormattedAddress:      "123 Main St, Seattle, WA",
		ShortFormattedAddress: "Capitol Hill",
		Rating:                4.7,
		UserRatingCount:       120,
		GoogleMapsUri:         "https://maps.google.com/example",
	})
	place.Name = "Daily Grind"

	assert.Equal(t, "abc123", place.ExternalPlaceID)
	assert.Equal(t, SourceProviderGooglePlaces, place.SourceProvider)
	assert.Equal(t, "discover", place.VisitStatus)
	assert.NotNil(t, place.Latitude)
	assert.NotNil(t, place.Longitude)
}
