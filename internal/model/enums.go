package model

import "fmt"

func (s Severity) IsValid() bool {
	switch s {
	case SeverityCritical, SeverityWarning, SeveritySuggestion:
		return true
	default:
		return false
	}
}

func ParseSeverity(v string) (Severity, error) {
	s := Severity(v)
	if !s.IsValid() {
		return "", fmt.Errorf("invalid severity: %s", v)
	}
	return s, nil
}

func (c Category) IsValid() bool {
	switch c {
	case CategoryDockerfile, CategoryImage, CategoryRuntime, CategoryEnv, CategoryDB, CategoryCloud:
		return true
	default:
		return false
	}
}

func ParseCategory(v string) (Category, error) {
	c := Category(v)
	if !c.IsValid() {
		return "", fmt.Errorf("invalid category: %s", v)
	}
	return c, nil
}
