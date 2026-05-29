package dockerfile

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"deploy-doctor/internal/model"
	parser "github.com/moby/buildkit/frontend/dockerfile/parser"
)

type Dockerfile struct {
	RawLines      []string
	Instructions  []Instruction
	BaseImageList []string
}

type Instruction struct {
	Keyword string
	Value   string
	Line    int
}

func ParseFile(path string) (Dockerfile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Dockerfile{}, err
	}
	return Parse(string(b))
}

func Parse(content string) (Dockerfile, error) {
	res, err := parser.Parse(bytes.NewBufferString(content))
	if err != nil {
		return Dockerfile{}, err
	}
	if res.AST == nil || len(res.AST.Children) == 0 {
		return Dockerfile{}, fmt.Errorf("empty dockerfile")
	}

	df := Dockerfile{RawLines: strings.Split(content, "\n")}
	for _, n := range res.AST.Children {
		k := strings.ToUpper(strings.TrimSpace(n.Value))
		v := strings.TrimSpace(n.Original)
		if len(v) >= len(k) {
			v = strings.TrimSpace(v[len(k):])
		}
		if k == "FROM" {
			parts := strings.Fields(v)
			if len(parts) > 0 {
				df.BaseImageList = append(df.BaseImageList, parts[0])
			}
		}
		df.Instructions = append(df.Instructions, Instruction{Keyword: k, Value: v, Line: n.StartLine})
	}
	if len(df.BaseImageList) == 0 {
		return Dockerfile{}, fmt.Errorf("no FROM found")
	}
	return df, nil
}

func RunChecks(df Dockerfile) []model.Issue {
	var issues []model.Issue
	issues = append(issues, checkBase(df)...)
	issues = append(issues, checkUser(df)...)
	issues = append(issues, checkHealth(df)...)
	issues = append(issues, checkCmd(df)...)
	issues = append(issues, checkCache(df)...)
	issues = append(issues, checkSecret(df)...)
	issues = append(issues, checkApt(df)...)
	issues = append(issues, checkCopy(df)...)
	return issues
}

func issue(id, title string, sev model.Severity, impact, fix string) model.Issue {
	return model.Issue{ID: id, Title: title, Severity: sev, Category: model.CategoryDockerfile, Impact: impact, Fix: fix}
}

func issueAt(id, title string, sev model.Severity, impact, fix string, line int, keyword string) model.Issue {
	is := issue(id, title, sev, impact, fix)
	is.Evidence = map[string]interface{}{"line": line, "column": 1, "keyword": keyword}
	return is
}

func checkBase(df Dockerfile) []model.Issue {
	for _, b := range df.BaseImageList {
		if strings.Contains(b, ":latest") || (!strings.Contains(b, "alpine") && !strings.Contains(b, "slim") && !strings.Contains(b, "distroless")) {
			return []model.Issue{issue("DF_BASE_0001", "Base image may be oversized or unpinned", model.SeverityWarning, "Larger/unpinned base increases risk", "Use pinned slim/alpine/distroless tag")}
		}
	}
	return nil
}

func checkUser(df Dockerfile) []model.Issue {
	for i := len(df.Instructions) - 1; i >= 0; i-- {
		ins := df.Instructions[i]
		if ins.Keyword == "USER" {
			v := strings.TrimSpace(ins.Value)
			if v == "root" || v == "0" {
				return []model.Issue{issueAt("DF_USER_0001", "Container runs as root", model.SeverityCritical, "Root user raises breakout risk", "Set non-root USER in final stage", ins.Line, ins.Keyword)}
			}
			return nil
		}
	}
	return []model.Issue{issueAt("DF_USER_0001", "Container user not specified", model.SeverityWarning, "Default user is root", "Set non-root USER in final stage", 1, "USER")}
}

func checkHealth(df Dockerfile) []model.Issue {
	for _, ins := range df.Instructions {
		if ins.Keyword == "HEALTHCHECK" {
			return nil
		}
	}
	return []model.Issue{issueAt("DF_HEALTH_0001", "Missing HEALTHCHECK", model.SeverityWarning, "Unhealthy container may go undetected", "Add HEALTHCHECK command", 1, "HEALTHCHECK")}
}

func checkCmd(df Dockerfile) []model.Issue {
	for _, ins := range df.Instructions {
		if ins.Keyword == "CMD" || ins.Keyword == "ENTRYPOINT" {
			return nil
		}
	}
	return []model.Issue{issueAt("DF_CMD_0001", "Missing CMD/ENTRYPOINT", model.SeverityCritical, "Container may not start correctly", "Set explicit CMD or ENTRYPOINT", 1, "CMD")}
}

func checkCache(df Dockerfile) []model.Issue {
	copyAll := -1
	deps := -1
	for i, ins := range df.Instructions {
		if ins.Keyword == "COPY" && strings.Contains(ins.Value, ". .") {
			copyAll = i
		}
		if ins.Keyword == "RUN" && (strings.Contains(ins.Value, "npm install") || strings.Contains(ins.Value, "pip install") || strings.Contains(ins.Value, "go mod download")) {
			deps = i
		}
	}
	if copyAll >= 0 && deps >= 0 && copyAll < deps {
		ins := df.Instructions[copyAll]
		return []model.Issue{issueAt("DF_CACHE_0001", "Poor cache ordering", model.SeveritySuggestion, "Dependency layer invalidates often", "Copy dependency manifests before full source copy", ins.Line, ins.Keyword)}
	}
	return nil
}

func checkSecret(df Dockerfile) []model.Issue {
	re := regexp.MustCompile(`(?i)(SECRET|PASSWORD|TOKEN|KEY)=`)
	for _, ins := range df.Instructions {
		if (ins.Keyword == "ENV" || ins.Keyword == "ARG") && re.MatchString(ins.Value) {
			return []model.Issue{issueAt("DF_SECRET_0001", "Potential secret in Dockerfile", model.SeverityCritical, "Secrets may leak into image layers", "Inject secrets at runtime, not ARG/ENV", ins.Line, ins.Keyword)}
		}
	}
	return nil
}

func checkApt(df Dockerfile) []model.Issue {
	for _, ins := range df.Instructions {
		if ins.Keyword == "RUN" && strings.Contains(ins.Value, "apt-get update") && !strings.Contains(ins.Value, "rm -rf /var/lib/apt/lists") {
			return []model.Issue{issueAt("DF_APT_0001", "apt cache not cleaned", model.SeveritySuggestion, "Image size increases from package index cache", "Clean apt lists in same RUN layer", ins.Line, ins.Keyword)}
		}
	}
	return nil
}

func checkCopy(df Dockerfile) []model.Issue {
	for i, ins := range df.Instructions {
		if ins.Keyword == "COPY" && strings.Contains(ins.Value, ". .") && i <= 1 {
			return []model.Issue{issueAt("DF_COPY_0001", "Broad COPY too early", model.SeverityWarning, "Large context copy early hurts caching and may include junk", "Copy only needed files first, broad copy later", ins.Line, ins.Keyword)}
		}
	}
	return nil
}
