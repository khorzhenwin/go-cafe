package discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildStaticMapParamsSinglePointUsesCenter(t *testing.T) {
	params := buildStaticMapParams([]StaticMapPoint{{Lat: 1.29, Lon: 103.85}}, nil, 800, 500)

	assert.Equal(t, "osm-bright", params.Get("style"))
	assert.Equal(t, "800", params.Get("width"))
	assert.Equal(t, "500", params.Get("height"))
	assert.Contains(t, params.Get("center"), "lonlat:103.850000,1.290000")
	assert.Equal(t, "15", params.Get("zoom"))
	assert.NotEmpty(t, params.Get("geometry"))
}

func TestBuildStaticMapParamsMultiplePointsUsesArea(t *testing.T) {
	points := []StaticMapPoint{
		{Lat: 1.29, Lon: 103.85},
		{Lat: 1.31, Lon: 103.89},
	}

	params := buildStaticMapParams(points, &points[0], 0, 0)

	assert.NotEmpty(t, params.Get("area"))
	assert.NotEmpty(t, params.Get("geometry"))
	assert.Equal(t, "1200", params.Get("width"))
	assert.Equal(t, "720", params.Get("height"))
}
