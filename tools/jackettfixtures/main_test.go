package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestClearFixtureDir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	jsonFile := filepath.Join(dir, "fixture.json")
	textFile := filepath.Join(dir, "keep.txt")
	subdir := filepath.Join(dir, "nested")

	if err := os.WriteFile(jsonFile, []byte("{}"), 0o644); err != nil {
		t.Fatalf("write json fixture: %v", err)
	}
	if err := os.WriteFile(textFile, []byte("keep"), 0o644); err != nil {
		t.Fatalf("write text fixture: %v", err)
	}
	if err := os.Mkdir(subdir, 0o755); err != nil {
		t.Fatalf("create nested dir: %v", err)
	}

	if err := clearFixtureDir(dir); err != nil {
		t.Fatalf("clearFixtureDir: %v", err)
	}

	if _, err := os.Stat(jsonFile); !os.IsNotExist(err) {
		t.Fatalf("expected json fixture to be removed, stat err = %v", err)
	}
	if _, err := os.Stat(textFile); err != nil {
		t.Fatalf("expected non-json file to remain: %v", err)
	}
	if _, err := os.Stat(subdir); err != nil {
		t.Fatalf("expected nested directory to remain: %v", err)
	}
}

func TestSanitizeResults(t *testing.T) {
	t.Parallel()

	link := "https://example.com/download"
	details := "https://example.com/details"
	magnet := "magnet:?xt=urn:btih:abc"

	got := sanitizeResults([]rawResult{
		{
			Title:        " Sample Title 2024 1080p WEB-DL ",
			Tracker:      " ExampleTracker ",
			TrackerID:    "",
			CategoryDesc: " TV ",
			Size:         json.RawMessage(`"12345"`),
			Seeders:      12,
			Peers:        18,
			PublishDate:  "2026-03-22T10:11:12.999999999Z",
			IMDb:         json.RawMessage(`"tt1234567"`),
			Link:         &link,
			Details:      &details,
			MagnetURI:    &magnet,
		},
		{
			Title:     "",
			Tracker:   "skip-empty-title",
			TrackerID: "skip-empty-title",
		},
		{
			Title:     "Skip Missing Tracker",
			Tracker:   "   ",
			TrackerID: "   ",
		},
	})

	if len(got) != 1 {
		t.Fatalf("sanitizeResults() len = %d, want 1", len(got))
	}

	want := fixtureResult{
		Title:       "Sample Title 2024 1080p WEB-DL",
		Tracker:     "ExampleTracker",
		TrackerID:   "exampletracker",
		Category:    "TV",
		Size:        12345,
		Seeders:     12,
		Peers:       18,
		PublishDate: "2026-03-22T10:11:12Z",
		IMDb:        "tt1234567",
		HasLink:     true,
		HasDetails:  true,
		HasMagnet:   true,
	}
	if got[0] != want {
		t.Fatalf("sanitizeResults()[0]\nwant: %#v\ngot:  %#v", want, got[0])
	}
}

func TestSanitizeIndexers(t *testing.T) {
	t.Parallel()

	errorText := " timed out "
	got := sanitizeIndexers([]rawIndexer{
		{ID: "b", Name: "B", Status: 1, Results: 2, ElapsedTime: 3},
		{ID: "a", Name: " A ", Status: 4, Results: 5, ElapsedTime: 6, Error: &errorText},
	})

	want := []fixtureIndexer{
		{ID: "a", Name: "A", Status: 4, Results: 5, ElapsedTime: 6, Error: "timed out"},
		{ID: "b", Name: "B", Status: 1, Results: 2, ElapsedTime: 3},
	}
	if len(got) != len(want) {
		t.Fatalf("sanitizeIndexers() len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("sanitizeIndexers()[%d]\nwant: %#v\ngot:  %#v", i, want[i], got[i])
		}
	}
}

func TestParseHelpers(t *testing.T) {
	t.Parallel()

	if got := parseSize(json.RawMessage(`"42.9"`)); got != 42 {
		t.Fatalf("parseSize() = %d, want 42", got)
	}
	if got := parseSize(json.RawMessage(`null`)); got != 0 {
		t.Fatalf("parseSize(null) = %d, want 0", got)
	}

	if got := parseIMDb(json.RawMessage(`"tt1234567"`)); got != "tt1234567" {
		t.Fatalf("parseIMDb(tt-prefixed) = %q, want %q", got, "tt1234567")
	}
	if got := parseIMDb(json.RawMessage(`1234567`)); got != "tt1234567" {
		t.Fatalf("parseIMDb(numeric) = %q, want %q", got, "tt1234567")
	}
	if got := parseIMDb(json.RawMessage(`"invalid"`)); got != "" {
		t.Fatalf("parseIMDb(invalid) = %q, want empty string", got)
	}

	if got := normalizeTimestamp("2026-03-22T10:11:12.999999999Z"); got != "2026-03-22T10:11:12Z" {
		t.Fatalf("normalizeTimestamp() = %q, want %q", got, "2026-03-22T10:11:12Z")
	}
	if got := normalizeTimestamp("not-a-time"); got != "" {
		t.Fatalf("normalizeTimestamp(invalid) = %q, want empty string", got)
	}
}

func TestFetchFixture(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("apikey"); got != "test-key" {
			t.Fatalf("apikey query = %q, want %q", got, "test-key")
		}
		if got := r.URL.Query().Get("configured"); got != "true" {
			t.Fatalf("configured query = %q, want %q", got, "true")
		}
		if got := r.URL.Query().Get("Query"); got != "safe query" {
			t.Fatalf("Query = %q, want %q", got, "safe query")
		}

		response := rawResponse{
			Results: []rawResult{
				{
					Title:        "Safe Sample 2026 1080p WEB-DL x264-GRP",
					Tracker:      "Tracker",
					TrackerID:    "tracker",
					CategoryDesc: "Movies",
					Size:         json.RawMessage(`"1000"`),
				},
			},
			Indexers: []rawIndexer{
				{ID: "idx", Name: "Indexer", Status: 1, Results: 1, ElapsedTime: 12},
			},
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	}))
	defer server.Close()

	fixture, err := fetchFixture(server.Client(), server.URL, "test-key", "safe query")
	if err != nil {
		t.Fatalf("fetchFixture() error = %v", err)
	}

	if fixture.Query != "safe query" {
		t.Fatalf("fixture.Query = %q, want %q", fixture.Query, "safe query")
	}
	if fixture.Source.BaseURL != server.URL {
		t.Fatalf("fixture.Source.BaseURL = %q, want %q", fixture.Source.BaseURL, server.URL)
	}
	if !fixture.Source.Configured {
		t.Fatalf("fixture.Source.Configured = false, want true")
	}
	if len(fixture.Results) != 1 {
		t.Fatalf("fixture.Results len = %d, want 1", len(fixture.Results))
	}
	if len(fixture.Indexers) != 1 {
		t.Fatalf("fixture.Indexers len = %d, want 1", len(fixture.Indexers))
	}
}
