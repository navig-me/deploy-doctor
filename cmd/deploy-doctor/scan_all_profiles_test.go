package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"deploy-doctor/internal/profiles"
)

func TestScanSupportsAllProfiles(t *testing.T) {
	t.Parallel()

	d := t.TempDir()
	if err := os.WriteFile(filepath.Join(d, "Dockerfile"), []byte("FROM alpine:3.20\nCMD [\"sh\"]\nUSER 1000\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, ".dockerignore"), []byte(".git\nnode_modules\n.env\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()
	_ = os.Chdir(d)

	for _, p := range profiles.List() {
		p := p
		t.Run(p.Name, func(t *testing.T) {
			t.Parallel()
			cmd := newRootCmd()
			out := &bytes.Buffer{}
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"scan", "--profile", p.Name, "--static-only"})
			if err := cmd.Execute(); err != nil {
				t.Fatalf("scan failed for profile %s: %v", p.Name, err)
			}
			got := out.String()
			if !strings.Contains(got, "profile="+p.Name) {
				t.Fatalf("missing profile in summary for %s: %s", p.Name, got)
			}
		})
	}
}
