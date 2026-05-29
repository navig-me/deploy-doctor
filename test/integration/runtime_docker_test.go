package integration

import (
	"os"
	"testing"
)

func TestDockerEnabledIntegrationSuite(t *testing.T) {
	if os.Getenv("DOCKER_RUNTIME_TEST") != "1" {
		t.Skip("set DOCKER_RUNTIME_TEST=1 to run Docker-capable integration suite")
	}
}
