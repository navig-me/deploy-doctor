package provider

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectSignalsConfidence(t *testing.T) {
	t.Parallel()
	d := t.TempDir()
	_ = os.WriteFile(filepath.Join(d, "render.yaml"), []byte("healthCheckPath: /health\nstartCommand: run"), 0o644)
	s := DetectSignals(d)
	if len(s) == 0 || s[0].Profile != "render" || s[0].Confidence != "high" { t.Fatalf("unexpected signals: %+v", s) }
}
