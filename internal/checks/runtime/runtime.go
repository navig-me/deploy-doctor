package runtime

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"deploy-doctor/internal/model"
)

type ContainerInfo struct { ID string; Running bool; ExitCode *int; BoundAddress string; BoundPort int; Logs string; HealthOK bool; MemoryMB int }
type Runner interface { Start(ctx context.Context, image string, env map[string]string) (string, error); Inspect(ctx context.Context, containerID string) (ContainerInfo, error); Stop(ctx context.Context, containerID string) error; Remove(ctx context.Context, containerID string) error; SendSIGTERM(ctx context.Context, containerID string) error }
type ProbeConfig struct { ExpectedPort int; HealthPath string; Timeout time.Duration; MemoryWarnMB int }
type ProbeResult struct { Issues []model.Issue; Errs []error; Duration time.Duration; CleanupStopOK bool; CleanupRemoveOK bool }

func RunProbes(ctx context.Context, r Runner, image string, cfg ProbeConfig) ProbeResult {
	start := time.Now()
	if cfg.Timeout <= 0 { cfg.Timeout = 20 * time.Second }
	if cfg.ExpectedPort == 0 { cfg.ExpectedPort = 8080 }
	if cfg.MemoryWarnMB == 0 { cfg.MemoryWarnMB = 512 }
	pctx, cancel := context.WithTimeout(ctx, cfg.Timeout); defer cancel()
	cid, err := r.Start(pctx, image, nil)
	if err != nil { return ProbeResult{Issues: []model.Issue{issue("RT_BOOT_0001","Container failed to start",model.SeverityCritical,"App cannot boot","Check startup command and required env")}, Errs: []error{err}, Duration: time.Since(start)} }
	stopOK, removeOK := cleanupWithRetry(cid, r)
	defer func() { _, _ = stopOK, removeOK }()
	info, err := r.Inspect(pctx, cid)
	if err != nil { return ProbeResult{Errs: []error{err}, Duration: time.Since(start), CleanupStopOK: stopOK, CleanupRemoveOK: removeOK} }
	issues := make([]model.Issue, 0)
	if !info.Running { issues = append(issues, issue("RT_EXIT_0001","Container exited early",model.SeverityCritical,"Process does not stay alive","Keep main process in foreground and investigate startup errors")) }
	if info.BoundPort != cfg.ExpectedPort { issues = append(issues, issue("RT_PORT_0001","Expected port not listening",model.SeverityCritical,"Ingress cannot reach app","Bind service to configured platform port")) }
	if strings.TrimSpace(info.BoundAddress) == "127.0.0.1" { issues = append(issues, issue("RT_BIND_0001","Service binds localhost only",model.SeverityCritical,"External traffic cannot reach service","Bind to 0.0.0.0")) }
	if !info.HealthOK { issues = append(issues, issue("RT_HEAL_0001","Health check failed",model.SeverityWarning,"Platform may restart instance","Expose healthy endpoint and update config")) }
	if strings.TrimSpace(info.Logs) == "" { issues = append(issues, issue("RT_LOG_0001","No stdout/stderr logs detected",model.SeveritySuggestion,"Observability degraded","Log app output to stdout/stderr")) }
	if info.MemoryMB > cfg.MemoryWarnMB { issues = append(issues, issue("RT_MEM_0001","High startup memory usage",model.SeverityWarning,"May cause OOM or slow cold starts","Reduce startup allocations and tune runtime")) }
	if err := r.SendSIGTERM(pctx, cid); err != nil { issues = append(issues, issue("RT_SIGT_0001","SIGTERM handling appears broken",model.SeverityWarning,"Graceful shutdown may fail","Handle SIGTERM and stop cleanly")) }
	return ProbeResult{Issues: issues, Duration: time.Since(start), CleanupStopOK: stopOK, CleanupRemoveOK: removeOK}
}

func cleanupWithRetry(cid string, r Runner) (bool, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second); defer cancel()
	stopOK, removeOK := false, false
	for i := 0; i < 3 && !stopOK; i++ { if err := r.Stop(ctx, cid); err == nil { stopOK = true } }
	for i := 0; i < 3 && !removeOK; i++ { if err := r.Remove(ctx, cid); err == nil { removeOK = true } }
	return stopOK, removeOK
}

func issue(id, title string, sev model.Severity, impact, fix string) model.Issue { return model.Issue{ID:id,Title:title,Severity:sev,Category:model.CategoryRuntime,Impact:impact,Fix:fix} }
var ErrDockerUnavailable = errors.New("docker unavailable")
func IsDockerAvailable(ctx context.Context, r Runner) error { _, err := r.Start(ctx, "", nil); if err != nil { return fmt.Errorf("%w: %v", ErrDockerUnavailable, err) }; return nil }
