package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunReportsFixtureMetrics(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	payload := `{
		"results": [
			{"title": "Sample Series S01E02 1080p WEB-DL x264-GRP"},
			{"title": "Movie 2024 HDRip XviD AC3-EVO"}
		]
	}`
	if err := os.WriteFile(filepath.Join(dir, "sample.json"), []byte(payload), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if err := run(dir); err != nil {
		t.Fatalf("run() error = %v", err)
	}
}

func TestRunRequiresFixtures(t *testing.T) {
	t.Parallel()

	err := run(t.TempDir())
	if err == nil {
		t.Fatalf("run() error = nil, want missing fixture error")
	}
}
