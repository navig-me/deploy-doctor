package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestScanStaticOnlyFlag(t *testing.T) {
	t.Parallel()
	cmd := newRootCmd()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"scan", "--static-only"})
	if err := cmd.Execute(); err != nil { t.Fatalf("execute failed: %v", err) }
	if !strings.Contains(out.String(), "Mode: static-only") { t.Fatalf("expected static-only marker, got: %s", out.String()) }
}
