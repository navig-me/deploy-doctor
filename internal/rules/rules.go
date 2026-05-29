package rules

import (
	"context"
	"sort"

	"docker-doctor/internal/model"
)

// ScanContext carries scan-time inputs for rule evaluation.
type ScanContext struct{}

// Rule is core check contract.
type Rule interface {
	ID() string
	Category() string
	SeverityDefault() string
	Run(ctx context.Context, scanContext ScanContext) ([]model.Issue, error)
}

// Registry stores rules by stable ID.
type Registry struct {
	rules map[string]Rule
}

func NewRegistry() *Registry {
	return &Registry{rules: map[string]Rule{}}
}

func (r *Registry) Register(rule Rule) error {
	id := rule.ID()
	if _, exists := r.rules[id]; exists {
		return ErrDuplicateRuleID(id)
	}
	r.rules[id] = rule
	return nil
}

func (r *Registry) Get(id string) (Rule, bool) {
	rule, ok := r.rules[id]
	return rule, ok
}

func (r *Registry) List() []Rule {
	ids := make([]string, 0, len(r.rules))
	for id := range r.rules {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	out := make([]Rule, 0, len(ids))
	for _, id := range ids {
		out = append(out, r.rules[id])
	}
	return out
}
