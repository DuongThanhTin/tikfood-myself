package ingest

import "context"

// SourcePlace is the source-agnostic representation of one ingested venue.
// Both the Overpass client (now) and a future Google client populate it.
type SourcePlace struct {
	ExternalID      string // "<type>/<id>" for OSM, e.g. "node/123"
	Name            string
	Address         string
	Lat             float64
	Lng             float64
	Phone           string
	Website         string
	Amenity         string
	Cuisines        []string
	OpeningHoursRaw string
}

// SearchParams describes a single bounded source query.
type SearchParams struct {
	Lat       float64
	Lng       float64
	RadiusM   int
	Amenities []string
	Limit     int
}

// VenueSource fetches places for one district query.
type VenueSource interface {
	Search(ctx context.Context, params SearchParams) ([]SourcePlace, error)
}
