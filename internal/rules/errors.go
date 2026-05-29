package rules

import "fmt"

type duplicateRuleIDError struct {
	id string
}

func (e duplicateRuleIDError) Error() string {
	return fmt.Sprintf("duplicate rule id: %s", e.id)
}

func ErrDuplicateRuleID(id string) error {
	return duplicateRuleIDError{id: id}
}
