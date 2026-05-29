package scanner

import (
	"docker-doctor/internal/model"

	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"docker-doctor/internal/rules"
)

type testRule struct {
	id     string
	delay  time.Duration
	issues []model.Issue
	err    error
	runs   *atomic.Int32
}

func (r testRule) ID() string              { return r.id }
func (r testRule) Category() string        { return "dockerfile" }
func (r testRule) SeverityDefault() string { return "warning" }
func (r testRule) Run(ctx context.Context, scanContext rules.ScanContext) ([]model.Issue, error) {
	if r.runs != nil {
		r.runs.Add(1)
	}
	if r.delay > 0 {
		time.Sleep(r.delay)
	}
	return r.issues, r.err
}

func TestOrchestratorRunsRulesAndSortsIssues(t *testing.T) {
	t.Parallel()

	reg := rules.NewRegistry()
	_ = reg.Register(testRule{id: "B", issues: []model.Issue{{ID: "RT_PORT_0001", Title: "b"}}})
	_ = reg.Register(testRule{id: "A", issues: []model.Issue{{ID: "DF_BASE_0001", Title: "a"}}})

	orch := NewOrchestrator(reg)
	issues, errs := orch.Run(context.Background(), rules.ScanContext{})

	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(issues) != 2 {
		t.Fatalf("unexpected issue count: %d", len(issues))
	}
	if issues[0].ID != "DF_BASE_0001" || issues[1].ID != "RT_PORT_0001" {
		t.Fatalf("issues not sorted deterministically: %+v", issues)
	}
}

func TestOrchestratorIsolatesRuleErrors(t *testing.T) {
	t.Parallel()

	reg := rules.NewRegistry()
	_ = reg.Register(testRule{id: "OK", issues: []model.Issue{{ID: "IMG_SIZE_0001", Title: "ok"}}})
	_ = reg.Register(testRule{id: "BAD", err: errors.New("boom")})

	orch := NewOrchestrator(reg)
	issues, errs := orch.Run(context.Background(), rules.ScanContext{})

	if len(issues) != 1 {
		t.Fatalf("expected successful rule issues preserved, got %d", len(issues))
	}
	if len(errs) != 1 || errs[0].RuleID != "BAD" {
		t.Fatalf("expected isolated BAD rule error, got %+v", errs)
	}
}

func TestOrchestratorRunsRulesConcurrently(t *testing.T) {
	t.Parallel()

	reg := rules.NewRegistry()
	runs := &atomic.Int32{}
	_ = reg.Register(testRule{id: "R1", delay: 120 * time.Millisecond, runs: runs})
	_ = reg.Register(testRule{id: "R2", delay: 120 * time.Millisecond, runs: runs})

	orch := NewOrchestrator(reg)
	start := time.Now()
	_, errs := orch.Run(context.Background(), rules.ScanContext{})
	elapsed := time.Since(start)

	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if runs.Load() != 2 {
		t.Fatalf("expected two rule runs, got %d", runs.Load())
	}
	if elapsed >= 220*time.Millisecond {
		t.Fatalf("expected concurrent execution, elapsed=%s", elapsed)
	}
}
