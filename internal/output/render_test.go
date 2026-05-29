package output

import (
	"os"
	"path/filepath"
	"testing"

	"deploy-doctor/internal/model"
)

func sample() model.ScanResult {
	return model.ScanResult{Score: 72, Status: "risky", Issues: []model.Issue{{ID: "DF_BASE_0001", Title: "Base image may be oversized", Severity: model.SeverityWarning}}}
}

func snap(t *testing.T, name, got string) {
	t.Helper()
	p := filepath.Join("testdata", name)
	want, err := os.ReadFile(p)
	if err != nil { t.Fatalf("read snapshot %s: %v", name, err) }
	if got != string(want) { t.Fatalf("snapshot mismatch for %s\n--- got ---\n%s\n--- want ---\n%s", name, got, string(want)) }
}

func TestTextSnapshot(t *testing.T) {
	t.Parallel()
	snap(t, "text.golden", RenderText(sample()))
}

func TestJSONSnapshot(t *testing.T) {
	t.Parallel()
	got, err := RenderJSON(sample())
	if err != nil { t.Fatal(err) }
	snap(t, "json.golden", got)
}

func TestSARIFSnapshot(t *testing.T) {
	t.Parallel()
	got, err := RenderSARIF(sample())
	if err != nil { t.Fatal(err) }
	snap(t, "sarif.golden", got)
}

func TestMarkdownSnapshot(t *testing.T) {
	t.Parallel()
	snap(t, "markdown.golden", RenderMarkdown(sample()))
}
