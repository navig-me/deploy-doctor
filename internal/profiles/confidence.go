package profiles

import "docker-doctor/internal/model"

func ConfidenceForProfile(profileName string) string {
	if profileName == "generic" {
		return "low"
	}
	return "high"
}

func ApplyConfidence(profileName string, issues []model.Issue) []model.Issue {
	conf := ConfidenceForProfile(profileName)
	out := make([]model.Issue, len(issues))
	copy(out, issues)
	for i := range out {
		out[i].Confidence = conf
	}
	return out
}
