package runtime

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRunner struct {
	startErr    error
	inspectInfo ContainerInfo
	inspectErr  error
	sigtermErr  error
	stopped     bool
	removed     bool
}

func (f *fakeRunner) Start(ctx context.Context, image string, env map[string]string) (string, error) {
	if f.startErr != nil { return "", f.startErr }
	return "c1", nil
}
func (f *fakeRunner) Inspect(ctx context.Context, containerID string) (ContainerInfo, error) {
	return f.inspectInfo, f.inspectErr
}
func (f *fakeRunner) Stop(ctx context.Context, containerID string) error { f.stopped = true; return nil }
func (f *fakeRunner) Remove(ctx context.Context, containerID string) error { f.removed = true; return nil }
func (f *fakeRunner) SendSIGTERM(ctx context.Context, containerID string) error { return f.sigtermErr }

func TestRunProbesCoversRuntimeRules(t *testing.T) {
	t.Parallel()
	r := &fakeRunner{inspectInfo: ContainerInfo{Running: false, BoundAddress: "127.0.0.1", BoundPort: 3000, HealthOK: false, MemoryMB: 1024}}
	res := RunProbes(context.Background(), r, "img", ProbeConfig{ExpectedPort: 8080, Timeout: 2 * time.Second, MemoryWarnMB: 512})
	seen := map[string]bool{}
	for _, is := range res.Issues { seen[is.ID] = true }
	for _, id := range []string{"RT_EXIT_0001","RT_PORT_0001","RT_BIND_0001","RT_HEAL_0001","RT_LOG_0001","RT_MEM_0001"} {
		if !seen[id] { t.Fatalf("missing %s in %+v", id, res.Issues) }
	}
	if !r.stopped || !r.removed { t.Fatalf("cleanup guarantees failed") }
}

func TestRunProbesStartupFailure(t *testing.T) {
	t.Parallel()
	r := &fakeRunner{startErr: errors.New("boom")}
	res := RunProbes(context.Background(), r, "img", ProbeConfig{})
	if len(res.Issues) == 0 || res.Issues[0].ID != "RT_BOOT_0001" { t.Fatalf("expected RT_BOOT_0001, got %+v", res.Issues) }
}

func TestRunProbesSigtermFailure(t *testing.T) {
	t.Parallel()
	r := &fakeRunner{inspectInfo: ContainerInfo{Running: true, BoundAddress: "0.0.0.0", BoundPort: 8080, HealthOK: true, Logs: "ok"}, sigtermErr: errors.New("sig")}
	res := RunProbes(context.Background(), r, "img", ProbeConfig{ExpectedPort: 8080})
	found := false
	for _, is := range res.Issues { if is.ID == "RT_SIGT_0001" { found = true } }
	if !found { t.Fatalf("expected RT_SIGT_0001, got %+v", res.Issues) }
}
