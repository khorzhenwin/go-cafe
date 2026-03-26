package discovery

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultStaticMapWidth  = 1200
	defaultStaticMapHeight = 720
	maxStaticMapDimension  = 2000
)

type StaticMapPoint struct {
	Lat float64
	Lon float64
}

type StaticMapClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewStaticMapClient(apiKey string) *StaticMapClient {
	return &StaticMapClient{
		apiKey:  strings.TrimSpace(apiKey),
		baseURL: "https://maps.geoapify.com/v1/staticmap",
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func NewStaticMapClientFromEnv() *StaticMapClient {
	apiKey := strings.TrimSpace(os.Getenv("GEOAPIFY_API_KEY"))
	if apiKey == "" {
		return nil
	}
	return NewStaticMapClient(apiKey)
}

func (c *StaticMapClient) GetMap(ctx context.Context, points []StaticMapPoint, selected *StaticMapPoint, width, height int) (*http.Response, error) {
	if c == nil || c.apiKey == "" {
		return nil, fmt.Errorf("geoapify static maps is not configured")
	}

	params := buildStaticMapParams(points, selected, width, height)
	params.Set("apiKey", c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("geoapify static map failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return resp, nil
}

func buildStaticMapParams(points []StaticMapPoint, selected *StaticMapPoint, width, height int) url.Values {
	if width <= 0 {
		width = defaultStaticMapWidth
	}
	if width > maxStaticMapDimension {
		width = maxStaticMapDimension
	}
	if height <= 0 {
		height = defaultStaticMapHeight
	}
	if height > maxStaticMapDimension {
		height = maxStaticMapDimension
	}

	params := url.Values{}
	params.Set("style", "osm-bright")
	params.Set("format", "png")
	params.Set("width", strconv.Itoa(width))
	params.Set("height", strconv.Itoa(height))
	params.Set("scaleFactor", "2")
	params.Set("lang", "en")

	if len(points) == 0 && selected != nil {
		points = append(points, *selected)
	}

	if len(points) <= 1 {
		point := StaticMapPoint{Lat: defaultDiscoveryLat, Lon: defaultDiscoveryLon}
		if len(points) == 1 {
			point = points[0]
		}
		params.Set("center", fmt.Sprintf("lonlat:%f,%f", point.Lon, point.Lat))
		params.Set("zoom", "15")
	} else {
		params.Set("area", buildStaticMapArea(points))
	}

	geometries := make([]string, 0, len(points)+1)
	for _, point := range points {
		geometries = append(geometries, buildStaticCircle(point, false))
	}
	if selected != nil {
		geometries = append(geometries, buildStaticCircle(*selected, true))
	}
	params.Set("geometry", strings.Join(geometries, "|"))

	return params
}

func buildStaticCircle(point StaticMapPoint, selected bool) string {
	fillColor := "#b7794e"
	lineColor := "#7c5234"
	radius := 11
	if selected {
		fillColor = "#2f6f62"
		lineColor = "#1f5147"
		radius = 16
	}

	return fmt.Sprintf(
		"circle:%f,%f,%d;linewidth:3;linecolor:%s;fillcolor:%s;lineopacity:0.95;fillopacity:0.82",
		point.Lon,
		point.Lat,
		radius,
		lineColor,
		fillColor,
	)
}

func buildStaticMapArea(points []StaticMapPoint) string {
	minLat := points[0].Lat
	maxLat := points[0].Lat
	minLon := points[0].Lon
	maxLon := points[0].Lon

	for _, point := range points[1:] {
		minLat = math.Min(minLat, point.Lat)
		maxLat = math.Max(maxLat, point.Lat)
		minLon = math.Min(minLon, point.Lon)
		maxLon = math.Max(maxLon, point.Lon)
	}

	latPadding := math.Max((maxLat-minLat)*0.18, 0.01)
	lonPadding := math.Max((maxLon-minLon)*0.18, 0.01)

	return fmt.Sprintf(
		"rect:%f,%f,%f,%f",
		minLon-lonPadding,
		maxLat+latPadding,
		maxLon+lonPadding,
		minLat-latPadding,
	)
}
