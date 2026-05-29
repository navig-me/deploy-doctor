package dockerfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"docker-doctor/internal/model"
)

type Dockerfile struct {
	RawLines []string
	Stages   []Stage
}

type Stage struct {
	BaseImage    string
	Instructions []Instruction
}

type Instruction struct {
	Keyword string
	Value   string
	Line    int
}

var instrRe = regexp.MustCompile(`^([A-Z]+)\s+(.*)$`)

func ParseFile(path string) (Dockerfile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Dockerfile{}, err
	}
	return Parse(string(b))
}

func Parse(content string) (Dockerfile, error) {
	lines := strings.Split(content, "\n")
	norm := normalizeLines(lines)

	df := Dockerfile{RawLines: norm}
	var cur *Stage
	for i, line := range norm {
		m := instrRe.FindStringSubmatch(strings.TrimSpace(line))
		if len(m) != 3 {
			continue
		}
		k, v := m[1], m[2]
		ins := Instruction{Keyword: k, Value: v, Line: i + 1}
		if k == "FROM" {
			base := strings.Fields(v)[0]
			s := Stage{BaseImage: base}
			df.Stages = append(df.Stages, s)
			cur = &df.Stages[len(df.Stages)-1]
		}
		if cur == nil {
			return Dockerfile{}, fmt.Errorf("instruction before FROM at line %d", i+1)
		}
		cur.Instructions = append(cur.Instructions, ins)
	}
	if len(df.Stages) == 0 {
		return Dockerfile{}, fmt.Errorf("no FROM found")
	}
	return df, nil
}

func normalizeLines(lines []string) []string {
	out := make([]string, 0, len(lines))
	buf := ""
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasSuffix(line, "\\") {
			buf += strings.TrimSuffix(line, "\\") + " "
			continue
		}
		line = buf + line
		buf = ""
		out = append(out, strings.TrimSpace(line))
	}
	if strings.TrimSpace(buf) != "" {
		out = append(out, strings.TrimSpace(buf))
	}
	return out
}

func RunChecks(df Dockerfile) []model.Issue {
	var issues []model.Issue
	issues = append(issues, checkBase(df)...)    // DF_BASE_0001
	issues = append(issues, checkUser(df)...)    // DF_USER_0001
	issues = append(issues, checkHealth(df)...)  // DF_HEALTH_0001
	issues = append(issues, checkCmd(df)...)     // DF_CMD_0001
	issues = append(issues, checkCache(df)...)   // DF_CACHE_0001
	issues = append(issues, checkSecret(df)...)  // DF_SECRET_0001
	issues = append(issues, checkApt(df)...)     // DF_APT_0001
	issues = append(issues, checkCopy(df)...)    // DF_COPY_0001
	return issues
}

func issue(id, title string, sev model.Severity, impact, fix string) model.Issue {
	return model.Issue{ID: id, Title: title, Severity: sev, Category: model.CategoryDockerfile, Impact: impact, Fix: fix}
}

func checkBase(df Dockerfile) []model.Issue { for _, s := range df.Stages { if strings.Contains(s.BaseImage, ":latest") || (!strings.Contains(s.BaseImage, "alpine") && !strings.Contains(s.BaseImage, "slim") && !strings.Contains(s.BaseImage, "distroless")) { return []model.Issue{issue("DF_BASE_0001", "Base image may be oversized or unpinned", model.SeverityWarning, "Larger/unpinned base increases risk", "Use pinned slim/alpine/distroless tag")}}}; return nil }
func checkUser(df Dockerfile) []model.Issue { for i := len(df.Stages)-1; i >= 0; i-- { for j := len(df.Stages[i].Instructions)-1; j >= 0; j-- { ins := df.Stages[i].Instructions[j]; if ins.Keyword=="USER" { if strings.TrimSpace(ins.Value)=="root" || strings.TrimSpace(ins.Value)=="0" { return []model.Issue{issue("DF_USER_0001","Container runs as root",model.SeverityCritical,"Root user raises breakout risk","Set non-root USER in final stage")}}; return nil } } }; return []model.Issue{issue("DF_USER_0001","Container user not specified",model.SeverityWarning,"Default user is root","Set non-root USER in final stage")} }
func checkHealth(df Dockerfile) []model.Issue { for _, st := range df.Stages { for _, ins := range st.Instructions { if ins.Keyword=="HEALTHCHECK" { return nil } } }; return []model.Issue{issue("DF_HEALTH_0001","Missing HEALTHCHECK",model.SeverityWarning,"Unhealthy container may go undetected","Add HEALTHCHECK command")} }
func checkCmd(df Dockerfile) []model.Issue { last:=df.Stages[len(df.Stages)-1]; for _,ins:= range last.Instructions { if ins.Keyword=="CMD"||ins.Keyword=="ENTRYPOINT" { return nil } }; return []model.Issue{issue("DF_CMD_0001","Missing CMD/ENTRYPOINT",model.SeverityCritical,"Container may not start correctly","Set explicit CMD or ENTRYPOINT")} }
func checkCache(df Dockerfile) []model.Issue { last:=df.Stages[len(df.Stages)-1]; copyAll:=-1; deps:=-1; for i,ins:= range last.Instructions { if ins.Keyword=="COPY" && strings.Contains(ins.Value,". .") { copyAll=i }; if ins.Keyword=="RUN" && (strings.Contains(ins.Value,"npm install")||strings.Contains(ins.Value,"pip install")||strings.Contains(ins.Value,"go mod download")) { deps=i } }; if copyAll>=0 && deps>=0 && copyAll<deps { return []model.Issue{issue("DF_CACHE_0001","Poor cache ordering",model.SeveritySuggestion,"Dependency layer invalidates often","Copy dependency manifests before full source copy")}}; return nil }
func checkSecret(df Dockerfile) []model.Issue { re:=regexp.MustCompile(`(?i)(SECRET|PASSWORD|TOKEN|KEY)=`); for _,st:= range df.Stages { for _,ins:= range st.Instructions { if (ins.Keyword=="ENV"||ins.Keyword=="ARG") && re.MatchString(ins.Value) { return []model.Issue{issue("DF_SECRET_0001","Potential secret in Dockerfile",model.SeverityCritical,"Secrets may leak into image layers","Inject secrets at runtime, not ARG/ENV")}} } }; return nil }
func checkApt(df Dockerfile) []model.Issue { for _,st:= range df.Stages { for _,ins:= range st.Instructions { if ins.Keyword=="RUN" && strings.Contains(ins.Value,"apt-get update") && !strings.Contains(ins.Value,"rm -rf /var/lib/apt/lists") { return []model.Issue{issue("DF_APT_0001","apt cache not cleaned",model.SeveritySuggestion,"Image size increases from package index cache","Clean apt lists in same RUN layer")}} } }; return nil }
func checkCopy(df Dockerfile) []model.Issue { last:=df.Stages[len(df.Stages)-1]; for i,ins:= range last.Instructions { if ins.Keyword=="COPY" && strings.Contains(ins.Value,". .") && i<=1 { return []model.Issue{issue("DF_COPY_0001","Broad COPY too early",model.SeverityWarning,"Large context copy early hurts caching and may include junk","Copy only needed files first, broad copy later")}} }; return nil }
