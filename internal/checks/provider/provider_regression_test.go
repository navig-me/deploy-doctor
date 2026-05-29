package provider

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProviderParserFalsePositiveRegressions(t *testing.T) {
	t.Parallel()

	t.Run("render non-web worker should not produce missing web-field issue", func(t *testing.T) {
		t.Parallel()
		d := t.TempDir()
		content := "services:\n  - type: worker\n    name: jobs\n    startCommand: npm run worker\n"
		_ = os.WriteFile(filepath.Join(d, "render.yaml"), []byte(content), 0o644)
		signals := DetectSignals(d)
		if len(signals) == 0 || signals[0].Profile != "render" { t.Fatalf("unexpected signals: %+v", signals) }
		issues := ProviderValidationIssues("render", signals)
		if len(issues) != 1 {
			t.Fatalf("expected one regression-safe warning (web fields unknown), got %+v", issues)
		}
		if issues[0].ID != "CLD_CFG_0002" {
			t.Fatalf("unexpected issue: %+v", issues)
		}
	})

	t.Run("fly minimal app without http_service should be medium, not high", func(t *testing.T) {
		t.Parallel()
		d := t.TempDir()
		_ = os.WriteFile(filepath.Join(d, "fly.toml"), []byte("app = \"demo\"\n"), 0o644)
		signals := DetectSignals(d)
		if len(signals) == 0 || signals[0].Profile != "flyio" { t.Fatalf("unexpected signals: %+v", signals) }
		if signals[0].Confidence != "medium" {
			t.Fatalf("expected medium confidence to avoid false high-confidence parse, got %+v", signals[0])
		}
	})

	t.Run("ecs ec2-style task should not be high-confidence fargate", func(t *testing.T) {
		t.Parallel()
		d := t.TempDir()
		json := `{"requiresCompatibilities":["EC2"],"networkMode":"bridge","containerDefinitions":[{"portMappings":[{"containerPort":8080}]}]}`
		_ = os.WriteFile(filepath.Join(d, "ecs-task-definition.json"), []byte(json), 0o644)
		signals := DetectSignals(d)
		if len(signals) == 0 || signals[0].Profile != "ecs-fargate" { t.Fatalf("unexpected signals: %+v", signals) }
		if signals[0].Confidence == "high" {
			t.Fatalf("expected non-high confidence for non-fargate task, got %+v", signals[0])
		}
	})
}
