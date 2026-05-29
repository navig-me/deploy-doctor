package runtime

import "testing"

func TestUniqueScanNameIsolation(t *testing.T) {
	t.Parallel()
	a := uniqueScanName()
	b := uniqueScanName()
	if a == b {
		t.Fatalf("expected unique names, got duplicate %q", a)
	}
	if len(a) < len("deploy-doctor-") || a[:14] != "deploy-doctor-" {
		t.Fatalf("unexpected prefix: %s", a)
	}
}
