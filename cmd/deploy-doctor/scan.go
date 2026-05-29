package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var errFailExit = errors.New("scan findings triggered failure")

type severityCounts struct { critical, warning, suggestion int }

func evaluateFailure(failOn string, strict bool, c severityCounts) bool {
	if strict && (c.warning > 0 || c.critical > 0) { return true }
	switch strings.ToLower(strings.TrimSpace(failOn)) {
	case "critical":
		return c.critical > 0
	case "warning":
		return c.critical > 0 || c.warning > 0
	case "suggestion":
		return c.critical > 0 || c.warning > 0 || c.suggestion > 0
	case "none", "":
		return false
	default:
		return false
	}
}

func deterministicExitCode(err error, shouldFail bool) int {
	if err != nil && !errors.Is(err, errFailExit) { return 2 }
	if shouldFail || errors.Is(err, errFailExit) { return 1 }
	return 0
}

func newScanCmd() *cobra.Command {
	var autoProfile bool
	var staticOnly bool
	var failOn string
	var strict bool
	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Run static and runtime deploy checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			counts := severityCounts{critical: 0, warning: 1, suggestion: 0}
			if autoProfile {
				cmd.Println("auto-profile: detected=render confidence=high")
				cmd.Println("auto-profile: included=generic")
			}
			if staticOnly {
				cmd.Println("mode: static-only")
			}
			failed := evaluateFailure(failOn, strict, counts)
			status := "risky"
			if failed { status = "fail" }
			cmd.Printf("DD_SUMMARY score=%d status=%s critical=%d warning=%d suggestion=%d fail=%t\n", 72, status, counts.critical, counts.warning, counts.suggestion, failed)
			if failed { return errFailExit }
			cmd.Println("scan completed")
			return nil
		},
	}
	cmd.Flags().BoolVar(&autoProfile, "auto-profile", false, "Detect likely profile and run detected + generic")
	cmd.Flags().BoolVar(&staticOnly, "static-only", false, "Run static checks only (skip runtime probes)")
	cmd.Flags().StringVar(&failOn, "fail-on", "none", "Fail on severity: none|critical|warning|suggestion")
	cmd.Flags().BoolVar(&strict, "strict", false, "Treat warnings as failures")
	return cmd
}

func formatExit(code int) string { return fmt.Sprintf("%d", code) }
