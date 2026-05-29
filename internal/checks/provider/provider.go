package provider

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"deploy-doctor/internal/model"
)

type Signal struct {
	Profile    string
	Confidence string
	Evidence   map[string]string
}

func DetectSignals(cwd string) []Signal {
	out := []Signal{}
	if b, err := os.ReadFile(filepath.Join(cwd, "render.yaml")); err == nil {
		s := string(b)
		out = append(out, Signal{Profile: "render", Confidence: confidenceByFields(2, countContains(s, "healthCheckPath", "startCommand")), Evidence: map[string]string{"file": "render.yaml"}})
	}
	if b, err := os.ReadFile(filepath.Join(cwd, "fly.toml")); err == nil {
		s := string(b)
		out = append(out, Signal{Profile: "flyio", Confidence: confidenceByFields(2, countContains(s, "internal_port", "[http_service]")), Evidence: map[string]string{"file": "fly.toml"}})
	}
	if b, err := os.ReadFile(filepath.Join(cwd, "railway.json")); err == nil {
		_ = b
		out = append(out, Signal{Profile: "railway", Confidence: "high", Evidence: map[string]string{"file": "railway.json"}})
	}
	if b, err := os.ReadFile(filepath.Join(cwd, "ecs-task-definition.json")); err == nil {
		var m map[string]interface{}
		_ = json.Unmarshal(b, &m)
		out = append(out, Signal{Profile: "ecs-fargate", Confidence: "high", Evidence: map[string]string{"file": "ecs-task-definition.json"}})
	}
	if len(out) == 0 {
		out = append(out, Signal{Profile: "generic", Confidence: "low", Evidence: map[string]string{"reason": "no provider config files"}})
	}
	return out
}

func countContains(s string, needles ...string) int { c := 0; low := strings.ToLower(s); for _, n := range needles { if strings.Contains(low, strings.ToLower(n)) { c++ } }; return c }
func confidenceByFields(total, matched int) string { if matched >= total { return "high" }; if matched > 0 { return "medium" }; return "low" }

func ProviderEvidenceIssues(profile string, signals []Signal) []model.Issue {
	for _, s := range signals {
		if s.Profile == profile && profile != "generic" {
			return []model.Issue{{ID: "CLD_CFG_0001", Title: "Provider config detected", Severity: model.SeveritySuggestion, Category: model.CategoryCloud, Confidence: s.Confidence, Evidence: map[string]interface{}{"signal": s.Evidence}, Impact: "Provider contracts can be validated more accurately", Fix: "Keep provider config in repo for higher-fidelity checks"}}
		}
	}
	return nil
}
