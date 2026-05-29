package provider

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFixture(t *testing.T, root, srcRel, dstName string) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", srcRel))
	if err != nil { t.Fatalf("read fixture %s: %v", srcRel, err) }
	if err := os.WriteFile(filepath.Join(root, dstName), b, 0o644); err != nil { t.Fatalf("write fixture: %v", err) }
}

func TestProviderParserFixturesValidInvalidPermutations(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		fixture         string
		dst             string
		profile         string
		expectConfidence string
		expectMissing   bool
	}{
		{name: "render valid", fixture: "render/valid.yaml", dst: "render.yaml", profile: "render", expectConfidence: "high", expectMissing: false},
		{name: "render missing health", fixture: "render/missing-health.yaml", dst: "render.yaml", profile: "render", expectConfidence: "medium", expectMissing: true},
		{name: "render missing start", fixture: "render/missing-start.yaml", dst: "render.yaml", profile: "render", expectConfidence: "medium", expectMissing: true},
		{name: "fly valid", fixture: "fly/valid.toml", dst: "fly.toml", profile: "flyio", expectConfidence: "high", expectMissing: false},
		{name: "fly missing port", fixture: "fly/missing-port.toml", dst: "fly.toml", profile: "flyio", expectConfidence: "medium", expectMissing: true},
		{name: "ecs valid", fixture: "ecs/valid.json", dst: "ecs-task-definition.json", profile: "ecs-fargate", expectConfidence: "high", expectMissing: false},
		{name: "ecs missing port", fixture: "ecs/missing-port.json", dst: "ecs-task-definition.json", profile: "ecs-fargate", expectConfidence: "medium", expectMissing: true},
		{name: "lightsail valid", fixture: "lightsail/valid.yml", dst: "lightsail.yml", profile: "lightsail", expectConfidence: "high", expectMissing: false},
		{name: "lightsail missing endpoint", fixture: "lightsail/missing-endpoint.yml", dst: "lightsail.yml", profile: "lightsail", expectConfidence: "medium", expectMissing: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			d := t.TempDir()
			writeFixture(t, d, tt.fixture, tt.dst)
			signals := DetectSignals(d)
			if len(signals) == 0 || signals[0].Profile != tt.profile {
				t.Fatalf("unexpected signals: %+v", signals)
			}
			if signals[0].Confidence != tt.expectConfidence {
				t.Fatalf("unexpected confidence: got %s want %s", signals[0].Confidence, tt.expectConfidence)
			}
			issues := ProviderValidationIssues(tt.profile, signals)
			if tt.expectMissing && len(issues) == 0 {
				t.Fatalf("expected validation issue for missing fields")
			}
			if !tt.expectMissing && len(issues) != 0 {
				t.Fatalf("did not expect validation issues, got %+v", issues)
			}
		})
	}
}
