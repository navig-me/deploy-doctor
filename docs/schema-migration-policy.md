# Schema Migration Policy
- Additive fields: keep same major schema version.
- Breaking changes: bump schema version and keep backward fixture tests.
- JSON output must include `schema_version`.
