package workflowpolicy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const workflowsDir = "../../.github/workflows"

func workflows(t *testing.T) map[string]string {
	t.Helper()
	entries, err := os.ReadDir(workflowsDir)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yml") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(workflowsDir, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		out[e.Name()] = string(b)
	}
	if len(out) == 0 {
		t.Fatal("no workflows found")
	}
	return out
}

// zizmor (mise run actions) owns action pinning and dangerous triggers.
func TestOrgWorkflowInvariants(t *testing.T) {
	for name, w := range workflows(t) {
		t.Run(name, func(t *testing.T) {
			if strings.Contains(w, "\n  schedule:\n") {
				t.Fatal("declares a GitHub schedule; Cloudflare owns recurring dispatch")
			}
			if !strings.Contains(w, "\npermissions:") {
				t.Fatal("no workflow-level permissions block")
			}
			if strings.Count(w, "actions/checkout@") != strings.Count(w, "persist-credentials: false") {
				t.Fatal("every actions/checkout must set persist-credentials: false")
			}
		})
	}
}
