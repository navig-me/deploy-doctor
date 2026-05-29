package dockerfile

import (
	"os"
	"path/filepath"
	"testing"
)

func readFixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return string(b)
}

func TestDockerfileRegressionFixtures(t *testing.T) {
	t.Parallel()

	t.Run("multistage parse and no critical defaults", func(t *testing.T) {
		t.Parallel()
		df, err := Parse(readFixture(t, "multistage.dockerfile"))
		if err != nil { t.Fatalf("parse: %v", err) }
		if len(df.BaseImageList) != 2 { t.Fatalf("expected 2 stages, got %d", len(df.BaseImageList)) }
		issues := RunChecks(df)
		for _, is := range issues {
			if is.ID == "DF_USER_0001" || is.ID == "DF_CMD_0001" {
				t.Fatalf("unexpected rule for good multistage fixture: %+v", issues)
			}
		}
	})

	t.Run("arg complex parse and secret detection", func(t *testing.T) {
		t.Parallel()
		df, err := Parse(readFixture(t, "arg-complex.dockerfile"))
		if err != nil { t.Fatalf("parse: %v", err) }
		issues := RunChecks(df)
		foundSecret := false
		for _, is := range issues {
			if is.ID == "DF_SECRET_0001" { foundSecret = true; break }
		}
		if !foundSecret { t.Fatalf("expected DF_SECRET_0001 in %+v", issues) }
	})

	t.Run("complex cache ordering should not trigger", func(t *testing.T) {
		t.Parallel()
		df, err := Parse(readFixture(t, "complex-copy-cache.dockerfile"))
		if err != nil { t.Fatalf("parse: %v", err) }
		issues := RunChecks(df)
		for _, is := range issues {
			if is.ID == "DF_CACHE_0001" {
				t.Fatalf("did not expect DF_CACHE_0001 for good ordering: %+v", issues)
			}
		}
	})
}
