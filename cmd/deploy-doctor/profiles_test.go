package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProfilesListCommand(t *testing.T) {
	t.Parallel()
	cmd := newRootCmd()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"profiles", "list"})
	if err := cmd.Execute(); err != nil { t.Fatalf("execute failed: %v", err) }
	got := out.String()
	for _, name := range []string{"generic", "lightsail", "render", "railway", "flyio"} {
		if !strings.Contains(got, name) { t.Fatalf("expected profile %q in output", name) }
	}
}

func TestProfilesExplainCommand(t *testing.T) {
	t.Parallel()
	cmd := newRootCmd()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"profiles", "explain", "render"})
	if err := cmd.Execute(); err != nil { t.Fatalf("execute failed: %v", err) }
	got := out.String()
	if !strings.Contains(got, "render") || !strings.Contains(got, "thresholds:") { t.Fatalf("unexpected explain output: %s", got) }
}

func TestProfilesListRecommendedCommand(t *testing.T) {
	d := t.TempDir()
	if err := os.WriteFile(filepath.Join(d, "lightsail.yml"), []byte("container:\npublicEndpoint: 80\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()
	_ = os.Chdir(d)

	cmd := newRootCmd()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"profiles", "list", "--recommended"})
	if err := cmd.Execute(); err != nil { t.Fatalf("execute failed: %v", err) }
	got := out.String()
	if !strings.Contains(got, "lightsail") || !strings.Contains(got, "confidence=") { t.Fatalf("expected lightsail recommendation output, got: %s", got) }
}

func TestScanAutoProfileFlag(t *testing.T) {
	t.Parallel()
	cmd := newRootCmd()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"scan", "--auto-profile"})
	if err := cmd.Execute(); err != nil { t.Fatalf("execute failed: %v", err) }
	got := out.String()
	if !strings.Contains(got, "Auto profile detected:") { t.Fatalf("expected auto-profile output, got: %s", got) }
}
