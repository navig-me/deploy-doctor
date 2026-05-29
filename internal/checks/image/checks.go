package image

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"deploy-doctor/internal/model"
)

type Metadata struct {
	SizeBytes int64    `json:"size_bytes"`
	Layers    []string `json:"layers"`
	Arch      string   `json:"arch"`
	OS        string   `json:"os"`
	Packages  []string `json:"packages"`
}

type Thresholds struct {
	SizeWarnMB     int
	SizeCriticalMB int
	LayerWarn      int
}

func DefaultThresholds() Thresholds {
	return Thresholds{SizeWarnMB: 800, SizeCriticalMB: 1500, LayerWarn: 25}
}

func ReadMetadata(path string) (Metadata, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Metadata{}, err
	}
	var m Metadata
	if err := json.Unmarshal(b, &m); err != nil {
		return Metadata{}, err
	}
	return m, nil
}

func CheckDockerignore(contextDir string) []model.Issue {
	p := filepath.Join(contextDir, ".dockerignore")
	b, err := os.ReadFile(p)
	if err != nil {
		return []model.Issue{issue("CTX_IGNR_0001", "Missing .dockerignore", model.SeverityWarning, "Large/sensitive context can be sent to build", "Add .dockerignore with common excludes")}
	}
	content := string(b)
	required := []string{".git", "node_modules", ".env"}
	for _, needle := range required {
		if !strings.Contains(content, needle) {
			return []model.Issue{issue("CTX_IGNR_0001", ".dockerignore missing common excludes", model.SeveritySuggestion, "Build context may include unnecessary files", "Add .git, node_modules, .env to .dockerignore")}
		}
	}
	return nil
}

func CheckContextJunk(contextDir string) []model.Issue {
	candidates := []string{".env", ".git", "node_modules", ".cache"}
	for _, name := range candidates {
		if _, err := os.Stat(filepath.Join(contextDir, name)); err == nil {
			return []model.Issue{issue("CTX_JUNK_0001", "Junk/sensitive files present in build context", model.SeverityWarning, "Context may leak sensitive or large files", "Exclude junk files via .dockerignore")}
		}
	}
	return nil
}

func CheckImageMetadata(m Metadata, t Thresholds) []model.Issue {
	var out []model.Issue
	sizeMB := int(m.SizeBytes / 1024 / 1024)
	if sizeMB >= t.SizeCriticalMB {
		out = append(out, issue("IMG_SIZE_0001", "Image size exceeds critical threshold", model.SeverityCritical, "Large image slows deploy/cold start", "Reduce base image and prune artifacts"))
	} else if sizeMB >= t.SizeWarnMB {
		out = append(out, issue("IMG_SIZE_0001", "Image size exceeds warning threshold", model.SeverityWarning, "Image may be slower to ship/start", "Use slimmer base and multi-stage builds"))
	}
	if len(m.Layers) > t.LayerWarn {
		out = append(out, issue("IMG_LAYR_0001", "Image has many layers", model.SeveritySuggestion, "Many layers may indicate inefficient build", "Consolidate RUN instructions where safe"))
	}

	hostArch := runtime.GOARCH
	if m.Arch != "" && m.Arch != hostArch {
		out = append(out, issue("IMG_ARCH_0001", "Image architecture differs from host", model.SeverityWarning, "Arch mismatch can hide runtime problems", "Build/test for target architecture explicitly"))
	}

	for _, pkg := range m.Packages {
		low := strings.ToLower(pkg)
		if strings.Contains(low, "gcc") || strings.Contains(low, "g++") || strings.Contains(low, "make") || strings.Contains(low, "build-essential") {
			out = append(out, issue("IMG_BUILD_0001", "Build tools found in runtime image", model.SeveritySuggestion, "Runtime image includes unnecessary build toolchain", "Use multi-stage build and copy only runtime artifacts"))
			break
		}
	}
	return out
}

func issue(id, title string, sev model.Severity, impact, fix string) model.Issue {
	return model.Issue{ID: id, Title: title, Severity: sev, Category: model.CategoryImage, Impact: impact, Fix: fix}
}
