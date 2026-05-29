package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScanAddsActionableIssueWhenDockerfileParseFails(t *testing.T) {
	d := t.TempDir()
	if err := os.Mkdir(filepath.Join(d, "Dockerfile"), 0o755); err != nil {
		t.Fatal(err)
	}
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()
	_ = os.Chdir(d)

	cmd := newRootCmd()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"scan", "--static-only", "--verbose"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "DF_PARSE_0001") {
		t.Fatalf("expected DF_PARSE_0001 fallback issue, got: %s", got)
	}
	if !strings.Contains(got, "Dockerfile parse failed") {
		t.Fatalf("expected actionable parse failure title, got: %s", got)
	}
}
