package model

// Severity of finding.
type Severity string

const (
	SeverityCritical   Severity = "critical"
	SeverityWarning    Severity = "warning"
	SeveritySuggestion Severity = "suggestion"
)

// Category of finding.
type Category string

const (
	CategoryDockerfile Category = "dockerfile"
	CategoryImage      Category = "image"
	CategoryRuntime    Category = "runtime"
	CategoryEnv        Category = "env"
	CategoryDB         Category = "db"
	CategoryCloud      Category = "cloud"
)

// Issue is single rule finding.
type Issue struct {
	ID       string                 `json:"id"`
	Title    string                 `json:"title"`
	Severity Severity               `json:"severity"`
	Category Category               `json:"category"`
	Confidence string               `json:"confidence,omitempty"`
	Evidence map[string]interface{} `json:"evidence,omitempty"`
	Impact   string                 `json:"impact"`
	Fix      string                 `json:"fix"`
	DocsURL  string                 `json:"docs_url"`
}

// Summary contains issue counts by severity.
type Summary struct {
	Critical   int `json:"critical"`
	Warning    int `json:"warning"`
	Suggestion int `json:"suggestion"`
}

// Metadata describes scan execution context.
type Metadata struct {
	Profile   string `json:"profile"`
	Duration  string `json:"duration"`
	Platform  string `json:"platform"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// ScanResult is top-level scan output payload.
type ScanResult struct {
	Score    int      `json:"score"`
	Status   string   `json:"status"`
	Issues   []Issue  `json:"issues"`
	Summary  Summary  `json:"summary"`
	Metadata Metadata `json:"metadata"`
}
