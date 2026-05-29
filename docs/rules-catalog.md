# Rules Catalog (Issue / Impact / Fix)

## Dockerfile
- `DF_BASE_0001`: oversized/unpinned base. Impact: larger risk/surface. Fix: pin slim/alpine/distroless.
- `DF_USER_0001`: root user. Impact: privilege risk. Fix: non-root `USER`.
- `DF_HEALTH_0001`: missing healthcheck. Impact: bad health visibility. Fix: add `HEALTHCHECK`.
- `DF_CMD_0001`: missing command. Impact: startup failure risk. Fix: explicit `CMD`/`ENTRYPOINT`.
- `DF_CACHE_0001`: cache-poor layering. Impact: slow rebuilds. Fix: copy manifests before source.
- `DF_SECRET_0001`: secret-like values in build args/env. Impact: secret leak. Fix: runtime secret injection.
- `DF_APT_0001`: apt cache not cleaned. Impact: image bloat. Fix: clean apt lists same layer.
- `DF_COPY_0001`: broad copy too early. Impact: cache invalidation + context risk. Fix: narrow copy first.

## Context/Image
- `CTX_IGNR_0001`: weak/missing `.dockerignore`. Impact: noisy/sensitive context. Fix: add excludes.
- `CTX_JUNK_0001`: junk/sensitive files present. Impact: leak/bloat risk. Fix: prune + ignore.
- `IMG_SIZE_0001`: image too large. Impact: slow pushes/cold starts. Fix: slim/multi-stage.
- `IMG_LAYR_0001`: too many layers. Impact: inefficiency. Fix: consolidate where safe.
- `IMG_ARCH_0001`: arch mismatch. Impact: hidden runtime issues. Fix: test/build target arch.
- `IMG_BUILD_0001`: build tools in runtime. Impact: attack surface/bloat. Fix: multi-stage runtime-only.

## Runtime
- `RT_BOOT_0001`: startup failed. Impact: deployment breakage. Fix: verify entrypoint/env.
- `RT_PORT_0001`: expected port not bound. Impact: unreachable service. Fix: bind platform port.
- `RT_BIND_0001`: localhost-only bind. Impact: unreachable externally. Fix: bind `0.0.0.0`.
- `RT_HEAL_0001`: health probe failed. Impact: restarts/unhealthy deploy. Fix: stable health path.
- `RT_EXIT_0001`: exits early. Impact: crash-loop. Fix: long-lived foreground process.
- `RT_LOG_0001`: no stdout/stderr logs. Impact: poor observability. Fix: log to stdout/stderr.
- `RT_SIGT_0001`: poor SIGTERM handling. Impact: unsafe shutdown. Fix: graceful signal handling.
- `RT_MEM_0001`: high startup memory. Impact: OOM/cold-start risk. Fix: reduce startup footprint.

## Env/DB
- `ENV_MISS_0001`: missing required env. Impact: startup failure. Fix: define required vars.
- `ENV_INSEC_0001`: unsafe env values. Impact: insecure runtime. Fix: production-safe values.
- `ENV_DB_0001`: DB URL SSL lint fail. Impact: managed DB compatibility issues. Fix: explicit SSL mode.
- `ENV_HOST_0001`: localhost misuse. Impact: cloud connectivity failure. Fix: service hostname.
- `DB_CONN_0001`: active DB probe failed. Impact: runtime DB outage risk. Fix: network/auth/SSL checks.
- `DB_MIGR_0001`: migration probe failed. Impact: schema drift risk. Fix: repair migration command.
