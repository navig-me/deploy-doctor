package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"deploy-doctor/internal/model"
)

func TestScanAutoProfileConfidenceFromSignals(t *testing.T) {
	d := t.TempDir()
	_ = os.WriteFile(filepath.Join(d, "render.yaml"), []byte("services:\n  - type: web\n    startCommand: run\n    healthCheckPath: /health\n"), 0o644)
	_ = os.WriteFile(filepath.Join(d, "Dockerfile"), []byte("FROM alpine\nCMD [\"sh\"]\nUSER 1000\n"), 0o644)
	cwd, _ := os.Getwd(); defer func(){ _ = os.Chdir(cwd) }(); _ = os.Chdir(d)
	cmd := newRootCmd(); out := &bytes.Buffer{}; cmd.SetOut(out); cmd.SetErr(out)
	cmd.SetArgs([]string{"scan","--auto-profile","--static-only"})
	if err := cmd.Execute(); err != nil { t.Fatalf("execute failed: %v", err) }
	if !strings.Contains(out.String(), "confidence: high") { t.Fatalf("expected high confidence output: %s", out.String()) }
}

func TestFalsePositiveTuningRenderBoundary(t *testing.T) {
	t.Parallel()
	issues := tuneFalsePositives("render", []model.Issue{{ID:"ENV_HOST_0001", Severity:model.SeverityCritical}})
	if issues[0].Severity != model.SeverityWarning { t.Fatalf("expected severity downgrade: %+v", issues[0]) }
}
