# Profile Selection Guide

- Start with `generic` if provider unknown.
- Use provider profile when platform known:
  - `render`, `railway`, `flyio`, `lightsail`, `ecs-fargate`, `digitalocean-app-platform`, `gcp-cloud-run`, `azure-container-apps`, `dokku`, `vps-systemd-docker`.
- Use auto-detection when uncertain:
```bash
deploy-doctor profiles list --recommended
deploy-doctor scan --auto-profile
```
- Rule of thumb: provider profile = higher-confidence findings for platform contracts.
