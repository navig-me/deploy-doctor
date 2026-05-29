package profiles

import "strings"

func RecommendFromRepoFiles(paths []string) []string {
	joined := strings.ToLower(strings.Join(paths, "\n"))
	seen := map[string]bool{}
	out := make([]string, 0)
	add := func(name string) {
		if !seen[name] {
			seen[name] = true
			out = append(out, name)
		}
	}
	if strings.Contains(joined, "render.yaml") {
		add("render")
	}
	if strings.Contains(joined, "fly.toml") {
		add("flyio")
	}
	if strings.Contains(joined, "railway.json") {
		add("railway")
	}
	if strings.Contains(joined, "app.yaml") || strings.Contains(joined, "cloudrun") {
		add("gcp-cloud-run")
	}
	if len(out) == 0 {
		add("generic")
	}
	return out
}
