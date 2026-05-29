package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestDoctorCommandPrintsDockerDiagnostic(t *testing.T) {
	t.Parallel()
	cmd := newRootCmd()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"doctor"})
	if err := cmd.Execute(); err != nil { t.Fatalf("execute failed: %v", err) }
	got := out.String()
	if !strings.Contains(got, "docker:") { t.Fatalf("expected docker diagnostic, got: %s", got) }
}
