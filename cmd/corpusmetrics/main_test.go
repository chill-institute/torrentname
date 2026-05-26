package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunReportsFixtureMetrics(t *testing.T) {
	t.Parallel()

	dir := writeSampleFixture(t)

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

func TestRunAcceptsPassingThreshold(t *testing.T) {
	t.Parallel()

	dir := writeSampleFixture(t)
	if err := run(dir, fieldThreshold{name: "title", percent: 100}); err != nil {
		t.Fatalf("run() error = %v", err)
	}
}

func TestRunFailsBelowThreshold(t *testing.T) {
	t.Parallel()

	dir := writeSampleFixture(t)
	err := run(dir, fieldThreshold{name: "source", percent: 1})
	if err == nil {
		t.Fatalf("run() error = nil, want threshold error")
	}
	if !strings.Contains(err.Error(), "source coverage") {
		t.Fatalf("run() error = %v, want source threshold error", err)
	}
}

func TestThresholdFlagsValidateFields(t *testing.T) {
	t.Parallel()

	var flags thresholdFlags
	if err := flags.Set("title=99.5%"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if err := flags.Set("unknown=1"); err == nil {
		t.Fatalf("Set() error = nil, want unknown field error")
	}
}

func writeSampleFixture(t *testing.T) string {
	t.Helper()

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
	return dir
}
