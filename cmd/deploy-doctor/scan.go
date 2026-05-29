package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"deploy-doctor/internal/checks/dockerfile"
	"deploy-doctor/internal/checks/envdb"
	imagecheck "deploy-doctor/internal/checks/image"
	providercheck "deploy-doctor/internal/checks/provider"
	runtimecheck "deploy-doctor/internal/checks/runtime"
	"deploy-doctor/internal/model"
	"deploy-doctor/internal/profiles"
	"deploy-doctor/internal/scoring"
	"github.com/spf13/cobra"
)

var errFailExit = errors.New("scan findings triggered failure")
type severityCounts struct{ critical, warning, suggestion int }

func evaluateFailure(failOn string, strict bool, c severityCounts) bool { if strict && (c.warning > 0 || c.critical > 0) { return true }; switch strings.ToLower(strings.TrimSpace(failOn)) { case "critical": return c.critical > 0; case "warning": return c.critical > 0 || c.warning > 0; case "suggestion": return c.critical > 0 || c.warning > 0 || c.suggestion > 0; default: return false } }
func collectCounts(issues []model.Issue) severityCounts { c:=severityCounts{}; for _,is:= range issues { switch is.Severity { case model.SeverityCritical: c.critical++; case model.SeverityWarning: c.warning++; case model.SeveritySuggestion: c.suggestion++ } }; return c }
func deterministicExitCode(err error, shouldFail bool) int { if err != nil && !errors.Is(err, errFailExit) { return 2 }; if shouldFail || errors.Is(err, errFailExit) { return 1 }; return 0 }

func tuneFalsePositives(profile string, issues []model.Issue) []model.Issue {
	out := make([]model.Issue, 0, len(issues))
	for _, is := range issues {
		if profile == "render" && is.ID == "ENV_HOST_0001" {
			is.Severity = model.SeverityWarning
			if is.Evidence == nil { is.Evidence = map[string]interface{}{} }
			is.Evidence["tuned"] = "render false-positive boundary"
		}
		out = append(out, is)
	}
	return out
}

func runStaticScan(cwd, profileName string, signals []providercheck.Signal) ([]model.Issue, error) {
	p, err := profiles.Get(profileName)
	if err != nil { return nil, err }
	issues := make([]model.Issue, 0)
	if df, err := dockerfile.ParseFile(filepath.Join(cwd, "Dockerfile")); err == nil { issues = append(issues, dockerfile.RunChecks(df)...)}
	issues = append(issues, imagecheck.CheckDockerignore(cwd)...)
	issues = append(issues, imagecheck.CheckContextJunk(cwd)...)
	if m, err := imagecheck.ReadMetadata(filepath.Join(cwd, ".deploy-doctor-image.json")); err == nil { t:=imagecheck.Thresholds{SizeWarnMB:p.Thresholds.ImageSizeWarnMB,SizeCriticalMB:p.Thresholds.ImageSizeCriticalMB,LayerWarn:p.Thresholds.LayerWarnCount}; issues=append(issues,imagecheck.CheckImageMetadata(m,t)...)}
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" { issues = append(issues, envdb.CheckDBURL(dbURL)...)}
	if os.Getenv("APP_DEBUG") == "true" || os.Getenv("DEBUG") == "true" { issues = append(issues, envdb.CheckUnsafeValues(map[string]string{"DEBUG":"true"})...)}
	issues = append(issues, providercheck.ProviderEvidenceIssues(profileName, signals)...)
	issues = profiles.ApplyConfidence(profileName, issues)
	issues = tuneFalsePositives(profileName, issues)
	sort.SliceStable(issues, func(i,j int) bool { if issues[i].ID==issues[j].ID { return issues[i].Title < issues[j].Title }; return issues[i].ID < issues[j].ID })
	return issues,nil
}

func newScanCmd() *cobra.Command {
	var autoProfile, staticOnly, verbose, runtimeEnabled bool
	var failOn, profileName string
	cmd := &cobra.Command{Use:"scan", Short:"Run static and runtime deploy checks", RunE: func(cmd *cobra.Command, args []string) error {
		cwd,_ := os.Getwd()
		signals := providercheck.DetectSignals(cwd)
		selected := profileName
		if autoProfile { selected = signals[0].Profile; cmd.Printf("Auto profile detected: %s (confidence: %s)\n", signals[0].Profile, signals[0].Confidence); cmd.Println("Auto profile mode includes baseline generic checks.") }
		if verbose { cmd.Printf("Verbose: cwd=%s\n", cwd); cmd.Printf("Verbose: selected profile=%s\n", selected); cmd.Printf("Verbose: provider signals=%v\n", signals) }
		issues, err := runStaticScan(cwd, selected, signals); if err != nil { return err }
		if !staticOnly && runtimeEnabled {
			rr := runtimecheck.RunProbes(context.Background(), runtimecheck.DockerCLIRunner{}, "", runtimecheck.ProbeConfig{Timeout: 15*time.Second})
			issues = append(issues, rr.Issues...)
			if verbose && len(rr.Errs) > 0 { cmd.Printf("Verbose: runtime probe errors=%v\n", rr.Errs) }
		} else if staticOnly { cmd.Println("Mode: static-only (runtime probes skipped).") }
		counts := collectCounts(issues)
		score := scoring.CalculateScore(issues)
		status := scoring.MapStatus(score, model.Summary{Critical:counts.critical,Warning:counts.warning,Suggestion:counts.suggestion})
		failed := evaluateFailure(failOn,false,counts)
		cmd.Printf("Scan summary: score=%d status=%s profile=%s findings=(critical:%d warning:%d suggestion:%d)\n", score, status, selected, counts.critical, counts.warning, counts.suggestion)
		cmd.Printf("DD_SUMMARY score=%d status=%s critical=%d warning=%d suggestion=%d fail=%t profile=%s\n", score, status, counts.critical, counts.warning, counts.suggestion, failed, selected)
		if verbose { for _,is := range issues { cmd.Printf("  - %s [%s] %s evidence=%v\n", is.ID, is.Severity, is.Title, is.Evidence) } }
		if failed { return errFailExit }
		cmd.Println("Scan completed successfully.")
		return nil
	}}
	cmd.Flags().BoolVar(&autoProfile,"auto-profile",false,"Detect likely profile and run detected + generic")
	cmd.Flags().BoolVar(&staticOnly,"static-only",false,"Run static checks only (skip runtime probes)")
	cmd.Flags().BoolVar(&runtimeEnabled,"runtime",true,"Enable runtime probes during scan")
	cmd.Flags().BoolVarP(&verbose,"verbose","v",false,"Enable verbose output")
	cmd.Flags().StringVar(&profileName,"profile","generic","Profile to apply (e.g. lightsail)")
	cmd.Flags().StringVar(&failOn,"fail-on","none","Fail on severity: none|critical|warning|suggestion")
	return cmd
}

func formatExit(code int) string { return fmt.Sprintf("%d", code) }
