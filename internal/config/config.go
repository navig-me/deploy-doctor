package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"docker-doctor/internal/model"
	"gopkg.in/yaml.v3"
)

const DefaultConfigPath = ".deploy-doctor.yml"
const DefaultProfile = "generic"

type Config struct {
	Version    int         `yaml:"version"`
	Profile    string      `yaml:"profile"`
	Timeouts   Timeouts    `yaml:"timeouts"`
	Thresholds Thresholds  `yaml:"thresholds"`
	Rules      RulesConfig `yaml:"rules"`
	Runtime    Runtime     `yaml:"runtime"`
}

type Timeouts struct {
	StartupSeconds int `yaml:"startup_seconds"`
	HealthSeconds  int `yaml:"health_seconds"`
}

type Thresholds struct {
	ImageSizeMBWarn     int `yaml:"image_size_mb_warn"`
	ImageSizeMBCritical int `yaml:"image_size_mb_critical"`
}

type RulesConfig struct {
	Ignore            []string          `yaml:"ignore"`
	SeverityOverrides map[string]string `yaml:"severity_overrides"`
}

type Runtime struct {
	ExpectedPort int    `yaml:"expected_port"`
	HealthPath   string `yaml:"health_path"`
}

func Load(path string) (Config, error) {
	if path == "" {
		path = DefaultConfigPath
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config yaml: %w", err)
	}
	if err := ValidateAndApplyDefaults(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func ValidateAndApplyDefaults(cfg *Config) error {
	if cfg.Version == 0 {
		cfg.Version = 1
	}
	if cfg.Version != 1 {
		return fmt.Errorf("unsupported config version: %d", cfg.Version)
	}
	if cfg.Profile == "" {
		cfg.Profile = DefaultProfile
	}
	if cfg.Timeouts.StartupSeconds <= 0 {
		cfg.Timeouts.StartupSeconds = 45
	}
	if cfg.Timeouts.HealthSeconds <= 0 {
		cfg.Timeouts.HealthSeconds = 20
	}
	if cfg.Runtime.ExpectedPort < 0 || cfg.Runtime.ExpectedPort > 65535 {
		return errors.New("runtime.expected_port must be between 0 and 65535")
	}
	if cfg.Runtime.ExpectedPort == 0 {
		cfg.Runtime.ExpectedPort = 8080
	}
	if cfg.Runtime.HealthPath == "" {
		cfg.Runtime.HealthPath = "/health"
	}
	if !strings.HasPrefix(cfg.Runtime.HealthPath, "/") {
		return errors.New("runtime.health_path must start with '/'")
	}
	if cfg.Thresholds.ImageSizeMBWarn <= 0 {
		cfg.Thresholds.ImageSizeMBWarn = 800
	}
	if cfg.Thresholds.ImageSizeMBCritical <= 0 {
		cfg.Thresholds.ImageSizeMBCritical = 1500
	}
	if cfg.Thresholds.ImageSizeMBWarn >= cfg.Thresholds.ImageSizeMBCritical {
		return errors.New("thresholds.image_size_mb_warn must be less than thresholds.image_size_mb_critical")
	}
	for ruleID, sev := range cfg.Rules.SeverityOverrides {
		if _, err := model.ParseSeverity(sev); err != nil {
			return fmt.Errorf("rules.severity_overrides[%s]: %w", ruleID, err)
		}
	}
	return nil
}

func ResolveProfile(profileFlag string, cfg Config) string {
	if strings.TrimSpace(profileFlag) != "" {
		return profileFlag
	}
	if strings.TrimSpace(cfg.Profile) != "" {
		return cfg.Profile
	}
	return DefaultProfile
}
