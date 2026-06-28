package overpass

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/ingest"
)

const userAgent = "tikfood-ingest/1.0 (OpenStreetMap venue ingestion)"

// Client queries the OpenStreetMap Overpass API.
type Client struct {
	endpoint string
	http     *http.Client
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		http:     &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *Client) Search(ctx context.Context, params ingest.SearchParams) ([]ingest.SourcePlace, error) {
	query := buildQuery(params)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint,
		strings.NewReader(url.Values{"data": {query}}.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build overpass request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("overpass request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read overpass response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("overpass status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return DecodeElements(data)
}

// buildQuery builds an Overpass QL query for the given amenities around a point.
func buildQuery(params ingest.SearchParams) string {
	amenities := params.Amenities
	if len(amenities) == 0 {
		amenities = []string{"restaurant"}
	}
	radius := params.RadiusM
	if radius <= 0 {
		radius = 1300
	}
	filter := strings.Join(amenities, "|")

	var b strings.Builder
	b.WriteString("[out:json][timeout:25];(")
	fmt.Fprintf(&b, `node["amenity"~"^(%s)$"](around:%d,%f,%f);`, filter, radius, params.Lat, params.Lng)
	b.WriteString(");out body;")
	return b.String()
}
