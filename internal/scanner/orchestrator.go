package scanner

import (
	"deploy-doctor/internal/model"

	"context"
	"fmt"
	"sort"
	"sync"

	"deploy-doctor/internal/rules"
)

// RuleError keeps rule-level failure isolated from overall scan flow.
type RuleError struct {
	RuleID string
	Err    error
}

type Orchestrator struct {
	registry *rules.Registry
}

func NewOrchestrator(registry *rules.Registry) *Orchestrator {
	return &Orchestrator{registry: registry}
}

func (o *Orchestrator) Run(ctx context.Context, scanCtx rules.ScanContext) ([]model.Issue, []RuleError) {
	ruleList := o.registry.List()

	issuesCh := make(chan []model.Issue, len(ruleList))
	errCh := make(chan RuleError, len(ruleList))

	var wg sync.WaitGroup
	for _, rule := range ruleList {
		rule := rule
		wg.Add(1)
		go func() {
			defer wg.Done()
			issues, err := rule.Run(ctx, scanCtx)
			if err != nil {
				errCh <- RuleError{RuleID: rule.ID(), Err: err}
				return
			}
			issuesCh <- issues
		}()
	}

	wg.Wait()
	close(issuesCh)
	close(errCh)

	allIssues := make([]model.Issue, 0)
	for chunk := range issuesCh {
		allIssues = append(allIssues, chunk...)
	}

	sort.SliceStable(allIssues, func(i, j int) bool {
		if allIssues[i].ID == allIssues[j].ID {
			return allIssues[i].Title < allIssues[j].Title
		}
		return allIssues[i].ID < allIssues[j].ID
	})

	ruleErrs := make([]RuleError, 0)
	for e := range errCh {
		ruleErrs = append(ruleErrs, e)
	}
	sort.SliceStable(ruleErrs, func(i, j int) bool {
		return ruleErrs[i].RuleID < ruleErrs[j].RuleID
	})

	return allIssues, ruleErrs
}

func (e RuleError) Error() string {
	return fmt.Sprintf("rule %s failed: %v", e.RuleID, e.Err)
}
