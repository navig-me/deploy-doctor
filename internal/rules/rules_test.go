package rules

import (
	"context"
	"reflect"
	"testing"

	"docker-doctor/internal/model"
)

type fakeRule struct {
	id string
}

func (f fakeRule) ID() string { return f.id }
func (f fakeRule) Category() string { return "dockerfile" }
func (f fakeRule) SeverityDefault() string { return "warning" }
func (f fakeRule) Run(ctx context.Context, scanContext ScanContext) ([]model.Issue, error) {
	return nil, nil
}

func TestRegistryRejectsDuplicateRuleID(t *testing.T) {
	t.Parallel()

	reg := NewRegistry()
	if err := reg.Register(fakeRule{id: "DF_BASE_0001"}); err != nil {
		t.Fatalf("unexpected register error: %v", err)
	}
	if err := reg.Register(fakeRule{id: "DF_BASE_0001"}); err == nil {
		t.Fatal("expected duplicate ID error")
	}
}

func TestRegistryListIsDeterministicByID(t *testing.T) {
	t.Parallel()

	reg := NewRegistry()
	_ = reg.Register(fakeRule{id: "RT_PORT_0001"})
	_ = reg.Register(fakeRule{id: "DF_BASE_0001"})
	_ = reg.Register(fakeRule{id: "IMG_SIZE_0001"})

	got := reg.List()
	ids := make([]string, 0, len(got))
	for _, rule := range got {
		ids = append(ids, rule.ID())
	}

	want := []string{"DF_BASE_0001", "IMG_SIZE_0001", "RT_PORT_0001"}
	if !reflect.DeepEqual(ids, want) {
		t.Fatalf("unexpected order: got %v want %v", ids, want)
	}
}
