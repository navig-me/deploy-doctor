package envdb

import (
	"context"
	"errors"
	"testing"
)

type fakeDB struct{ err error }
func (f fakeDB) Ping(ctx context.Context, dbURL string) error { return f.err }

type fakeMig struct{ err error }
func (f fakeMig) RunMigrationCheck(ctx context.Context) error { return f.err }

func TestDetectEnvReferences(t *testing.T) {
	t.Parallel()
	content := `process.env.DATABASE_URL\nos.Getenv("SECRET_KEY")\n${PORT}\nos.environ.get("REDIS_URL")`
	refs := DetectEnvReferences(content)
	if len(refs) != 4 { t.Fatalf("expected 4 refs, got %v", refs) }
}

func TestMissingRequiredEnvBoundary(t *testing.T) {
	t.Parallel()
	issues := CheckMissingRequiredEnv([]string{"DATABASE_URL"}, map[string]string{"DATABASE_URL": ""})
	if len(issues) == 0 || issues[0].ID != "ENV_MISS_0001" { t.Fatalf("expected ENV_MISS_0001") }
	if got := CheckMissingRequiredEnv([]string{"DATABASE_URL"}, map[string]string{"DATABASE_URL": "postgres://x"}); len(got) != 0 { t.Fatalf("unexpected false positive: %+v", got) }
}

func TestUnsafeValueBoundaries(t *testing.T) {
	t.Parallel()
	if got := CheckUnsafeValues(map[string]string{"DEBUG": "false"}); len(got) != 0 { t.Fatalf("unexpected false positive: %+v", got) }
	if got := CheckUnsafeValues(map[string]string{"APP_DEBUG": "true"}); len(got) == 0 || got[0].ID != "ENV_INSEC_0001" { t.Fatalf("expected ENV_INSEC_0001") }
}

func TestDBURLChecksAndLocalhostMisuse(t *testing.T) {
	t.Parallel()
	issues := CheckDBURL("postgres://u:p@localhost:5432/db")
	seen := map[string]bool{}
	for _, is := range issues { seen[is.ID] = true }
	if !seen["ENV_DB_0001"] || !seen["ENV_HOST_0001"] { t.Fatalf("expected both DB lint + localhost issues, got %+v", issues) }
	if got := CheckDBURL("postgres://u:p@db:5432/db?sslmode=require"); len(got) != 0 { t.Fatalf("unexpected false positive: %+v", got) }
}

func TestOptionalProbes(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	if got := OptionalActiveDBProbe(ctx, false, "postgres://x", fakeDB{}); len(got) != 0 { t.Fatalf("expected disabled no issues") }
	if got := OptionalActiveDBProbe(ctx, true, "", fakeDB{}); len(got) != 0 { t.Fatalf("expected empty-url no issues") }
	if got := OptionalActiveDBProbe(ctx, true, "postgres://x", fakeDB{err: errors.New("nope")}); len(got) == 0 || got[0].ID != "DB_CONN_0001" { t.Fatalf("expected DB_CONN_0001") }

	if got := OptionalMigrationProbe(ctx, false, fakeMig{}); len(got) != 0 { t.Fatalf("expected disabled no issues") }
	if got := OptionalMigrationProbe(ctx, true, fakeMig{err: errors.New("bad")}); len(got) == 0 || got[0].ID != "DB_MIGR_0001" { t.Fatalf("expected DB_MIGR_0001") }
}
