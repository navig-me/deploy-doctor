# Deploy Doctor

`deploy-doctor` is a command-line tool that evaluates containerized applications for deployment readiness across common platforms.

It runs static checks (Dockerfile, build context, env/DB heuristics), optional runtime probes, profile-aware thresholds, and outputs machine-readable summaries for CI.

## Install

### Homebrew (recommended)
```bash
brew tap navig-me/deploy-doctor
brew install deploy-doctor
```

### Curl installer
```bash
curl -fsSL https://raw.githubusercontent.com/navig-me/deploy-doctor/main/scripts/install.sh | bash
```

### Curl installer (pinned version)
```bash
VERSION=v0.1.0 curl -fsSL https://raw.githubusercontent.com/navig-me/deploy-doctor/main/scripts/install.sh | bash
```

### Build from source
```bash
make build
./bin/deploy-doctor version
```

## Command List

- `deploy-doctor scan`: run deployment checks in the current repository.
- `deploy-doctor doctor`: verify local prerequisites (for example Docker availability).
- `deploy-doctor profiles list`: list supported deployment profiles.
- `deploy-doctor profiles list --recommended`: show profile recommendations from repository signals.
- `deploy-doctor profiles explain <profile>`: print profile details and thresholds.
- `deploy-doctor version`: print installed CLI version.

## How `scan` Works

`scan` executes in phases:
1. Determine profile (`--profile` or `--auto-profile`).
2. Run static checks:
   - Dockerfile checks
   - context/image metadata checks
   - env/DB heuristics
   - provider config signals (evidence-driven confidence)
3. Optionally run runtime probes (unless `--static-only`).
4. Aggregate findings and compute:
   - severity counts
   - score
   - status
5. Emit:
   - human-readable summary
   - machine-readable CI line (`DD_SUMMARY ...`)

## `scan` Flags

- `--profile <name>`: force profile (default: `generic`).
- `--auto-profile`: detect likely provider profile from repo files.
- `--static-only`: skip runtime probes.
- `--runtime=false`: disable runtime probes explicitly.
- `--fail-on <none|critical|warning|suggestion>`: fail by severity policy.
- `--strict`: treat warnings as failures.
- `--verbose` / `-v`: print detailed execution and finding evidence.

## Output and Exit Behavior

- Human-readable summary line:
  - `Scan summary: score=... status=... profile=... findings=(...)`
- Machine-readable CI summary line:
  - `DD_SUMMARY score=... status=... critical=... warning=... suggestion=... fail=... profile=...`

Exit semantics:
- `0`: scan completed without failure policy trigger.
- `1`: findings violated `--fail-on` or `--strict` policy.
- `2`: tool/runtime error.

## Profiles

Examples:
- `generic`
- `lightsail`
- `render`
- `railway`
- `flyio`
- `ecs-fargate`
- `digitalocean-app-platform`
- `gcp-cloud-run`
- `azure-container-apps`
- `dokku`
- `vps-systemd-docker`

Use:
```bash
deploy-doctor profiles list
deploy-doctor profiles explain lightsail
deploy-doctor scan --profile lightsail --static-only
```

## Quickstart

```bash
deploy-doctor doctor
deploy-doctor scan --static-only
deploy-doctor scan --auto-profile -v
```

## CI Usage Example

```bash
deploy-doctor scan --fail-on warning --static-only
```

Parse `DD_SUMMARY` in CI logs for machine-readable gating or metrics.

## Documentation

- [First Scan in 3 Minutes](docs/first-scan-3-minutes.md)
- [Rules Catalog](docs/rules-catalog.md)
- [Profile Selection Guide](docs/profile-selection-guide.md)
- [Ignore/Suppression Best Practices](docs/suppression-best-practices.md)
- [Troubleshooting Playbook](docs/troubleshooting-playbook.md)
- [Release Playbook](docs/release-playbook.md)
