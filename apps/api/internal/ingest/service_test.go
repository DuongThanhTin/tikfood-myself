package ingest

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
)

type fakeSource struct {
	places []SourcePlace
	err    error
}

func (f fakeSource) Search(_ context.Context, _ SearchParams) ([]SourcePlace, error) {
	return f.places, f.err
}

type mockRepo struct {
	upserts    []UpsertVenueInput
	finished   bool
	finalState string
	finalCount Summary
}

func (m *mockRepo) ResolveOrCreateDistrict(_ context.Context, _, _ string) (string, error) {
	return "district-loc-1", nil
}
func (m *mockRepo) RecordIngestionRun(_ context.Context, _ IngestionRun) (string, error) {
	return "run-1", nil
}
func (m *mockRepo) UpsertVenueFromSource(_ context.Context, input UpsertVenueInput) error {
	m.upserts = append(m.upserts, input)
	return nil
}
func (m *mockRepo) FinishIngestionRun(_ context.Context, _ string, counts Summary, status string) error {
	m.finished = true
	m.finalCount = counts
	m.finalState = status
	return nil
}

func testLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

func baseOptions() IngestOptions {
	return IngestOptions{
		DistrictSlug: "quan-1", DistrictName: "Quận 1", Category: "restaurant",
		Lat: 10.7756, Lng: 106.7019, RadiusM: 1300, Amenities: []string{"restaurant"}, Limit: 20,
	}
}

func TestIngestDistrictUpsertsEachPlace(t *testing.T) {
	src := fakeSource{places: []SourcePlace{
		{ExternalID: "node/1", Name: "Alpha", Lat: 10, Lng: 106},
		{ExternalID: "node/2", Name: "Beta", Lat: 11, Lng: 107},
	}}
	repo := &mockRepo{}
	svc := NewService(src, repo, testLogger())

	summary, err := svc.IngestDistrict(context.Background(), baseOptions())
	if err != nil {
		t.Fatalf("ingest: %v", err)
	}
	if summary.Fetched != 2 || summary.Upserted != 2 || summary.Failed != 0 {
		t.Fatalf("summary: %+v", summary)
	}
	if len(repo.upserts) != 2 {
		t.Fatalf("expected 2 upserts, got %d", len(repo.upserts))
	}
	if repo.upserts[0].Source != "openstreetmap" || repo.upserts[0].ExternalID != "node/1" {
		t.Fatalf("upsert[0]: %+v", repo.upserts[0])
	}
	if repo.upserts[0].DistrictLocationID != "district-loc-1" {
		t.Fatalf("district id not set: %+v", repo.upserts[0])
	}
	if !repo.finished || repo.finalState != "completed" {
		t.Fatalf("run not finished completed: %+v", repo)
	}
}

func TestIngestDistrictRespectsLimit(t *testing.T) {
	src := fakeSource{places: []SourcePlace{
		{ExternalID: "node/1", Name: "Alpha"},
		{ExternalID: "node/2", Name: "Beta"},
		{ExternalID: "node/3", Name: "Gamma"},
	}}
	repo := &mockRepo{}
	svc := NewService(src, repo, testLogger())

	opts := baseOptions()
	opts.Limit = 2
	summary, err := svc.IngestDistrict(context.Background(), opts)
	if err != nil {
		t.Fatalf("ingest: %v", err)
	}
	if summary.Fetched != 2 || len(repo.upserts) != 2 {
		t.Fatalf("limit not applied: fetched=%d upserts=%d", summary.Fetched, len(repo.upserts))
	}
}

func TestIngestDistrictDryRunWritesNothing(t *testing.T) {
	src := fakeSource{places: []SourcePlace{{ExternalID: "node/1", Name: "Alpha"}}}
	repo := &mockRepo{}
	svc := NewService(src, repo, testLogger())

	opts := baseOptions()
	opts.DryRun = true
	summary, err := svc.IngestDistrict(context.Background(), opts)
	if err != nil {
		t.Fatalf("ingest: %v", err)
	}
	if len(repo.upserts) != 0 {
		t.Fatalf("dry run must not upsert, got %d", len(repo.upserts))
	}
	if summary.Fetched != 1 {
		t.Fatalf("summary fetched: %+v", summary)
	}
}

func TestIngestDistrictSourceErrorFailsRun(t *testing.T) {
	src := fakeSource{err: errors.New("overpass down")}
	repo := &mockRepo{}
	svc := NewService(src, repo, testLogger())

	_, err := svc.IngestDistrict(context.Background(), baseOptions())
	if err == nil {
		t.Fatal("expected error")
	}
	if !repo.finished || repo.finalState != "failed" {
		t.Fatalf("run should be marked failed: %+v", repo)
	}
}
