# Release Playbook

## Overview

This project releases with:
- `release-please` on pushes to `main` (version PR automation)
- GoReleaser on tag push (`v*`) for binary assets + Homebrew tap update

## Required GitHub Secrets

In `navig-me/deploy-doctor` repository:
- `HOMEBREW_TAP_GITHUB_TOKEN`: PAT with write access to `navig-me/homebrew-deploy-doctor`

## Normal Release Flow

1. Merge feature/fix PRs into `main`.
2. `Release Please` workflow updates/opens release PR.
3. Merge release PR.
4. Tag (for example `v0.1.1`) is created.
5. `Release` workflow runs:
   - runs tests
   - builds archives/checksums
   - publishes GitHub Release assets
   - updates Homebrew formula in tap repo

## Manual Tag Release (Fallback)

If needed:

```bash
git checkout main
git pull
git tag v0.1.1
git push origin v0.1.1
```

Then monitor:
- `.github/workflows/release.yml`
- release page in `navig-me/deploy-doctor`

## Verify Artifacts

After release:
- GitHub release contains platform tarballs + `checksums.txt`
- Homebrew tap formula points to new version + valid SHA256

Quick checks:

```bash
gh release view v0.1.1 --repo navig-me/deploy-doctor
brew update
brew reinstall deploy-doctor
deploy-doctor version
```

## Rollback / Hotfix

If release broken:
1. Patch `main`.
2. Cut next patch tag (e.g. `v0.1.2`) instead of mutating existing tag.
3. Re-run install verification.

Avoid deleting/replacing published tags unless absolutely necessary.

## Troubleshooting

- Homebrew 404:
  - confirm release asset filenames match formula URLs.
- Homebrew checksum mismatch:
  - update formula SHA256 in tap and push.
- Tap update fails in release workflow:
  - verify `HOMEBREW_TAP_GITHUB_TOKEN` permissions and repo access.
