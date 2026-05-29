# Troubleshooting Playbook

## Docker unavailable
- Symptom: `doctor` shows `docker: unavailable`.
- Action: run static checks first:
```bash
deploy-doctor scan --static-only
```

## Fail-on exits CI with code 1
- Symptom: pipeline fails with findings.
- Action: lower policy (`--fail-on critical`) temporarily, fix warnings quickly.

## Profile mismatch noise
- Symptom: findings look irrelevant.
- Action: switch profile or use `--auto-profile`.

## Large image warnings
- Action: move to multi-stage builds, slim base images, remove build tools.

## DB connectivity/migration probe failures
- Action: verify `DATABASE_URL`, SSL mode, network egress, migration command.
