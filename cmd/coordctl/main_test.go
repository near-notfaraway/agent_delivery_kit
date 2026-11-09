package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunClaimAndSummary(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "project.json")
	initial := `{
  "version": 1,
  "name": "demo",
  "tasks": [
    {"id": "T1", "repo": "api", "status": "pending", "depends_on": []},
    {"id": "T2", "repo": "service", "status": "pending", "depends_on": ["T1"]}
  ]
}`
	if err := os.WriteFile(projectPath, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if code := run([]string{"claim", "--project", projectPath, "--task", "T1", "--owner", "api-agent"}, &stdout, &stderr); code != 0 {
		t.Fatalf("claim exit = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "claimed T1") {
		t.Fatalf("claim stdout = %q", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"summary", "--project", projectPath}, &stdout, &stderr); code != 0 {
		t.Fatalf("summary exit = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"claimed": 1`) || !strings.Contains(stdout.String(), `"pending": 1`) {
		t.Fatalf("summary stdout = %q", stdout.String())
	}
}

func TestRunCompleteRequiresEvidenceFlags(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"complete", "--project", "project.json", "--task", "T1", "--owner", "agent"}, &stdout, &stderr)
	if code == 0 {
		t.Fatalf("complete exit = 0, want failure")
	}
	if !strings.Contains(stderr.String(), "commit, validation, and handoff are required") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunContextPrintsUpstreamHandoff(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "project.json")
	initial := `{
  "version": 1,
  "name": "demo",
  "tasks": [
    {"id": "T1", "repo": "api", "status": "completed", "depends_on": [], "result": {"commit": "v1", "validation": ["passed"], "handoff": "Item schema v1"}},
    {"id": "T2", "repo": "sdk", "status": "pending", "depends_on": ["T1"]}
  ]
}`
	if err := os.WriteFile(projectPath, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if code := run([]string{"context", "--project", projectPath, "--task", "T2"}, &stdout, &stderr); code != 0 {
		t.Fatalf("context exit = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Item schema v1") {
		t.Fatalf("context stdout = %q", stdout.String())
	}
}

func TestRunSkipAndDecision(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "project.json")
	initial := `{"version":1,"name":"demo","tasks":[{"id":"T1","repo":"sdk","status":"pending","depends_on":[]}]}`
	if err := os.WriteFile(projectPath, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if code := run([]string{"skip", "--project", projectPath, "--task", "T1", "--reason", "not required by v1"}, &stdout, &stderr); code != 0 {
		t.Fatalf("skip exit = %d, stderr = %q", code, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	args := []string{"decide", "--project", projectPath, "--id", "D1", "--task", "T1", "--original", "generate SDK", "--problem", "no public client", "--choice", "skip SDK", "--needs-approval"}
	if code := run(args, &stdout, &stderr); code != 0 {
		t.Fatalf("decide exit = %d, stderr = %q", code, stderr.String())
	}
	data, err := os.ReadFile(filepath.Join(dir, "decisions.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"needs_approval":true`) || !strings.Contains(string(data), `"choice":"skip SDK"`) {
		t.Fatalf("decision log = %s", data)
	}
}
