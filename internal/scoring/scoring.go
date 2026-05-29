package scoring

import "docker-doctor/internal/model"

const (
	maxScore              = 100
	criticalDeduction     = 30
	warningDeduction      = 10
	suggestionDeduction   = 3
	passScoreThreshold    = 80
	failScoreThreshold    = 60
	passWarningMaxAllowed = 1
)

func CalculateScore(issues []model.Issue) int {
	score := maxScore
	for _, issue := range issues {
		switch issue.Severity {
		case model.SeverityCritical:
			score -= criticalDeduction
		case model.SeverityWarning:
			score -= warningDeduction
		case model.SeveritySuggestion:
			score -= suggestionDeduction
		}
	}
	if score < 0 {
		return 0
	}
	if score > maxScore {
		return maxScore
	}
	return score
}

func MapStatus(score int, summary model.Summary) string {
	if summary.Critical > 0 || score < failScoreThreshold {
		return "fail"
	}
	if score >= passScoreThreshold && summary.Warning <= passWarningMaxAllowed {
		return "pass"
	}
	return "risky"
}
