package provider

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"deploy-doctor/internal/model"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type Signal struct {
	Profile    string
	Confidence string
	Evidence   map[string]string
}

type renderConfig struct {
	Services []struct {
		Type            string `yaml:"type"`
		Name            string `yaml:"name"`
		StartCommand    string `yaml:"startCommand"`
		HealthCheckPath string `yaml:"healthCheckPath"`
	} `yaml:"services"`
}

type flyConfig struct {
	App         string `toml:"app"`
	HTTPService struct {
		InternalPort int `toml:"internal_port"`
	} `toml:"http_service"`
}

type ecsTaskDef struct {
	RequiresCompatibilities []string `json:"requiresCompatibilities"`
	NetworkMode             string   `json:"networkMode"`
	ContainerDefinitions    []struct {
		PortMappings []struct {
			ContainerPort int `json:"containerPort"`
		} `json:"portMappings"`
	} `json:"containerDefinitions"`
}

func DetectSignals(cwd string) []Signal {
	out := []Signal{}

	if s, ok := parseLightsail(cwd); ok { out = append(out, s) }
	if s, ok := parseRender(cwd); ok { out = append(out, s) }
	if s, ok := parseFly(cwd); ok { out = append(out, s) }
	if s, ok := parseRailway(cwd); ok { out = append(out, s) }
	if s, ok := parseECS(cwd); ok { out = append(out, s) }

	if len(out) == 0 {
		out = append(out, Signal{Profile: "generic", Confidence: "low", Evidence: map[string]string{"reason": "no provider config files"}})
	}
	return out
}

func parseLightsail(cwd string) (Signal, bool) {
	if b, err := os.ReadFile(filepath.Join(cwd, "lightsail.yml")); err == nil {
		var raw map[string]interface{}
		if err := yaml.Unmarshal(b, &raw); err == nil {
			_, hasContainer := raw["container"]
			_, hasEndpoint := raw["publicEndpoint"]
			conf := "low"
			if hasContainer && hasEndpoint { conf = "high" } else if hasContainer || hasEndpoint { conf = "medium" }
			return Signal{Profile: "lightsail", Confidence: conf, Evidence: map[string]string{"file": "lightsail.yml", "hasContainer": strconv.FormatBool(hasContainer), "hasPublicEndpoint": strconv.FormatBool(hasEndpoint)}}, true
		}
	}
	if _, err := os.Stat(filepath.Join(cwd, "lightsail.json")); err == nil {
		return Signal{Profile: "lightsail", Confidence: "high", Evidence: map[string]string{"file": "lightsail.json"}}, true
	}
	return Signal{}, false
}

func parseRender(cwd string) (Signal, bool) {
	b, err := os.ReadFile(filepath.Join(cwd, "render.yaml"))
	if err != nil { return Signal{}, false }
	var cfg renderConfig
	if err := yaml.Unmarshal(b, &cfg); err != nil { return Signal{}, false }
	matched := 0
	hasStart, hasHealth := false, false
	for _, s := range cfg.Services {
		if s.Type == "web" && s.StartCommand != "" { hasStart = true }
		if s.Type == "web" && s.HealthCheckPath != "" { hasHealth = true }
	}
	if hasStart { matched++ }
	if hasHealth { matched++ }
	conf := confidenceByFields(2, matched)
	return Signal{Profile: "render", Confidence: conf, Evidence: map[string]string{"file": "render.yaml", "hasStartCommand": strconv.FormatBool(hasStart), "hasHealthCheckPath": strconv.FormatBool(hasHealth)}}, true
}

func parseFly(cwd string) (Signal, bool) {
	var cfg flyConfig
	if _, err := toml.DecodeFile(filepath.Join(cwd, "fly.toml"), &cfg); err != nil { return Signal{}, false }
	matched := 0
	hasApp := cfg.App != ""
	hasPort := cfg.HTTPService.InternalPort > 0
	if hasApp { matched++ }
	if hasPort { matched++ }
	return Signal{Profile: "flyio", Confidence: confidenceByFields(2, matched), Evidence: map[string]string{"file": "fly.toml", "hasApp": strconv.FormatBool(hasApp), "hasInternalPort": strconv.FormatBool(hasPort)}}, true
}

func parseRailway(cwd string) (Signal, bool) {
	if _, err := os.Stat(filepath.Join(cwd, "railway.json")); err == nil {
		return Signal{Profile: "railway", Confidence: "high", Evidence: map[string]string{"file": "railway.json"}}, true
	}
	return Signal{}, false
}

func parseECS(cwd string) (Signal, bool) {
	b, err := os.ReadFile(filepath.Join(cwd, "ecs-task-definition.json"))
	if err != nil { return Signal{}, false }
	var td ecsTaskDef
	if err := json.Unmarshal(b, &td); err != nil { return Signal{}, false }
	matched := 0
	hasFargate := false
	for _, c := range td.RequiresCompatibilities {
		if c == "FARGATE" { hasFargate = true; break }
	}
	hasAwsvpc := td.NetworkMode == "awsvpc"
	hasPortMapping := len(td.ContainerDefinitions) > 0 && len(td.ContainerDefinitions[0].PortMappings) > 0 && td.ContainerDefinitions[0].PortMappings[0].ContainerPort > 0
	if hasFargate { matched++ }
	if hasAwsvpc { matched++ }
	if hasPortMapping { matched++ }
	return Signal{Profile: "ecs-fargate", Confidence: confidenceByFields(3, matched), Evidence: map[string]string{"file": "ecs-task-definition.json", "hasFargateCompat": strconv.FormatBool(hasFargate), "hasAwsvpc": strconv.FormatBool(hasAwsvpc), "hasPortMapping": strconv.FormatBool(hasPortMapping)}}, true
}

func confidenceByFields(total, matched int) string { if matched >= total { return "high" }; if matched > 0 { return "medium" }; return "low" }

func ProviderEvidenceIssues(profile string, signals []Signal) []model.Issue {
	for _, s := range signals {
		if s.Profile == profile && profile != "generic" {
			return []model.Issue{{ID: "CLD_CFG_0001", Title: "Provider config detected", Severity: model.SeveritySuggestion, Category: model.CategoryCloud, Confidence: s.Confidence, Evidence: map[string]interface{}{"signal": s.Evidence}, Impact: "Provider contracts can be validated more accurately", Fix: "Keep provider config in repo for higher-fidelity checks"}}
		}
	}
	return nil
}

func ProviderValidationIssues(profile string, signals []Signal) []model.Issue {
	for _, s := range signals {
		if s.Profile != profile || profile == "generic" {
			continue
		}
		missing := missingRequiredFields(profile, s.Evidence)
		if len(missing) == 0 {
			return nil
		}
		return []model.Issue{{
			ID:         "CLD_CFG_0002",
			Title:      "Provider config missing required fields",
			Severity:   model.SeverityWarning,
			Category:   model.CategoryCloud,
			Confidence: "high",
			Evidence: map[string]interface{}{
				"profile":        profile,
				"missing_fields": missing,
				"signal":         s.Evidence,
			},
			Impact: "Provider-specific deployment contract may be incomplete or invalid.",
			Fix:    "Add the required provider fields shown in evidence for accurate deploy behavior.",
		}}
	}
	return nil
}

func missingRequiredFields(profile string, ev map[string]string) []string {
	req := map[string][]string{
		"lightsail":   {"hasContainer", "hasPublicEndpoint"},
		"render":      {"hasStartCommand", "hasHealthCheckPath"},
		"flyio":       {"hasApp", "hasInternalPort"},
		"ecs-fargate": {"hasFargateCompat", "hasAwsvpc", "hasPortMapping"},
	}
	required := req[profile]
	if len(required) == 0 {
		return nil
	}
	missing := make([]string, 0)
	for _, k := range required {
		if ev[k] != "true" {
			missing = append(missing, k)
		}
	}
	return missing
}
