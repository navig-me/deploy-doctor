package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()
	d := t.TempDir()
	p := filepath.Join(d, ".deploy-doctor.yml")
	yml := "version: 1\nprofile: render\nrules:\n  ignore:\n    - DF_HEALTH_0001\n  severity_overrides:\n    IMG_SIZE_0001: critical\nruntime:\n  expected_port: 3000\n  health_path: /readyz\nthresholds:\n  image_size_mb_warn: 700\n  image_size_mb_critical: 1200\n"
	if err := os.WriteFile(p, []byte(yml), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Profile != "render" || cfg.Runtime.ExpectedPort != 3000 || cfg.Runtime.HealthPath != "/readyz" {
		t.Fatalf("unexpected loaded config: %+v", cfg)
	}
	if len(cfg.Rules.Ignore) != 1 || cfg.Rules.Ignore[0] != "DF_HEALTH_0001" {
		t.Fatalf("ignore rules not loaded: %+v", cfg.Rules.Ignore)
	}
}

func TestValidateDefaults(t *testing.T) {
	t.Parallel()
	cfg := Config{Version: 1}
	if err := ValidateAndApplyDefaults(&cfg); err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	if cfg.Profile != "generic" || cfg.Timeouts.StartupSeconds != 45 || cfg.Timeouts.HealthSeconds != 20 {
		t.Fatalf("defaults not applied: %+v", cfg)
	}
	if cfg.Runtime.ExpectedPort != 8080 || cfg.Runtime.HealthPath != "/health" {
		t.Fatalf("runtime defaults not applied: %+v", cfg.Runtime)
	}
}

func TestValidateRejectsInvalidSeverityOverride(t *testing.T) {
	t.Parallel()
	cfg := Config{Version: 1, Rules: RulesConfig{SeverityOverrides: map[string]string{"DF_BASE_0001": "severe"}}}
	if err := ValidateAndApplyDefaults(&cfg); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestResolveProfilePrecedence(t *testing.T) {
	t.Parallel()
	cfg := Config{Profile: "render"}
	if got := ResolveProfile("flyio", cfg); got != "flyio" {
		t.Fatalf("flag precedence failed: %q", got)
	}
	if got := ResolveProfile("", cfg); got != "render" {
		t.Fatalf("config precedence failed: %q", got)
	}
	if got := ResolveProfile("", Config{}); got != "generic" {
		t.Fatalf("default precedence failed: %q", got)
	}
}
