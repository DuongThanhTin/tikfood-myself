package overpass

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDecodeElements(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "elements.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	places, err := DecodeElements(data)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	// The third element has no name and must be skipped.
	if len(places) != 2 {
		t.Fatalf("expected 2 named places, got %d", len(places))
	}

	first := places[0]
	if first.ExternalID != "node/1770305118" {
		t.Fatalf("external id: got %q", first.ExternalID)
	}
	if first.Name != "Cà Phê Ciao" {
		t.Fatalf("name: got %q", first.Name)
	}
	if first.Lat != 10.7751001 || first.Lng != 106.7029693 {
		t.Fatalf("coords: got %v/%v", first.Lat, first.Lng)
	}
	if first.Address != "74-76 Nguyễn Huệ" {
		t.Fatalf("address: got %q", first.Address)
	}
	if first.Phone != "+8438231130" || first.Website != "https://ciaocafe.vn/" {
		t.Fatalf("contact: got %q / %q", first.Phone, first.Website)
	}
	if len(first.Cuisines) != 3 || first.Cuisines[0] != "vietnamese" {
		t.Fatalf("cuisines: got %v", first.Cuisines)
	}
	if first.OpeningHoursRaw != "Mo-Su 07:00-23:00" {
		t.Fatalf("hours raw: got %q", first.OpeningHoursRaw)
	}
	if first.Amenity != "cafe" {
		t.Fatalf("amenity: got %q", first.Amenity)
	}
}
