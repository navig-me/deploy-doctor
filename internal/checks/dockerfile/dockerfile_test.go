package dockerfile

import "testing"

func TestParseDockerfileReliably(t *testing.T) {
	t.Parallel()
	df, err := Parse("FROM golang:1.22\nRUN echo hi")
	if err != nil || len(df.BaseImageList) != 1 || len(df.Instructions) < 2 { t.Fatalf("parse failed: %v", err) }
}

func TestRules(t *testing.T) {
	t.Parallel()
	tests := []struct{name, dockerfile, ruleID string}{
		{"DF_BASE_0001","FROM ubuntu:latest\nCMD [\"app\"]","DF_BASE_0001"},
		{"DF_USER_0001","FROM alpine\nCMD [\"app\"]","DF_USER_0001"},
		{"DF_HEALTH_0001","FROM alpine\nCMD [\"app\"]\nUSER app","DF_HEALTH_0001"},
		{"DF_CMD_0001","FROM alpine\nUSER app","DF_CMD_0001"},
		{"DF_CACHE_0001","FROM node:20\nCOPY . .\nRUN npm install\nCMD [\"node\"]\nUSER app","DF_CACHE_0001"},
		{"DF_SECRET_0001","FROM alpine\nENV API_KEY=abc\nCMD [\"app\"]\nUSER app","DF_SECRET_0001"},
		{"DF_APT_0001","FROM ubuntu:22.04\nRUN apt-get update && apt-get install -y curl\nCMD [\"bash\"]\nUSER app","DF_APT_0001"},
		{"DF_COPY_0001","FROM node:20\nCOPY . .\nRUN echo ok\nCMD [\"node\"]\nUSER app","DF_COPY_0001"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			df, err := Parse(tt.dockerfile)
			if err != nil { t.Fatalf("parse: %v", err) }
			issues := RunChecks(df)
			found := false
			for _, is := range issues { if is.ID == tt.ruleID { found = true; break } }
			if !found { t.Fatalf("expected %s in issues: %+v", tt.ruleID, issues) }
		})
	}
}

func TestRuleEvidenceIncludesLineColumnMapping(t *testing.T) {
	t.Parallel()
	df, err := Parse("FROM ubuntu:22.04\nRUN apt-get update && apt-get install -y curl\nCMD [\"bash\"]\nUSER app\n")
	if err != nil { t.Fatalf("parse: %v", err) }
	issues := RunChecks(df)
	var aptIssueFound bool
	for _, is := range issues {
		if is.ID == "DF_APT_0001" {
			aptIssueFound = true
			if is.Evidence == nil {
				t.Fatalf("expected evidence for DF_APT_0001")
			}
			if got := is.Evidence["line"]; got != 2 {
				t.Fatalf("expected line=2, got %v", got)
			}
			if got := is.Evidence["column"]; got != 1 {
				t.Fatalf("expected column=1, got %v", got)
			}
			if got := is.Evidence["keyword"]; got != "RUN" {
				t.Fatalf("expected keyword=RUN, got %v", got)
			}
		}
	}
	if !aptIssueFound {
		t.Fatalf("expected DF_APT_0001 issue")
	}
}
