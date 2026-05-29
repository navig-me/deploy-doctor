package scoring

import (
	"testing"

	"docker-doctor/internal/model"
)

func TestCalculateScore(t *testing.T) {
	t.Parallel()

	issues := []model.Issue{
		{Severity: model.SeverityCritical},
		{Severity: model.SeverityWarning},
		{Severity: model.SeveritySuggestion},
	}

	got := CalculateScore(issues)
	if got != 57 {
		t.Fatalf("unexpected score: got %d want %d", got, 57)
	}
}

func TestCalculateScoreClampedToZero(t *testing.T) {
	t.Parallel()

	issues := make([]model.Issue, 5)
	for i := range issues {
		issues[i] = model.Issue{Severity: model.SeverityCritical}
	}

	if got := CalculateScore(issues); got != 0 {
		t.Fatalf("score should clamp to zero, got %d", got)
	}
}

func TestMapStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		score   int
		summary model.Summary
		want    string
	}{
		{name: "fail due to critical", score: 95, summary: model.Summary{Critical: 1}, want: "fail"},
		{name: "fail due to low score", score: 59, summary: model.Summary{}, want: "fail"},
		{name: "pass", score: 85, summary: model.Summary{Warning: 1}, want: "pass"},
		{name: "risky due to warnings", score: 85, summary: model.Summary{Warning: 3}, want: "risky"},
		{name: "risky due to mid score", score: 75, summary: model.Summary{Warning: 1}, want: "risky"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := MapStatus(tt.score, tt.summary); got != tt.want {
				t.Fatalf("unexpected status: got %q want %q", got, tt.want)
			}
		})
	}
}
