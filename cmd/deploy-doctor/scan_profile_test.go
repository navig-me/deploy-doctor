package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScanWithLightsailProfile(t *testing.T) {
	t.Parallel()
	d := t.TempDir()
	if err := os.WriteFile(filepath.Join(d, "Dockerfile"), []byte("FROM ubuntu:latest\nCMD [\"app\"]\n"), 0o644); err != nil { t.Fatal(err) }
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()
	_ = os.Chdir(d)

	cmd := newRootCmd()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"scan", "--profile", "lightsail", "--static-only"})
	if err := cmd.Execute(); err != nil { t.Fatalf("execute failed: %v", err) }
	got := out.String()
	if !strings.Contains(got, "profile=lightsail") { t.Fatalf("expected lightsail profile in summary: %s", got) }
}
