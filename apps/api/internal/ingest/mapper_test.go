package ingest

import (
	"testing"
)

func samplePlace() SourcePlace {
	return SourcePlace{
		ExternalID:      "node/1770305118",
		Name:            "Cà Phê Ciao",
		Address:         "74-76 Nguyễn Huệ",
		Lat:             10.7751001,
		Lng:             106.7029693,
		Phone:           "+8438231130",
		Website:         "https://ciaocafe.vn/",
		Amenity:         "cafe",
		Cuisines:        []string{"vietnamese", "international", "coffee_shop"},
		OpeningHoursRaw: "Mo-Su 07:00-23:00",
	}
}

func TestMapPlaceCoreFields(t *testing.T) {
	m := MapPlace(samplePlace())

	if m.Venue.Name != "Cà Phê Ciao" {
		t.Fatalf("name: got %q", m.Venue.Name)
	}
	if m.Venue.Slug != "ca-phe-ciao" { // textutil.Slugish handles Vietnamese
		t.Fatalf("slug: got %q", m.Venue.Slug)
	}
	if m.Venue.Latitude != 10.7751001 || m.Venue.Longitude != 106.7029693 {
		t.Fatalf("coords: got %v/%v", m.Venue.Latitude, m.Venue.Longitude)
	}
	if m.Venue.Address != "74-76 Nguyễn Huệ" {
		t.Fatalf("address: got %q", m.Venue.Address)
	}
	if m.Venue.Currency != "VND" {
		t.Fatalf("currency: got %q", m.Venue.Currency)
	}
	if m.Venue.Phone() != "+8438231130" || m.Venue.Website() != "https://ciaocafe.vn/" {
		t.Fatalf("contact: got %q / %q", m.Venue.Phone(), m.Venue.Website())
	}
}

func TestMapPlaceOpeningHours(t *testing.T) {
	m := MapPlace(samplePlace())
	if len(m.OpeningHours) != 7 {
		t.Fatalf("hours len: got %d", len(m.OpeningHours))
	}
}

func TestMapPlaceTagSlugs(t *testing.T) {
	m := MapPlace(samplePlace())

	// amenity "cafe" + cuisines; coffee_shop -> cafe (dedup); vietnamese stays.
	want := map[string]bool{"cafe": true, "vietnamese": true, "international": true}
	got := map[string]bool{}
	for _, s := range m.TagSlugs {
		got[s] = true
	}
	for slug := range want {
		if !got[slug] {
			t.Fatalf("missing tag %q in %v", slug, m.TagSlugs)
		}
	}
	if got["coffee_shop"] {
		t.Fatalf("coffee_shop should map to cafe, got %v", m.TagSlugs)
	}
}

func TestMapPlaceNoHoursWhenAbsent(t *testing.T) {
	p := samplePlace()
	p.OpeningHoursRaw = ""
	m := MapPlace(p)
	if m.OpeningHours != nil {
		t.Fatalf("expected no hours, got %+v", m.OpeningHours)
	}
}
