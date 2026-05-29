package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestScanVerboseMode(t *testing.T) {
	t.Parallel()
	cmd := newRootCmd()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"scan", "--static-only", "--verbose"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "Verbose: selected profile") {
		t.Fatalf("expected verbose output, got: %s", got)
	}
}
