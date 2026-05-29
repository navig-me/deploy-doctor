# Deploy Doctor

`deploy-doctor` is a Go CLI for diagnosing container deploy readiness.

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
```

## Quickstart

```bash
deploy-doctor doctor
deploy-doctor scan --static-only
deploy-doctor scan --auto-profile
```

## Documentation

- [First Scan in 3 Minutes](docs/first-scan-3-minutes.md)
- [Rules Catalog](docs/rules-catalog.md)
- [Profile Selection Guide](docs/profile-selection-guide.md)
- [Ignore/Suppression Best Practices](docs/suppression-best-practices.md)
- [Troubleshooting Playbook](docs/troubleshooting-playbook.md)
