package model

import (
	"encoding/json"
	"testing"
)

func TestScanResultJSONShape(t *testing.T) {
	t.Parallel()

	result := ScanResult{Score: 72, Status: "risky", Issues: []Issue{{ID: "RT_BIND_0001", Title: "Binds", Severity: SeverityCritical, Category: CategoryRuntime}}, Summary: Summary{Critical: 1}, Metadata: Metadata{Profile: "generic"}}
	b, err := json.Marshal(result)
	if err != nil { t.Fatalf("marshal failed: %v", err) }
	var got map[string]interface{}
	if err := json.Unmarshal(b, &got); err != nil { t.Fatalf("unmarshal failed: %v", err) }
	for _, key := range []string{"score", "status", "issues", "summary", "metadata"} {
		if _, ok := got[key]; !ok { t.Fatalf("missing key %q", key) }
	}
}

func TestEnumValidation(t *testing.T) {
	t.Parallel()
	if _, err := ParseSeverity("bad"); err == nil { t.Fatal("expected severity error") }
	if _, err := ParseCategory("bad"); err == nil { t.Fatal("expected category error") }
}
