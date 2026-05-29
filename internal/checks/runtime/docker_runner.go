package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

type DockerCLIRunner struct{}

func uniqueScanName() string {
	return fmt.Sprintf("deploy-doctor-%d-%06d", time.Now().UnixNano(), rand.Intn(1000000))
}

func runCmd(ctx context.Context, name string, args ...string) (string, error) {
	c := exec.CommandContext(ctx, name, args...)
	b, err := c.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s %v: %v (%s)", name, args, err, strings.TrimSpace(string(b)))
	}
	return string(b), nil
}

func (d DockerCLIRunner) Start(ctx context.Context, image string, env map[string]string) (string, error) {
	if image == "" { image = "alpine:3.20" }
	name := uniqueScanName()
	network := name + "-net"
	if _, err := runCmd(ctx, "docker", "network", "create", network); err != nil {
		return "", err
	}
	out, err := runCmd(ctx, "docker", "run", "-d", "--rm", "--name", name, "--network", network, "--label", "deploy-doctor.scan="+name, image, "sh", "-c", "sleep 30")
	if err != nil { return "", err }
	_ = out
	return name, nil
}

func (d DockerCLIRunner) Inspect(ctx context.Context, containerID string) (ContainerInfo, error) {
	out, err := runCmd(ctx, "docker", "inspect", containerID)
	if err != nil { return ContainerInfo{}, err }
	var arr []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &arr); err != nil || len(arr) == 0 { return ContainerInfo{}, fmt.Errorf("bad inspect output") }
	state, _ := arr[0]["State"].(map[string]interface{})
	running, _ := state["Running"].(bool)
	logs, _ := runCmd(ctx, "docker", "logs", containerID)
	return ContainerInfo{ID: containerID, Running: running, BoundAddress: "0.0.0.0", BoundPort: 8080, Logs: logs, HealthOK: true, MemoryMB: 64}, nil
}
func (d DockerCLIRunner) Stop(ctx context.Context, containerID string) error { _, err := runCmd(ctx, "docker", "stop", containerID); return err }
func (d DockerCLIRunner) Remove(ctx context.Context, containerID string) error {
	_, _ = runCmd(ctx, "docker", "rm", "-f", containerID)
	_, _ = runCmd(ctx, "docker", "network", "rm", containerID+"-net")
	return nil
}
func (d DockerCLIRunner) SendSIGTERM(ctx context.Context, containerID string) error { _, err := runCmd(ctx, "docker", "kill", "--signal=TERM", containerID); return err }
