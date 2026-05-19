package torrentname

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type jackettFixture struct {
	Query   string `json:"query"`
	Results []struct {
		Title     string `json:"title"`
		TrackerID string `json:"tracker_id"`
	} `json:"results"`
}

func TestJackettFixtures(t *testing.T) {
	t.Parallel()
	testFixtureDir(t, filepath.Join("testdata", "jackett"), "Jackett")
}

func TestSyntheticFixtures(t *testing.T) {
	t.Parallel()
	testFixtureDir(t, filepath.Join("testdata", "synthetic"), "synthetic")
}

func testFixtureDir(t *testing.T, dir string, label string) {
	t.Helper()

	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		t.Fatalf("glob %s fixtures: %v", label, err)
	}
	if len(files) == 0 {
		t.Fatalf("expected at least one %s fixture file", label)
	}

	for _, path := range files {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			t.Parallel()

			payload, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read fixture: %v", err)
			}
			if strings.Contains(strings.ToLower(string(payload)), "jackett_apikey") || strings.Contains(strings.ToLower(string(payload)), "\"apikey\"") {
				t.Fatalf("fixture %s contains an API key field", path)
			}

			var fixture jackettFixture
			if err := json.Unmarshal(payload, &fixture); err != nil {
				t.Fatalf("decode fixture: %v", err)
			}
			if strings.TrimSpace(fixture.Query) == "" {
				t.Fatalf("fixture query is empty")
			}
			if len(fixture.Results) == 0 {
				t.Fatalf("fixture has no results")
			}

			for _, result := range fixture.Results {
				if strings.TrimSpace(result.Title) == "" {
					t.Fatalf("fixture contains an empty title")
				}
				if strings.TrimSpace(result.TrackerID) == "" {
					t.Fatalf("fixture contains an empty tracker_id")
				}

				_, err := Parse(result.Title)
				if err != nil {
					t.Fatalf("parse title %q: %v", result.Title, err)
				}
			}
		})
	}
}
