# Performance Targets Validation

## Target
- Static scan completes under 3s on medium repo.
- Memory usage remains stable without sustained growth.

## Command
```bash
time ./bin/deploy-doctor scan --static-only
```

## Sample Repo Matrix
- Go web service
- Node API
- Python Flask service

Record results in `docs/quality/perf-results.md`.
