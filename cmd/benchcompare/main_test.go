package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestComparisonRows(t *testing.T) {
	t.Parallel()

	baseline := benchmarkMetrics{
		"BenchmarkParse/tv_basic": {
			"ns/op":     100,
			"B/op":      50,
			"allocs/op": 5,
		},
	}
	current := benchmarkMetrics{
		"BenchmarkParse/tv_basic": {
			"ns/op":     90,
			"B/op":      55,
			"allocs/op": 5,
		},
	}

	rows := comparisonRows(baseline, current)
	if len(rows) != 3 {
		t.Fatalf("comparisonRows len = %d, want 3", len(rows))
	}
}

func TestFormatDeltaWithZeroBaseline(t *testing.T) {
	t.Parallel()

	if got := formatDelta(0, 1); got != "n/a" {
		t.Fatalf("formatDelta(0, 1) = %q, want n/a", got)
	}
	if got := formatDelta(0, 0); got != "+0.0%" {
		t.Fatalf("formatDelta(0, 0) = %q, want +0.0%%", got)
	}
}

func TestParseBenchmarkFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "bench.txt")
	payload := strings.Join([]string{
		"BenchmarkParse/tv_basic-10 12345 100.0 ns/op 50 B/op 5 allocs/op",
		"BenchmarkParse/tv_basic-10 12345 101.0 ns/op 52 B/op 5 allocs/op",
		"",
	}, "\n")
	if err := os.WriteFile(path, []byte(payload), 0o600); err != nil {
		t.Fatalf("write benchmark: %v", err)
	}

	metrics, err := parseBenchmarkFile(path)
	if err != nil {
		t.Fatalf("parseBenchmarkFile() error = %v", err)
	}
	got := metrics["BenchmarkParse/tv_basic"]["ns/op"]
	if got != 100.5 {
		t.Fatalf("ns/op = %v, want 100.5", got)
	}
}
