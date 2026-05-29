package profiles

import "fmt"

type Thresholds struct {
	ImageSizeWarnMB     int
	ImageSizeCriticalMB int
	LayerWarnCount      int
	StartupMemoryWarnMB int
}

type Definition struct {
	Name            string
	Description     string
	EnabledFamilies []string
	Thresholds      Thresholds
}

var base = Definition{
	Name:            "base",
	Description:     "Base defaults shared by all profiles",
	EnabledFamilies: []string{"dockerfile", "context-image", "runtime", "env-db"},
	Thresholds:      Thresholds{ImageSizeWarnMB: 800, ImageSizeCriticalMB: 1500, LayerWarnCount: 25, StartupMemoryWarnMB: 512},
}

var generic = Definition{
	Name:            "generic",
	Description:     "Baseline portability checks with no provider-specific assumptions",
	EnabledFamilies: nil,
	Thresholds:      Thresholds{},
}

var lightsail = Definition{
	Name:            "lightsail",
	Description:     "AWS Lightsail container constraints and small-service defaults",
	EnabledFamilies: []string{"dockerfile", "context-image", "runtime", "env-db", "cloud"},
	Thresholds:      Thresholds{ImageSizeWarnMB: 600, ImageSizeCriticalMB: 1200, LayerWarnCount: 22, StartupMemoryWarnMB: 384},
}

var render = Definition{
	Name:            "render",
	Description:     "Render web service defaults including PORT contract and health checks",
	EnabledFamilies: []string{"dockerfile", "context-image", "runtime", "env-db", "cloud"},
	Thresholds:      Thresholds{ImageSizeWarnMB: 700, ImageSizeCriticalMB: 1300, LayerWarnCount: 24, StartupMemoryWarnMB: 448},
}

var railway = Definition{
	Name:            "railway",
	Description:     "Railway deploy defaults for PORT, service DNS, and runtime behavior",
	EnabledFamilies: []string{"dockerfile", "context-image", "runtime", "env-db", "cloud"},
	Thresholds:      Thresholds{ImageSizeWarnMB: 700, ImageSizeCriticalMB: 1250, LayerWarnCount: 24, StartupMemoryWarnMB: 448},
}

var flyio = Definition{
	Name:            "flyio",
	Description:     "Fly.io defaults for internal port alignment, graceful shutdown, and runtime health",
	EnabledFamilies: []string{"dockerfile", "context-image", "runtime", "env-db", "cloud"},
	Thresholds:      Thresholds{ImageSizeWarnMB: 650, ImageSizeCriticalMB: 1200, LayerWarnCount: 23, StartupMemoryWarnMB: 384},
}

var ecsFargate = Definition{
	Name:            "ecs-fargate",
	Description:     "ECS Fargate defaults for task sizing, port mapping, and runtime behavior",
	EnabledFamilies: []string{"dockerfile", "context-image", "runtime", "env-db", "cloud"},
	Thresholds:      Thresholds{ImageSizeWarnMB: 700, ImageSizeCriticalMB: 1300, LayerWarnCount: 24, StartupMemoryWarnMB: 448},
}

var digitaloceanAppPlatform = Definition{
	Name:            "digitalocean-app-platform",
	Description:     "DigitalOcean App Platform defaults for port, health checks, and app runtime constraints",
	EnabledFamilies: []string{"dockerfile", "context-image", "runtime", "env-db", "cloud"},
	Thresholds:      Thresholds{ImageSizeWarnMB: 700, ImageSizeCriticalMB: 1250, LayerWarnCount: 24, StartupMemoryWarnMB: 448},
}

var gcpCloudRun = Definition{
	Name:            "gcp-cloud-run",
	Description:     "GCP Cloud Run defaults for strict PORT listening, stateless runtime, and startup behavior",
	EnabledFamilies: []string{"dockerfile", "context-image", "runtime", "env-db", "cloud"},
	Thresholds:      Thresholds{ImageSizeWarnMB: 650, ImageSizeCriticalMB: 1200, LayerWarnCount: 23, StartupMemoryWarnMB: 384},
}

var azureContainerApps = Definition{
	Name:            "azure-container-apps",
	Description:     "Azure Container Apps defaults for ingress target port, probes, and runtime behavior",
	EnabledFamilies: []string{"dockerfile", "context-image", "runtime", "env-db", "cloud"},
	Thresholds:      Thresholds{ImageSizeWarnMB: 700, ImageSizeCriticalMB: 1250, LayerWarnCount: 24, StartupMemoryWarnMB: 448},
}

var dokku = Definition{
	Name:            "dokku",
	Description:     "Dokku defaults for reverse-proxy port behavior, restart model, and runtime expectations",
	EnabledFamilies: []string{"dockerfile", "context-image", "runtime", "env-db", "cloud"},
	Thresholds:      Thresholds{ImageSizeWarnMB: 750, ImageSizeCriticalMB: 1350, LayerWarnCount: 24, StartupMemoryWarnMB: 512},
}

var vpsSystemdDocker = Definition{
	Name:            "vps-systemd-docker",
	Description:     "VPS + systemd + Docker defaults for host-managed runtime and restart behavior",
	EnabledFamilies: []string{"dockerfile", "context-image", "runtime", "env-db", "cloud"},
	Thresholds:      Thresholds{ImageSizeWarnMB: 800, ImageSizeCriticalMB: 1400, LayerWarnCount: 25, StartupMemoryWarnMB: 512},
}

func Get(name string) (Definition, error) {
	switch name {
	case "generic":
		return inherit(base, generic), nil
	case "lightsail":
		return inherit(base, lightsail), nil
	case "render":
		return inherit(base, render), nil
	case "railway":
		return inherit(base, railway), nil
	case "flyio":
		return inherit(base, flyio), nil
	case "ecs-fargate":
		return inherit(base, ecsFargate), nil
	case "digitalocean-app-platform":
		return inherit(base, digitaloceanAppPlatform), nil
	case "gcp-cloud-run":
		return inherit(base, gcpCloudRun), nil
	case "azure-container-apps":
		return inherit(base, azureContainerApps), nil
	case "dokku":
		return inherit(base, dokku), nil
	case "vps-systemd-docker":
		return inherit(base, vpsSystemdDocker), nil
	default:
		return Definition{}, fmt.Errorf("unknown profile: %s", name)
	}
}

func List() []Definition {
	return []Definition{
		inherit(base, generic),
		inherit(base, lightsail),
		inherit(base, render),
		inherit(base, railway),
		inherit(base, flyio),
		inherit(base, ecsFargate),
		inherit(base, digitaloceanAppPlatform),
		inherit(base, gcpCloudRun),
		inherit(base, azureContainerApps),
		inherit(base, dokku),
		inherit(base, vpsSystemdDocker),
	}
}

func inherit(base Definition, override Definition) Definition {
	out := base
	out.Name = override.Name
	out.Description = override.Description
	if len(override.EnabledFamilies) > 0 {
		out.EnabledFamilies = append([]string{}, override.EnabledFamilies...)
	}
	if override.Thresholds.ImageSizeWarnMB > 0 {
		out.Thresholds.ImageSizeWarnMB = override.Thresholds.ImageSizeWarnMB
	}
	if override.Thresholds.ImageSizeCriticalMB > 0 {
		out.Thresholds.ImageSizeCriticalMB = override.Thresholds.ImageSizeCriticalMB
	}
	if override.Thresholds.LayerWarnCount > 0 {
		out.Thresholds.LayerWarnCount = override.Thresholds.LayerWarnCount
	}
	if override.Thresholds.StartupMemoryWarnMB > 0 {
		out.Thresholds.StartupMemoryWarnMB = override.Thresholds.StartupMemoryWarnMB
	}
	return out
}
