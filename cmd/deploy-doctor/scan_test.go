package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestEvaluateFailure(t *testing.T) {
	t.Parallel()
	c := severityCounts{critical: 0, warning: 1, suggestion: 0}
	if !evaluateFailure("warning", false, c) { t.Fatal("warning threshold should fail") }
	if evaluateFailure("critical", false, c) { t.Fatal("critical threshold should not fail") }
	if !evaluateFailure("none", true, c) { t.Fatal("strict should fail on warning") }
}

func TestDeterministicExitCodeMatrix(t *testing.T) {
	t.Parallel()
	if deterministicExitCode(nil, false) != 0 { t.Fatal("expected 0") }
	if deterministicExitCode(errFailExit, true) != 1 { t.Fatal("expected 1") }
	if deterministicExitCode(errors.New("boom"), false) != 2 { t.Fatal("expected 2") }
}

func TestScanMachineReadableSummaryAndFailOn(t *testing.T) {
	t.Parallel()
	cmd := newRootCmd()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"scan", "--fail-on", "warning"})
	err := cmd.Execute()
	if err == nil { t.Fatal("expected fail-on warning to fail") }
	if !strings.Contains(out.String(), "DD_SUMMARY") { t.Fatalf("missing summary line: %s", out.String()) }
}
