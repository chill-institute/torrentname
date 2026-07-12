package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

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
	}, 0)

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
		{ID: "a", Name: "A", Status: 4, Results: 5, ElapsedTime: 6, HasError: true},
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

func TestRefreshFixturesPreservesExistingCorpusOnFetchFailure(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	original := []byte("{\"original\":true}\n")
	path := filepath.Join(dir, "existing.json")
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatalf("write existing fixture: %v", err)
	}

	calls := 0
	fetch := func(*http.Client, string, string, string, int) (fixtureFile, error) {
		calls++
		if calls == 2 {
			return fixtureFile{}, errors.New("injected fetch failure")
		}
		return fixtureFile{Results: []fixtureResult{{Title: "Safe", TrackerID: "safe"}}}, nil
	}
	err := refreshFixtures(http.DefaultClient, "http://localhost:9117", "neutral-key", dir, []fixtureQuery{
		{slug: "first", query: "first"},
		{slug: "second", query: "second"},
	}, fetch)
	if err == nil {
		t.Fatal("refreshFixtures() error = nil, want injected failure")
	}
	got, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("read preserved fixture: %v", readErr)
	}
	if string(got) != string(original) {
		t.Fatalf("existing fixture changed after failed refresh: %q", got)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "first.json")); !os.IsNotExist(statErr) {
		t.Fatalf("partial fixture exists after failed refresh: %v", statErr)
	}
}

func TestRefreshFixturesReplacesJSONAndPreservesOtherEntries(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "old.json"), []byte("{}"), 0o600); err != nil {
		t.Fatalf("write old fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "keep.txt"), []byte("keep"), 0o600); err != nil {
		t.Fatalf("write preserved file: %v", err)
	}
	if err := os.Mkdir(filepath.Join(dir, "nested"), 0o700); err != nil {
		t.Fatalf("create nested directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "nested", "keep.txt"), []byte("nested"), 0o600); err != nil {
		t.Fatalf("write nested file: %v", err)
	}

	fetch := func(*http.Client, string, string, string, int) (fixtureFile, error) {
		return fixtureFile{
			Query:   "safe query",
			Source:  fixtureSource{Configured: true},
			Results: []fixtureResult{{Title: "Safe", TrackerID: "safe"}},
		}, nil
	}
	if err := refreshFixtures(http.DefaultClient, "http://localhost:9117", "neutral-key", dir, []fixtureQuery{{slug: "new", query: "safe"}}, fetch); err != nil {
		t.Fatalf("refreshFixtures() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "old.json")); !os.IsNotExist(err) {
		t.Fatalf("old fixture still exists: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "new.json")); err != nil {
		t.Fatalf("new fixture missing: %v", err)
	}
	if info, err := os.Stat(dir); err != nil {
		t.Fatalf("stat refreshed directory: %v", err)
	} else if got := info.Mode().Perm(); got != 0o700 {
		t.Fatalf("refreshed directory mode = %o, want 700", got)
	}
	for _, preserved := range []string{"keep.txt", filepath.Join("nested", "keep.txt")} {
		if _, err := os.Stat(filepath.Join(dir, preserved)); err != nil {
			t.Fatalf("preserved entry %s missing: %v", preserved, err)
		}
	}
}

func TestFetchFixtureRedactsAPIKeyFromTransportErrors(t *testing.T) {
	t.Parallel()

	const credential = "neutral-placeholder-key"
	client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		return nil, errors.New("transport failed for " + request.URL.String())
	})}
	_, err := fetchFixture(client, "http://localhost:9117", credential, "safe query", 1)
	if err == nil {
		t.Fatal("fetchFixture() error = nil, want transport error")
	}
	if strings.Contains(err.Error(), credential) {
		t.Fatalf("fetchFixture() error contains credential: %v", err)
	}
	if !strings.Contains(err.Error(), "localhost:9117") || !strings.Contains(err.Error(), "transport failed") {
		t.Fatalf("fetchFixture() error lost useful context: %v", err)
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

	fixture, err := fetchFixture(server.Client(), server.URL, "test-key", "safe query", 0)
	if err != nil {
		t.Fatalf("fetchFixture() error = %v", err)
	}

	if fixture.Query != "safe query" {
		t.Fatalf("fixture.Query = %q, want %q", fixture.Query, "safe query")
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
