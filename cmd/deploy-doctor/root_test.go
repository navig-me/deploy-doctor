package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootHasExpectedSubcommands(t *testing.T) {
	t.Parallel()

	cmd := newRootCmd()
	names := map[string]bool{}
	for _, c := range cmd.Commands() {
		names[c.Name()] = true
	}

	for _, name := range []string{"scan", "version", "doctor"} {
		if !names[name] {
			t.Fatalf("missing subcommand %q", name)
		}
	}
}

func TestVersionCommandPrintsVersion(t *testing.T) {
	t.Parallel()

	cmd := newRootCmd()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	if got := strings.TrimSpace(out.String()); got != version {
		t.Fatalf("unexpected version output: got %q want %q", got, version)
	}
}
