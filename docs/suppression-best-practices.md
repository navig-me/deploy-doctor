# Ignore/Suppression Best Practices

- Prefer fixing issue before suppressing.
- Scope suppression to specific stable rule IDs.
- Keep suppression list short and reviewed.
- Add reason and owner in config review notes.
- Revisit suppressions on major deploy/runtime changes.

Example config:
```yaml
rules:
  ignore:
    - DF_HEALTH_0001
  severity_overrides:
    IMG_SIZE_0001: critical
```
