package provider

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectSignalsRenderSchemaAware(t *testing.T) {
	t.Parallel()
	d := t.TempDir()
	_ = os.WriteFile(filepath.Join(d, "render.yaml"), []byte("services:\n  - type: web\n    startCommand: npm start\n    healthCheckPath: /health\n"), 0o644)
	s := DetectSignals(d)
	if len(s) == 0 || s[0].Profile != "render" || s[0].Confidence != "high" { t.Fatalf("unexpected signals: %+v", s) }
	if s[0].Evidence["hasStartCommand"] != "true" { t.Fatalf("missing schema evidence: %+v", s[0].Evidence) }
}

func TestDetectSignalsFlySchemaAware(t *testing.T) {
	t.Parallel()
	d := t.TempDir()
	_ = os.WriteFile(filepath.Join(d, "fly.toml"), []byte("app = \"demo\"\n[http_service]\ninternal_port = 8080\n"), 0o644)
	s := DetectSignals(d)
	if len(s) == 0 || s[0].Profile != "flyio" || s[0].Confidence != "high" { t.Fatalf("unexpected signals: %+v", s) }
}

func TestDetectSignalsECSSchemaAware(t *testing.T) {
	t.Parallel()
	d := t.TempDir()
	json := `{"requiresCompatibilities":["FARGATE"],"networkMode":"awsvpc","containerDefinitions":[{"portMappings":[{"containerPort":8080}]}]}`
	_ = os.WriteFile(filepath.Join(d, "ecs-task-definition.json"), []byte(json), 0o644)
	s := DetectSignals(d)
	if len(s) == 0 || s[0].Profile != "ecs-fargate" || s[0].Confidence != "high" { t.Fatalf("unexpected signals: %+v", s) }
}

func TestProviderValidationIssuesMissingFields(t *testing.T) {
	t.Parallel()
	signals := []Signal{{
		Profile:    "render",
		Confidence: "medium",
		Evidence: map[string]string{
			"file":               "render.yaml",
			"hasStartCommand":    "true",
			"hasHealthCheckPath": "false",
		},
	}}
	issues := ProviderValidationIssues("render", signals)
	if len(issues) != 1 || issues[0].ID != "CLD_CFG_0002" {
		t.Fatalf("expected CLD_CFG_0002, got %+v", issues)
	}
	missing, ok := issues[0].Evidence["missing_fields"].([]string)
	if !ok || len(missing) == 0 {
		t.Fatalf("expected missing_fields evidence, got %+v", issues[0].Evidence)
	}
}

func TestProviderValidationIssuesNoMissingFields(t *testing.T) {
	t.Parallel()
	signals := []Signal{{
		Profile:    "flyio",
		Confidence: "high",
		Evidence: map[string]string{
			"file":            "fly.toml",
			"hasApp":          "true",
			"hasInternalPort": "true",
		},
	}}
	if issues := ProviderValidationIssues("flyio", signals); len(issues) != 0 {
		t.Fatalf("expected no validation issues, got %+v", issues)
	}
}
