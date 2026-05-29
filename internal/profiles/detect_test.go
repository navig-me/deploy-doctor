package profiles

import "testing"

func TestRecommendFromRepoFiles(t *testing.T) {
	t.Parallel()
	got := RecommendFromRepoFiles([]string{"render.yaml", "fly.toml"})
	if len(got) < 2 || got[0] != "render" { t.Fatalf("unexpected recommendations: %+v", got) }
}

func TestProfileSpecificEnablement(t *testing.T) {
	t.Parallel()
	g, _ := Get("generic")
	r, _ := Get("render")
	if len(g.EnabledFamilies) == 0 || len(r.EnabledFamilies) == 0 { t.Fatal("families should be non-empty") }
	if r.Thresholds.ImageSizeWarnMB >= g.Thresholds.ImageSizeWarnMB { t.Fatal("render should have stricter warn threshold than generic") }
}
