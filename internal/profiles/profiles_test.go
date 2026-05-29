package profiles

import "testing"

func TestGenericProfileDefinition(t *testing.T) {
	t.Parallel()
	p, err := Get("generic")
	if err != nil || p.Name != "generic" { t.Fatalf("bad generic: %+v %v", p, err) }
}

func TestLightsailProfileDefinition(t *testing.T) {
	t.Parallel()
	p, err := Get("lightsail")
	if err != nil { t.Fatalf("expected lightsail: %v", err) }
	g, _ := Get("generic")
	if p.Thresholds.ImageSizeWarnMB >= g.Thresholds.ImageSizeWarnMB { t.Fatalf("lightsail should be stricter") }
}

func TestRenderProfileDefinition(t *testing.T) {
	t.Parallel()
	p, err := Get("render")
	if err != nil { t.Fatalf("expected render profile: %v", err) }
	if p.Name != "render" { t.Fatalf("unexpected name: %s", p.Name) }
	g, _ := Get("generic")
	if p.Thresholds.ImageSizeWarnMB > g.Thresholds.ImageSizeWarnMB { t.Fatalf("render should be <= generic threshold") }
}

func TestRailwayProfileDefinition(t *testing.T) {
	t.Parallel()
	p, err := Get("railway")
	if err != nil { t.Fatalf("expected railway profile: %v", err) }
	if p.Name != "railway" { t.Fatalf("unexpected name: %s", p.Name) }
	if p.Thresholds.ImageSizeCriticalMB > render.Thresholds.ImageSizeCriticalMB { t.Fatalf("railway should be <= render critical threshold") }
}

func TestFlyioProfileDefinition(t *testing.T) {
	t.Parallel()
	p, err := Get("flyio")
	if err != nil { t.Fatalf("expected flyio profile: %v", err) }
	if p.Name != "flyio" { t.Fatalf("unexpected name: %s", p.Name) }
	if p.Thresholds.ImageSizeWarnMB > render.Thresholds.ImageSizeWarnMB { t.Fatalf("flyio should be <= render warn threshold") }
}

func TestECSFargateProfileDefinition(t *testing.T) {
	t.Parallel()
	p, err := Get("ecs-fargate")
	if err != nil { t.Fatalf("expected ecs-fargate profile: %v", err) }
	if p.Name != "ecs-fargate" { t.Fatalf("unexpected name: %s", p.Name) }
	g, _ := Get("generic")
	if p.Thresholds.ImageSizeCriticalMB > g.Thresholds.ImageSizeCriticalMB { t.Fatalf("ecs-fargate should be <= generic critical threshold") }
}

func TestDigitalOceanAppPlatformProfileDefinition(t *testing.T) {
	t.Parallel()
	p, err := Get("digitalocean-app-platform")
	if err != nil { t.Fatalf("expected digitalocean-app-platform profile: %v", err) }
	if p.Name != "digitalocean-app-platform" { t.Fatalf("unexpected name: %s", p.Name) }
	g, _ := Get("generic")
	if p.Thresholds.ImageSizeCriticalMB > g.Thresholds.ImageSizeCriticalMB { t.Fatalf("digitalocean-app-platform should be <= generic critical threshold") }
}

func TestGCPCloudRunProfileDefinition(t *testing.T) {
	t.Parallel()
	p, err := Get("gcp-cloud-run")
	if err != nil { t.Fatalf("expected gcp-cloud-run profile: %v", err) }
	if p.Name != "gcp-cloud-run" { t.Fatalf("unexpected name: %s", p.Name) }
	if p.Thresholds.ImageSizeWarnMB > render.Thresholds.ImageSizeWarnMB { t.Fatalf("gcp-cloud-run should be <= render warn threshold") }
}

func TestAzureContainerAppsProfileDefinition(t *testing.T) {
	t.Parallel()
	p, err := Get("azure-container-apps")
	if err != nil { t.Fatalf("expected azure-container-apps profile: %v", err) }
	if p.Name != "azure-container-apps" { t.Fatalf("unexpected name: %s", p.Name) }
	g, _ := Get("generic")
	if p.Thresholds.ImageSizeCriticalMB > g.Thresholds.ImageSizeCriticalMB { t.Fatalf("azure-container-apps should be <= generic critical threshold") }
}

func TestDokkuProfileDefinition(t *testing.T) {
	t.Parallel()
	p, err := Get("dokku")
	if err != nil { t.Fatalf("expected dokku profile: %v", err) }
	if p.Name != "dokku" { t.Fatalf("unexpected name: %s", p.Name) }
	if p.Thresholds.ImageSizeCriticalMB < render.Thresholds.ImageSizeCriticalMB { t.Fatalf("dokku expected to allow >= render critical threshold") }
}

func TestVPSSystemdDockerProfileDefinition(t *testing.T) {
	t.Parallel()
	p, err := Get("vps-systemd-docker")
	if err != nil { t.Fatalf("expected vps-systemd-docker profile: %v", err) }
	if p.Name != "vps-systemd-docker" { t.Fatalf("unexpected name: %s", p.Name) }
	if p.Thresholds.ImageSizeCriticalMB < dokku.Thresholds.ImageSizeCriticalMB { t.Fatalf("vps-systemd-docker expected to allow >= dokku critical threshold") }
}

func TestListStableOrder(t *testing.T) {
	t.Parallel()
	all := List()
	if len(all) < 11 { t.Fatalf("expected 11 profiles") }
	if all[0].Name != "generic" || all[1].Name != "lightsail" || all[2].Name != "render" || all[3].Name != "railway" || all[4].Name != "flyio" || all[5].Name != "ecs-fargate" || all[6].Name != "digitalocean-app-platform" || all[7].Name != "gcp-cloud-run" || all[8].Name != "azure-container-apps" || all[9].Name != "dokku" || all[10].Name != "vps-systemd-docker" {
		t.Fatalf("unexpected list order: %+v", all)
	}
}

func TestProfileInheritanceBasePlusOverride(t *testing.T) {
	t.Parallel()

	p, err := Get("generic")
	if err != nil { t.Fatalf("get generic: %v", err) }
	if p.Thresholds.ImageSizeWarnMB != 800 || p.Thresholds.ImageSizeCriticalMB != 1500 {
		t.Fatalf("generic should inherit base thresholds: %+v", p.Thresholds)
	}
	if len(p.EnabledFamilies) == 0 {
		t.Fatalf("generic should inherit base families")
	}

	r, err := Get("render")
	if err != nil { t.Fatalf("get render: %v", err) }
	if r.Thresholds.ImageSizeWarnMB == p.Thresholds.ImageSizeWarnMB {
		t.Fatalf("render should override base warn threshold")
	}
	if len(r.EnabledFamilies) == 0 {
		t.Fatalf("render families should be present")
	}
}
