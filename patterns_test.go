package torrentname

import "testing"

func TestPatternsHaveExpectedCaptureGroups(t *testing.T) {
	t.Parallel()

	for _, pattern := range patterns {
		if got := pattern.re.NumSubexp(); got != 2 {
			t.Fatalf("pattern %q capture groups = %d, want 2", pattern.name, got)
		}
	}
}
