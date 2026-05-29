package profiles

import (
	"testing"

	"docker-doctor/internal/model"
)

func TestApplyConfidenceLabels(t *testing.T) {
	t.Parallel()

	issues := []model.Issue{{ID: "DF_BASE_0001"}}
	g := ApplyConfidence("generic", issues)
	if g[0].Confidence != "low" { t.Fatalf("expected low confidence for generic, got %q", g[0].Confidence) }

	r := ApplyConfidence("render", issues)
	if r[0].Confidence != "high" { t.Fatalf("expected high confidence for provider profile, got %q", r[0].Confidence) }
}
