package envdb

import (
	"context"
	"regexp"
	"strings"

	"deploy-doctor/internal/model"
)

var envRefPatterns = []*regexp.Regexp{
	regexp.MustCompile(`process\.env\.([A-Z0-9_]+)`),
	regexp.MustCompile(`os\.Getenv\("([A-Z0-9_]+)"\)`),
	regexp.MustCompile(`os\.environ\.get\("([A-Z0-9_]+)"\)`),
	regexp.MustCompile(`\$\{([A-Z0-9_]+)\}`),
}

func DetectEnvReferences(content string) []string {
	seen := map[string]bool{}
	for _, re := range envRefPatterns {
		for _, m := range re.FindAllStringSubmatch(content, -1) {
			if len(m) > 1 {
				seen[m[1]] = true
			}
		}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}

func CheckMissingRequiredEnv(required []string, provided map[string]string) []model.Issue {
	for _, k := range required {
		if strings.TrimSpace(provided[k]) == "" {
			return []model.Issue{issue("ENV_MISS_0001", "Missing required environment variable", model.SeverityCritical, "App may fail at runtime", "Set required env vars in deployment config")}
		}
	}
	return nil
}

func CheckUnsafeValues(provided map[string]string) []model.Issue {
	for k, v := range provided {
		x := strings.ToLower(v)
		if strings.Contains(k, "TLS") && x == "0" {
			return []model.Issue{issue("ENV_INSEC_0001", "Unsafe TLS-related env value", model.SeverityWarning, "Security may be weakened in production", "Enable TLS verification in production")}
		}
		if strings.Contains(k, "DEBUG") && x == "true" {
			return []model.Issue{issue("ENV_INSEC_0001", "Debug mode enabled", model.SeveritySuggestion, "Verbose debug may expose sensitive internals", "Disable debug mode in production")}
		}
	}
	return nil
}

func CheckDBURL(url string) []model.Issue {
	if url == "" { return nil }
	low := strings.ToLower(url)
	var out []model.Issue
	if strings.Contains(low, "postgres://") || strings.Contains(low, "postgresql://") {
		if !strings.Contains(low, "sslmode=") {
			out = append(out, issue("ENV_DB_0001", "DB URL missing sslmode", model.SeverityWarning, "Managed DBs often require explicit SSL mode", "Set sslmode=require/verify-full as appropriate"))
		}
	}
	if strings.Contains(low, "localhost") || strings.Contains(low, "127.0.0.1") {
		out = append(out, issue("ENV_HOST_0001", "DB URL points to localhost", model.SeverityCritical, "Containerized cloud deployments cannot reach host localhost DB", "Use managed DB hostname/service DNS"))
	}
	return out
}

type DBProber interface { Ping(ctx context.Context, dbURL string) error }

type MigrationProber interface { RunMigrationCheck(ctx context.Context) error }

func OptionalActiveDBProbe(ctx context.Context, enabled bool, dbURL string, p DBProber) []model.Issue {
	if !enabled || p == nil || strings.TrimSpace(dbURL) == "" { return nil }
	if err := p.Ping(ctx, dbURL); err != nil {
		return []model.Issue{issue("DB_CONN_0001", "Active DB connectivity probe failed", model.SeverityWarning, "Connectivity issues may block startup", "Verify network access, credentials, and SSL settings")}
	}
	return nil
}

func OptionalMigrationProbe(ctx context.Context, enabled bool, p MigrationProber) []model.Issue {
	if !enabled || p == nil { return nil }
	if err := p.RunMigrationCheck(ctx); err != nil {
		return []model.Issue{issue("DB_MIGR_0001", "Migration command probe failed", model.SeverityWarning, "Schema migration may fail during deploy", "Fix migration command and DB connectivity")}
	}
	return nil
}

func issue(id, title string, sev model.Severity, impact, fix string) model.Issue {
	return model.Issue{ID: id, Title: title, Severity: sev, Category: model.CategoryEnv, Impact: impact, Fix: fix}
}
