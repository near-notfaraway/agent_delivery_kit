package control

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
)

func TestClaimMarksPendingTaskClaimed(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "project.json")
	initial := `{
  "version": 1,
  "name": "demo",
  "tasks": [
    {"id": "T1", "repo": "api", "status": "pending", "depends_on": []}
  ]
}`
	if err := os.WriteFile(projectPath, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}

	store := NewStore(projectPath)
	if err := store.Claim("T1", "api-agent"); err != nil {
		t.Fatalf("Claim() error = %v", err)
	}

	project, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got := project.Tasks[0].Status; got != "claimed" {
		t.Fatalf("task status = %q, want claimed", got)
	}
	if got := project.Tasks[0].Owner; got != "api-agent" {
		t.Fatalf("task owner = %q, want api-agent", got)
	}
}

func TestClaimRejectsTaskAlreadyClaimed(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "project.json")
	initial := `{
  "version": 1,
  "name": "demo",
  "tasks": [
    {"id": "T1", "repo": "api", "status": "claimed", "depends_on": [], "owner": "first-agent"}
  ]
}`
	if err := os.WriteFile(projectPath, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}

	store := NewStore(projectPath)
	if err := store.Claim("T1", "second-agent"); err == nil {
		t.Fatal("Claim() error = nil, want already claimed error")
	}

	project, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got := project.Tasks[0].Owner; got != "first-agent" {
		t.Fatalf("task owner = %q, want first-agent", got)
	}
}

func TestClaimRejectsTaskWithIncompleteDependency(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "project.json")
	initial := `{
  "version": 1,
  "name": "demo",
  "tasks": [
    {"id": "T1", "repo": "api", "status": "in_progress", "depends_on": [], "owner": "api-agent"},
    {
      "id": "T2",
      "repo": "service",
      "title": "Consume item contract",
      "objective": "Map the API item into the service model",
      "acceptance": ["service tests pass"],
      "handoff_requirements": ["publish service revision"],
      "status": "pending",
      "depends_on": ["T1"]
    }
  ]
}`
	if err := os.WriteFile(projectPath, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}

	store := NewStore(projectPath)
	if err := store.Claim("T2", "service-agent"); err == nil {
		t.Fatal("Claim() error = nil, want incomplete dependency error")
	}
}

func TestCompleteRejectsMissingEvidence(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "project.json")
	initial := `{
  "version": 1,
  "name": "demo",
  "tasks": [
    {"id": "T1", "repo": "api", "status": "claimed", "depends_on": [], "owner": "api-agent"}
  ]
}`
	if err := os.WriteFile(projectPath, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}

	store := NewStore(projectPath)
	if err := store.Complete("T1", "api-agent", Completion{}); err == nil {
		t.Fatal("Complete() error = nil, want missing evidence error")
	}
}

func TestCompleteStoresEvidenceAndMarksTaskCompleted(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "project.json")
	initial := `{
  "version": 1,
  "name": "demo",
  "tasks": [
    {"id": "T1", "repo": "api", "status": "claimed", "depends_on": [], "owner": "api-agent"}
  ]
}`
	if err := os.WriteFile(projectPath, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}

	store := NewStore(projectPath)
	evidence := Completion{
		Commit:     "abc123",
		Validation: []string{"go test ./..."},
		Handoff:    "API contract v2 is ready",
	}
	if err := store.Complete("T1", "api-agent", evidence); err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	project, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	task := project.Tasks[0]
	if task.Status != "completed" {
		t.Fatalf("task status = %q, want completed", task.Status)
	}
	if task.Result == nil || task.Result.Commit != "abc123" {
		t.Fatalf("task result = %#v, want commit abc123", task.Result)
	}
}

func TestSkipRequiresReasonAndRecordsResolution(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "project.json")
	initial := `{
  "version": 1,
  "name": "demo",
  "tasks": [
    {"id": "T3", "repo": "sdk", "status": "pending", "depends_on": []}
  ]
}`
	if err := os.WriteFile(projectPath, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}

	store := NewStore(projectPath)
	if err := store.Skip("T3", ""); err == nil {
		t.Fatal("Skip() error = nil, want missing reason error")
	}
	if err := store.Skip("T3", "existing contract is sufficient"); err != nil {
		t.Fatalf("Skip() error = %v", err)
	}

	project, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	task := project.Tasks[0]
	if task.Status != "skipped" || task.Resolution != "existing contract is sufficient" {
		t.Fatalf("task = %#v, want skipped with reason", task)
	}
}

func TestAddDecisionAppendsDurableRecord(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "project.json")
	if err := os.WriteFile(projectPath, []byte(`{"version":1,"name":"demo","tasks":[]}`), 0o600); err != nil {
		t.Fatal(err)
	}

	store := NewStore(projectPath)
	decision := Decision{
		ID:            "D1",
		TaskID:        "T2",
		Original:      "poll until ready",
		Problem:       "polling loses state after restart",
		Choice:        "persist readiness in workflow state",
		NeedsApproval: true,
	}
	if err := store.AddDecision(decision); err != nil {
		t.Fatalf("AddDecision() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "decisions.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	var got Decision
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decision file is not JSON: %v", err)
	}
	if got.ID != "D1" || got.Choice != decision.Choice {
		t.Fatalf("decision = %#v, want %#v", got, decision)
	}
}

func TestSummaryReportsCountsAndReadyTasks(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "project.json")
	initial := `{
  "version": 1,
  "name": "demo",
  "tasks": [
    {"id": "T1", "repo": "api", "status": "completed", "depends_on": []},
    {"id": "T2", "repo": "service", "status": "pending", "depends_on": ["T1"]},
    {"id": "T3", "repo": "sdk", "status": "pending", "depends_on": ["T2"]},
    {"id": "T4", "repo": "web", "status": "blocked", "depends_on": []}
  ]
}`
	if err := os.WriteFile(projectPath, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}

	summary, err := NewStore(projectPath).Summary()
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if summary.Counts["pending"] != 2 || summary.Counts["completed"] != 1 || summary.Counts["blocked"] != 1 {
		t.Fatalf("counts = %#v", summary.Counts)
	}
	if len(summary.Ready) != 1 || summary.Ready[0] != "T2" {
		t.Fatalf("ready = %#v, want [T2]", summary.Ready)
	}
}

func TestSummaryUsesEmptyReadyListWhenNothingCanStart(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "project.json")
	initial := `{"version":1,"name":"done","tasks":[{"id":"T1","repo":"api","status":"completed","depends_on":[]}]}`
	if err := os.WriteFile(projectPath, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}

	summary, err := NewStore(projectPath).Summary()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"counts":{"completed":1},"ready":[]}` {
		t.Fatalf("summary JSON = %s", data)
	}
}

func TestConcurrentClaimAllowsOnlyOneOwner(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "project.json")
	initial := `{
  "version": 1,
  "name": "demo",
  "tasks": [
    {"id": "T1", "repo": "api", "status": "pending", "depends_on": []}
  ]
}`
	if err := os.WriteFile(projectPath, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}

	const contenders = 20
	start := make(chan struct{})
	var wg sync.WaitGroup
	var successes atomic.Int32
	for i := 0; i < contenders; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			if err := NewStore(projectPath).Claim("T1", string(rune('A'+i))); err == nil {
				successes.Add(1)
			}
		}(i)
	}
	close(start)
	wg.Wait()

	if got := successes.Load(); got != 1 {
		t.Fatalf("successful claims = %d, want 1", got)
	}
}

func TestContextIncludesDependencyHandoff(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "project.json")
	initial := `{
  "version": 1,
  "name": "demo",
  "tasks": [
    {
      "id": "T1",
      "repo": "api",
      "status": "completed",
      "depends_on": [],
      "owner": "api-agent",
      "result": {"commit": "api-v1", "validation": ["contract test passed"], "handoff": "GET /items returns Item[]"}
    },
    {
      "id": "T2",
      "repo": "service",
      "title": "Consume item contract",
      "objective": "Map the API item into the service model",
      "acceptance": ["service tests pass"],
      "handoff_requirements": ["publish service revision"],
      "status": "pending",
      "depends_on": ["T1"]
    }
  ]
}`
	if err := os.WriteFile(projectPath, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}

	context, err := NewStore(projectPath).Context("T2")
	if err != nil {
		t.Fatalf("Context() error = %v", err)
	}
	if context.Task.ID != "T2" || len(context.Dependencies) != 1 {
		t.Fatalf("context = %#v", context)
	}
	if context.Task.Objective == "" || len(context.Task.Acceptance) != 1 || len(context.Task.HandoffRequirements) != 1 {
		t.Fatalf("task instructions = %#v", context.Task)
	}
	dependency := context.Dependencies[0]
	if dependency.ID != "T1" || dependency.Result == nil || dependency.Result.Handoff != "GET /items returns Item[]" {
		t.Fatalf("dependency = %#v", dependency)
	}
}
