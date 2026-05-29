package image

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func hasRule(issues []interface{}, _ string) bool { return false }

func containsRule(issues []struct{ID string}, id string) bool { return false }

func TestFixtureBasedChecks(t *testing.T) {
	t.Parallel()
	d := t.TempDir()

	if err := os.WriteFile(filepath.Join(d, ".dockerignore"), []byte(".git\n"), 0o644); err != nil { t.Fatal(err) }
	if err := os.WriteFile(filepath.Join(d, ".env"), []byte("SECRET=1\n"), 0o644); err != nil { t.Fatal(err) }

	ign := CheckDockerignore(d)
	if len(ign) == 0 || ign[0].ID != "CTX_IGNR_0001" { t.Fatalf("expected CTX_IGNR_0001, got %+v", ign) }

	junk := CheckContextJunk(d)
	if len(junk) == 0 || junk[0].ID != "CTX_JUNK_0001" { t.Fatalf("expected CTX_JUNK_0001, got %+v", junk) }
}

func TestReadMetadataAndImageChecks(t *testing.T) {
	t.Parallel()
	d := t.TempDir()
	arch := "amd64"
	if runtime.GOARCH == "amd64" { arch = "arm64" }
	fixture := `{"size_bytes":1700000000,"layers":["l1","l2","l3","l4","l5","l6","l7","l8","l9","l10","l11","l12","l13","l14","l15","l16","l17","l18","l19","l20","l21","l22","l23","l24","l25","l26"],"arch":"` + arch + `","os":"linux","packages":["gcc","curl"]}`
	p := filepath.Join(d, "image.json")
	if err := os.WriteFile(p, []byte(fixture), 0o644); err != nil { t.Fatal(err) }

	m, err := ReadMetadata(p)
	if err != nil { t.Fatalf("read metadata: %v", err) }
	issues := CheckImageMetadata(m, DefaultThresholds())

	seen := map[string]bool{}
	for _, is := range issues { seen[is.ID] = true }
	for _, id := range []string{"IMG_SIZE_0001","IMG_LAYR_0001","IMG_ARCH_0001","IMG_BUILD_0001"} {
		if !seen[id] { t.Fatalf("expected %s in %+v", id, issues) }
	}
}
