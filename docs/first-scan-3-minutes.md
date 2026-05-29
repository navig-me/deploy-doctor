# First Scan in 3 Minutes

## 1. Install

### Download binary
- Download release artifact for your OS/arch from GitHub Releases.
- Place `deploy-doctor` in your `PATH`.

### Homebrew
```bash
brew tap navig-me/deploy-doctor
brew install deploy-doctor
```

### Curl installer
```bash
curl -fsSL https://raw.githubusercontent.com/navig-me/docker-doctor/main/scripts/install.sh | bash
```

## 2. Run first scan
```bash
deploy-doctor scan --static-only
```

## 3. Optional profile-aware scan
```bash
deploy-doctor scan --auto-profile
```
