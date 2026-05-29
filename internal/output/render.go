package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"deploy-doctor/internal/model"
)

const JSONSchemaVersion = "v1"

func RenderText(r model.ScanResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Score: %d/100 (%s)\n", r.Score, r.Status)
	for _, is := range r.Issues {
		fmt.Fprintf(&b, "- [%s] %s: %s\n", is.Severity, is.ID, is.Title)
	}
	return b.String()
}

func RenderJSON(r model.ScanResult) (string, error) {
	type payload struct {
		SchemaVersion string           `json:"schema_version"`
		Result        model.ScanResult `json:"result"`
	}
	b, err := json.MarshalIndent(payload{SchemaVersion: JSONSchemaVersion, Result: r}, "", "  ")
	if err != nil { return "", err }
	return string(b) + "\n", nil
}

func RenderSARIF(r model.ScanResult) (string, error) {
	type sarifResult struct { RuleID, Level, Message string }
	out := map[string]interface{}{
		"version": "2.1.0",
		"runs": []interface{}{map[string]interface{}{
			"tool": map[string]interface{}{"driver": map[string]interface{}{"name": "deploy-doctor"}},
			"results": func() []sarifResult { rr := []sarifResult{}; for _, i := range r.Issues { rr = append(rr, sarifResult{RuleID: i.ID, Level: string(i.Severity), Message: i.Title}) }; return rr }(),
		}},
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil { return "", err }
	return string(b) + "\n", nil
}

func RenderMarkdown(r model.ScanResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Deploy Doctor Report\n\n")
	fmt.Fprintf(&b, "- Score: **%d/100**\n", r.Score)
	fmt.Fprintf(&b, "- Status: **%s**\n\n", r.Status)
	fmt.Fprintf(&b, "## Issues\n")
	for _, is := range r.Issues {
		fmt.Fprintf(&b, "- `%s` **%s**: %s\n", is.ID, is.Severity, is.Title)
	}
	return b.String()
}
